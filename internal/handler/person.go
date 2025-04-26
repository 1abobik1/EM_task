package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/1abobik1/EM_task/internal/dto"
	"github.com/1abobik1/EM_task/internal/models"
	"github.com/1abobik1/EM_task/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type PersonService interface {
	EnrichPersonInfo(ctx context.Context, person *models.Person) error
	SavePerson(ctx context.Context, person *models.Person) error
	DeletePerson(ctx context.Context, id int) error
	UpdatePerson(ctx context.Context, req dto.UpdatePersonRequest, id int) (dto.PersonResponse, error)
	ListPersons(ctx context.Context, req dto.ListPersonsRequest) (dto.ListPersonsResponse, error)
	GetPersonByID(ctx context.Context, id int) (dto.PersonResponse, error)
}

type PersonHandler struct {
	personService PersonService
}

func NewPersonHandler(personService PersonService) *PersonHandler {
	return &PersonHandler{personService: personService}
}

// @Summary     Add person
// @Description Создаёт нового человека, обогащает информацию с помощью внешнего API и сохраняет в PostgreSQL
// @Tags        persons
// @Accept      json
// @Produce     json
// @Param       person  body     dto.CreatePersonRequest    true  "Параметры для создания"
// @Success     201     {object} dto.PersonResponse        "Сущность человека"
// @Failure     400     {object} dto.BadRequest            "Некорректный запрос"
// @Failure     500     {object} dto.InternalServerError   "Внутренняя ошибка"
// @Failure     503     {object} dto.ServiceUnavailable    "Сервис недоступен"
// @Router      /persons [post]
func (h *PersonHandler) CreatePerson(c *gin.Context) {
	var req dto.CreatePersonRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		handleBindError(c, err)
		return
	}

	domainPerson := models.Person{
		Name:       req.Name,
		Surname:    req.Surname,
		Patronymic: req.Patronymic,
	}

	if err := h.personService.EnrichPersonInfo(c, &domainPerson); err != nil {
		logrus.WithError(err).Error("failed to enrich person data")
		c.JSON(http.StatusServiceUnavailable, dto.ServiceUnavailable{Error: "service unavailable"})
		return
	}

	if err := h.personService.SavePerson(c, &domainPerson); err != nil {
		logrus.WithError(err).Error("failed to save person")
		c.JSON(http.StatusInternalServerError, dto.InternalServerError{Error: "internal server error"})
		return
	}

	response := dto.PersonResponse{
		ID:          domainPerson.ID,
		Name:        domainPerson.Name,
		Surname:     domainPerson.Surname,
		Patronymic:  domainPerson.Patronymic,
		Age:         domainPerson.Age,
		Gender:      domainPerson.Gender,
		Nationality: domainPerson.Nationality,
		CreatedAt:   domainPerson.CreatedAt,
		UpdatedAt:   domainPerson.CreatedAt,
	}

	c.JSON(http.StatusCreated, response)
}

// @Summary     List persons
// @Description Получение списка людей при помощи фильтров и пагинацией. Примечание name и surname работает через поиск подстроки
// @Tags        persons
// @Accept      json
// @Produce     json
// @Param       name        query     string  false  "Фильтр по имени"
// @Param       surname     query     string  false  "Фильтр по фамилии"
// @Param       age         query     integer false  "Фильтр по возрасту"
// @Param       gender      query     string  false  "Фильтр по полу"
// @Param       nationality query     string  false  "Фильтр по национальности"
// @Param       page        query     integer false  "Номер страницы"  default(1)
// @Param       limit       query     integer false  "Размер страницы"  default(10)
// @Success     200         {object}  dto.ListPersonsResponse
// @Failure     400         {object}  dto.BadRequest
// @Failure     404         {object}  dto.NotFound
// @Failure     500         {object}  dto.InternalServerError
// @Router      /persons [get]
func (h *PersonHandler) ListPersons(c *gin.Context) {
	var req dto.ListPersonsRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		handleBindError(c, err)
		return
	}

	resp, err := h.personService.ListPersons(c, req)
	if err != nil {
		if errors.Is(err, repository.ErrPersonNotFound) {
			logrus.Error(err)
			c.JSON(http.StatusNotFound, dto.NotFound{Error: "person does not exist according to these filters"})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.InternalServerError{Error: "internal server error"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Summary     Get person by ID
// @Description Возвращает информацию о человеке по его уникальному идентификатору
// @Tags        persons
// @Accept      json
// @Produce     json
// @Param       id   path      int  true  "Идентификатор человека"
// @Success     200  {object}  dto.PersonResponse           "Найдённая сущность Person"
// @Failure     400  {object}  dto.BadRequest               "Неверный формат ID"
// @Failure     404  {object}  dto.NotFound                 "Person не найден"
// @Failure     500  {object}  dto.InternalServerError      "Внутренняя ошибка сервера"
// @Router      /persons/{id} [get]
func (h *PersonHandler) GetPersonByID(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.WithError(err).Warnf("invalid id param, id: %s", idStr)
		c.JSON(http.StatusBadRequest, dto.BadRequest{Error: "invalid id " + idStr})
		return
	}

	resp, err := h.personService.GetPersonByID(c, id)
	if err != nil {
		handlePersonNotFoundError(c, idStr, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Summary     Update person
// @Description Обновляет информацию о человеке (любые поля: name, surname, patronymic, age, gender, nationality)
// @Tags        persons
// @Accept      json
// @Produce     json
// @Param       id      path      int                       true  "Идентификатор человека"
// @Param       person  body      dto.UpdatePersonRequest   true  "Поля для обновления"
// @Success     200     {object} dto.PersonResponse        "Обновлённая сущность человека"
// @Failure     400     {object} dto.BadRequest            "Некорректный запрос или параметры"
// @Failure     404     {object} dto.NotFound              "Человек не найден"
// @Failure     500     {object} dto.InternalServerError   "Внутренняя ошибка"
// @Router      /persons/{id} [put]
// @Router      /persons/{id} [patch]
func (h *PersonHandler) UpdatePerson(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.WithError(err).Warnf("invalid id param, id: %s", idStr)
		c.JSON(http.StatusBadRequest, dto.BadRequest{Error: "invalid id " + idStr})
		return
	}

	var req dto.UpdatePersonRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		handleBindError(c, err)
		return
	}

	resp, err := h.personService.UpdatePerson(c, req, id)
	if err != nil {
		handlePersonNotFoundError(c, idStr, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Summary     Delete person
// @Description Удаляет человека по идентификатору
// @Tags        persons
// @Param       id   path      int  true  "Person ID"
// @Success     204   "No Content"
// @Failure     400  {object}  dto.BadRequest
// @Failure     404  {object}  dto.NotFound
// @Failure     500  {object}  dto.InternalServerError
// @Router      /persons/{id} [delete]
func (h *PersonHandler) DeletePerson(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.WithError(err).Warnf("invalid id param, id: %s", idStr)
		c.JSON(http.StatusBadRequest, dto.BadRequest{Error: "invalid id " + idStr})
		return
	}

	if err := h.personService.DeletePerson(c, id); err != nil {
		handlePersonNotFoundError(c, idStr, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func handleBindError(c *gin.Context, err error) {

	if verrs, ok := err.(validator.ValidationErrors); ok {
		out := make(map[string]string, len(verrs))

		for _, fe := range verrs {
			out[fe.Field()] = fmt.Sprintf("must satisfy %s", fe.Tag())
		}

		logrus.WithError(err).Warn(out)
		c.JSON(http.StatusBadRequest, gin.H{"errors": out})
		return
	}

	logrus.WithError(err).Warn("invalid request data")
	c.JSON(http.StatusBadRequest, dto.BadRequest{Error: "invalid request data"})
}

func handlePersonNotFoundError(c *gin.Context, idStr string, err error) {

	if errors.Is(err, repository.ErrPersonNotFound) {
		logrus.WithError(err).Warnf(" id: %s", idStr)
		c.JSON(http.StatusNotFound, dto.NotFound{Error: "person was not found by ID " + idStr})
		return
	}

	logrus.Error(err)
	c.JSON(http.StatusInternalServerError, dto.InternalServerError{Error: "internal server error"})
}

package service

import (
	"context"
	"time"

	"github.com/1abobik1/EM_task/internal/dto"
	"github.com/1abobik1/EM_task/internal/models"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type PostgresRepo interface {
	SavePerson(ctx context.Context, p *models.Person) error
	DeletePerson(ctx context.Context, id int) error
	UpdatePerson(ctx context.Context, p *models.Person) error
	GetPersonByID(ctx context.Context, id int) (models.Person, error)
	ListPersons(ctx context.Context, filter models.PersonFilter) ([]models.Person, int, error)
}

type AgifyClient interface {
	GetAge(ctx context.Context, name string) (int, error)
}

type GenderizeClient interface {
	GetGender(ctx context.Context, name string) (string, error)
}

type NationalizeClient interface {
	GetNationality(ctx context.Context, name string) (string, error)
}

type PersonService struct {
	postgresRepo PostgresRepo
	ageClient    AgifyClient
	genderClient GenderizeClient
	natClient    NationalizeClient
}

func NewPersonService(postgresRepo PostgresRepo, ageClient AgifyClient, genderClient GenderizeClient, natClient NationalizeClient) *PersonService {

	return &PersonService{postgresRepo: postgresRepo, ageClient: ageClient,
		genderClient: genderClient, natClient: natClient}
}

func (s *PersonService) EnrichPersonInfo(ctx context.Context, person *models.Person) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var (
		eg  errgroup.Group
		age int
		gen string
		nat string
	)

	eg.Go(func() error {
		var err error
		age, err = s.ageClient.GetAge(ctx, person.Name)
		return err
	})

	eg.Go(func() error {
		var err error
		gen, err = s.genderClient.GetGender(ctx, person.Name)
		return err
	})

	eg.Go(func() error {
		var err error
		nat, err = s.natClient.GetNationality(ctx, person.Name)
		return err
	})

	if err := eg.Wait(); err != nil {
		return err
	}

	logrus.Infof("the information collection was successfully completed: age:%v, gender:%v, nationality: %v", age, gen, nat)

	person.Age = age
	person.Gender = gen
	person.Nationality = nat

	return nil
}

func (s *PersonService) SavePerson(ctx context.Context, person *models.Person) error {
	return s.postgresRepo.SavePerson(ctx, person)
}

func (s *PersonService) DeletePerson(ctx context.Context, id int) error {
	return s.postgresRepo.DeletePerson(ctx, id)
}

func (s *PersonService) UpdatePerson(ctx context.Context, req dto.UpdatePersonRequest, id int) (dto.PersonResponse, error) {
	existing, err := s.postgresRepo.GetPersonByID(ctx, id)
	if err != nil {
		return dto.PersonResponse{}, err
	}

	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.Surname != nil {
		existing.Surname = *req.Surname
	}
	if req.Patronymic != nil {
		existing.Patronymic = *req.Patronymic
	}
	if req.Age != nil {
		existing.Age = *req.Age
	}
	if req.Gender != nil {
		existing.Gender = *req.Gender
	}
	if req.Nationality != nil {
		existing.Nationality = *req.Nationality
	}

	if err := s.postgresRepo.UpdatePerson(ctx, &existing); err != nil {
		return dto.PersonResponse{}, err
	}

	resp := dto.PersonResponse{
		ID:          existing.ID,
		Name:        existing.Name,
		Surname:     existing.Surname,
		Patronymic:  existing.Patronymic,
		Age:         existing.Age,
		Gender:      existing.Gender,
		Nationality: existing.Nationality,
		CreatedAt:   existing.CreatedAt,
		UpdatedAt:   existing.UpdatedAt,
	}

	return resp, nil
}

func (s *PersonService) GetPersonByID(ctx context.Context, id int) (dto.PersonResponse, error) {
	person, err := s.postgresRepo.GetPersonByID(ctx, id)
	if err != nil {
		return dto.PersonResponse{}, err
	}

	resp := dto.PersonResponse{
		ID:          person.ID,
		Name:        person.Name,
		Surname:     person.Surname,
		Patronymic:  person.Patronymic,
		Age:         person.Age,
		Gender:      person.Gender,
		Nationality: person.Nationality,
		CreatedAt:   person.CreatedAt,
		UpdatedAt:   person.UpdatedAt,
	}

	return resp, nil
}

func (s *PersonService) ListPersons(ctx context.Context, req dto.ListPersonsRequest) (dto.ListPersonsResponse, error) {
	filter := models.PersonFilter{
		Name:        req.Name,
		Surname:     req.Surname,
		Age:         req.Age,
		Gender:      req.Gender,
		Nationality: req.Nationality,
		Page:        req.Page,
		Limit:       req.Limit,
	}

	rows, total, err := s.postgresRepo.ListPersons(ctx, filter)
	if err != nil {
		return dto.ListPersonsResponse{}, err
	}

	persons := make([]dto.PersonResponse, len(rows))
	for i, p := range rows {
		persons[i] = dto.PersonResponse{
			ID:          p.ID,
			Name:        p.Name,
			Surname:     p.Surname,
			Patronymic:  p.Patronymic,
			Age:         p.Age,
			Gender:      p.Gender,
			Nationality: p.Nationality,
			CreatedAt:   p.CreatedAt,
			UpdatedAt:   p.UpdatedAt,
		}
	}

	listPersonRespons := dto.ListPersonsResponse{
		Persons: persons,
		Pagination: dto.PaginationInfo{
			Page:  req.Page,
			Limit: req.Limit,
			Total: total,
		},
	}

	return listPersonRespons, nil
}

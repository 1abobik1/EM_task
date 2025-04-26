package dto

type CreatePersonRequest struct {
	Name       string `json:"name" binding:"required,min=2,max=50"`
	Surname    string `json:"surname" binding:"required,min=2,max=50"`
	Patronymic string `json:"patronymic,omitempty" binding:"max=50"`
}

type UpdatePersonRequest struct {
    Name        *string `json:"name,omitempty" binding:"omitempty,min=2,max=50"`
    Surname     *string `json:"surname,omitempty" binding:"omitempty,min=2,max=50"`
    Patronymic  *string `json:"patronymic,omitempty" binding:"omitempty,max=50"`
    Age         *int    `json:"age,omitempty" binding:"omitempty,gte=0,lte=110"`
    Gender      *string `json:"gender,omitempty" binding:"omitempty,oneof=male female"`
    Nationality *string `json:"nationality,omitempty" binding:"omitempty,len=2,alpha"`
}

// ListPersonsRequest содержит параметры фильтрации и пагинации для списка людей
// swagger:model
// @description Параметры запроса для получения списка Person
// @param name query string false "Фильтр по имени"
// @param surname query string false "Фильтр по фамилии"
// @param age query integer false "Фильтр по возрасту"
// @param gender query string false "Фильтр по полу (male/female)"
// @param nationality query string false "Фильтр по национальности"
// @param page query integer false "Номер страницы, начинается с 1" default(1)
// @param limit query integer false "Размер страницы" default(10)
type ListPersonsRequest struct {
    Name        string `form:"name" binding:"omitempty,min=2,max=50"`
    Surname     string `form:"surname" binding:"omitempty,min=2,max=50"`
    Age         int   `form:"age" binding:"omitempty,gte=0,lte=110"`
    Gender      string `form:"gender" binding:"omitempty,oneof=male female"`
    Nationality string `form:"nationality" binding:"omitempty,len=2,alpha"`
    Page        int    `form:"page,default=1" binding:"gte=1"`
    Limit       int    `form:"limit,default=10" binding:"gte=1,lte=100"`
}
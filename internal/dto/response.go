package dto

import "time"

// PersonResponse основная информация о пользователе
// swagger:model
// @description основная информация о пользователе в json
type PersonResponse struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Surname     string    `json:"surname"`
	Patronymic  string    `json:"patronymic,omitempty"`
	Age         int       `json:"age"`
	Gender      string    `json:"gender"`
	Nationality string    `json:"nationality,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

// PaginationInfo метаинформация о пагинации
// swagger:model
// @description Метаданные пагинации
type PaginationInfo struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}

// ListPersonsResponse ответ с массивом PersonResponse и пагинацией
// swagger:model
// @description Ответ списка Person
type ListPersonsResponse struct {
	Persons    []PersonResponse `json:"persons"`
	Pagination PaginationInfo   `json:"pagination"`
}

package persistenceentities

import "time"

type Person struct {
	ID          int
	Name        string
	Surname     string
	Patronymic  string
	Age         int
	Gender      string
	Nationality string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
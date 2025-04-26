package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/1abobik1/EM_task/internal/models"
	"github.com/1abobik1/EM_task/internal/repository"
	_ "github.com/lib/pq"

	sq "github.com/Masterminds/squirrel"
)

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorageProd(storagePath string) (*PostgresStorage, error) {
	db, err := sql.Open("postgres", storagePath)
	if err != nil {
		return nil, err
	}

	return &PostgresStorage{db: db}, nil
}

func (ps *PostgresStorage) SavePerson(ctx context.Context, p *models.Person) error {
	const q = `
    INSERT INTO persons (name, surname, patronymic, age, gender, nationality)
    VALUES ($1,$2,$3,$4,$5,$6)
    RETURNING id, created_at
    `
	return ps.db.QueryRowContext(ctx, q,
		p.Name,
		p.Surname,
		p.Patronymic,
		p.Age,
		p.Gender,
		p.Nationality,
	).Scan(&p.ID, &p.CreatedAt)
}

func (ps *PostgresStorage) DeletePerson(ctx context.Context, id int) error {
	const q = `DELETE FROM persons WHERE id = $1`
	res, err := ps.db.ExecContext(ctx, q, id)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return repository.ErrPersonNotFound
	}
	return nil
}

func (ps *PostgresStorage) UpdatePerson(ctx context.Context, p *models.Person) error {
	const q = `
    UPDATE persons
    SET name        = $1,
        surname     = $2,
        patronymic  = $3,
        age         = $4,
        gender      = $5,
        nationality = $6,
        updated_at  = NOW()
    WHERE id = $7
    RETURNING created_at, updated_at
    `
	err := ps.db.QueryRowContext(ctx, q,
		p.Name,
		p.Surname,
		p.Patronymic,
		p.Age,
		p.Gender,
		p.Nationality,
		p.ID,
	).Scan(&p.CreatedAt, &p.UpdatedAt)

	if err == sql.ErrNoRows {
		return repository.ErrPersonNotFound
	}

	return err
}

func (ps *PostgresStorage) GetPersonByID(ctx context.Context, id int) (models.Person, error) {
	const q = `
    SELECT id, name, surname, patronymic, age, gender, nationality, created_at, updated_at
    FROM persons
    WHERE id = $1
    `
	p := models.Person{}
	err := ps.db.QueryRowContext(ctx, q, id).Scan(
		&p.ID,
		&p.Name,
		&p.Surname,
		&p.Patronymic,
		&p.Age,
		&p.Gender,
		&p.Nationality,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Person{}, repository.ErrPersonNotFound
		}
		return models.Person{}, err
	}
	return p, nil
}

// ListPersons возвращает срез Person и общее количество по фильтрам и пагинации
func (ps *PostgresStorage) ListPersons(ctx context.Context, filter models.PersonFilter) ([]models.Person, int, error) {

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	// базовый селект для получения данных
	query := builder.Select(
		"id", "name", "surname", "patronymic", "age", "gender", "nationality", "created_at", "updated_at",
	).From("persons")

	// селект для подсчёта total
	countQuery := builder.Select("COUNT(*)").From("persons")

	// фильтры
	if filter.Name != "" {
		query = query.Where(sq.ILike{"name": fmt.Sprintf("%%%s%%", filter.Name)})
		countQuery = countQuery.Where(sq.ILike{"name": fmt.Sprintf("%%%s%%", filter.Name)})
	}
	if filter.Surname != "" {
		query = query.Where(sq.ILike{"surname": fmt.Sprintf("%%%s%%", filter.Surname)})
		countQuery = countQuery.Where(sq.ILike{"surname": fmt.Sprintf("%%%s%%", filter.Surname)})
	}
	if filter.Age != 0 {
		query = query.Where(sq.Eq{"age": filter.Age})
		countQuery = countQuery.Where(sq.Eq{"age": filter.Age})
	}
	if filter.Gender != "" {
		query = query.Where(sq.Eq{"gender": filter.Gender})
		countQuery = countQuery.Where(sq.Eq{"gender": filter.Gender})
	}
	if filter.Nationality != "" {
		query = query.Where(sq.Eq{"nationality": filter.Nationality})
		countQuery = countQuery.Where(sq.Eq{"nationality": filter.Nationality})
	}

	// считываем total
	countSql, countArgs, err := countQuery.ToSql()
	if err != nil {
		return nil, 0, err
	}
	var total int
	if err := ps.db.QueryRowContext(ctx, countSql, countArgs...).Scan(&total); err != nil {
		return nil, 0, err
	}
 	if total == 0 {
		return nil, 0, repository.ErrPersonNotFound
	}
	
	// пагинация LIMIT и OFFSET
	query = query.Limit(uint64(filter.Limit)).Offset(uint64(filter.Offset()))
	query = query.OrderBy("id")

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, 0, err
	}

	rows, err := ps.db.QueryContext(ctx, sqlStr, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var persons []models.Person
	for rows.Next() {
		var p models.Person
		if err := rows.Scan(
			&p.ID, &p.Name, &p.Surname, &p.Patronymic,
			&p.Age, &p.Gender, &p.Nationality,
			&p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		persons = append(persons, p)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return persons, total, nil
}

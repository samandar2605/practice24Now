package postgres

import (
	"github.com/jmoiron/sqlx"
	"github.com/practice2311/storage/repo"

	_ "github.com/lib/pq"
)

type userRepo struct {
	DB *sqlx.DB
}

func NewUser(db *sqlx.DB) repo.UserStorageI {
	return &userRepo{
		DB: db,
	}
}

func (ur *userRepo) Create(user *repo.User) (*repo.User, error) {
	query :=`
		INSERT INTO users(
			name,
			email,
			password
		) VALUES($1, $2, $3)
		RETURNING id
	`

	row := ur.DB.QueryRow(
		query,
		user.Name,
		user.Email,
		user.Password,
	)

	if err := row.Scan(&user.Id); err != nil {
		return nil, err
	}

	return user, nil
}

func (ur *userRepo) GetByEmail(email string) (*repo.User, error) {
	var result repo.User
	query := `
		select 
			id,
			name,
			email,
			password
		from users where email=$1
	`

	row := ur.DB.QueryRow(query, email)
	if err := row.Scan(
		&result.Id,
		&result.Name,
		&result.Email,
		&result.Password,
	); err != nil {
		return nil, err
	}

	return &result, nil
}

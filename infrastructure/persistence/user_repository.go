package persistence

import (
	"api_blog/domain/entities"
	"database/sql"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (r *UserRepository) FindByID(id int) *entities.User {
	query := `
		select id, name, email
		from users
		where id = $1
	`
	var user entities.User
	err := r.DB.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		return nil
	}
	return &user
}

func (r *UserRepository) FindByEmail(email string) *entities.User {
	query := `
		select id, name, email, password
		from users
		where email = $1
	`
	var user entities.User
	err := r.DB.QueryRow(query, email).Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return nil
	}
	return &user
}

func (r *UserRepository) Create(user entities.User) (*entities.User, error) {
	query := `
		insert into users (name, email, password, created_at, updated_at)
		values ($1, $2, $3, now(), now())
		returning id, name, email
	`
	var newUser entities.User
	err := r.DB.QueryRow(query, user.Name, user.Email, user.Password).Scan(&newUser.ID, &newUser.Name, &newUser.Email)
	if err != nil {
		return nil, err
	}
	return &newUser, nil
}

package repository

import (
	"database/sql"
	"jwt-auth/model"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(username, email, password string) (*model.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{}
	query := `INSERT INTO users (username, email, password, created_at) 
			  VALUES ($1, $2, $3, $4) 
			  RETURNING id, username, email, created_at`

	err = r.db.QueryRow(query, username, email, string(hashedPassword), time.Now()).
		Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt)

	return user, err
}

func (r *UserRepository) GetByEmail(email string) (*model.User, error) {
	user := &model.User{}
	query := `SELECT id, username, email, password, created_at FROM users WHERE email = $1`

	err := r.db.QueryRow(query, email).
		Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return user, err
}

func (r *UserRepository) GetByID(id int) (*model.User, error) {
	user := &model.User{}
	query := `SELECT id, username, email, password, created_at FROM users WHERE id = $1`

	err := r.db.QueryRow(query, id).
		Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return user, err
}

func (r *UserRepository) ValidatePassword(user *model.User, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) == nil
}
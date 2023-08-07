package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const dbTimeout = time.Second * 3

var db *sql.DB

func New(poolDB *sql.DB) Models {
	db = poolDB

	return Models{
		User: User{},
	}
}

type Models struct {
	User User
}

type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty"`
	Password  string    `json:"-"`
	Active    int       `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (usr *User) GetAll() ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `SELECT id, email, first_name, last_name, password, user_active, created_at, updated_at
	FROM users ORDER by last_name`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User

	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.Password,
			&user.Active,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			log.Println("error user GetAll: ", err)
			return nil, err
		}

		users = append(users, &user)
	}

	return users, nil
}

func (usr *User) GetByEmail(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `SELECT id, email, first_name, last_name, password, user_active, created_at, updated_at
	FROM users WHERE email = $1`

	var user User
	row := db.QueryRowContext(ctx, query, email)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		log.Println("error user GetByEmail: ", err)
		return nil, err
	}

	return &user, nil
}

func (usr *User) GetById(id string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `SELECT id, email, first_name, last_name, password, user_active, created_at, updated_at
	FROM users WHERE id = $1`

	var user User
	row := db.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		log.Println("error user GetById: ", err)
		return nil, err
	}

	return &user, nil
}

func (usr *User) Update() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `UPDATE users SET 
		email = $1, 
		first_name = $2, 
		last_name = $3, 
		user_active = $4, 
		updated_at = $5
		WHERE id = $6
	`

	_, err := db.ExecContext(ctx, query,
		usr.Email,
		usr.FirstName,
		usr.LastName,
		usr.Active,
		time.Now,
		usr.ID,
	)

	if err != nil {
		log.Println("error user Update: ", err)
		return err
	}

	return nil
}

func (usr *User) Delete() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `DELETE from users WHERE id = $1`

	_, err := db.ExecContext(ctx, query, usr.ID)
	if err != nil {
		log.Println("error user Delete: ", err)
		return err
	}

	return nil
}

func (usr *User) DeleteByID(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `DELETE from users WHERE id = $1`

	_, err := db.ExecContext(ctx, query, id)
	if err != nil {
		log.Println("error user Delete: ", err)
		return err
	}

	return nil
}

func (usr *User) Insert(user User) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		log.Println("error user GenerateFromPassword: ", err)
		return 0, err
	}

	var newID int
	query := `INSERT INTO users (email, first_name, last_name, password, user_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7) returning id`

	err = db.QueryRowContext(ctx, query,
		user.Email,
		user.FirstName,
		user.LastName,
		hashedPassword,
		user.Active,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	if err != nil {
		log.Println("error user Insert: ", err)
		return 0, err
	}

	return newID, nil
}

func (usr *User) ResetPassword(password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Println("error user GenerateFromPassword: ", err)
		return err
	}

	query := `UPDATE users SET password = $1 WHERE id = $2`
	_, err = db.ExecContext(ctx, query, hashedPassword, usr.ID)
	if err != nil {
		log.Println("error user ResetPassword: ", err)
		return err
	}

	return nil
}

func (usr *User) PasswordMatch(inputPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(inputPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

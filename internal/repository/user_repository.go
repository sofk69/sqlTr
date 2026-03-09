package repository

import (
	"database/sql"
	"fmt"
	"sqlTr/internal/domain"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *domain.User) (int, error) {
	sqlStatement := `
	INSERT INTO users (age, email, first_name, last_name)
	VALUES ($1, $2, $3, $4)
	RETURNING id`

	var id int
	err := r.db.QueryRow(sqlStatement, user.Age, user.Email, user.FirstName, user.LastName).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("error creating user: %w", err)
	}
	return id, nil
}

func (r *UserRepository) CreatePayment(payment *domain.Payment) (int, error) {
	sqlStatement := `
	INSERT INTO users_payments (user_id, amount, payment_name)
	VALUES ($1, $2, $3)
	RETURNING payment_id`

	var id int
	err := r.db.QueryRow(sqlStatement, payment.UserID, payment.Amount, payment.PaymentName).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("error creating payment: %w", err)
	}
	return id, nil
}

func (r *UserRepository) List(limit int) ([]domain.User, error) {
	rows, err := r.db.Query("SELECT id, first_name, last_name FROM users LIMIT $1", limit)
	if err != nil {
		return nil, fmt.Errorf("error querying users: %w", err)
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		err = rows.Scan(&user.ID, &user.FirstName, &user.LastName)
		if err != nil {
			return nil, fmt.Errorf("error scanning user: %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error after scanning: %w", err)
	}

	return users, nil
}

func (r *UserRepository) Update(id int, firstName, lastName string) error {
	sqlStatement := `
	UPDATE users
	SET first_name = $2, last_name = $3
	WHERE id = $1;`

	result, err := r.db.Exec(sqlStatement, id, firstName, lastName)
	if err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("user with id %d not found", id)
	}

	fmt.Printf("updated user with id %d\n", id)
	return nil
}

func (r *UserRepository) FindByFirstName(firstName string) (*domain.User, error) {
	sqlStatement := `SELECT id, age, email, first_name, last_name FROM users WHERE first_name = $1;`

	var user domain.User
	err := r.db.QueryRow(sqlStatement, firstName).Scan(
		&user.ID, &user.Age, &user.Email, &user.FirstName, &user.LastName,
	)

	if err == sql.ErrNoRows {
		return nil, nil // пользователь не найден
	}
	if err != nil {
		return nil, fmt.Errorf("error finding user: %w", err)
	}

	return &user, nil
}

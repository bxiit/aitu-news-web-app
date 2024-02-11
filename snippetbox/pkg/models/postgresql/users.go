package postgresql

import (
	"alexedwards.net/snippetbox/pkg/models"
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type UserModel struct {
	DB *sql.DB
}

// We'll use the Insert method to add a new record to the users table.
func (m *UserModel) Insert(name, email, password string) error {
	// Create a bcrypt hash of the plain-text password.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (name, email, hashed_password, created, role)
	VALUES($1, $2, $3, CURRENT_TIMESTAMP, $4);`

	_, err = m.DB.Exec(stmt, name, email, string(hashedPassword), "student")
	if err != nil {
		// If this returns an error, we use the errors.As() function to check
		// whether the error has the type *mysql.MySQLError. If it does, the
		// error will be assigned to the mySQLError variable. We can then check
		// whether or not the error relates to our users_uc_email key by
		// checking the contents of the message string. If it does, we return
		// an ErrDuplicateEmail error.
		var pgErr *pq.Error

		if err != nil {
			// PostgreSQL specific error handling for duplicate email
			ok := errors.As(err, &pgErr)
			if ok && pgErr.Code == "23505" && strings.Contains(pgErr.Message, "users_uc_email") {
				return models.ErrDuplicateEmail
			}
			return err
		}
		return err
	}

	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	// Retrieve the id and hashed password associated with the given email. If no
	// matching email exists, or the user is not active, we return the
	// ErrInvalidCredentials error.
	var id int
	var hashedPassword []byte
	stmt := "SELECT id, hashed_password FROM users WHERE email = $1 AND active = TRUE;"
	row := m.DB.QueryRow(stmt, email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	// Check whether the hashed password and plain-text password provided match.
	// If they don't, we return the ErrInvalidCredentials error.
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	// Otherwise, the password is correct. Return the user ID.
	return id, nil
}

// We'll use the Get method to fetch details for a specific user based
// on their user ID.
func (m *UserModel) Get(id int) (*models.User, error) {
	stmt := `SELECT * FROM users WHERE id = $1`

	userRow := m.DB.QueryRow(stmt, id)

	u := &models.User{}

	err := userRow.Scan(&u.ID, &u.Name, &u.Email, &u.HashedPassword, &u.Created, &u.Active, &u.Role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return u, nil
}

func (m *UserModel) GetAll() ([]*models.User, error) {
	stmt := `SELECT * FROM users`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	users := []*models.User{}

	for rows.Next() {
		user := &models.User{}
		err = rows.Scan(&user.ID, &user.Name, &user.Email, &user.HashedPassword, &user.Created, &user.Active, &user.Role)
		if err != nil {
			return nil, err
		}
		// Append it to the slice of snippets.
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (m *UserModel) GetRole(id int) string {
	stmt := `SELECT role FROM users WHERE id = $1`
	var role string
	err := m.DB.QueryRow(stmt, id).Scan(&role)
	if err != nil {
		return ""
	}

	return role
}

func (m *UserModel) ChangeRole(id int, newRole string) {
	stmt := `UPDATE users SET role = $1 WHERE id = $2`
	_, err := m.DB.Exec(stmt, newRole, id)
	if err != nil {
		return
	}
}

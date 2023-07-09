package mysql

import (
	"database/sql"
	"github.com/csn2002/Snippetbox/pkg/models"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type Usermodel struct {
	DB *sql.DB
}

func (m *Usermodel) Insert(name, email, password string) error {
	hashedpassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (name, email, hashed_password,created ) 
VALUES(?, ?, ?, UTC_TIMESTAMP())`
	_, err = m.DB.Exec(stmt, name, email, string(hashedpassword))
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1062 && strings.Contains(mysqlErr.Message, "users_uc_email") {
				return models.ErrDuplicateEmail
			}
		}
	}
	return err

}
func (m *Usermodel) Authenticate(email, password string) (int, error) {
	stmt := `SELECT id, Hashed_password FROM users where email = ?`
	row := m.DB.QueryRow(stmt, email)
	s := models.User{}
	err := row.Scan(&s.ID, &s.HashedPassword)
	if err == sql.ErrNoRows {
		return 0, models.ErrInvalidCredential
	} else if err != nil {
		return 0, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(s.HashedPassword), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, models.ErrInvalidCredential
	} else if err != nil {
		return 0, err
	}
	return s.ID, nil
}
func (m *Usermodel) Get(id int) (*models.User, error) {
	return nil, nil
}

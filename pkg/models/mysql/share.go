package mysql

import (
	"database/sql"
	"github.com/csn2002/Snippetbox/pkg/models"
	"github.com/go-sql-driver/mysql"
	"strings"
	"time"
)

type ShareModel struct {
	DB *sql.DB
}

func (m *ShareModel) Insert(user_id, shared_user_id, snippet_id int) error {

	stmt := `INSERT INTO share (userid, shared_user_id, snippet_id ) 
VALUES(?, ?, ?)`
	_, err := m.DB.Exec(stmt, user_id, shared_user_id, snippet_id)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1062 && strings.Contains(mysqlErr.Message, "users_uc_email") {
				return models.ErrDuplicateEmail
			}
		}
	}
	return err
}

func (m *ShareModel) Authenticate(email string) (int, error) {
	stmt := `SELECT id  FROM users where email = ?`
	row := m.DB.QueryRow(stmt, email)
	s := models.Share{}
	err := row.Scan(&s.ID)
	if err == sql.ErrNoRows {
		return 0, models.ErrInvalidCredential
	} else if err != nil {
		return 0, err
	}
	return s.ID, nil
}
func (m *ShareModel) LatestSharedWithYou(userid int) ([]*models.Snippet, error) {
	stmt := `SELECT a.id, a.title, a.content, a.created, a.expires
FROM snippets AS a
INNER JOIN share AS c ON a.id = c.snippet_id
WHERE c.shared_user_id = ? AND a.expires > UTC_TIMESTAMP()
ORDER BY a.created DESC
LIMIT 10`
	//stmt := `SELECT a.id, a.title, a.content, a.created, a.expires FROM
	// (snippets AS a
	//     INNER JOIN (SELECT snippet_id FROM users WHERE shared_user_id=?) as c ON a.id=c.snippet_id)
	//     WHERE expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT 10`
	snippets := []*models.Snippet{}
	rows, err := m.DB.Query(stmt, userid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		s := &models.Snippet{}
		var createdStr string
		var expiredStr string
		err := rows.Scan(&s.ID, &s.Title, &s.Content, &createdStr, &expiredStr)
		if err != nil {
			return nil, err
		}
		createdTime, parseErr := time.Parse("2006-01-02 15:04:05", createdStr)
		if parseErr != nil {
			return nil, parseErr
		}
		expireTime, parseErr := time.Parse("2006-01-02 15:04:05", expiredStr)
		if parseErr != nil {
			return nil, parseErr
		}
		s.Expires = expireTime
		s.Created = createdTime
		snippets = append(snippets, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return snippets, nil
	return nil, nil
}
func (m *ShareModel) Sharedsnippets(userid int) ([]*models.Snippet, error) {
	stmt := `SELECT a.id, a.title, a.content, a.created, a.expires
FROM snippets AS a
INNER JOIN share AS c ON a.id = c.snippet_id
WHERE c.userid = ? AND a.expires > UTC_TIMESTAMP()
ORDER BY a.created DESC
LIMIT 10`
	//stmt := `SELECT a.id, a.title, a.content, a.created, a.expires FROM
	// (snippets AS a
	//     INNER JOIN (SELECT snippet_id FROM users WHERE shared_user_id=?) as c ON a.id=c.snippet_id)
	//     WHERE expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT 10`
	snippets := []*models.Snippet{}
	rows, err := m.DB.Query(stmt, userid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		s := &models.Snippet{}
		var createdStr string
		var expiredStr string
		err := rows.Scan(&s.ID, &s.Title, &s.Content, &createdStr, &expiredStr)
		if err != nil {
			return nil, err
		}
		createdTime, parseErr := time.Parse("2006-01-02 15:04:05", createdStr)
		if parseErr != nil {
			return nil, parseErr
		}
		expireTime, parseErr := time.Parse("2006-01-02 15:04:05", expiredStr)
		if parseErr != nil {
			return nil, parseErr
		}
		s.Expires = expireTime
		s.Created = createdTime
		snippets = append(snippets, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return snippets, nil
	return nil, nil
}

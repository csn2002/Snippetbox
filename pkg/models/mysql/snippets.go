package mysql

import (
	"database/sql"
	"github.com/csn2002/Snippetbox/pkg/models"
	"time"
)

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title, content, expires string, userid int) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires, userid) 
VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY),?)`
	result, err := m.DB.Exec(stmt, title, content, expires, userid)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}
func (m *SnippetModel) Get(id, userid int) (*models.Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() AND userid = ? AND id = ?`
	row := m.DB.QueryRow(stmt, userid, id)
	s := &models.Snippet{}
	var createdStr string
	var expiredStr string
	err := row.Scan(&s.ID, &s.Title, &s.Content, &createdStr, &expiredStr)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
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
	return s, nil
}
func (m *SnippetModel) Latest(userid int) ([]*models.Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() AND userid = ? ORDER BY created DESC LIMIT 10`
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

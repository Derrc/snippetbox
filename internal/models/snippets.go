package models

import (
	"database/sql"
	"errors"
	"time"
)

type Snippet struct {
	ID int
	Title string
	Content string
	Created time.Time
	Expires time.Time
}

// interface for Snippet CRUD methods
type SnippetModelInterface interface {
	Insert(title string, content string, expires int) (int, error)
	Get(id int) (Snippet, error)
	Latest() ([]Snippet, error)
}

// implements SnippetModelInterface
type SnippetModel struct {
		DB *sql.DB
}

// inserts snippet into 'snippets' table
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// returns snippet with corresponding id
func (m *SnippetModel) Get(id int) (Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() AND id = ?`

	// sql.Row object contains results from query execution
	row := m.DB.QueryRow(stmt, id)

	var s Snippet;

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		// row.Scan returns sql.ErrNoRows if query returns no rows
		if errors.Is(err, sql.ErrNoRows) {
			return Snippet{}, ErrNoRecord
		} else {
			return Snippet{}, err
		}
	}


	return s, nil;
}

// returns 10 most recently created snippets
func (m *SnippetModel) Latest() ([]Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	// defer rows.Close() to ensure that sql.Rows resultset is properly closed
	// if not closed, keeps underyling db connection open -> uses up all of the connections in the pool
	defer rows.Close()

	var snippets []Snippet

	for rows.Next() {
		var s Snippet

		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	// retrieve any error that was encountered during rows.Next()
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
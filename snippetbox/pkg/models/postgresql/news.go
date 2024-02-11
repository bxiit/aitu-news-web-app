package postgresql

import (
	"alexedwards.net/snippetbox/pkg/models"
	"database/sql"
	"errors"
	"fmt"
)

type NewsModel struct {
	DB *sql.DB
}

func (m *NewsModel) LatestTen() ([]*models.News, error) {
	// Write the SQL statement we want to execute.
	stmt := `SELECT id, title, content, created, expires, category FROM news
    WHERE expires > CURRENT_TIMESTAMP ORDER BY created DESC LIMIT 10`

	// Use the Query() method on the connection pool to execute our
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	snippets := []*models.News{}

	for rows.Next() {
		s := &models.News{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires, &s.Category)
		if err != nil {
			return nil, err
		}
		// Append it to the slice of snippets.
		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}

func (m *NewsModel) Get(id int) (*models.News, error) {
	stmt := `SELECT id, title, content, created, expires, category FROM news
    WHERE expires > CURRENT_TIMESTAMP AND id = $1`

	row := m.DB.QueryRow(stmt, id)

	s := &models.News{}

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires, &s.Category)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	// If everything went OK then return the News object.
	return s, nil
}

func (m *NewsModel) Update(id int, title, content, category string) (*models.News, error) {
	stmt := `
		UPDATE news
		SET title = $1, 
		    content = $2, 
		    created = CURRENT_TIMESTAMP, 
		    expires = CURRENT_TIMESTAMP + INTERVAL '7 days', 
		    category = $3
		WHERE id = $4;
	`
	_, err := m.DB.Exec(stmt, title, content, category, id)
	if err != nil {
		return nil, err
	}
	news, err := m.Get(id)
	if err != nil {
		return nil, err
	}
	return news, nil
}

func (m *NewsModel) Insert(title, content, expires, category string) (int, error) {
	stmt := `INSERT INTO news (title, content, created, expires, category)
    VALUES($1, $2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP + INTERVAL '1 day' * $3, $4)
    RETURNING id`

	var id int
	err := m.DB.QueryRow(stmt, title, content, expires, category).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (m *NewsModel) GetCategory(category string) ([]*models.News, error) {
	stmt := `SELECT id, title, content, created, expires FROM news
    WHERE category = $1 AND expires > CURRENT_TIMESTAMP
    ORDER BY created DESC LIMIT 10`

	rows, err := m.DB.Query(stmt, category)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	snippets := []*models.News{}

	for rows.Next() {
		s := &models.News{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
		fmt.Println(snippets)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}

func (m *NewsModel) Delete(id int) error {
	stmt := `DELETE FROM news WHERE id = $1`
	_, err := m.DB.Exec(stmt, id)
	if err != nil {
		return err
	}
	return nil
}

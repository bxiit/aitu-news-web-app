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

	// We defer rows.Close() to ensure the sql.Rows resultset is
	// always properly closed before the LatestTen() method returns. This defer
	// statement should come *after* you check for an error from the Query()
	// method. Otherwise, if Query() returns an error, you'll get a panic
	// trying to close a nil resultset.
	defer rows.Close()

	// Initialize an empty slice to hold the models.Snippets objects.
	snippets := []*models.News{}

	// Use rows.Next to iterate through the rows in the resultset. This
	// prepares the first (and then each subsequent) row to be acted on by the
	// rows.Scan() method. If iteration over all the rows completes then the
	// resultset automatically closes itself and frees-up the underlying
	// database connection.
	for rows.Next() {
		// Create a pointer to a new zeroed News struct.
		s := &models.News{}
		// Use rows.Scan() to copy the values from each field in the row to the
		// new News object that we created. Again, the arguments to row.Scan()
		// must be pointers to the place you want to copy the data into, and the
		// number of arguments must be exactly the same as the number of
		// columns returned by your statement.
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires, &s.Category)
		if err != nil {
			return nil, err
		}
		// Append it to the slice of snippets.
		snippets = append(snippets, s)
	}

	// When the rows.Next() loop has finished we call rows.Err() to retrieve any
	// error that was encountered during the iteration. It's important to
	// call this - don't assume that a successful iteration was completed
	// over the whole resultset.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// If everything went OK then return the Snippets slice.
	return snippets, nil
}

func (m *NewsModel) Get(id int) (*models.News, error) {
	// Write the SQL statement we want to execute. Again, I've split it over two
	// lines for readability.
	stmt := `SELECT id, title, content, created, expires, category FROM news
    WHERE expires > CURRENT_TIMESTAMP AND id = $1`

	// Use the QueryRow() method on the connection pool to execute our
	// SQL statement, passing in the untrusted id variable as the value for the
	// placeholder parameter. This returns a pointer to a sql.Row object which
	// holds the result from the database.
	row := m.DB.QueryRow(stmt, id)

	// Initialize a pointer to a new zeroed News struct.
	s := &models.News{}

	// Use row.Scan() to copy the values from each field in sql.Row to the
	// corresponding field in the News struct. Notice that the arguments
	// to row.Scan are *pointers* to the place you want to copy the data into,
	// and the number of arguments must be exactly the same as the number of
	// columns returned by your statement.
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires, &s.Category)
	if err != nil {
		// If the query returns no rows, then row.Scan() will return a
		// sql.ErrNoRows error. We use the errors.Is() function check for that
		// error specifically, and return our own models.ErrNoRecord error
		// instead.
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
	/// for test git
}

func (m *NewsModel) Insert(title, content, expires, category string) (int, error) {
	// Write the SQL statement we want to execute. Using the RETURNING clause to
	// get the last inserted ID.
	stmt := `INSERT INTO news (title, content, created, expires, category)
    VALUES($1, $2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP + INTERVAL '1 day' * $3, $4)
    RETURNING id`

	// Use the QueryRow() method on the connection pool to execute the
	// SQL statement and retrieve the last inserted ID.
	var id int
	err := m.DB.QueryRow(stmt, title, content, expires, category).Scan(&id)
	if err != nil {
		return 0, err
	}

	// Return the last inserted ID.
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

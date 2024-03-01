package postgresql

import (
	"alexedwards.net/snippetbox/pkg/models"
	"database/sql"
	_ "github.com/lib/pq"
)

type CommentModel struct {
	DB *sql.DB
}

func (c *CommentModel) GetAllByNewsId(newsId int) ([]*models.Comment, error) {
	stmt := `SELECT * FROM comments WHERE news_id = $1`

	rows, err := c.DB.Query(stmt, newsId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	comments := []*models.Comment{}

	for rows.Next() {
		c := &models.Comment{}
		err = rows.Scan(&c.ID, &c.UserId, &c.NewsID, &c.Text)
		if err != nil {
			return nil, err
		}

		comments = append(comments, c)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func (c *CommentModel) InsertComment(userId, newsId int, text string) (int, error) {
	stmt := `INSERT INTO comments (user_id, news_id, text) VALUES($1, $2, $3)`

	var id int
	err := c.DB.QueryRow(stmt, userId, newsId, text).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (c *CommentModel) DeleteComment(commentId int) error {
	stmt := `DELETE FROM comments WHERE id = $1`
	_, err := c.DB.Exec(stmt, commentId)
	if err != nil {
		return err
	}
	return nil
}

func (c *CommentModel) GetOwnerId(newsId int) (int, error) {
	stmt := `SELECT user_id FROM comments WHERE news_id = $1`

	row := c.DB.QueryRow(stmt, newsId)

	var ownerId int

	err := row.Scan(&ownerId)
	if err != nil {
		return 0, err
	}
	return ownerId, err
}

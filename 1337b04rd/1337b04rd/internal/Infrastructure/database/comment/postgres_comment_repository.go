package comment

import (
	"1337b04rd/internal/domain/models"
	"database/sql"
)

type PostgresCommentRepository struct {
	db *sql.DB
}

func NewPostgresCommentRepository(db *sql.DB) CommentRepository {
	return &PostgresCommentRepository{db: db}
}

func (r *PostgresCommentRepository) CreateComment(comment *models.Comment) error {
	query := `
		INSERT INTO comments (id, post_id, parent_id, content, avatar, name, image, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.Exec(
		query,
		comment.ID,
		comment.PostID,
		comment.ParentID,
		comment.Content,
		comment.Avatar,
		comment.Name,
		comment.Image,
		comment.CreatedAt,
	)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(
		`UPDATE posts SET last_comment_at = $1 WHERE id = $2`,
		comment.CreatedAt,
		comment.PostID,
	)
	return err
}

func (r *PostgresCommentRepository) GetCommentByID(id int) (*models.Comment, error) {
	query := `
		SELECT id, post_id, parent_id, content, avatar, name, image, created_at
		FROM comments WHERE id = $1
	`
	row := r.db.QueryRow(query, id)

	comment := &models.Comment{}
	err := row.Scan(
		&comment.ID,
		&comment.PostID,
		&comment.ParentID,
		&comment.Content,
		&comment.Avatar,
		&comment.Name,
		&comment.Image,
		&comment.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return comment, err
}

func (r *PostgresCommentRepository) GetCommentsByPostID(postID int) ([]*models.Comment, error) {
	query := `
		SELECT id, post_id, parent_id, content, avatar, name, image, created_at
		FROM comments WHERE post_id = $1 ORDER BY created_at ASC
	`
	rows, err := r.db.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*models.Comment
	for rows.Next() {
		c := &models.Comment{}
		err := rows.Scan(
			&c.ID,
			&c.PostID,
			&c.ParentID,
			&c.Content,
			&c.Avatar,
			&c.Name,
			&c.Image,
			&c.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, nil
}

func (r *PostgresCommentRepository) DeleteComment(id int) error {
	_, err := r.db.Exec(`DELETE FROM comments WHERE id = $1`, id)
	return err
}

func (r *PostgresCommentRepository) GetReplies(commentID int) ([]*models.Comment, error) {
	query := `
		SELECT id, post_id, parent_id, content, avatar, name, image, created_at
		FROM comments WHERE parent_id = $1 ORDER BY created_at ASC
	`
	rows, err := r.db.Query(query, commentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var replies []*models.Comment
	for rows.Next() {
		c := &models.Comment{}
		err := rows.Scan(
			&c.ID,
			&c.PostID,
			&c.ParentID,
			&c.Content,
			&c.Avatar,
			&c.Name,
			&c.Image,
			&c.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		replies = append(replies, c)
	}
	return replies, nil
}

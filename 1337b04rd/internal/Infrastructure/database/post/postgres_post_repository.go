package post

import (
	"1337b04rd/internal/domain/models"
	"database/sql"
	"fmt"
	"time"
)

type PostgresPostRepository struct {
	db *sql.DB
}

func NewPostgresPostRepository(db *sql.DB) PostRepository {
	return &PostgresPostRepository{db: db}
}

func (r *PostgresPostRepository) CreatePost(post *models.Post) error {
	query := `
		INSERT INTO posts (title, content, avatar, name, image, created_at, last_comment_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	err := r.db.QueryRow(
		query,
		post.Title,
		post.Content,
		post.Avatar,
		post.Name,
		post.Image,
		post.CreatedAt,
		post.LastCommentAt,
	).Scan(&post.ID)
	if err != nil {
		return fmt.Errorf("error creating post: %v", err)
	}
	return nil
}

func (r *PostgresPostRepository) GetPostByID(id int) (*models.Post, error) {
	query := `
		SELECT id, title, content, avatar, name, image, created_at, last_comment_at
		FROM posts
		WHERE id = $1
	`

	var post models.Post
	err := r.db.QueryRow(query, id).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.Avatar,
		&post.Name,
		&post.Image,
		&post.CreatedAt,
		&post.LastCommentAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("post with id %d not found", id)
		}
		return nil, fmt.Errorf("error fetching post by ID: %v", err)
	}

	return &post, nil
}

func (r *PostgresPostRepository) DeletePost(id int) error {
	query := `
		DELETE FROM posts WHERE id = $1
	`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting post: %v", err)
	}

	return nil
}

func (r *PostgresPostRepository) GetActivePosts() ([]*models.Post, error) {
	now := time.Now()
	noCommentsSince := now.Add(-10 * time.Minute)
	withCommentsSince := now.Add(-15 * time.Minute)

	query := `
		SELECT id, title, content, avatar, name, image, created_at, last_comment_at
		FROM posts
		WHERE 
			(last_comment_at IS NULL AND created_at > $1) OR
			(last_comment_at IS NOT NULL AND last_comment_at > $2)
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, noCommentsSince, withCommentsSince)
	if err != nil {
		return nil, fmt.Errorf("error getting active posts: %v", err)
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.Avatar,
			&post.Name,
			&post.Image,
			&post.CreatedAt,
			&post.LastCommentAt,
		); err != nil {
			return nil, fmt.Errorf("error scanning post: %v", err)
		}
		posts = append(posts, &post)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return posts, nil
}

func (r *PostgresPostRepository) GetArchivedPosts() ([]*models.Post, error) {
	now := time.Now()
	noCommentsBefore := now.Add(-10 * time.Minute)
	withCommentsBefore := now.Add(-15 * time.Minute)

	query := `
		SELECT id, title, content, avatar, name, image, created_at, last_comment_at
		FROM posts
		WHERE 
			(last_comment_at IS NULL AND created_at <= $1) OR
			(last_comment_at IS NOT NULL AND last_comment_at <= $2)
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, noCommentsBefore, withCommentsBefore)
	if err != nil {
		return nil, fmt.Errorf("error getting archived posts: %v", err)
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.Avatar,
			&post.Name,
			&post.Image,
			&post.CreatedAt,
			&post.LastCommentAt,
		); err != nil {
			return nil, fmt.Errorf("error scanning post: %v", err)
		}
		posts = append(posts, &post)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return posts, nil
}

// func (r *PostgresPostRepository) UpdatePost(post *models.Post) error {
// 	query := `
// 		UPDATE posts
// 		SET 
// 			title = COALESCE($1, title),
// 			content = COALESCE($2, content),
// 			avatar = COALESCE($3, avatar),
// 			name = COALESCE($4, name),
// 			image = COALESCE($5, image),
// 			last_comment_at = COALESCE($6, last_comment_at)
// 		WHERE id = $7
// 	`
// 	_, err := r.db.Exec(
// 		query,
// 		post.Title,
// 		post.Content,
// 		post.Avatar,
// 		post.Name,
// 		post.Image,
// 		post.LastCommentAt,
// 		post.ID,
// 	)
// 	if err != nil {
// 		return fmt.Errorf("error updating post: %v", err)
// 	}

// 	return nil
// }

func (r *PostgresPostRepository) UpdateLastCommentAt(postID int) error {
	query := `
		UPDATE posts
		SET last_comment_at = $1
		WHERE id = $2
	`

	_, err := r.db.Exec(query, time.Now(), postID)
	if err != nil {
		return fmt.Errorf("error updating last_comment_at for post %d: %v", postID, err)
	}

	return nil
}

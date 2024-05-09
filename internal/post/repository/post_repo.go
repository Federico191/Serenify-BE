package repository

import (
	"FindIt/internal/entity"
	customError "FindIt/pkg/error"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PostRepoItf interface {
	CreatePost(post *entity.Post) error
	GetPostById(id uuid.UUID) (*entity.Post, error)
	GetAllPosts(limit, offset int) ([]entity.PostWithLikeCount, error)
	UpdatePost(post *entity.Post) error
	IsPostOwner(userId, postId uuid.UUID) (bool, error)
	DeletePost(id uuid.UUID) error
}

type PostRepo struct {
	db *sqlx.DB
}

func NewPostRepo(db *sqlx.DB) PostRepoItf {
	return &PostRepo{db: db}
}

// CreatePost implements PostRepoItf.
func (p *PostRepo) CreatePost(post *entity.Post) error {
	result, err := p.db.Exec(CreatePostQuery, post.ID, post.Content, post.PhotoLink, post.UserID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows != 1 {
		return fmt.Errorf("unexpected number of rows affected: %d", rows)
	}

	return nil
}

// DeletePost implements PostRepoItf.
func (p *PostRepo) DeletePost(id uuid.UUID) error {
	result, err := p.db.Exec(DeletePostQuery, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows != 1 {
		return fmt.Errorf("unexpected number of rows affected: %d", rows)
	}

	return nil
}

// GetPostById implements PostRepoItf.
func (p *PostRepo) GetPostById(id uuid.UUID) (*entity.Post, error) {
	var post entity.Post

	err := p.db.QueryRowx(GetPostByIDQuery, id).StructScan(&post)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, customError.ErrRecordNotFound
		}
		return nil, err
	}

	return &post, nil
}

// GetAllPosts implements PostRepoItf.
func (p *PostRepo) GetAllPosts(limit, offset int) ([]entity.PostWithLikeCount, error) {
	var posts []entity.PostWithLikeCount

	rows, err := p.db.Queryx(GetAllPostsQuery, limit, offset)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var post entity.Post
		var like_count int
		if err := rows.Scan(&post.ID, &post.UserID, &post.Content,
			&post.PhotoLink, &post.CreatedAt, &post.UpdatedAt, &like_count); err != nil {
			return nil, err
		}

		posts = append(posts, entity.PostWithLikeCount{
			Post:      post,
			LikeCount: like_count,
		})
	}

	if len(posts) == 0 {
		return nil, customError.ErrRecordNotFound
	}

	return posts, nil
}

// UpdatePost implements PostRepoItf.
func (p *PostRepo) UpdatePost(post *entity.Post) error {
	result, err := p.db.Exec(UpdatePostQuery, post.Content, post.PhotoLink, post.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows != 1 {
		return fmt.Errorf("unexpected number of rows affected: %d", rows)
	}

	return nil
}

// IsPostOwner implements PostRepoItf.
func (p *PostRepo) IsPostOwner(userId uuid.UUID, postId uuid.UUID) (bool, error) {
	var exists bool

	err := p.db.QueryRowx(IsOwnerQuery , postId, userId).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

package repository

import (
	"FindIt/internal/entity"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type CommentRepoItf interface {
	CreateComment(uc *entity.Comment) error
	GetTotalComments(postId uuid.UUID) (int, error)
	GetAllCommentsByPostId(postId uuid.UUID) ([]entity.Comment, error)
	GetCommentById(commentId int) (*entity.Comment, error)
	IsOwner(commentId int, userId uuid.UUID) (bool, error)
	UpdateComment(uc *entity.Comment) error
	DeleteComment(commentId int) error
}

func NewCommentRepo(db *sqlx.DB) CommentRepoItf {
	return &CommentRepo{db: db}
}

type CommentRepo struct {
	db *sqlx.DB
}

// CreateComment implements CommentRepoItf.
func (c *CommentRepo) CreateComment(uc *entity.Comment) error {
	result, err := c.db.Exec(CreateCommentQuery, uc.UserID, uc.PostID, uc.Comment)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows != 1 {
		return fmt.Errorf("expected single row affected, got %d rows affected", rows)
	}

	return nil
}

// GetCommentById implements CommentRepoItf.
func (c *CommentRepo) GetCommentById(commentId int) (*entity.Comment, error) {
	var comment entity.Comment

	err := c.db.QueryRowx(GetCommentByIdQuery, commentId).StructScan(&comment)
	if err != nil {
		return nil, err
	}

	return &comment, nil
}

// GetAllComments implements CommentRepoItf.
func (c *CommentRepo) GetAllCommentsByPostId(postId uuid.UUID) ([]entity.Comment, error) {
	var comments []entity.Comment

    rows, err := c.db.Queryx(GetAllCommentByPostIdQuery, postId)
    if err != nil {
        return nil, err
    }

    for rows.Next() {
        var comment entity.Comment
        err = rows.StructScan(&comment)
        if err != nil {
            return nil, err
        }
        comments = append(comments, comment)
    }

    return comments, nil
}

// GetTotalComments implements CommentRepoItf.
func (c *CommentRepo) GetTotalComments(postId uuid.UUID) (int, error) {
	var total int

	err := c.db.QueryRowx(GetTotalCommentsQuery, postId).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}

// IsOwner implements CommentRepoItf.
func (c *CommentRepo) IsOwner(commentId int, userId uuid.UUID) (bool, error) {
	var exists bool

	err := c.db.QueryRowx(IsOwnerQuery, commentId, userId).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// UpdateComment implements CommentRepoItf.
func (c *CommentRepo) UpdateComment(uc *entity.Comment) error {
	result, err := c.db.Exec(UpdateCommentQuery, uc.Comment, uc.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows != 1 {
		return fmt.Errorf("expected single row affected, got %d rows affected", rows)
	}

	return nil
}

// DeleteComment implements CommentRepoItf.
func (c *CommentRepo) DeleteComment(commentId int) error {
	result, err := c.db.Exec(DeleteCommentQuery, commentId)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows != 1 {
		return fmt.Errorf("expected single row affected, got %d rows affected", rows)
	}

	return nil
}

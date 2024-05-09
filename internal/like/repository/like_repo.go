package repository

import (
	"FindIt/internal/entity"
	"FindIt/model"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type LikeRepoItf interface {
	CreatePostLike(postLike *entity.PostLike) error
	CreateCommentLike(commentLike *entity.CommentLike) error
	GetTotalPostLikes(postId uuid.UUID) (int, error)
	GetTotalCommentLikes(commentId int) (int, error)
	IsPostOwner(userId, postId uuid.UUID) (bool, error)
	IsCommentOwner(userId uuid.UUID, commentId int) (bool, error)
	DeletePostLike(req model.PostLikeReq) error
	DeleteCommentLike(req model.CommentLikeReq) error
}

type LikeRepo struct {
	db *sqlx.DB
}

func NewLikeRepo(db *sqlx.DB) LikeRepoItf {
	return &LikeRepo{db: db}
}

// CreateLike implements LikeRepoItf.
func (l *LikeRepo) CreatePostLike(postLike *entity.PostLike) error {
	result, err := l.db.Exec(CreatePostLikeQuery, postLike.UserID, postLike.PostID)
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

// CreateCommentLike implements LikeRepoItf.
func (l *LikeRepo) CreateCommentLike(commentLike *entity.CommentLike) error {
	result, err := l.db.Exec(CreateCommentLikeQuery, commentLike.UserID, commentLike.CommentID)
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

// IsPostOwner implements LikeRepoItf.
func (l *LikeRepo) IsPostOwner(userId uuid.UUID, postId uuid.UUID) (bool, error) {
	var exists bool

	err := l.db.QueryRowx(IsPostOwnerQuery, postId, userId).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// IsCommentOwner implements LikeRepoItf.
func (l *LikeRepo) IsCommentOwner(userId uuid.UUID, commentId int) (bool, error) {
	var exists bool

	err := l.db.QueryRowx(IsCommentOwnerQuery, commentId , userId).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// DeleteLike implements LikeRepoItf.
func (l *LikeRepo) DeletePostLike(req model.PostLikeReq) error {
	result, err := l.db.Exec(DeletePostLikeQuery, req.UserID, req.PostID)
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

// DeleteCommentLike implements LikeRepoItf.
func (l *LikeRepo) DeleteCommentLike(req model.CommentLikeReq) error {
	result, err := l.db.Exec(DeleteCommentLikeQuery, req.UserID, req.CommentID)
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

// GetTotalLikes implements LikeRepoItf.
func (l *LikeRepo) GetTotalPostLikes(postId uuid.UUID) (int, error) {
	var total int

	err := l.db.QueryRowx(GetTotalPostLikesQuery, postId).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}

// GetTotalCommentLikes implements LikeRepoItf.
func (l *LikeRepo) GetTotalCommentLikes(commentId int) (int, error) {
	var total int

	err := l.db.QueryRowx(GetTotalCommentLikesQuery, commentId).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}

package usecase

import (
	"FindIt/internal/entity"
	"FindIt/internal/like/repository"
	"FindIt/model"
	customError "FindIt/pkg/error"
	"strings"
)

type LikeUCItf interface {
	CreatePostLike(req model.PostLikeReq) error
	CreateCommentLike(req model.CommentLikeReq) error
	DeletePostLike(req model.PostLikeReq) error
	DeleteCommentLike(req model.CommentLikeReq) error
}

type LikeUC struct {
	likeRepo repository.LikeRepoItf
}

func NewLikeUC(likeRepo repository.LikeRepoItf) LikeUCItf {
	return &LikeUC{likeRepo: likeRepo}
}

// CreateLike implements LikeUCItf.
func (l *LikeUC) CreatePostLike(req model.PostLikeReq) error {
	like := &entity.PostLike{
		UserID: req.UserID,
		PostID: req.PostID,
	}

	err := l.likeRepo.CreatePostLike(like)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return customError.ErrRecordAlreadyExists
		}
		return err
	}

	return nil
}

// CreateCommentLike implements LikeUCItf.
func (l *LikeUC) CreateCommentLike(req model.CommentLikeReq) error {
	like := &entity.CommentLike{
		UserID: req.UserID,
		CommentID: req.CommentID,
	}

	err := l.likeRepo.CreateCommentLike(like)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return customError.ErrRecordAlreadyExists
		}
		return err
	}

	return nil
}

// DeleteCommentLike implements LikeUCItf.
func (l *LikeUC) DeleteCommentLike(req model.CommentLikeReq) error {
	exist, err := l.likeRepo.IsCommentOwner(req.UserID, req.CommentID)
	if err != nil {
		return err
	}

	if !exist {
		return customError.ErrNotAuthorize
	}

	return l.likeRepo.DeleteCommentLike(req)
}

// DeleteLike implements LikeUCItf.
func (l *LikeUC) DeletePostLike(req model.PostLikeReq) error {
	exist, err := l.likeRepo.IsPostOwner(req.UserID, req.PostID)
	if err != nil {
		return err
	}

	if !exist {
		return customError.ErrNotAuthorize
	}

	return l.likeRepo.DeletePostLike(req)
}

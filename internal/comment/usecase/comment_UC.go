package usecase

import (
	commentRepo "FindIt/internal/comment/repository"
	"FindIt/internal/entity"
	likeRepo "FindIt/internal/like/repository"
	userRepo "FindIt/internal/user/repository"
	"FindIt/model"
	customError "FindIt/pkg/error"
	"strings"

	"github.com/google/uuid"
)

type CommentUCItf interface {
	CreateComment(req model.CreateCommentReq) (*model.CommentResp, error)
	UpdateComment(req model.UpdateCommentReq) (*model.CommentResp, error)
	DeleteComment(userId uuid.UUID, commentId int) error
}

type CommentUC struct {
	commentRepo commentRepo.CommentRepoItf
	userRepo    userRepo.UserRepoItf
	likeRepo    likeRepo.LikeRepoItf
}

func NewCommentUC(commentRepo commentRepo.CommentRepoItf,
	userRepo userRepo.UserRepoItf, likeRepo likeRepo.LikeRepoItf) CommentUCItf {
	return &CommentUC{commentRepo: commentRepo, userRepo: userRepo, likeRepo: likeRepo}
}

// CreateComment implements CommentUCItf.
func (c *CommentUC) CreateComment(req model.CreateCommentReq) (*model.CommentResp, error) {
	user, err := c.userRepo.GetUserById(req.UserID)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return nil, customError.ErrRecordAlreadyExists
		}
		return nil, err
	}

	comment := &entity.Comment{
		UserID:  req.UserID,
		PostID:  req.PostID,
		Comment: req.Comment,
	}

	err = c.commentRepo.CreateComment(comment)
	if err != nil {
		return nil, err
	}

	totalLikes, err := c.likeRepo.GetTotalCommentLikes(comment.ID)
	if err != nil {
		return nil, err
	}

	return convertToCommentResp(comment, user, totalLikes), nil
}

// UpdateComment implements CommentUCItf.
func (c *CommentUC) UpdateComment(req model.UpdateCommentReq) (*model.CommentResp, error) {
	comment, err := c.commentRepo.GetCommentById(req.ID)
	if err != nil {
		return nil, err
	}

	user, err := c.userRepo.GetUserById(req.UserID)
	if err != nil {
		return nil, err
	}

	exist, err := c.commentRepo.IsOwner(req.ID, req.UserID)
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, customError.ErrNotAuthorize
	}

	comment.Comment = req.Comment

	err = c.commentRepo.UpdateComment(comment)
	if err != nil {
		return nil, err
	}

	totalLikes, err := c.likeRepo.GetTotalCommentLikes(comment.ID)
	if err != nil {
		return nil, err
	}

	return convertToCommentResp(comment, user, totalLikes), nil
}

// DeleteComment implements CommentUCItf.
func (c *CommentUC) DeleteComment(userId uuid.UUID, commentId int) error {
	exist, err := c.commentRepo.IsOwner(commentId, userId)
	if err != nil {
		return err
	}

	if !exist {
		return customError.ErrNotAuthorize
	}

	return c.commentRepo.DeleteComment(commentId)
}

func convertToCommentResp(comment *entity.Comment, user *entity.User, totalLikes int) *model.CommentResp {
	return &model.CommentResp{
		ID:         comment.ID,
		UserID:     comment.UserID,
		UserName:   user.FullName,
		UserPhoto:  user.PhotoLink.String,
		PostID:     comment.PostID,
		Comment:    comment.Comment,
		CreatedAt:  comment.CreatedAt,
		TotalLikes: totalLikes,
	}
}

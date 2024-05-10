package usecase

import (
	commentRepo "FindIt/internal/comment/repository"
	"FindIt/internal/entity"
	likeRepo "FindIt/internal/like/repository"
	postRepo "FindIt/internal/post/repository"
	userRepo "FindIt/internal/user/repository"
	"FindIt/model"
	customError "FindIt/pkg/error"
	supabaseStorage "FindIt/pkg/supabase"
	"database/sql"
	"os"
	"strings"

	"github.com/google/uuid"
)

type PostUCItf interface {
	CreatePost(req *model.CreatePostReq) (*model.PostResponse, error)
	GetAllPosts() ([]model.PostResponse, error)
	GetPostByPostId(post_id uuid.UUID) (*model.PostDetailResponse, error)
	UpdatePost(req *model.UpdatePostReq) (*model.PostResponse, error)
	DeletePost(user_id, post_id uuid.UUID) error
}

type PostUC struct {
	postRepo    postRepo.PostRepoItf
	likeRepo    likeRepo.LikeRepoItf
	commentRepo commentRepo.CommentRepoItf
	userRepo    userRepo.UserRepoItf
	supabase    supabaseStorage.SupabaseStorageItf
}

func NewPostUC(postRepo postRepo.PostRepoItf, userRepo userRepo.UserRepoItf,
	likeRepo likeRepo.LikeRepoItf, commentRepo commentRepo.CommentRepoItf,
	supabase supabaseStorage.SupabaseStorageItf) PostUCItf {
	return &PostUC{postRepo: postRepo, userRepo: userRepo,
		likeRepo: likeRepo, commentRepo: commentRepo,
        supabase: supabase}
}

// CreatePost implements PostUCItf.
func (p *PostUC) CreatePost(req *model.CreatePostReq) (*model.PostResponse, error) {
	user, err := p.userRepo.GetUserById(req.UserID)
	if err != nil {
		return nil, err
	}

	var photoLink string
	if req.Photo != nil && req.Photo.Size > 0 {
		photoLink, err = p.supabase.Upload(os.Getenv("SUPABASE_BUCKET_POST"), req.Photo)
		if err != nil {
			return nil, err
		}
	}

	post := &entity.Post{
		ID:      uuid.New(),
		Content: req.Content,
		UserID:  user.ID,
		PhotoLink: sql.NullString{
			String: photoLink,
			Valid:  true,
		},
	}

	err = p.postRepo.CreatePost(post)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return nil, customError.ErrRecordAlreadyExists
		}
		return nil, err
	}

	totalLikes, err := p.likeRepo.GetTotalPostLikes(post.ID)
	if err != nil {
		return nil, err
	}

	totalComments, err := p.commentRepo.GetTotalComments(post.ID)
	if err != nil {
		return nil, err
	}

	return convertToPostResponse(post, user, totalLikes, totalComments), nil
}

// GetAllPosts implements PostUCItf.
func (p *PostUC) GetAllPosts() ([]model.PostResponse, error) {
	posts, err := p.postRepo.GetAllPosts()
	if err != nil {
		return nil, err
	}

	var responses []model.PostResponse
	for _, post := range posts {
		user, err := p.userRepo.GetUserById(post.Post.UserID)
		if err != nil {
			return nil, err
		}

		totalComments, err := p.commentRepo.GetTotalComments(post.Post.ID)
		if err != nil {
			return nil, err
		}

		response := &model.PostResponse{
			ID:            post.Post.ID,
			UserID:        user.ID,
			UserName:      user.FullName,
			UserPhoto:     user.PhotoLink.String,
			Content:       post.Post.Content,
			PhotoLink:     post.Post.PhotoLink.String,
			TotalLikes:    post.LikeCount,
			TotalComments: totalComments,
		}

		responses = append(responses, *response)
	}

	return responses, nil
}

// UpdatePost implements PostUCItf.
func (p *PostUC) UpdatePost(req *model.UpdatePostReq) (*model.PostResponse, error) {
	exist, err := p.postRepo.IsPostOwner(req.UserID, req.ID)
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, customError.ErrNotAuthorize
	}

    user, err := p.userRepo.GetUserById(req.UserID)
    if err != nil {
        return nil, err
    }

	post, err := p.postRepo.GetPostById(req.ID)
	if err != nil {
		return nil, err
	}

	post.Content = req.Content

	if req.Photo != nil {
		photoLink, err := p.supabase.Upload(os.Getenv("SUPABASE_BUCKET_POST"), req.Photo)
		if err != nil {
			return nil, err
		}

		post.PhotoLink = sql.NullString{
			String: photoLink,
			Valid:  true,
		}
	}

    totalLikes, err := p.likeRepo.GetTotalPostLikes(post.ID)
	if err != nil {
		return nil, err
	}

	totalComments, err := p.commentRepo.GetTotalComments(post.ID)
	if err != nil {
		return nil, err
	}

	err = p.postRepo.UpdatePost(post)
	if err != nil {
		return nil, err
	}

	return convertToPostResponse(post, user, totalLikes, totalComments) ,nil
}

// DeletePost implements PostUCItf.
func (p *PostUC) DeletePost(user_id, post_id uuid.UUID) error {
	exist, err := p.postRepo.IsPostOwner(user_id, post_id)
	if err != nil {
		return err
	}

	if !exist {
		return customError.ErrNotAuthorize
	}

	err = p.postRepo.DeletePost(post_id)
	if err != nil {
		return err
	}

	return nil
}

// GetPostByPostId implements PostUCItf.
func (p *PostUC) GetPostByPostId(post_id uuid.UUID) (*model.PostDetailResponse, error) {
	post, err := p.postRepo.GetPostById(post_id)
	if err != nil {
		return nil, err
	}

	user, err := p.userRepo.GetUserById(post.UserID)
	if err != nil {
		return nil, err
	}

	totalLikes, err := p.likeRepo.GetTotalPostLikes(post.ID)
	if err != nil {
		return nil, err
	}

	totalComments, err := p.commentRepo.GetTotalComments(post.ID)
	if err != nil {
		return nil, err
	}

	comments, err := p.commentRepo.GetAllCommentsByPostId(post.ID)
	if err != nil {
		return nil, err
	}


    var commentsResp []model.CommentResp
    for _, comment := range comments {
        user, err := p.userRepo.GetUserById(comment.UserID)
        if err != nil {
            return nil, err
        }
        totalLikes, err := p.likeRepo.GetTotalCommentLikes(comment.ID)
        if err != nil {
            return nil, err
        }
        commentsResp = append(commentsResp, *convertToCommentResp(&comment, user, totalLikes))
    }

	resp := &model.PostDetailResponse{
		ID:            post.ID,
		UserID:        user.ID,
		UserName:      user.FullName,
		UserPhoto:     user.PhotoLink.String,
		Content:       post.Content,
		PhotoLink:     post.PhotoLink.String,
		TotalLikes:    totalLikes,
		TotalComments: totalComments,
		Comments:      commentsResp,
	}

	return resp, nil
}

func convertToPostResponse(post *entity.Post, user *entity.User, totalLikes, totalComment int) *model.PostResponse {
	return &model.PostResponse{
		ID:            post.ID,
		UserID:        post.UserID,
		UserName:      user.FullName,
		UserPhoto:     user.PhotoLink.String,
		Content:       post.Content,
		PhotoLink:     post.PhotoLink.String,
		TotalLikes:    totalLikes,
		TotalComments: totalComment,
	}
}

func convertToCommentResp(comment *entity.Comment, user *entity.User, totalLikes int) *model.CommentResp {
	return &model.CommentResp{
		ID:        comment.ID,
		UserID:    comment.UserID,
		UserName:  user.FullName,
		UserPhoto: user.PhotoLink.String,
		PostID:    comment.PostID,
		Comment:   comment.Comment,
		CreatedAt: comment.CreatedAt,
        TotalLikes: totalLikes,
	}
}
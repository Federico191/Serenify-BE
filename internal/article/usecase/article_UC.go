package usecase

import (
	"FindIt/internal/article/repository"
	"FindIt/internal/entity"
	"FindIt/model"
    customError "FindIt/pkg/error"
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

type ArticleUCItf interface {
	GetAllArticles() ([]*model.ArticleResponse, error)
	GetArticleById(id uuid.UUID) (*model.ArticleResponse, error)
}

type ArticleUC struct {
	articleRepo repository.ArticleRepoItf
}

func NewArticleUC(articleRepo repository.ArticleRepoItf) ArticleUCItf {
	return &ArticleUC{
		articleRepo: articleRepo,
	}
}

// GetAllArticles implements ArticleUCItf.
func (a *ArticleUC) GetAllArticles() ([]*model.ArticleResponse, error) {
    articles, err := a.articleRepo.GetAllArticles()
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, customError.ErrRecordNotFound
        }
        return nil, err
    }

    var responses []*model.ArticleResponse
    for _, article := range articles {
        response := &entity.Article{
            ID: article.ID,
            Title: article.Title,
            Content: article.Content,
            PhotoLink: article.PhotoLink,
        }

        responses = append(responses, convertToArticleResponse(response))
    }

    return responses, nil
}

// GetArticleById implements ArticleUCItf.
func (a *ArticleUC) GetArticleById(id uuid.UUID) (*model.ArticleResponse, error) {
    article, err := a.articleRepo.GetArticleById(id)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, customError.ErrRecordNotFound
        }
        return nil, err
    }

    return convertToArticleResponse(article), nil
}

func convertToArticleResponse(article *entity.Article) *model.ArticleResponse {
    return &model.ArticleResponse{
        ID: article.ID,
        Title: article.Title,
        Content: article.Content,
        PhotoLink: article.PhotoLink.String,
    }
}
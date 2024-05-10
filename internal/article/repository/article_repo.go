package repository

import (
	"FindIt/internal/entity"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ArticleRepoItf interface {
	GetAllArticles() ([]*entity.Article, error)
	GetArticleById(id uuid.UUID) (*entity.Article, error)
}

type ArticleRepo struct {
	db *sqlx.DB
}

func NewArticleRepo(db *sqlx.DB) ArticleRepoItf {
	return &ArticleRepo{db: db}
}


// GetAllArticles implements ArticleRepoItf.
func (a *ArticleRepo) GetAllArticles() ([]*entity.Article, error) {
	var articles []*entity.Article

    rows, err := a.db.Queryx(getAllArticlesQuery)
    if err != nil {
        return nil, err
    }

    for rows.Next() {
        var article entity.Article
        err = rows.StructScan(&article)
        if err != nil {
            return nil, err
        }
        articles = append(articles, &article)
    }

    return articles, nil
}

// GetArticleById implements ArticleRepoItf.
func (a *ArticleRepo) GetArticleById(id uuid.UUID) (*entity.Article, error) {
    var article entity.Article

    err := a.db.QueryRowx(getArticleByIdQuery, id).StructScan(&article)
    if err != nil {
        return nil, err
    }

    return &article, nil
}

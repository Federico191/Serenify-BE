package repository

const (
    getAllArticlesQuery = `SELECT * FROM articles`

    getArticleByIdQuery = `SELECT * FROM articles WHERE id = $1`
)
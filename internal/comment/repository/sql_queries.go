package repository

const (
    CreateCommentQuery = `INSERT INTO comments (user_id, post_id, comment) VALUES ($1, $2, $3) RETURNING *`

    GetCommentByIdQuery = `SELECT * FROM comments WHERE id = $1`

    GetAllCommentByPostIdQuery = `SELECT * FROM comments WHERE post_id = $1`
    
    GetTotalCommentsQuery = `SELECT COUNT(*) FROM comments WHERE post_id = $1`

    IsOwnerQuery = `SELECT EXISTS(
        SELECT 1 FROM comments
        WHERE id = $1 AND user_id = $2
    )`

    UpdateCommentQuery = `UPDATE comments SET comment = $1 WHERE id = $2 RETURNING *`

    DeleteCommentQuery = `DELETE FROM comments WHERE id = $1 RETURNING *`
)
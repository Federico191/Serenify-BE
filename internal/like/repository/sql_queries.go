package repository

const (
    CreatePostLikeQuery = `INSERT INTO post_likes (user_id, post_id) VALUES ($1, $2) RETURNING *`

    CreateCommentLikeQuery = `INSERT INTO comment_likes (user_id, comment_id) VALUES ($1, $2) RETURNING *`

    GetTotalCommentLikesQuery  = `SELECT COUNT(*) FROM comment_likes WHERE comment_id = $1`

    GetTotalPostLikesQuery = `SELECT COUNT(*) FROM post_likes WHERE post_id = $1`

    IsPostOwnerQuery = `SELECT EXISTS(
        SELECT 1 FROM post_likes
        WHERE post_id = $1 AND user_id = $2
    )`

    IsCommentOwnerQuery = `SELECT EXISTS(
        SELECT 1 FROM comment_likes
        WHERE comment_id = $1 AND user_id = $2
    )`

    DeletePostLikeQuery = `DELETE FROM post_likes WHERE user_id = $1 AND post_id = $2 RETURNING *`

    DeleteCommentLikeQuery = `DELETE FROM comment_likes WHERE user_id = $1 AND comment_id = $2 RETURNING *`
)
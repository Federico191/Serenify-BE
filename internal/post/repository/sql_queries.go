package repository

const (
    CreatePostQuery = `INSERT INTO posts (id ,content, photo_link, user_id) VALUES ($1, $2, $3, $4) RETURNING *`

    GetPostByIDQuery = `SELECT * FROM posts WHERE id = $1`

    GetAllPostsQuery = `SELECT p.*, COUNT(l.user_id) as like_count FROM posts p 
    LEFT JOIN post_likes l ON p.id = l.post_id
    GROUP BY p.id
    ORDER BY like_count DESC`

    UpdatePostQuery = `UPDATE posts 
    SET content = $1, photo_link = $2, updated_at = NOW() 
    WHERE id = $3 RETURNING *`

    IsOwnerQuery = `SELECT EXISTS(
        SELECT 1 FROM posts
        WHERE id = $1 AND user_id = $2
    )`

    DeletePostQuery = `DELETE FROM posts WHERE id = $1 RETURNING *`
)
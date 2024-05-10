package repository

const (
    GetUserByIdQuery = `SELECT * FROM users WHERE id = $1`

    GetScoreTestQuery = `SELECT score_test FROM users WHERE id = $1`

    UpdateUserQuery = `UPDATE users 
	SET full_name = $1, email = $2, password = $3, birth_date = $4, photo_link = $5, score_test = $6, updated_at = NOW() 
	WHERE id = $7 RETURNING *`
)
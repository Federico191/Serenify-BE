package repository

const (
	createUserQuery = `INSERT INTO users (id, full_name, birth_date, email, password, verification_code) 
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING *`

	getUserByIdQuery = `SELECT * FROM users WHERE id = $1`

	getUserByEmailQuery = `SELECT * FROM users WHERE email = $1`

	getUserByVerificationCodeQuery = `SELECT * FROM users WHERE verification_code = $1`

	getExpiredVerificationCodeQuery = `SELECT * FROM users WHERE verification_code IS NOT NULL AND created_at < NOW() - INTERVAL '15 minutes' - INTERVAL '1 second'`

	updateUserQuery = `UPDATE users 
	SET full_name = $1, email = $2, password = $3, birth_date = $4, is_verified = $5, photo_link = $6, updated_at = NOW() 
	WHERE id = $7 RETURNING *`

	deleteVerificationCodeQuery = `UPDATE users SET verification_code = NULL WHERE email = $1 RETURNING *`
)

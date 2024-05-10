package repository

const (
    getAllSeminarsQuery = `SELECT * FROM seminars`

    getSeminarByIdQuery = `SELECT * FROM seminars WHERE id = $1`
)
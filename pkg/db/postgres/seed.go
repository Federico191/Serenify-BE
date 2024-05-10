package postgres

import (
	"FindIt/internal/entity"
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func SeedInit(db *sqlx.DB) {
    generateArticles(db)
    generateSeminars(db)
}

func generateArticles(db *sqlx.DB) {
	articles := []entity.Article{
		{
			ID:        uuid.New(),
			Title:     "Kesehatan Mental: Mengatasi Stres dan Kecemasan di Era Pandemi",
			Content:   "Kesehatan mental menjadi kunci penting dalam menghadapi tantangan pandemi. Artikel ini menjelaskan pentingnya menjaga kesehatan mental dan memberikan tips praktis untuk menghadapi stres dan kecemasan.",
			PhotoLink: sql.NullString{String: "https://djuvkqyzbqemumxpcnxr.supabase.co/storage/v1/object/public/article/3803041.jpg", Valid: true},
			CreatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			Title:     "Mengatasi Stres di Tempat Kerja",
			Content:   "Stres di tempat kerja dapat mempengaruhi produktivitas dan kesejahteraan karyawan. Temukan strategi praktis untuk mengelola stres dan meningkatkan kesejahteraan di lingkungan kerja.",
			PhotoLink: sql.NullString{String: "https://djuvkqyzbqemumxpcnxr.supabase.co/storage/v1/object/public/article/man-513529_640.jpg", Valid: true},
			CreatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			Title:     "Mengatasi Stres dan Kecemasan di Era Pandemi",
			Content:   "Stres dan kecemasan di era pandemi mempengaruhi kehidupan sehari-hari. Cari solusi yang tepat untuk mengatasi stres dan kecemasan di era pandemi.",
			PhotoLink: sql.NullString{String: "https://djuvkqyzbqemumxpcnxr.supabase.co/storage/v1/object/public/article/man-513529_640.jpg", Valid: true},
			CreatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			Title:     "Mengenali macam-macam Kesehatan Mental",
			Content:   "Kesehatan mental adalah kesehatan yang melibatkan emosi, pikuan, dan perasaan. Ini adalah bagian dari kesehatan yang harus dijaga.",
			PhotoLink: sql.NullString{String: "https://djuvkqyzbqemumxpcnxr.supabase.co/storage/v1/object/public/article/3803041.jpg", Valid: true},
			CreatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			Title:     "Edukasi mengenai self diagnosis kesehatan mental",
			Content:   "Bahaya self diagnosis kesehatan mental sering kali terbukti. Banyak sekali anak muda yang merasa terganggunya kesehatan mental. Hal tersebut mengancam kehidupan mereka.",
			PhotoLink: sql.NullString{String: "https://djuvkqyzbqemumxpcnxr.supabase.co/storage/v1/object/public/article/3803041.jpg", Valid: true},
			CreatedAt: time.Now(),
		},
	}

	for _, article := range articles {
		_, err := db.Exec("INSERT INTO articles (id, title, content, photo_link, created_at) VALUES ($1, $2, $3, $4, $5)",
			article.ID, article.Title, article.Content, article.PhotoLink, article.CreatedAt)
		if err != nil {
			panic(err)
		}
	}
    log.Println("Seed articles success")
}

func generateSeminars(db *sqlx.DB) {
	seminars := []entity.Seminar{
		{
			ID:          uuid.New(),
			Title:       "Muda, Berdaya, dan Bahagia: Menjawab Tantangan Kesehatan Mental di Era Perubahan",
			Time:        "2023-06-15 09:00:00",
			Place:       "Zoom Meeting",
			Price:       50000,
			Description: "Discover how to cope with stress and anxiety in a healthy and productive way. Join us for a 1:1 seminar on Muda, Berdaya, dan Bahagia.",
			PhotoLink:   sql.NullString{String: "https://djuvkqyzbqemumxpcnxr.supabase.co/storage/v1/object/public/seminar/seminar.jpg", Valid: true},
			CreatedAt:   time.Now(),
		},
		{
			ID:          uuid.New(),
			Title:       "Emotional Intelligence in the Workplace",
			Time:        "2023-07-20 14:00:00",
			Place:       "Hotel Bintang Lima, Malang",
			Price:       75000,
			Description: "Discover how to cultivate emotional intelligence and improve interpersonal relationships in a professional setting.",
			PhotoLink:   sql.NullString{String: "https://djuvkqyzbqemumxpcnxr.supabase.co/storage/v1/object/public/seminar/8b06b74e2acc4e47e01c4759d5ed4470.jpg", Valid: true},
			CreatedAt:   time.Now(),
		},
        {
            ID:          uuid.New(),
            Title:       "Mengenal Kesehatan Mental",
            Time:        "2021-06-15 09:00:00",
            Place:       "Zoom Meeting",
            Price:       50000,
            Description: "Belajar mengenal kesehatan mental dan bahayanya di era pandemi. Temukan solusi yang tepat untuk mengatasi stres dan kecemasan di era pandemi.",
            PhotoLink:   sql.NullString{String: "https://djuvkqyzbqemumxpcnxr.supabase.co/storage/v1/object/public/seminar/290f1b7df3b2653ac1cc8dc8bf16507c.jpg", Valid: true},
            CreatedAt:   time.Now(),
        },
        {
            ID:          uuid.New(),
            Title:       "Manajemen Kesehatan Mental di Era Media Sosial", 
            Time:        "2021-06-15 09:00:00",
            Place:       "Zoom Meeting",
            Price:       30000,
            Description: "Belajar tentang manajemen kesehatan mental di era media sosial. Mengenal strategi dan metode yang efektif untuk mengatasi stres dan kecemasan.",
            PhotoLink:   sql.NullString{String: "https://djuvkqyzbqemumxpcnxr.supabase.co/storage/v1/object/public/seminar/8b06b74e2acc4e47e01c4759d5ed4470.jpg", Valid: true},
            CreatedAt:   time.Now(),
        },
        {
            ID:          uuid.New(),
            Title:       "Menemukan kedamaian kesehatan mental dengan berolahraga",
            Time:        "2021-06-15 09:00:00",
            Place:       "Lapangan Rampal, Malang",
            Price:       20000,
            Description: "Bagaimana menemukan kedamaian kesehatan mental dengan berolahraga? Bergabunglah sekarang dan temukan solusi yang tepat.",
            PhotoLink:   sql.NullString{String: "https://djuvkqyzbqemumxpcnxr.supabase.co/storage/v1/object/public/seminar/290f1b7df3b2653ac1cc8dc8bf16507c.jpg", Valid: true},
            CreatedAt:   time.Now(),
        },
	}

    for _, seminar := range seminars {
        _, err := db.Exec("INSERT INTO seminars (id, title, time, place, price, description, photo_link, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
            seminar.ID, seminar.Title, seminar.Time, seminar.Place, seminar.Price, seminar.Description, seminar.PhotoLink, seminar.CreatedAt)
        if err != nil {
            panic(err)
        }
    }
    log.Println("Seed seminars success")
}

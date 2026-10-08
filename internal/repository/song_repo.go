package repository

import(
	"database/sql"
	"errors"
	"ecommerce/internal/models"
)

type SongRepository struct{
	DB *sql.DB
}

func NewSongRepository(db *sql.DB) *SongRepository{
	return &SongRepository{DB: db}
}

func(r *SongRepository) GetAll()([]*models.Song, error){
	rows, err := r.DB.Query(`
		SELECT id, title, artist, file_url, created_at
		FROM songs ORDER BY ASC
	`)

	if err != nil{
		return nil, err
	}

	defer rows.Close()

	var songs []*models.Song
	for rows.Next(){
		var s models.Song
		if err:= rows.Scan(&s.ID, &s.Title, &s.Artist, &s.FileURL, &s.CreatedAt) ;err != nil{
			return nil, err
		}
		songs = append(songs, &s)
	}
	return songs, rows.Err()
}


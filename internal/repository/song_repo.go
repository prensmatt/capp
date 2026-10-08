package repository

import(
	"database/sql"
	
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

func(r *SongRepository) Insert(s *models.Song) error{
	return r.DB.QueryRow(`INSERT INTO songs (title, artist, file_url) VALUES ($1, $2, $3) RETURNING id, created_at`,
	  s.Title, s.Artist, s.FileURL,
	).Scan(&s.ID, &s.CreatedAt)

}

func(r *SongRepository) Delete(id int) error{
	result, err := r.DB.Exec(`DELETE FROM songs WHERE id = $1`, id)
	if err != nil{
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil{
		return err
	}
	if rows == 0 {
		return models.ErrNotFound
	}
	return nil
}
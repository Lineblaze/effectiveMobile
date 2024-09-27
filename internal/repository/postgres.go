package postgresql

import (
	"effectiveMobile/pkg/storage/postgres"
	"fmt"
	openapi "github.com/Lineblaze/effective_gen"
	"strings"
)

//go:generate ifacemaker -f postgres.go -o ../repository.go -i Repository -s PostgresRepository -p internal -y "Controller describes methods, implemented by the repository package."
type PostgresRepository struct {
	db postgres.Postgres
}

func NewPostgresRepository(db postgres.Postgres) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (p *PostgresRepository) GetSongDetail(group, song string) (*openapi.SongDetail, error) {
	var songDetail openapi.SongDetail
	err := p.db.QueryRow(
		`SELECT release_date, text, link FROM songs_detail WHERE "group" = $1 AND song = $2`,
		group, song,
	).Scan(&songDetail.ReleaseDate, &songDetail.Text, &songDetail.Link)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch song detail: %w", err)
	}

	return &songDetail, nil
}

func (p *PostgresRepository) CreateSong(song *openapi.Song) (*openapi.Song, error) {
	query := `
		INSERT INTO songs ("group", song, release_date, text, link)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, "group", song, release_date, text, link
	`

	err := p.db.QueryRow(query, song.Group, song.Song, song.ReleaseDate, song.Text, song.Link).Scan(
		&song.Id,
		&song.Group,
		&song.Song,
		&song.ReleaseDate,
		&song.Text,
		&song.Link,
	)

	if err != nil {
		return nil, fmt.Errorf("inserting song: %w", err)
	}

	return song, nil
}

func (p *PostgresRepository) UpdateSong(songID string, req *openapi.UpdateSongBody) (*openapi.Song, error) {
	var args []any
	var fields []string
	argID := 1

	if req.Group != nil {
		fields = append(fields, fmt.Sprintf(`"group" = $%d`, argID))
		args = append(args, *req.Group)
		argID++
	}
	if req.Song != nil {
		fields = append(fields, fmt.Sprintf(`song = $%d`, argID))
		args = append(args, *req.Song)
		argID++
	}
	if req.ReleaseDate != nil {
		fields = append(fields, fmt.Sprintf(`release_date = $%d`, argID))
		args = append(args, *req.ReleaseDate)
		argID++
	}
	if req.Text != nil {
		fields = append(fields, fmt.Sprintf(`text = $%d`, argID))
		args = append(args, *req.Text)
		argID++
	}
	if req.Link != nil {
		fields = append(fields, fmt.Sprintf(`link = $%d`, argID))
		args = append(args, *req.Link)
		argID++
	}

	if len(fields) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	query := fmt.Sprintf(`
		UPDATE songs
		SET %s
		WHERE id = $%d
		RETURNING id, COALESCE("group", ''), COALESCE(song, ''), COALESCE(release_date, ''), COALESCE(text, ''), COALESCE(link, '')
	`, strings.Join(fields, ", "), argID)

	args = append(args, songID)

	var updatedSong openapi.Song
	err := p.db.QueryRow(query, args...).Scan(
		&updatedSong.Id,
		&updatedSong.Group,
		&updatedSong.Song,
		&updatedSong.ReleaseDate,
		&updatedSong.Text,
		&updatedSong.Link,
	)

	if err != nil {
		return nil, fmt.Errorf("updating song: %w", err)
	}

	return &updatedSong, nil
}

func (p *PostgresRepository) DeleteSong(songID string) error {
	_, err := p.db.Exec("DELETE FROM songs WHERE id = $1", songID)
	if err != nil {
		return fmt.Errorf("deleting song: %w", err)
	}
	return nil
}

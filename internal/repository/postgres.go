package postgresql

import (
	"effectiveMobile/pkg/logger"
	"effectiveMobile/pkg/storage/postgres"
	"fmt"
	openapi "github.com/Lineblaze/effective_mobile_gen"
	"strings"
)

//go:generate ifacemaker -f postgres.go -o ../repository.go -i Repository -s PostgresRepository -p internal -y "Controller describes methods, implemented by the repository package."
type PostgresRepository struct {
	db     postgres.Postgres
	logger *logger.ApiLogger
}

func NewPostgresRepository(db postgres.Postgres, logger *logger.ApiLogger) *PostgresRepository {
	return &PostgresRepository{db: db, logger: logger}
}

func (p *PostgresRepository) GetSongDetail(group, song string) (*openapi.SongDetail, error) {
	p.logger.Debugf("Fetching song detail for group: %s, song: %s", group, song)
	var songDetail openapi.SongDetail
	err := p.db.QueryRow(
		`SELECT release_date, text, link FROM songs_detail WHERE "group" = $1 AND song = $2`,
		group, song,
	).Scan(&songDetail.ReleaseDate, &songDetail.Text, &songDetail.Link)

	if err != nil {
		p.logger.Errorf("failed to fetch song detail: %v", err)
		return nil, fmt.Errorf("failed to fetch song detail: %v", err)
	}

	p.logger.Infof("Successfully fetched song detail for group: %s, song: %s", group, song)
	return &songDetail, nil
}

func (p *PostgresRepository) GetSongs(body *openapi.GetSongsBody) ([]*openapi.Song, error) {
	p.logger.Debug("Getting songs with filter parameters")

	query := `SELECT id, "group", song, release_date, "text", link FROM songs WHERE 1=1`
	var params []interface{}
	paramIdx := 1

	if body.Id != nil {
		query += fmt.Sprintf(" AND id = $%d", paramIdx)
		params = append(params, *body.Id)
		paramIdx++
	}
	if body.Group != nil {
		query += fmt.Sprintf(" AND \"group\" ILIKE $%d", paramIdx)
		params = append(params, "%"+*body.Group+"%")
		paramIdx++
	}
	if body.Song != nil {
		query += fmt.Sprintf(" AND song ILIKE $%d", paramIdx)
		params = append(params, "%"+*body.Song+"%")
		paramIdx++
	}
	if body.ReleaseDate != nil {
		query += fmt.Sprintf(" AND release_date = $%d", paramIdx)
		params = append(params, *body.ReleaseDate)
		paramIdx++
	}
	if body.Text != nil {
		query += fmt.Sprintf(" AND \"text\" ILIKE $%d", paramIdx)
		params = append(params, "%"+*body.Text+"%")
		paramIdx++
	}
	if body.Link != nil {
		query += fmt.Sprintf(" AND link ILIKE $%d", paramIdx)
		params = append(params, "%"+*body.Link+"%")
		paramIdx++
	}
	if body.Limit != nil {
		query += fmt.Sprintf(" LIMIT %d", *body.Limit)
	}
	if body.Offset != nil {
		query += fmt.Sprintf(" OFFSET %d", *body.Offset)
	}

	rows, err := p.db.Query(query, params...)
	if err != nil {
		p.logger.Errorf("failed to get songs: %v", err)
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	var songs []*openapi.Song
	for rows.Next() {
		var song openapi.Song
		if err = rows.Scan(&song.Id, &song.Group, &song.Song, &song.ReleaseDate, &song.Text, &song.Link); err != nil {
			p.logger.Errorf("failed to scan song: %v", err)
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		songs = append(songs, &song)
	}

	if err = rows.Err(); err != nil {
		p.logger.Errorf("error reading rows: %v", err)
		return nil, fmt.Errorf("row iteration error: %v", err)
	}

	p.logger.Infof("Successfully retrieved %d songs", len(songs))
	return songs, nil
}

func (p *PostgresRepository) GetSongText(group, song string) (string, error) {
	p.logger.Debugf("Fetching song text for group: %s, song: %s", group, song)
	var songText string
	query := `SELECT "text" FROM songs WHERE "group" = $1 AND song = $2`
	err := p.db.QueryRow(query, group, song).Scan(&songText)
	if err != nil {
		p.logger.Errorf("failed to fetch song text: %v", err)
		return "", fmt.Errorf("failed to get song text: %v", err)
	}

	p.logger.Infof("Successfully fetched song text for group: %s, song: %s", group, song)
	return songText, nil
}

func (p *PostgresRepository) CreateSong(song *openapi.Song) (*openapi.Song, error) {
	p.logger.Debugf("Creating song for group: %s, song: %s", song.Group, song.Song)

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
		p.logger.Errorf("failed to create song: %v", err)
		return nil, fmt.Errorf("inserting song: %v", err)
	}

	p.logger.Infof("Successfully created song for group: %s, song: %s", song.Group, song.Song)
	return song, nil
}

func (p *PostgresRepository) UpdateSong(songID string, req *openapi.UpdateSongBody) (*openapi.Song, error) {
	p.logger.Debugf("Updating song with ID: %s", songID)
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
		p.logger.Errorf("failed to update song: %v", err)
		return nil, fmt.Errorf("updating song: %v", err)
	}

	p.logger.Infof("Successfully updated song with ID: %s", songID)
	return &updatedSong, nil
}

func (p *PostgresRepository) DeleteSong(songID string) error {
	p.logger.Debugf("Deleting song with ID: %s", songID)
	_, err := p.db.Exec("DELETE FROM songs WHERE id = $1", songID)
	if err != nil {
		p.logger.Errorf("failed to delete song: %v", err)
		return fmt.Errorf("deleting song: %v", err)
	}

	p.logger.Infof("Successfully deleted song with ID: %s", songID)
	return nil
}

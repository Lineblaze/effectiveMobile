package usecase

import (
	repository "effectiveMobile/internal"
	"effectiveMobile/pkg/logger"
	"encoding/json"
	"fmt"
	openapi "github.com/Lineblaze/effective_mobile_gen"
	"net/http"
	"regexp"
	"strings"
)

//go:generate ifacemaker -f *.go -o ../usecase.go -i UseCase -s UseCase -p internal -y "Controller describes methods, implemented by the usecase package."
type UseCase struct {
	repo   repository.Repository
	logger *logger.ApiLogger
}

func NewUseCase(repo repository.Repository, logger *logger.ApiLogger) *UseCase {
	return &UseCase{repo: repo, logger: logger}
}

func (u *UseCase) FetchSongDetail(group, song string) (*openapi.SongDetail, error) {
	u.logger.Debugf("Fetching song detail for group: %s, song: %s", group, song)
	apiURL := fmt.Sprintf("http://localhost:8080/info?group=%s&song=%s", group, song)
	resp, err := http.Get(apiURL)
	if err != nil {
		u.logger.Errorf("failed to fetch song detail: %v", err)
		return nil, fmt.Errorf("failed to fetch song detail: %v", err)
	}
	defer resp.Body.Close()

	var detail openapi.SongDetail
	if err = json.NewDecoder(resp.Body).Decode(&detail); err != nil {
		u.logger.Errorf("failed to decode song detail: %v", err)
		return nil, fmt.Errorf("failed to decode song detail: %v", err)
	}

	u.logger.Infof("Successfully fetched song detail for group: %s, song: %s", group, song)
	return &detail, nil
}

func (u *UseCase) GetSongDetail(group, song string) (*openapi.SongDetail, error) {
	u.logger.Debugf("Getting song detail for group: %s, song: %s", group, song)
	songDetail, err := u.repo.GetSongDetail(group, song)
	if err != nil {
		u.logger.Infof("Successfully retrieved song detail for group: %s, song: %s", group, song)
		return nil, fmt.Errorf("getting bid by ID: %v", err)
	}

	u.logger.Infof("Successfully retrieved song detail for group: %s, song: %s", group, song)
	return songDetail, nil
}

func (u *UseCase) GetSongs(body *openapi.GetSongsBody) ([]*openapi.Song, error) {
	u.logger.Debug("Getting songs with filter parameters")
	songs, err := u.repo.GetSongs(body)
	if err != nil {
		u.logger.Errorf("error getting songs: %v", err)
		return nil, fmt.Errorf("getting songs: %v", err)
	}

	u.logger.Infof("Successfully retrieved %d songs", len(songs))
	return songs, nil
}

func (u *UseCase) GetSongText(body *openapi.GetSongTextBody) ([][]string, error) {
	u.logger.Debugf("Getting song text for group: %s, song: %s", body.Group, body.Song)
	songText, err := u.repo.GetSongText(body.Group, body.Song)
	if err != nil {
		u.logger.Errorf("error getting song text: %v", err)
		return nil, fmt.Errorf("getting song text: %v", err)
	}

	re := regexp.MustCompile("\\\\n\\\\n")
	split := re.Split(songText, -1)
	var verses [][]string

	for _, verse := range split {
		trimmedVerse := strings.TrimSpace(verse)
		if trimmedVerse != "" {
			lines := strings.Split(trimmedVerse, "\n")
			for i, line := range lines {
				lines[i] = strings.TrimSuffix(line, "\\n")
			}
			verses = append(verses, lines)
		}
	}

	start := 0
	if body.Offset != nil {
		start = int(*body.Offset)
	}
	if start > len(verses) {
		start = len(verses)
	}

	limit := 5
	if body.Limit != nil {
		limit = int(*body.Limit)
	}

	end := start + limit
	if end > len(verses) {
		end = len(verses)
	}

	u.logger.Infof("Returning %d verses starting from %d", len(verses[start:end]), start)
	return verses[start:end], nil
}

func (u *UseCase) CreateSong(req openapi.CreateSongBody, detail *openapi.SongDetail) (*openapi.Song, error) {
	u.logger.Debugf("Creating song for group: %s, song: %s", req.Group, req.Song)
	song := &openapi.Song{
		Group:       req.Group,
		Song:        req.Song,
		ReleaseDate: detail.ReleaseDate,
		Text:        detail.Text,
		Link:        detail.Link,
	}

	createdSong, err := u.repo.CreateSong(song)
	if err != nil {
		u.logger.Errorf("error creating song: %v", err)
		return nil, fmt.Errorf("creating song: %v", err)
	}

	u.logger.Infof("Successfully created song for group: %s, song: %s", req.Group, req.Song)
	return createdSong, nil
}

func (u *UseCase) UpdateSong(songID string, body *openapi.UpdateSongBody) (*openapi.Song, error) {
	u.logger.Debugf("Updating song with ID: %s", songID)
	updatedSong, err := u.repo.UpdateSong(songID, body)
	if err != nil {
		u.logger.Errorf("error updating song: %v", err)
		return nil, fmt.Errorf("updating song: %v", err)
	}

	u.logger.Infof("Successfully updated song with ID: %s", songID)
	return updatedSong, nil
}

func (u *UseCase) DeleteSong(songID string) error {
	u.logger.Debugf("Deleting song with ID: %s", songID)
	err := u.repo.DeleteSong(songID)
	if err != nil {
		u.logger.Errorf("error deleting song: %v", err)
		return fmt.Errorf("deleting song: %v", err)
	}

	u.logger.Infof("Successfully deleted song with ID: %s", songID)
	return nil
}

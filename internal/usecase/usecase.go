package usecase

import (
	repository "effectiveMobile/internal"
	"encoding/json"
	"fmt"
	openapi "github.com/Lineblaze/effective_mobile_gen"
	"net/http"
	"regexp"
	"strings"
)

//go:generate ifacemaker -f *.go -o ../usecase.go -i UseCase -s UseCase -p internal -y "Controller describes methods, implemented by the usecase package."
type UseCase struct {
	repo repository.Repository
}

func NewUseCase(repo repository.Repository) *UseCase {
	return &UseCase{repo: repo}
}

func (u *UseCase) FetchSongDetail(group, song string) (*openapi.SongDetail, error) {
	apiURL := fmt.Sprintf("http://localhost:8080/info?group=%s&song=%s", group, song)
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch song detail: %v", err)
	}
	defer resp.Body.Close()
	var detail openapi.SongDetail
	if err = json.NewDecoder(resp.Body).Decode(&detail); err != nil {
		return nil, fmt.Errorf("failed to decode song detail: %v", err)
	}
	return &detail, nil
}

func (u *UseCase) GetSongDetail(group, song string) (*openapi.SongDetail, error) {
	songDetail, err := u.repo.GetSongDetail(group, song)
	if err != nil {
		return nil, fmt.Errorf("getting bid by ID: %v", err)
	}
	return songDetail, nil
}

func (u *UseCase) GetSongs(body *openapi.GetSongsBody) ([]*openapi.Song, error) {
	songs, err := u.repo.GetSongs(body)
	if err != nil {
		return nil, fmt.Errorf("getting songs: %v", err)
	}
	return songs, nil
}

func (u *UseCase) GetSongText(body *openapi.GetSongTextBody) ([][]string, error) {
	songText, err := u.repo.GetSongText(body.Group, body.Song)
	if err != nil {
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

	return verses[start:end], nil
}

func (u *UseCase) CreateSong(req openapi.CreateSongBody, detail *openapi.SongDetail) (*openapi.Song, error) {
	song := &openapi.Song{
		Group:       req.Group,
		Song:        req.Song,
		ReleaseDate: detail.ReleaseDate,
		Text:        detail.Text,
		Link:        detail.Link,
	}

	createdSong, err := u.repo.CreateSong(song)
	if err != nil {
		return nil, fmt.Errorf("creating song: %v", err)
	}

	return createdSong, nil
}

func (u *UseCase) UpdateSong(songID string, body *openapi.UpdateSongBody) (*openapi.Song, error) {
	updatedSong, err := u.repo.UpdateSong(songID, body)
	if err != nil {
		return nil, fmt.Errorf("updating song: %v", err)
	}

	return updatedSong, nil
}

func (u *UseCase) DeleteSong(songID string) error {
	err := u.repo.DeleteSong(songID)
	if err != nil {
		return fmt.Errorf("deleting song: %v", err)
	}

	return nil
}

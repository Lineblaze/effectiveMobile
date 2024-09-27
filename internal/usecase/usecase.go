package usecase

import (
	repository "effectiveMobile/internal"
	"encoding/json"
	"fmt"
	openapi "github.com/Lineblaze/effective_gen"
	"net/http"
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
		return nil, fmt.Errorf("failed to fetch song detail: %w", err)
	}
	defer resp.Body.Close()
	var detail openapi.SongDetail
	if err = json.NewDecoder(resp.Body).Decode(&detail); err != nil {
		return nil, fmt.Errorf("failed to decode song detail: %w", err)
	}
	return &detail, nil
}

func (u *UseCase) GetSongDetail(group, song string) (*openapi.SongDetail, error) {
	songDetail, err := u.repo.GetSongDetail(group, song)
	if err != nil {
		return nil, fmt.Errorf("getting bid by ID: %w", err)
	}
	return songDetail, nil
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
		return nil, fmt.Errorf("creating song: %w", err)
	}

	return createdSong, nil
}

func (u *UseCase) UpdateSong(songID string, body *openapi.UpdateSongBody) (*openapi.Song, error) {
	updatedSong, err := u.repo.UpdateSong(songID, body)
	if err != nil {
		return nil, fmt.Errorf("updating song: %w", err)
	}

	return updatedSong, nil
}

func (u *UseCase) DeleteSong(songID string) error {
	err := u.repo.DeleteSong(songID)
	if err != nil {
		return fmt.Errorf("deleting song: %w", err)
	}

	return nil
}

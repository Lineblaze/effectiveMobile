// Code generated by ifacemaker; DO NOT EDIT.

package internal

import (
	openapi "github.com/Lineblaze/effective_mobile_gen"
)

// Controller describes methods, implemented by the repository package.
type Repository interface {
	GetSongDetail(group, song string) (*openapi.SongDetail, error)
	GetSongs(body *openapi.GetSongsBody) ([]*openapi.Song, error)
	GetSongText(group, song string) (string, error)
	CreateSong(song *openapi.Song) (*openapi.Song, error)
	UpdateSong(songID string, req *openapi.UpdateSongBody) (*openapi.Song, error)
	DeleteSong(songID string) error
}

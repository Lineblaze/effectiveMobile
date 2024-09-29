package http

import (
	handler "effectiveMobile/internal"
	"github.com/gofiber/fiber/v3"
)

func MapRoutes(r fiber.Router, h handler.Handler) {
	r.Get("/info", h.GetSongDetail())

	r.Get(`songs`, h.GetSongs())
	r.Get(`songs/text`, h.GetSongText())
	r.Post(`songs`, h.CreateSong())
	r.Patch(`songs/:songId`, h.UpdateSong())
	r.Delete(`songs/:songId`, h.DeleteSong())
}

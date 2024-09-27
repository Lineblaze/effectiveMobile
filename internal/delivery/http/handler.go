package http

import (
	useCase "effectiveMobile/internal"
	"effectiveMobile/pkg/logger"
	openapi "github.com/Lineblaze/effective_gen"
	"github.com/gofiber/fiber/v3"
)

//go:generate ifacemaker -f handler.go -o ../../handler.go -i Handler -s Handler -p internal -y "Controller describes methods, implemented by the http package."
type Handler struct {
	useCase useCase.UseCase
	logger  *logger.ApiLogger
}

func NewHandler(useCase useCase.UseCase, logger *logger.ApiLogger) *Handler {
	return &Handler{useCase: useCase, logger: logger}
}

func (h *Handler) GetSongDetail() fiber.Handler {
	return func(c fiber.Ctx) error {
		group := c.Query("group")
		if group == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "group is required"})
		}
		song := c.Query("song")
		if song == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "song is required"})
		}
		songDetail, err := h.useCase.GetSongDetail(group, song)
		if err != nil {
			h.logger.Errorf("Failed to get songDetail %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "InternalServerError"})
		}

		return c.Status(fiber.StatusOK).JSON(songDetail)
	}
}

//func (h Handler) GetSongs() fiber.Handler {
//	return func(ctx fiber.Ctx) error {
//		var body openapi.GetSongsBody
//		if err := ctx.Bind().Body(&body); err != nil {
//			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
//		}
//		songs, err := h.useCase.GetSongs(&body)
//		if err != nil {
//			h.logger.Errorf("Failed to get songs %v", err)
//			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "InternalServerError"})
//		}
//		return ctx.Status(fiber.StatusOK).JSON(songs)
//	}
//}

func (h Handler) CreateSong() fiber.Handler {
	return func(ctx fiber.Ctx) error {
		var req openapi.CreateSongBody
		if err := ctx.Bind().Body(&req); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
		}

		songDetail, err := h.useCase.FetchSongDetail(req.Group, req.Song)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch song detail"})
		}

		createdSong, err := h.useCase.CreateSong(req, songDetail)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "InternalServerError"})
		}

		return ctx.Status(fiber.StatusOK).JSON(createdSong)
	}
}

func (h Handler) UpdateSong() fiber.Handler {
	return func(ctx fiber.Ctx) error {
		songID := ctx.Params("songId")
		var req openapi.UpdateSongBody
		if err := ctx.Bind().Body(&req); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
		}
		updatedSong, err := h.useCase.UpdateSong(songID, &req)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "InternalServerError"})
		}
		return ctx.Status(fiber.StatusOK).JSON(updatedSong)
	}
}

func (h Handler) DeleteSong() fiber.Handler {
	return func(ctx fiber.Ctx) error {
		songID := ctx.Params("songId")
		err := h.useCase.DeleteSong(songID)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "InternalServerError"})
		}
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Song deleted successfully"})
	}
}

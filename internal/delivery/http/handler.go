package http

import (
	useCase "effectiveMobile/internal"
	"effectiveMobile/pkg/logger"
	openapi "github.com/Lineblaze/effective_mobile_gen"
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
			h.logger.Debug("group query param is missing")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "group is required"})
		}
		song := c.Query("song")
		if song == "" {
			h.logger.Debug("song query param is missing")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "song is required"})
		}

		songDetail, err := h.useCase.GetSongDetail(group, song)
		if err != nil {
			h.logger.Errorf("Failed to get songDetail %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "InternalServerError"})
		}

		h.logger.Infof("Successfully fetched song detail for group: %s, song: %s", group, song)
		return c.Status(fiber.StatusOK).JSON(songDetail)
	}
}

func (h Handler) GetSongs() fiber.Handler {
	return func(ctx fiber.Ctx) error {
		var body openapi.GetSongsBody
		if err := ctx.Bind().Body(&body); err != nil {
			h.logger.Debug("Failed to parse GetSongs request body")
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
		}

		h.logger.Infof("Fetching songs with body: %+v", body)
		songs, err := h.useCase.GetSongs(&body)
		if err != nil {
			h.logger.Errorf("Failed to get songs %v", err)
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "InternalServerError"})
		}

		h.logger.Infof("Successfully fetched songs, count: %d", len(songs))
		return ctx.Status(fiber.StatusOK).JSON(songs)
	}
}

func (h *Handler) GetSongText() fiber.Handler {
	return func(ctx fiber.Ctx) error {
		var body openapi.GetSongTextBody
		if err := ctx.Bind().Body(&body); err != nil {
			h.logger.Debug("Failed to parse GetSongText request body")
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
		}
		if body.Group == "" || body.Song == "" {
			h.logger.Debug("Missing group or song in GetSongText request")
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "group and song are required"})
		}

		h.logger.Infof("Fetching text for group: %s, song: %s", body.Group, body.Song)
		verses, err := h.useCase.GetSongText(&body)
		if err != nil {
			h.logger.Errorf("Failed to get song verses: %v", err)
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
		}

		h.logger.Infof("Successfully fetched verses for group: %s, song: %s", body.Group, body.Song)
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"verses": verses})
	}
}

func (h Handler) CreateSong() fiber.Handler {
	return func(ctx fiber.Ctx) error {
		var req openapi.CreateSongBody
		if err := ctx.Bind().Body(&req); err != nil {
			h.logger.Debug("Failed to parse CreateSong request body")
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
		}

		h.logger.Infof("Creating song for group: %s, song: %s", req.Group, req.Song)
		songDetail, err := h.useCase.FetchSongDetail(req.Group, req.Song)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch song detail"})
		}

		createdSong, err := h.useCase.CreateSong(req, songDetail)
		if err != nil {
			h.logger.Errorf("Failed to create song: %v", err)
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "InternalServerError"})
		}

		h.logger.Infof("Successfully created song for group: %s, song: %s", req.Group, req.Song)
		return ctx.Status(fiber.StatusOK).JSON(createdSong)
	}
}

func (h Handler) UpdateSong() fiber.Handler {
	return func(ctx fiber.Ctx) error {
		songID := ctx.Params("songId")
		var req openapi.UpdateSongBody
		h.logger.Debug("Failed to parse UpdateSong request body")
		if err := ctx.Bind().Body(&req); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
		}

		h.logger.Infof("Updating song with ID: %s", songID)
		updatedSong, err := h.useCase.UpdateSong(songID, &req)
		if err != nil {
			h.logger.Errorf("Failed to update song: %v", err)
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "InternalServerError"})
		}

		h.logger.Infof("Successfully updated song with ID: %s", songID)
		return ctx.Status(fiber.StatusOK).JSON(updatedSong)
	}
}

func (h Handler) DeleteSong() fiber.Handler {
	return func(ctx fiber.Ctx) error {
		songID := ctx.Params("songId")

		h.logger.Infof("Deleting song with ID: %s", songID)
		err := h.useCase.DeleteSong(songID)
		if err != nil {
			h.logger.Errorf("Failed to delete song: %v", err)
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "InternalServerError"})
		}

		h.logger.Infof("Successfully deleted song with ID: %s", songID)
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Song deleted successfully"})
	}
}

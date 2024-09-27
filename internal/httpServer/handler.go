package httpServer

import (
	"effectiveMobile/internal/delivery/http"
	repository "effectiveMobile/internal/repository"
	useCase "effectiveMobile/internal/usecase"
	"effectiveMobile/pkg/logger"
	storage "effectiveMobile/pkg/storage/postgres"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	serverLogger "github.com/gofiber/fiber/v3/middleware/logger"
)

func (s *Server) MapHandlers(app *fiber.App, logger *logger.ApiLogger) error {
	db, err := storage.InitPsqlDB(s.cfg)
	if err != nil {
		return err
	}

	repo := repository.NewPostgresRepository(db)
	useCase := useCase.NewUseCase(repo)
	handler := http.NewHandler(useCase, logger)

	app.Use(serverLogger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{},
		AllowHeaders: []string{},
	}))

	group := app.Group("")
	http.MapRoutes(group, handler)

	return nil
}

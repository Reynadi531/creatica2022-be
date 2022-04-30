package main

import (
	"flag"
	"os"

	"github.com/Reynadi531/creatica2022-be/pkg/database"
	"github.com/Reynadi531/creatica2022-be/pkg/middleware"
	"github.com/Reynadi531/creatica2022-be/pkg/routes"
	"github.com/Reynadi531/creatica2022-be/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	DEV_MODE = flag.Bool("dev", false, "development mode")
)

func init() {
	flag.Parse()
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	if *DEV_MODE {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

		log.Info().Msg("Running in development mode")
	}
}

func main() {
	app := fiber.New()

	middleware.RegisterMiddleware(app)

	db := database.InitDatabase()

	routes.RegisterRouteAuth(app, db)
	routes.RegisterRoutePost(app, db)

	if *DEV_MODE {
		utils.StartServer(app)
	} else {
		utils.StartServerWithGracefulShutdown(app)
	}
}

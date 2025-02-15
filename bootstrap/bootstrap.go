package bootstrap

import (
	"io"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
	"github.com/sgitwhyd/jagong/app/ws"
	"github.com/sgitwhyd/jagong/pkg/database"
	"github.com/sgitwhyd/jagong/pkg/env"
	"github.com/sgitwhyd/jagong/pkg/router"
)

func NewApplication() *fiber.App {
	env.SetupEnvFile()
	SetupLogFile()

	database.SetupDatabase()
	database.SetupMongoDB()


	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{Views: engine})
	app.Use(recover.New())
	app.Use(logger.New())
	app.Get("/dashboard", monitor.New())

	go ws.ServeWsMessaging(app)

	router.InstallRouter(app)
	

	return app
}

func SetupLogFile(){
	logFile, err := os.OpenFile("./logs/jagong_app.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}

	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
}

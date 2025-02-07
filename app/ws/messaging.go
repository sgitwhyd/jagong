package ws

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/sgitwhyd/jagong/app/models"
	"github.com/sgitwhyd/jagong/app/repository"
	"github.com/sgitwhyd/jagong/pkg/env"
)

func ServeWsMessaging(app *fiber.App) {
	var clients = make(map[*websocket.Conn]bool)
	var broadcast = make(chan models.MessagePayload)

	app.Get("/message/v1/send", websocket.New(func(conn *websocket.Conn) {
		defer func() {
			conn.Close()
			delete(clients, conn)
		}()

		clients[conn] = true
		for {
			var msg models.MessagePayload
			if err := conn.ReadJSON(&msg); err != nil {
				log.Printf("msg payload %v", msg)
				log.Printf("error payload %v", err.Error())
				break
			}
			msg.Date = time.Now()
			err := repository.InsertMessage(context.Background(), msg)
			if err != nil {
				log.Printf("error insert message %v", err.Error())
			}
			broadcast <- msg
		}
	}))

	go func() {
		for {
			msg := <-broadcast
			for client := range clients {
				err := client.WriteJSON(msg)
				if err != nil {
					log.Printf("error sending message %v", err.Error())
					client.Close()
					delete(clients, client)
				}
			}
		}
	}()

	log.Fatal(app.Listen(fmt.Sprintf("%s:%s", env.GetEnv("APP_HOST_SOCKET", "localhost"), env.GetEnv("APP_PORT_SOCKET", "8080"))))
}

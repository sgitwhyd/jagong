package ws

import (
	"fmt"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/sgitwhyd/jagong/pkg/env"
	"log"
)

type MessagePayload struct {
	From    string `json:"from"`
	Message string `json:"message"`
}

func ServeWsMessaging(app *fiber.App) {
	var clients = make(map[*websocket.Conn]bool)
	var broadcast = make(chan MessagePayload)

	app.Get("/message/v1/send", websocket.New(func(conn *websocket.Conn) {
		defer func() {
			conn.Close()
			delete(clients, conn)
		}()

		clients[conn] = true
		for {
			var msg MessagePayload
			if err := conn.ReadJSON(&msg); err != nil {
				fmt.Printf("error payload %v", err.Error())
				break
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
					fmt.Printf("error sending message %v", err.Error())
					client.Close()
					delete(clients, client)
				}
			}
		}
	}()

	log.Fatal(app.Listen(fmt.Sprintf("%s:%s", env.GetEnv("APP_HOST_SOCKET", "localhost"), env.GetEnv("APP_PORT_SOCKET", "8080"))))
}

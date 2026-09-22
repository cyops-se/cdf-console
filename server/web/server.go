package web

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"server/logger"
	"server/routes"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/websocket/v2"
)

//go:embed static/index.html
var admin string

//go:embed all:static/*
var static embed.FS

// websocket connections
type WebSocketMessage struct {
	Topic   string      `json:"topic"`
	Message interface{} `json:"message"`
}

type WebSocketClient struct {
	connection *websocket.Conn
	close      chan string
}

var dropList []int
var ws []*websocket.Conn
var wsMutex sync.Mutex
var clients = make(map[*websocket.Conn]*WebSocketClient)
var register = make(chan *WebSocketClient, 5)
var unregister = make(chan *websocket.Conn, 50)
var broadcast = make(chan *WebSocketMessage, 5)

func handlePanic() {
	if r := recover(); r != nil {
		logger.Log("error", "Server panic recovery", fmt.Sprintf("%#v", r))
		return
	}
}

func RunWeb() {
	defer handlePanic()

	go runSocketActions()

	// http.FS can be used to create a http Filesystem
	subFS2, _ := fs.Sub(static, "static")
	var staticFS = http.FS(subFS2)

	// Set a file transfer limit to 50MB
	app := fiber.New(fiber.Config{StrictRouting: true, BodyLimit: 50 * 1024 * 1024})
	app.Use("/", filesystem.New(filesystem.Config{
		Root:   staticFS,
		Browse: true,
	}))

	app.Get("/help/*", func(ctx *fiber.Ctx) error {
		ctx.Status(200)
		ctx.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
		return ctx.SendString(admin)
	})

	app.Get("/ui/*", func(ctx *fiber.Ctx) error {
		ctx.Status(200)
		ctx.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
		return ctx.SendString(admin)
	})

	// WebSocket registration
	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		client := &WebSocketClient{c, make(chan string)}

		register <- client
		close := <-client.close
		logger.Log("trace", "Websocket close signal received", fmt.Sprintf("%v", close))
	}))

	api := app.Group("/api")
	routes.RegisterRoutes(api)

	app.Listen(fmt.Sprintf(":%d", 3080))

	select {}
}

func runSocketActions() {
	for {
		select {
		case client := <-register:
			clients[client.connection] = client
			logger.Log("trace", "New websocket connection registered", "")

		case msg := <-broadcast:
			for connection, client := range clients {
				if err := connection.WriteJSON(msg); err != nil {
					client.close <- "timetoexit"
					unregister <- connection
					connection.WriteMessage(websocket.CloseMessage, []byte{})
					connection.Close()
				}
			}

		case connection := <-unregister:
			delete(clients, connection)
			logger.Log("trace", "Websocket connection unregistered", "")
		}
	}
}

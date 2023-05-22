package routes

import (
	"bufio"
	"fmt"

	"github.com/JustGritt/go-grpc/broadcast"
	"github.com/JustGritt/go-grpc/database"
	"github.com/JustGritt/go-grpc/models"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

type Stream struct {
	Notifier chan []models.Payment
}

// var once sync.Once
var b = broadcast.NewBroker[models.Payment]()

func GetStream(c *fiber.Ctx) error {
	ctx := c.Context()
	ctx.SetContentType("text/event-stream")
	ctx.Response.Header.Set("Cache-Control", "no-cache")
	ctx.Response.Header.Set("Connection", "keep-alive")
	ctx.Response.Header.Set("Transfer-Encoding", "chunked")
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	ctx.Response.Header.Set("Access-Control-Allow-Headers", "Cache-Control")
	ctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")

	var payments []models.Payment
	database.Database.Db.Last(&payments)
	fmt.Println("payment ", payments)

	if len(payments) == 0 {
		return c.Status(404).JSON("No payments found")
	}

	go b.Start()
	ctx.SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		msgCh := b.Subscribe()
		var i int
		for {
			i++
			fmt.Fprintf(w, "%d\n", <-msgCh)
			w.Flush()
		}
	}))

	return nil
}

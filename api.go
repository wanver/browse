package browse

import (
	"context"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/sync/semaphore"
)

// Global semaphore to ensure only one browser at a time
var browserSemaphore = semaphore.NewWeighted(1)

func App() error {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(400).JSON(map[string]string{"error": err.Error()})
		},
	})

	app.Post("/browse", handleBrowseRequest)

	return app.Listen(":3000")
}

func handleBrowseRequest(c *fiber.Ctx) error {
	var req BrowseRequest
	if err := c.BodyParser(&req); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Minute)
	defer cancel()

	if err := browserSemaphore.Acquire(ctx, 1); err != nil {
		return errors.New("server busy")
	}
	defer browserSemaphore.Release(1)

	page, err := New(&req, c.Context())
	if err != nil {
		return err
	}
	defer page.Close()
	defer page.Browser().Close()

	resp, err := req.Hijack(page)
	if err != nil {
		return err
	}

	err = page.Navigate(req.PageURL)
	if err != nil {
		return err
	}
	page.WaitLoad()

	for _, bri := range req.Instructions {
		_, err := bri.Act(page)
		if err != nil && bri.Fatal {
			return err
		}
	}
	time.Sleep(30 * time.Second)

	return c.JSON(resp)
}

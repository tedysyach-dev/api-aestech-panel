package middlewares

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func RequestLogger(log *logrus.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// lanjut ke handler berikutnya
		err := c.Next()

		// log setelah response selesai
		log.Infof("[%s] %s | Status: %d | IP: %s | Duration: %v",
			c.Method(),
			c.OriginalURL(),
			c.Response().StatusCode(),
			c.IP(),
			time.Since(start),
		)

		return err
	}
}

package main

import (
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})

	log.SetLevel(logrus.InfoLevel)

	_ = godotenv.Load()

	dsn := os.Getenv("SENTRY_DSN")

	if dsn == "" {
		log.WithField("reason", "SENTRY_DSN not set").
			Warn("sentry_disabled")
	} else {
		log.WithField("dsn_present", true).
			Info("initializing_sentry")

		if err := sentry.Init(sentry.ClientOptions{
			Dsn:              dsn,
			TracesSampleRate: 1.0,
			Debug:            true,
		}); err != nil {
			log.WithError(err).
				Error("sentry_init_failed")
		}
	}
	defer sentry.Flush(2 * time.Second)

	router := setupRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.WithField("default_port", port).
			Warn("port_not_set")
	}

	log.WithField("port", port).
		Info("server_starting")

	if err := router.Run(":" + port); err != nil {
		log.WithError(err).
			Fatal("server_failed")
	}
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	router.Use(sentrygin.New(sentrygin.Options{
		Repanic: true,
	}))

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	return router
}

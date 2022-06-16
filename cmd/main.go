package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.uber.org/dig"
	"random-numbers/internal/adapters"
	"random-numbers/internal/common"
	"random-numbers/internal/logic"
	"random-numbers/internal/presenter"
)

func main() {
	c := dig.New()

	err := buildContainer(c)
	if err != nil {
		panic(err)
	}

	if err := c.Invoke(func(e *gin.Engine) error {
		return e.Run("0.0.0.0:8080")
	}); err != nil {
		panic(err)
	}
}

func buildContainer(c *dig.Container) error {
	if err := c.Provide(common.NewConfig); err != nil {
		return err
	}

	if err := c.Provide(logic.NewService); err != nil {
		return err
	}

	if err := c.Provide(logic.NewSessionWorker); err != nil {
		return err
	}

	if err := c.Provide(presenter.NewREST); err != nil {
		return err
	}

	if err := c.Provide(func(r *presenter.RestAPI) (*gin.Engine, error) {
		router := gin.New()

		router.Use(r.LoggingMiddleware)

		router.GET("/generate", r.AuthMiddleware, r.StatsMiddleware, r.Generate)
		router.GET("/login", r.Login)
		router.GET("/register", r.Register)
		router.GET("/details", r.Details)

		return router, nil
	}); err != nil {

	}

	if err := c.Provide(func(cfg *common.Config) (*logrus.Entry, error) {
		logger := logrus.New()

		logger.SetFormatter(&logrus.JSONFormatter{
			DisableTimestamp: true,
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyLevel: "severity",
			},
		})

		level, err := logrus.ParseLevel(cfg.LogLevel)
		if err != nil {
			return nil, err
		}

		logger.SetLevel(level)

		return logrus.NewEntry(logger), nil
	}); err != nil {
		return err
	}

	if err := c.Provide(adapters.NewRandomOrg); err != nil {
		return err
	}

	if err := c.Provide(adapters.NewUserPersistence); err != nil {
		return err
	}

	if err := c.Provide(adapters.NewSessionPersistence); err != nil {
		return err
	}

	return nil
}

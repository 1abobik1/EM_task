// cmd/server/main.go
package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/1abobik1/EM_task/config"
	"github.com/1abobik1/EM_task/internal/api"
	"github.com/1abobik1/EM_task/internal/handler"
	"github.com/1abobik1/EM_task/internal/middleware"
	"github.com/1abobik1/EM_task/internal/repository/postgres"
	"github.com/1abobik1/EM_task/internal/service"

	_ "github.com/1abobik1/EM_task/docs"
)

func init() {
	binding.EnableDecoderDisallowUnknownFields = true // отвергает лишние поля у DTO при запросе
}

// @title           Persons API
// @version         1.0
// @description     Сервис выполненный по тз https://drive.google.com/file/d/1zUU44O1ye5-3yYRdhLMEhpG9JXz-TWeQ/view

// @host      localhost:8080
// @BasePath  /

func main() {
	cfg := config.MustLoad()

	storage, err := postgres.NewPostgresStorageProd(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}

	httpClient := &http.Client{
		Timeout: 3 * time.Second,
	}

	// внешние апи
	ageCl := api.NewAgifyClient(cfg.AgifyURL, httpClient)
	genCl := api.NewGenderizeClient(cfg.GenderizeURL, httpClient)
	natCl := api.NewNationalizeClient(cfg.NationalizeURL, httpClient)

	// сервисный слой
	svc := service.NewPersonService(storage, ageCl, genCl, natCl)

	// хендлерный слой
	personHandler := handler.NewPersonHandler(svc)

	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	persons := router.Group("/persons")
	{
		persons.POST("", personHandler.CreatePerson)
		persons.GET("", middleware.StrictQueryParamsMiddleware(middleware.AllowedParams), personHandler.ListPersons)
		persons.GET("/:id", personHandler.GetPersonByID)
		persons.PUT("/:id", personHandler.UpdatePerson)
		persons.PATCH("/:id", personHandler.UpdatePerson)
		persons.DELETE("/:id", personHandler.DeletePerson)
	}

	addr := ":" + cfg.Port
	logrus.Infof("starting server on %s", addr)
	if err := router.Run(addr); err != nil {
		logrus.Fatalf("server failed: %v", err)
	}
}

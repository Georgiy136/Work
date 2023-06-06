package app

import (
	"fmt"
	"log"
	"myapp/config"
	"myapp/internal/handler"
	"myapp/internal/usecase"
	"myapp/internal/usecase/repository"
	"myapp/pkg/postgres"
	"myapp/pkg/redis"

	"github.com/gin-gonic/gin"
)

func Run(cfg *config.Config) {

	// Repository
	pg, err := postgres.New(cfg.Postgres)
	if err != nil {
		log.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	rdb, err := redis.New(cfg.Redis)
	if err != nil {
		log.Fatal(fmt.Errorf("app - Run - redis.New: %w", err))
	}
	defer rdb.Close()

	authRepositoryPostgres := repository.NewAuthPostgres(pg)
	authRepositoryRedis := repository.NewAuthRedis(rdb)

	projectRepository := repository.NewProject(pg)
	operatorRepository := repository.NewOperator(pg)

	// Use case
	authUseCases := usecase.NewAuthUsecases(authRepositoryPostgres, authRepositoryRedis, cfg.Auth)

	projectUseCases := usecase.NewProjectUsecases(projectRepository)
	operatorUseCases := usecase.NewOperatorUsecases(operatorRepository)

	// HTTP Server
	router := gin.Default()

	handler.NewRouter(router, *authUseCases, *operatorUseCases, *projectUseCases)

	router.Run(fmt.Sprintf("localhost:%d", cfg.Http.Port))
}

package handler

import (
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "myapp/docs"
	"myapp/internal/usecase"
)

func NewRouter(router *gin.Engine, as usecase.AuthUseCases, os usecase.OperatorUseCases, ps usecase.ProjectUseCases) {
	authHandlers := &AuthHandler{
		us: as,
	}
	operatorHandlers := &OperatorHandler{
		us: os,
	}
	projectHandlers := &ProjectHandler{
		us: ps,
	}

	// Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Routers
	operator := router.Group("/operator", authHandlers.clientIdentity(), authHandlers.clientHasRights())
	{
		operator.POST("/", operatorHandlers.PostOperator)
		operator.GET("/", operatorHandlers.GetAllOperators)
		operator.GET("/:id", operatorHandlers.GetOneOperator)
		operator.PUT("/:id", operatorHandlers.UpdateOperator)
		operator.DELETE("/:id", operatorHandlers.DeleteOperator)
	}

	project := router.Group("/project", authHandlers.clientIdentity(), authHandlers.clientHasRights())
	{
		project.POST("/", projectHandlers.PostProject)
		project.GET("/", projectHandlers.GetAllProjects)
		project.GET("/:id", projectHandlers.GetOneProject)
		project.PUT("/:id", projectHandlers.UpdateProject)
		project.DELETE("/:id", projectHandlers.DeleteProject)

	}
	router.PUT("/AddOperatorToProject/:project_id/:operator_id", authHandlers.clientIdentity(), authHandlers.clientHasRights(), projectHandlers.AddOperatorToProject)
	router.PUT("/DelOperatorFromProject/:project_id/:operator_id", authHandlers.clientIdentity(), authHandlers.clientHasRights(), projectHandlers.DeleteOperatorFromProject)

	router.POST("/sign-up", authHandlers.RegistClient)
	router.POST("/sign-in", authHandlers.AuthClient)
}

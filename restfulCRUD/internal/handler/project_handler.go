package handler

import (
	"fmt"
	"log"
	"myapp/internal/models"
	"myapp/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProjectHandler struct {
	us usecase.ProjectUseCases
}

// PostProject godoc
//
//	@Security		ApiKeyAuth
//	@Summary		Add Project to database
//	@Tags			Projects
//	@Description	create project
//	@ID				create-project
//	@Accept			json
//	@Produce		json
//	@Param			input	body		models.Project	true	"Project info"
//	@Success		201		{object}	models.Project
//	@Router			/project [post]
func (h *ProjectHandler) PostProject(c *gin.Context) {

	type PostProjectRequest struct {
		Id          uuid.UUID          `json:"project_id"`
		ProjectName string             `json:"project_Name" binding:"required"`
		ProjectType string             `json:"project_Type" binding:"required"`
		Operators   []*models.Operator `json:"operators" bun:"-"`
	}

	PostProject := &PostProjectRequest{}

	if err := c.Bind(PostProject); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	project := &models.Project{
		Id:          PostProject.Id,
		ProjectName: PostProject.ProjectName,
		ProjectType: PostProject.ProjectType,
		Operators:   PostProject.Operators,
	}

	project, err := h.us.AddProject(c.Request.Context(), *project)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, project)
}

// GetOneProject godoc
//
//	@Security		ApiKeyAuth
//	@Summary		Retrieves Project based on given ID
//	@Tags			Projects
//	@Description	get project by id
//	@ID				get-project-by-id
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Project ID"
//	@Success		202	{object}	models.Project
//	@Router			/project/{id} [get]
func (h *ProjectHandler) GetOneProject(c *gin.Context) {
	id := c.Param("id")
	projects, err := h.us.GetOneProject(c.Request.Context(), id)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, projects)
}

// GetAllProjects godoc
//
//	@Security		ApiKeyAuth
//	@Summary		Retrieves All Projects
//	@Tags			Projects
//	@Description	get all projects
//	@ID				get-all-projects
//	@Accept			json
//	@Produce		json
//	@Success		202	{array}	[]models.Project
//	@Router			/project [get]
func (h *ProjectHandler) GetAllProjects(c *gin.Context) {
	projects, err := h.us.GetAllProjects(c.Request.Context())
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, projects)
}

// DeleteProject godoc
//
//	@Security		ApiKeyAuth
//	@Summary		Delete Project based on given ID
//	@Tags			Projects
//	@Description	delete project by id
//	@ID				delete-project-by-id
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"Project ID"
//	@Success		200
//	@Router			/project/{id} [delete]
func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	id := c.Param("id")
	err := h.us.DeleteProject(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, fmt.Sprintf("Запись проекта с id = %s успешно удалена", id))
}

// UpdateProject godoc
//
//	@Security		ApiKeyAuth
//	@Summary		Update Project based on given ID
//	@Tags			Projects
//	@Description	update project by id
//	@ID				update-project-by-id
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string			true	"Project ID"
//	@Param			input	body		models.Project	true	"Project info"
//	@Success		201		{object}	models.Project
//	@Router			/project/{id} [put]
func (h *ProjectHandler) UpdateProject(c *gin.Context) {

	id := c.Param("id")

	type PostProjectRequest struct {
		Id          uuid.UUID          `json:"project_id"`
		ProjectName string             `json:"project_Name" binding:"required"`
		ProjectType string             `json:"project_Type" binding:"required"`
		Operators   []*models.Operator `json:"operators" bun:"-"`
	}

	PostProject := &PostProjectRequest{}

	if err := c.ShouldBindJSON(PostProject); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	project := &models.Project{
		Id:          PostProject.Id,
		ProjectName: PostProject.ProjectName,
		ProjectType: PostProject.ProjectType,
		Operators:   PostProject.Operators,
	}

	p, err := h.us.UpdateProject(c.Request.Context(), id, *project)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, p)
}

// Add operator to project godoc
//
//	@Security		ApiKeyAuth
//	@Summary		Add Operator To Project
//	@Tags			Projects
//	@Description	add operator to project
//	@ID				add-operator-to-project
//	@Accept			json
//	@Produce		json
//	@Param			project_id	path		string	true	"Project ID"
//	@Param			operator_id	path		string	true	"Operator ID"
//	@Success		201			{object}	models.Project
//	@Router			/AddOperatorToProject/{project_id}/{operator_id} [put]
func (h *ProjectHandler) AddOperatorToProject(c *gin.Context) {
	project_id := c.Param("project_id")
	operator_id := c.Param("operator_id")
	p, err := h.us.AddOperatorToProject(c.Request.Context(), project_id, operator_id)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, p)
}

// Delete operator from project godoc
//
//	@Security		ApiKeyAuth
//	@Summary		Delete Operator From Project
//	@Tags			Projects
//	@Description	delete operator from project
//	@ID				delete-operator-from-project
//	@Accept			json
//	@Produce		json
//	@Param			project_id	path		string	true	"Project ID"
//	@Param			operator_id	path		string	true	"Operator ID"
//	@Success		201			{object}	models.Project
//	@Router			/DelOperatorFromProject/{project_id}/{operator_id} [put]
func (h *ProjectHandler) DeleteOperatorFromProject(c *gin.Context) {
	project_id := c.Param("project_id")
	operator_id := c.Param("operator_id")
	p, err := h.us.DeleteOperatorFromProject(c.Request.Context(), project_id, operator_id)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, p)
}

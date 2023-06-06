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

type OperatorHandler struct {
	us usecase.OperatorUseCases
}

// PostOperator godoc
//
//	@Security		ApiKeyAuth
//	@Summary		Add Operator to database
//	@Tags			Operators
//	@Description	create Operator
//	@ID				create-operator
//	@Accept			json
//	@Produce		json
//	@Param			input	body		PostOperatorRequest	true	"Operator info"
//	@Success		201		{object}	models.Operator
//	@Router			/operator [post]
func (h *OperatorHandler) PostOperator(c *gin.Context) {

	type PostOperatorRequest struct {
		Id         uuid.UUID `json:"id"`
		FirstName  string    `json:"firstName" binding:"required"`
		LastName   string    `json:"lastName" binding:"required"`
		Patronymic string    `json:"patronymic" binding:"required"`
		City       string    `json:"city" binding:"required"`
		Phone      string    `json:"phone" binding:"required"`
		Email      string    `json:"email" binding:"required"`
		Password   string    `json:"password"`
	}

	postOperatorRequest := &PostOperatorRequest{}

	if err := c.Bind(postOperatorRequest); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	operator := &models.Operator{
		Id:         postOperatorRequest.Id,
		FirstName:  postOperatorRequest.FirstName,
		LastName:   postOperatorRequest.LastName,
		Patronymic: postOperatorRequest.Patronymic,
		City:       postOperatorRequest.City,
		Phone:      postOperatorRequest.Phone,
		Email:      postOperatorRequest.Email,
		Password:   postOperatorRequest.Password,
	}

	operator, err := h.us.AddOperator(c.Request.Context(), *operator)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, operator)
}

// GetOneOperator godoc
//
//	@Security		ApiKeyAuth
//	@Summary		Retrieves Operator based on given ID
//	@Tags			Operators
//	@Description	get operator by id
//	@ID				get-operator-by-id
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Operator ID"
//	@Success		202	{object}	models.Operator
//	@Router			/operator/{id} [get]
func (h *OperatorHandler) GetOneOperator(c *gin.Context) {
	id := c.Param("id")
	projects, err := h.us.GetOneOperator(c.Request.Context(), id)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, projects)
}

// GetAllOperators godoc
//
//	@Security		ApiKeyAuth
//	@Summary		Retrieves All Operators
//	@Tags			Operators
//	@Description	get all operators
//	@ID				get-all-operators
//	@Accept			json
//	@Produce		json
//	@Success		202	{array}	[]models.Operator
//	@Router			/operator [get]
func (h *OperatorHandler) GetAllOperators(c *gin.Context) {
	operators, err := h.us.GetAllOperators(c.Request.Context())
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, operators)
}

// DeleteOperator godoc
//
//	@Security		ApiKeyAuth
//	@Summary		Delete Operator based on given ID
//	@Tags			Operators
//	@Description	delete operator by id
//	@ID				delete-operator-by-id
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"Operator ID"
//	@Success		200
//	@Router			/operator/{id} [delete]
func (h *OperatorHandler) DeleteOperator(c *gin.Context) {
	id := c.Param("id")
	err := h.us.DeleteOperator(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, fmt.Sprintf("Запись оператора с id = %s успешно удалена", id))
}

// UpdateOperator godoc
//
//	@Security		ApiKeyAuth
//	@Summary		Update Operator based on given ID
//	@Tags			Operators
//	@Description	update operator by id
//	@ID				update-operator-by-id
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string				true	"Operator ID"
//	@Param			input	body		PutOperatorRequest	true	"Operator info"
//	@Success		201		{object}	models.Operator
//	@Router			/operator/{id} [put]
func (h *OperatorHandler) UpdateOperator(c *gin.Context) {

	type PutOperatorRequest struct {
		Id         uuid.UUID `json:"id"`
		FirstName  string    `json:"firstName" binding:"required"`
		LastName   string    `json:"lastName" binding:"required"`
		Patronymic string    `json:"patronymic" binding:"required"`
		City       string    `json:"city" binding:"required"`
		Phone      string    `json:"phone" binding:"required"`
		Email      string    `json:"email" binding:"required"`
		Password   string    `json:"password"`
	}

	id := c.Param("id")

	postOperatorRequest := &PutOperatorRequest{}

	if err := c.Bind(postOperatorRequest); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	operator := &models.Operator{
		Id:         postOperatorRequest.Id,
		FirstName:  postOperatorRequest.FirstName,
		LastName:   postOperatorRequest.LastName,
		Patronymic: postOperatorRequest.Patronymic,
		City:       postOperatorRequest.City,
		Phone:      postOperatorRequest.Phone,
		Email:      postOperatorRequest.Email,
		Password:   postOperatorRequest.Password,
	}

	operator, err := h.us.UpdateOperator(c.Request.Context(), id, *operator)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, operator)
}

package handler

import (
	"myapp/internal/models"
	"myapp/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	us usecase.AuthUseCases
}

// Регистрация клиента
//	@Summary		SignUp
//	@Tags			auth
//	@Description	create account
//	@ID				create-account
//	@Accept			json
//	@Produce		json
//	@Param			input	body		RegistClientRequest	true	"account info"
//	@Success		200		{integer}	integer				1
//	@Router			/sign-up [post]
func (h *AuthHandler) RegistClient(c *gin.Context) {

	type RegistClientRequest struct {
		Id         uuid.UUID `json:"id"`
		ClientRole string    `json:"role" `
		Username   string    `json:"username" binding:"required,min=1"`
		Login      string    `json:"login" binding:"required,min=1"`
		Password   string    `json:"password" binding:"required,min=4"`
	}

	RegistClient := &RegistClientRequest{}

	if err := c.Bind(RegistClient); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	client := &models.Client{
		Id:         RegistClient.Id,
		ClientRole: RegistClient.ClientRole,
		Username:   RegistClient.Username,
		Login:      RegistClient.Login,
		Password:   RegistClient.Password,
	}

	err := h.us.RegistClient(c.Request.Context(), *client)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, nil)
}

// Аутентификация клиента
//	@Summary		SignIn
//	@Tags			auth
//	@Description	login
//	@ID				login
//	@Accept			json
//	@Produce		json
//	@Param			input	body		Client	true	"credentials"
//	@Success		200		{string}	string	"token"
//	@Router			/sign-in [post]
func (h *AuthHandler) AuthClient(c *gin.Context) {

	type Client struct {
		Login    string `json:"login" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	client := &Client{}

	if err := c.BindJSON(&client); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	token, err := h.us.GenerateToken(c.Request.Context(), client.Login, client.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//c.Header("Authorization", token)

	c.JSON(http.StatusOK, gin.H{"token": token})
}

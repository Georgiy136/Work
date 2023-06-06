package handler

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Аутентификация пользователя
func (h *AuthHandler) clientIdentity() gin.HandlerFunc {
	return func(c *gin.Context) {
		headerValue := c.GetHeader("Authorization")

		if headerValue == "" {
			log.Printf("AuthHandler - clientIdentity - c.GetHeader: %s", "empty auth header")
			c.AbortWithStatusJSON(http.StatusUnauthorized, fmt.Errorf("AuthHandler - clientIdentity - c.GetHeader: %s", "empty auth header"))
			return
		}

		headerParts := strings.Split(headerValue, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			log.Printf("AuthHandler - clientIdentity - c.GetHeader: %s", "invalid auth header")
			c.AbortWithStatusJSON(http.StatusUnauthorized, fmt.Errorf("AuthHandler - clientIdentity - c.GetHeader: %s", "invalid auth header"))
			return
		}

		if len(headerParts[1]) == 0 {
			log.Printf("AuthHandler - clientIdentity - c.GetHeader: %s", "token is empty")
			c.AbortWithStatusJSON(http.StatusUnauthorized, fmt.Errorf("AuthHandler - clientIdentity - c.GetHeader: %s", "token is empty"))
			return
		}

		userId, tokenLifetime, err := h.us.ParseToken(headerParts[1])
		if err != nil {
			log.Printf("AuthHandler - clientIdentity - h.us.ParseToken: %v \n", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, fmt.Errorf("AuthHandler - clientIdentity - h.us.ParseToken: %v", err))
			return
		}

		if time.Now().Unix() > tokenLifetime {
			c.AbortWithStatusJSON(http.StatusUnauthorized, fmt.Errorf("AuthHandler - clientIdentity - h.us.ParseToken: %s", "token has expired"))
			return
		}

		_, err = h.us.GetOneClientById(c.Request.Context(), userId)
		if err != nil {
			log.Printf("AuthHandler - clientIdentity - GetOneClientById: %v \n", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, fmt.Errorf("AuthHandler - clientIdentity - GetOneClientById: %v", err))
			return
		}

		log.Println("id:", userId)

		c.Set("userId", userId)

		c.Next()
	}
}

func (h *AuthHandler) clientHasRights() gin.HandlerFunc {
	return func(c *gin.Context) {

		userId, exist := c.Get("userId")

		if !exist {
			log.Printf("AuthHandler - clientIdentity - GetOneClientById: %s \n", "can not find authorized user")
			c.AbortWithStatusJSON(http.StatusForbidden, fmt.Errorf("AuthHandler - clientIdentity - GetOneClientById: %s", "can not find authorized user"))
			return
		}

		client, err := h.us.GetOneClientById(c.Request.Context(), fmt.Sprint(userId))
		if err != nil {
			log.Printf("AuthHandler - clientIdentity - GetOneClientById: %v \n", err)
			c.AbortWithStatusJSON(http.StatusForbidden, fmt.Errorf("AuthHandler - clientIdentity - GetOneClientById: %v", err))
			return
		}

		rights, err := h.us.GetRoleRights(c.Request.Context(), client.ClientRole)
		if err != nil {
			log.Printf("AuthHandler - clientIdentity - h.us.GetRoleRights: %v \n", err)
			c.AbortWithStatusJSON(http.StatusForbidden, fmt.Errorf("AuthHandler - clientIdentity - h.us.GetRoleRights: %v", err))
			return
		}

		log.Println(rights)

		for _, method := range rights {
			if method == c.Request.Method {
				c.Next()
				return
			}
		}
		log.Printf("AuthHandler - clientHasRights - c.Request.Method: %s", "client has no right")
		c.AbortWithStatusJSON(http.StatusForbidden, fmt.Errorf("AuthHandler - clientHasRights - c.Request.Method: %s", "client has no right"))
		return

	}
}

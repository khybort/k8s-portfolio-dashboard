package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"github.com/gin-gonic/gin"
)

func Auth(authServiceURL string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Extract token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		token := parts[1]

		// Verify token with auth service
		reqBody, _ := json.Marshal(map[string]string{"token": token})
		resp, err := http.Post(authServiceURL+"/api/v1/auth/verify", "application/json", bytes.NewBuffer(reqBody))
		if err != nil || resp.StatusCode != http.StatusOK {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		var verifyResp struct {
			Valid   bool   `json:"valid"`
			UserID  string `json:"user_id"`
			Role    string `json:"role"`
		}
		json.NewDecoder(resp.Body).Decode(&verifyResp)

		if !verifyResp.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("user_id", verifyResp.UserID)
		c.Set("role", verifyResp.Role)
		c.Set("token", token)

		c.Next()
	}
}


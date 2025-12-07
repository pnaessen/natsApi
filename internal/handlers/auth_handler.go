package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	"golang.org/x/oauth2"

	//"natsApi/internal/models"
	"natsApi/internal/config"
	"natsApi/internal/models"
	"natsApi/internal/utils"
)

type AuthHandler struct {
	Config   *oauth2.Config
	NatsConn *nats.Conn
}

func NewAuthHandler(nc *nats.Conn) *AuthHandler {
	env, err := config.LoadEnv()
	if err != nil {
		fmt.Println("Error loading env:", err)
	}

	return &AuthHandler{
		NatsConn: nc,
		Config: &oauth2.Config{
			ClientID:     env.ClientID,
			ClientSecret: env.ClientSecret,
			RedirectURL:  env.RedirectURL,
			Scopes:       []string{"public"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://api.intra.42.fr/oauth/authorize",
				TokenURL: "https://api.intra.42.fr/oauth/token",
			},
		},
	}
}

func (h *AuthHandler) Login(c *gin.Context) {

	url := h.Config.AuthCodeURL("random") // need to change
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *AuthHandler) CallBack(c *gin.Context) {

	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Code not found"})
		return
	}

	token, err := h.Config.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("failed to exchange token: %w", err).Error()})
		return
	}

	client := h.Config.Client(context.Background(), token)
	resp, err := client.Get("https://api.intra.42.fr/v2/me")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("failed to get user info: %w", err).Error()})
		return
	}
	defer resp.Body.Close()

	var user42Data models.User42

	if err := json.NewDecoder(resp.Body).Decode(&user42Data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("failed to decode user")})
	}

	userMessage := models.UserMessage{
		Username:   user42Data.Username,
		Email:      user42Data.Email,
		IntraID:    user42Data.ID,
		SchoolYear: user42Data.School_year,
		IsActive:   user42Data.Is_active,
	}

	reqBytes, err := json.Marshal(userMessage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "JSON marshal error"})
		return
	}

	msg, err := h.NatsConn.Request("user.login", reqBytes, 2*time.Second)
	if err != nil {
		fmt.Printf("NATS Error: %v\n", err)
		c.JSON(http.StatusGatewayTimeout, gin.H{"error": "No worker on user"})
		return
	}

	var workerResponse models.UserMessage
	if err := json.Unmarshal(msg.Data, &workerResponse); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid response from worker"})
		return
	}

	tokenString, err := utils.GenerateJWT(workerResponse.Db_id, workerResponse.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "JWT generation failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"db_id":   workerResponse.Db_id,
		"user":    workerResponse.Username,
		"token":   tokenString,
	})
}

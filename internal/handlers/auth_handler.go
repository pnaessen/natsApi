package handlers

import (
	"context"
	"fmt"
	"net/http"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"

	//"natsApi/internal/models"
	"natsApi/internal/config"
	"natsApi/internal/models"
)

type AuthHandler struct {
	Config *oauth2.Config
}

func NewAuthHandler() *AuthHandler {
	env, err := config.LoadEnv()
	if err != nil {
		// return nil
	}

	return &AuthHandler{
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

	// return en json le user42Data
}
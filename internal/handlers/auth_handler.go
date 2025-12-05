package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"

	//"natsApi/internal/models"
	"natsApi/internal/config"
)

type AuthHandler struct {
	Config *oauth2.Config
}

func NewAuthHandler() *AuthHandler {
	env, err := config.LoadEnv()
	if err != nil {

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

	url := h.Config.AuthCodeURL("random")
	c.Redirect(http.StatusTemporaryRedirect, url)
}

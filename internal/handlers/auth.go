package handlers

import (
	"net/http"

	"goauth/internal/auth"
	requestmodel "goauth/internal/models/request"
	"goauth/internal/services"

	"github.com/gin-gonic/gin"
)

type IAuthHandler interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	Refresh(c *gin.Context)
	MainAuthPage(c *gin.Context)
	StartGoogleOAuth(c *gin.Context)
	OAuthGoogleCallback(c *gin.Context)
}

type authHandler struct {
	authService services.IAuthService
}

func NewAuthHandler(authService services.IAuthService) IAuthHandler {
	return &authHandler{
		authService: authService,
	}
}

func (ah authHandler) Register(c *gin.Context) {
	var (
		payload requestmodel.RegisterPayload
		conn    = connection(c)
	)

	if err := c.ShouldBind(&payload); err != nil {
		c.JSON(400, gin.H{"error": err})
		return
	}

	// check email password
	if err := ah.authService.Register(conn, payload); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "register success"})
}

func (ah authHandler) Login(c *gin.Context) {
	var payload requestmodel.LoginPayload

	if err := c.ShouldBind(&payload); err != nil {
		c.JSON(400, gin.H{"error bind payload": err.Error()})
		return
	}

	conn := connection(c)
	token, err := ah.authService.Login(conn, payload)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "login success", "data": token})
}

func (ah authHandler) Refresh(c *gin.Context) {
	var (
		payload requestmodel.RefreshPayload
		conn    = connection(c)
	)

	if err := c.ShouldBind(&payload); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	token, err := ah.authService.RefreshToken(conn, payload)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "resfresh token success",
		"data":    gin.H{"token": token},
	})
}

func (ah authHandler) MainAuthPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func (ah authHandler) StartGoogleOAuth(c *gin.Context) {
	//  Get Google OAuth configuration
	config := auth.GetGoogleOauthConfig()

	// Generate the OAuth and redirect to the user to it
	url := config.AuthCodeURL(auth.OauthStateString)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (ah authHandler) OAuthGoogleCallback(c *gin.Context) {
	var (
		conn  = connection(c)
		state = c.Query("state")
		code  = c.Query("code")
	)
	res, err := ah.authService.HandleGoogleOAuthCallback(conn, state, code)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "login success", "data": res})
}

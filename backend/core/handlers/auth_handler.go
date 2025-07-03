package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gsarmaonline/goiter/config"
	"github.com/gsarmaonline/goiter/core/models"
)

type (
	TokenResponse struct {
		AccessToken  string `json:"access_token"`
		IDToken      string `json:"id_token"`
		ExpiresIn    int    `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
		TokenType    string `json:"token_type"`
		Scope        string `json:"scope"`
	}
)

func (h *Handler) handleShortCircuitLogin(c *gin.Context) {
	type (
		ShortCircuitLogin struct {
			Email string `json:"email"`
		}
	)
	if h.cfg.Mode != config.ModeDev {
		c.JSON(403, gin.H{"error": "Short circuit login is only allowed in development mode"})
		return
	}
	req := &ShortCircuitLogin{}
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	someRandomNumber := strconv.Itoa(rand.Int())
	// Create or update user in database
	user := models.User{
		GoogleID:   someRandomNumber,
		Email:      req.Email,
		Name:       fmt.Sprintf("User %s", req.Email),
		UserStatus: models.ActiveUser,
	}
	modUser := &models.User{}

	if result := h.db.Where(models.User{Email: req.Email}).FirstOrCreate(&modUser); result.Error != nil {
		c.JSON(500, gin.H{"error": "Failed to save user"})
		return
	}
	user.ID = modUser.ID

	if err := h.db.Save(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to update user status"})
		return
	}

	// Set session cookie with more permissive settings for development
	c.SetCookie(
		"session",
		someRandomNumber,
		3600*24*7, // 7 days
		"/",
		"",
		false, // Set to false for development
		false, // Set to false for development
	)

	c.JSON(200, gin.H{
		"message": "Short circuit login successful",
		"code":    someRandomNumber,
	})

}

func (h *Handler) handleGoogleLogin(c *gin.Context) {
	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	if googleClientID == "" {
		c.JSON(500, gin.H{"error": "Google client ID not configured"})
		return
	}

	callbackURL := os.Getenv("GOOGLE_CALLBACK_URL")
	if callbackURL == "" {
		c.JSON(500, gin.H{"error": "Google callback URL not configured"})
		return
	}

	redirectURL := fmt.Sprintf(
		"https://accounts.google.com/o/oauth2/v2/auth?client_id=%s&redirect_uri=%s&response_type=code&scope=email profile&access_type=offline",
		googleClientID,
		callbackURL,
	)

	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

func (h *Handler) handleGoogleCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(400, gin.H{"error": "No code provided"})
		return
	}

	// Exchange code for tokens
	token, err := h.exchangeCodeForToken(code)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to exchange code for token"})
		return
	}

	// Get user info from Google
	userInfo, err := h.getGoogleUserInfo(token.AccessToken)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get user info"})
		return
	}

	// Create or update user in database
	user := models.User{
		GoogleID:    userInfo.ID,
		Email:       userInfo.Email,
		Name:        userInfo.Name,
		Picture:     userInfo.Picture,
		AccessToken: token.AccessToken,
		UserStatus:  models.ActiveUser,
	}
	modUser := &models.User{}

	if result := h.db.Where(models.User{Email: userInfo.Email}).FirstOrCreate(&modUser); result.Error != nil {
		c.JSON(500, gin.H{"error": "Failed to save user"})
		return
	}
	user.ID = modUser.ID

	if err := h.db.Save(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to update user status"})
		return
	}

	// Set session cookie with more permissive settings for development
	c.SetCookie(
		"session",
		userInfo.ID,
		3600*24*7, // 7 days
		"/",
		"",
		false, // Set to false for development
		false, // Set to false for development
	)

	// Redirect to frontend with success parameter
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		c.JSON(500, gin.H{"error": "Frontend URL not configured"})
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, frontendURL+"?code=success")
}

func (h *Handler) handleGetUser(c *gin.Context) {
	user := h.GetUserFromContext(c)

	c.JSON(200, gin.H{
		"id":      user.ID,
		"email":   user.Email,
		"name":    user.Name,
		"picture": user.Picture,
		"profile": user.Profile,
	})
}

func (h *Handler) handleLogout(c *gin.Context) {
	// Clear the session cookie
	c.SetCookie(
		"session",
		"",
		-1, // Expire immediately
		"/",
		"",
		false, // Set to false for development
		false, // Set to false for development
	)
	c.JSON(200, gin.H{"message": "Logged out successfully"})
}

func (h *Handler) exchangeCodeForToken(code string) (*TokenResponse, error) {
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	callbackURL := os.Getenv("GOOGLE_CALLBACK_URL")

	if callbackURL == "" {
		return nil, fmt.Errorf("Google callback URL not configured")
	}

	url := "https://oauth2.googleapis.com/token"
	data := fmt.Sprintf(
		"client_id=%s&client_secret=%s&code=%s&grant_type=authorization_code&redirect_uri=%s",
		clientID,
		clientSecret,
		code,
		callbackURL,
	)

	resp, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var token TokenResponse
	if err := json.Unmarshal(body, &token); err != nil {
		return nil, err
	}

	return &token, nil
}

func (h *Handler) getGoogleUserInfo(accessToken string) (*models.GoogleUserInfo, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + accessToken)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userInfo models.GoogleUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}

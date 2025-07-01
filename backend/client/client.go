package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"
)

type GoiterClient struct {
	BaseURL    string
	httpClient *http.Client
	sessionID  string
}

type User struct {
	ID      uint   `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

type Project struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	AccountID   uint   `json:"account_id"`
	UserID      uint   `json:"user_id"`
}

type Account struct {
	ID                 uint   `json:"id"`
	Name               string `json:"name"`
	Description        string `json:"description"`
	PlanID             uint   `json:"plan_id"`
	SubscriptionStatus string `json:"subscription_status"`
}

// NewGoiterClient creates a new client instance
func NewGoiterClient(baseURL string) *GoiterClient {
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	return &GoiterClient{
		BaseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Login initiates the Google OAuth flow
func (c *GoiterClient) Login() error {
	fmt.Println("🔐 Goiter Client Login")
	fmt.Println("======================")
	fmt.Printf("Please follow these steps to authenticate:\n\n")

	fmt.Printf("1. Open your browser and visit: %s/auth/google\n", c.BaseURL)
	fmt.Printf("2. Complete the Google OAuth flow\n")
	fmt.Printf("3. After successful login, you'll be redirected to the frontend\n")
	fmt.Printf("4. Open your browser's developer tools (F12)\n")
	fmt.Printf("5. Go to Application/Storage -> Cookies\n")
	fmt.Printf("6. Find the 'session' cookie and copy its value\n")
	fmt.Printf("7. Paste the session cookie value here: ")

	// Read the session cookie from user input
	var sessionCookie string
	fmt.Scanln(&sessionCookie)

	if sessionCookie == "" {
		return fmt.Errorf("no session cookie provided")
	}

	// Set the session cookie
	c.sessionID = sessionCookie

	// Test the authentication by making a request to /me
	fmt.Println("\n🔄 Testing authentication...")
	user, err := c.GetUser()
	if err != nil {
		return fmt.Errorf("authentication failed: %v", err)
	}

	fmt.Printf("✅ Login successful! Welcome, %s (%s)\n", user.Name, user.Email)
	return nil
}

// openBrowser opens the default browser to the given URL
func (c *GoiterClient) openBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

// makeRequest makes an authenticated HTTP request
func (c *GoiterClient) makeRequest(method, endpoint string, body io.Reader) (*http.Response, error) {
	if c.sessionID == "" {
		return nil, fmt.Errorf("not authenticated - please login first")
	}

	req, err := http.NewRequest(method, c.BaseURL+endpoint, body)
	if err != nil {
		return nil, err
	}

	// Add session cookie
	req.AddCookie(&http.Cookie{
		Name:  "session",
		Value: c.sessionID,
	})

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return c.httpClient.Do(req)
}

// GetUser retrieves the current user information
func (c *GoiterClient) GetUser() (*User, error) {
	resp, err := c.makeRequest("GET", "/me", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get user: %s", string(body))
	}

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

// GetProjects retrieves all projects for the current user
func (c *GoiterClient) GetProjects() ([]Project, error) {
	resp, err := c.makeRequest("GET", "/projects", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get projects: %s", string(body))
	}

	var projects []Project
	if err := json.NewDecoder(resp.Body).Decode(&projects); err != nil {
		return nil, err
	}

	return projects, nil
}

// GetAccount retrieves the current user's account
func (c *GoiterClient) GetAccount() (*Account, error) {
	resp, err := c.makeRequest("GET", "/account", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get account: %s", string(body))
	}

	var account Account
	if err := json.NewDecoder(resp.Body).Decode(&account); err != nil {
		return nil, err
	}

	return &account, nil
}

// Ping tests the connection to the server
func (c *GoiterClient) Ping() error {
	resp, err := http.Get(c.BaseURL + "/ping")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned status: %d", resp.StatusCode)
	}

	return nil
}

// Logout clears the session
func (c *GoiterClient) Logout() error {
	if c.sessionID == "" {
		return nil
	}

	resp, err := c.makeRequest("POST", "/logout", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	c.sessionID = ""
	fmt.Println("Logged out successfully!")
	return nil
}

// Example usage and CLI interface
func Run() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	// Initialize client
	baseURL := os.Getenv("GOITER_BASE_URL")
	client := NewGoiterClient(baseURL)

	command := os.Args[1]

	switch command {
	case "ping":
		if err := client.Ping(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Server is running!")

	case "login":
		if err := client.Login(); err != nil {
			fmt.Printf("Login failed: %v\n", err)
			os.Exit(1)
		}

	case "user":
		if err := client.Login(); err != nil {
			fmt.Printf("Login failed: %v\n", err)
			os.Exit(1)
		}

		user, err := client.GetUser()
		if err != nil {
			fmt.Printf("Error getting user: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("User: %s (%s)\n", user.Name, user.Email)

	case "projects":
		if err := client.Login(); err != nil {
			fmt.Printf("Login failed: %v\n", err)
			os.Exit(1)
		}

		projects, err := client.GetProjects()
		if err != nil {
			fmt.Printf("Error getting projects: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Projects (%d):\n", len(projects))
		for _, project := range projects {
			fmt.Printf("- %s: %s\n", project.Name, project.Description)
		}

	case "account":
		if err := client.Login(); err != nil {
			fmt.Printf("Login failed: %v\n", err)
			os.Exit(1)
		}

		account, err := client.GetAccount()
		if err != nil {
			fmt.Printf("Error getting account: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Account: %s\n", account.Name)
		fmt.Printf("Description: %s\n", account.Description)
		fmt.Printf("Subscription: %s\n", account.SubscriptionStatus)

	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: go run client.go <command>")
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("  ping     - Test server connection")
	fmt.Println("  login    - Login with Google OAuth")
	fmt.Println("  user     - Get current user info")
	fmt.Println("  projects - List user's projects")
	fmt.Println("  account  - Get account information")
	fmt.Println("")
	fmt.Println("Environment variables:")
	fmt.Println("  GOITER_BASE_URL - Server URL (default: http://localhost:8080)")
}

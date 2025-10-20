package testutils

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

// MockStripeServer provides a mock Stripe server for testing
type MockStripeServer struct {
	Server   *httptest.Server
	Requests []MockStripeRequest
}

// MockStripeRequest captures the details of requests made to the mock server
type MockStripeRequest struct {
	Method string
	Path   string
	Body   string
}

// NewMockStripeServer creates a new mock Stripe server
func NewMockStripeServer(t *testing.T) *MockStripeServer {
	mock := &MockStripeServer{
		Requests: make([]MockStripeRequest, 0),
	}

	mock.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Capture the request
		bodyBytes := make([]byte, 0)
		if r.Body != nil {
			bodyBytes, _ = io.ReadAll(r.Body)
		}

		mock.Requests = append(mock.Requests, MockStripeRequest{
			Method: r.Method,
			Path:   r.URL.Path,
			Body:   string(bodyBytes),
		})

		// Mock responses based on the path
		switch r.URL.Path {
		case "/v1/products":
			w.Header().Set("Content-Type", "application/json")
			response := map[string]interface{}{
				"id":   "prod_test123",
				"name": "Test Product",
			}
			json.NewEncoder(w).Encode(response)

		case "/v1/prices":
			w.Header().Set("Content-Type", "application/json")
			response := map[string]interface{}{
				"id":       "price_test123",
				"product":  "prod_test123",
				"currency": "usd",
				"recurring": map[string]interface{}{
					"interval": "month",
				},
			}
			json.NewEncoder(w).Encode(response)

		case "/v1/subscriptions":
			w.Header().Set("Content-Type", "application/json")
			response := map[string]interface{}{
				"id":     "sub_test123",
				"status": "active",
			}
			json.NewEncoder(w).Encode(response)

		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))

	return mock
}

// Close closes the mock server
func (m *MockStripeServer) Close() {
	if m.Server != nil {
		m.Server.Close()
	}
}

// GetRequestCount returns the number of requests made to the mock server
func (m *MockStripeServer) GetRequestCount() int {
	return len(m.Requests)
}

// GetRequestsForPath returns all requests made to a specific path
func (m *MockStripeServer) GetRequestsForPath(path string) []MockStripeRequest {
	requests := make([]MockStripeRequest, 0)
	for _, req := range m.Requests {
		if req.Path == path {
			requests = append(requests, req)
		}
	}
	return requests
}

// MockGoogleServer provides a mock Google OAuth server for testing
type MockGoogleServer struct {
	Server *httptest.Server
}

// NewMockGoogleServer creates a new mock Google OAuth server
func NewMockGoogleServer(t *testing.T) *MockGoogleServer {
	mock := &MockGoogleServer{}

	mock.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/token":
			// Mock OAuth token exchange
			w.Header().Set("Content-Type", "application/json")
			response := map[string]interface{}{
				"access_token":  "mock_access_token",
				"token_type":    "Bearer",
				"expires_in":    3600,
				"refresh_token": "mock_refresh_token",
			}
			json.NewEncoder(w).Encode(response)

		case "/oauth2/v2/userinfo":
			// Mock user info endpoint
			w.Header().Set("Content-Type", "application/json")
			response := map[string]interface{}{
				"id":             "mock_google_id",
				"email":          "mockuser@gmail.com",
				"verified_email": true,
				"name":           "Mock User",
				"picture":        "https://example.com/picture.jpg",
			}
			json.NewEncoder(w).Encode(response)

		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))

	return mock
}

// Close closes the mock server
func (m *MockGoogleServer) Close() {
	if m.Server != nil {
		m.Server.Close()
	}
}

// HTTPTestHelper provides utilities for HTTP testing
type HTTPTestHelper struct{}

// NewHTTPTestHelper creates a new HTTP test helper
func NewHTTPTestHelper() *HTTPTestHelper {
	return &HTTPTestHelper{}
}

// AssertContentType asserts that the response has the expected content type
func (h *HTTPTestHelper) AssertContentType(t *testing.T, w *httptest.ResponseRecorder, expectedType string) {
	actualType := w.Header().Get("Content-Type")
	require.Contains(t, actualType, expectedType)
}

// AssertRedirect asserts that the response is a redirect to the expected URL
func (h *HTTPTestHelper) AssertRedirect(t *testing.T, w *httptest.ResponseRecorder, expectedURL string) {
	require.Equal(t, http.StatusTemporaryRedirect, w.Code)
	location := w.Header().Get("Location")
	require.Contains(t, location, expectedURL)
}

// ParseJSONResponse parses the JSON response body into a map
func (h *HTTPTestHelper) ParseJSONResponse(t *testing.T, w *httptest.ResponseRecorder) map[string]interface{} {
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	return response
}
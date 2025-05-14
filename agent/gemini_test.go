package agent

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/vhybZApp/api/config"
	"github.com/vhybZApp/api/models"
	"google.golang.org/genai"
)

func setupTestRouter() *gin.Engine {
	config.LoadConfig()
	Init()
	log.Println(config.AppConfig.GeminiAPIKey)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/html", MakeHTML)
	return router
}

func TestMakeHTML_Success(t *testing.T) {
	router := setupTestRouter()

	// Create a test request
	reqBody := HTMLRequest{
		Contents: []*genai.Content{
			genai.Text("Create a simple button")[0],
		},
	}
	jsonBody, _ := json.Marshal(reqBody)

	// Create a test request
	req, _ := http.NewRequest("POST", "/html", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder
	w := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, w.Code)

	var response HTMLResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response.HTML)
}

func TestMakeHTML_EmptyContents(t *testing.T) {
	router := setupTestRouter()

	// Create a test request with empty contents
	reqBody := HTMLRequest{
		Contents: []*genai.Content{},
	}
	jsonBody, _ := json.Marshal(reqBody)

	// Create a test request
	req, _ := http.NewRequest("POST", "/html", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder
	w := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response models.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "No contents provided", response.Error)
}

func TestMakeHTML_InvalidJSON(t *testing.T) {
	router := setupTestRouter()

	// Create a test request with invalid JSON
	invalidJSON := []byte(`{"invalid": json}`)

	// Create a test request
	req, _ := http.NewRequest("POST", "/html", bytes.NewBuffer(invalidJSON))
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder
	w := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response models.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response.Error)
}

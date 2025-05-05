package azure

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vhybZApp/api.git/config"
	"github.com/vhybZApp/api.git/database"
	"github.com/vhybZApp/api.git/models"
	"github.com/vhybZApp/api.git/services"
)

type ChatCompletionRequest struct {
	Messages         []Message `json:"messages"`
	MaxTokens        int       `json:"max_tokens,omitempty"`
	Temperature      float64   `json:"temperature,omitempty"`
	TopP             float64   `json:"top_p,omitempty"`
	FrequencyPenalty float64   `json:"frequency_penalty,omitempty"`
	PresencePenalty  float64   `json:"presence_penalty,omitempty"`
	Stop             []string  `json:"stop,omitempty"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// @Summary Proxy Azure OpenAI Chat Completion
// @Description Proxy requests to Azure OpenAI Chat Completion API
// @Tags azure
// @Accept json
// @Produce json
// @Param request body ChatCompletionRequest true "Chat completion request"
// @Success 200 {object} ChatCompletionResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /azure/chat/completions [post]
func ChatCompletion(c *gin.Context) {
	// Get user ID from context (assuming it's set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.NewErrorResponse("User not authenticated"))
		return
	}

	// Initialize token quota service
	tokenQuotaService := services.NewTokenQuotaService(database.GetDB())

	// Validate Azure OpenAI configuration
	if config.AppConfig.AzureOpenAIEndpoint == "" || config.AppConfig.AzureOpenAIKey == "" || config.AppConfig.AzureOpenAIDeployment == "" {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Azure OpenAI configuration is incomplete"))
		return
	}

	// Parse request body
	var req ChatCompletionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(err.Error()))
		return
	}

	// Create HTTP client
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Marshal request body
	reqBody, err := json.Marshal(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Error marshaling request body"))
		return
	}

	// Create request
	url := config.AppConfig.AzureOpenAIEndpoint + "/openai/deployments/" + config.AppConfig.AzureOpenAIDeployment + "/chat/completions?api-version=2023-05-15"
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Error creating request"))
		return
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("api-key", config.AppConfig.AzureOpenAIKey)

	// Send request
	resp, err := client.Do(httpReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Error sending request to Azure OpenAI"))
		return
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Error reading response body"))
		return
	}

	// Check for error response
	if resp.StatusCode != http.StatusOK {
		c.JSON(resp.StatusCode, models.NewErrorResponse(string(body)))
		return
	}

	// Parse response
	var chatResp ChatCompletionResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Error parsing response"))
		return
	}

	// Update token usage
	if err := tokenQuotaService.UpdateUsage(userID.(string), chatResp.Usage.TotalTokens); err != nil {
		c.JSON(http.StatusTooManyRequests, models.NewErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, chatResp)
}

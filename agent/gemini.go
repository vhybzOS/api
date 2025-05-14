package agent

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vhybZApp/api/config"
	"github.com/vhybZApp/api/models"
	"google.golang.org/genai"
)

var client *genai.Client

func Init() {
	c, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey:  config.AppConfig.GeminiAPIKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		log.Fatalf("Failed to create Gemini client: %v", err)
	}
	client = c
}

// HTMLRequest represents the request body for HTML generation
// @Description Request body for generating HTML using Gemini AI
type HTMLRequest struct {
	// Contents contains the parts of the content to be processed
	Contents []*genai.Content `json:"parts" binding:"required"`
}

// HTMLResponse represents the response from HTML generation
// @Description Response containing the generated HTML
type HTMLResponse struct {
	// HTML contains the generated HTML code
	HTML string `json:"html"`
}

// MakeHTML godoc
// @Summary Generate HTML using Gemini AI
// @Description Generates HTML code based on the provided content using Gemini AI
// @Tags agent
// @Accept json
// @Produce json
// @Param request body HTMLRequest true "Request body containing content parts"
// @Success 200 {object} HTMLResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /agent/make-html [post]
func MakeHTML(c *gin.Context) {
	var req HTMLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(err.Error()))
		return
	}
	if len(req.Contents) == 0 {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("No contents provided"))
		return
	}
	if len(req.Contents) == 1 {
		req.Contents = append(req.Contents, genai.Text("You are a ultimate software engineer that generates HTML code.")[0])
	}

	resp, err := client.Models.GenerateContent(context.Background(), "gemini-2.0-flash", req.Contents, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, HTMLResponse{HTML: resp.Text()})
}

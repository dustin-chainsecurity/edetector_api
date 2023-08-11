package analysis

import (
	"edetector_API/internal/errhandler"
	"edetector_API/internal/template"
	"edetector_API/pkg/mariadb/query"
	"net/http"

	"github.com/gin-gonic/gin"
)

type templateResponse struct {
	IsSuccess  bool           `json:"isSuccess"`
	Data       interface{}    `json:"template"`
}

func AddTemplate(c *gin.Context) {
	
}

func DeleteTemplate(c *gin.Context) {
	
}

func UpdateTemplate(c *gin.Context) {
	
}

func GetTemplate(c *gin.Context) {
	// Get the template id from the request
	id := c.Param("id")
	raw, err := query.LoadRawTemplate(id)
	if err != nil {
		errhandler.Handler(c, err, "Error loading raw template data")
	}
	// Process the raw template data
	template, err := template.ProcessRawTemplate(raw)
	if err != nil {
		errhandler.Handler(c, err, "Error processing template data")
	}
	// Create the response object
	res := templateResponse{
		IsSuccess: true,
		Data:      template,
	}
	c.JSON(http.StatusOK, res)
}

func GetTemplateList(c *gin.Context) {
	raws, err := query.LoadAllRawTemplate()
	if err != nil {
		errhandler.Handler(c, err, "Error loading raw template data")
	}
	// Process the raw template data
	templates := []template.Template{}
	for _, raw := range raws {
		t, err := template.ProcessRawTemplate(raw)
		if err != nil {
			errhandler.Handler(c, err, "Error processing template data")
		}
		templates = append(templates, t)
	}
	// Create the response object
	res := templateResponse{
		IsSuccess: true,
		Data:      templates,
	}
	c.JSON(http.StatusOK, res)
}
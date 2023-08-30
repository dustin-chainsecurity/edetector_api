package analysis

import (
	"edetector_API/internal/template"
	"edetector_API/pkg/errhandler"
	"edetector_API/pkg/mariadb/query"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddTemplate(c *gin.Context) {
	// Get the template data from the request
	var req template.Template
	if err := c.ShouldBindJSON(&req); err != nil {
		errhandler.Handler(c, err, "Error binding template data")
		return
	}
	// Check if the name exists
	if exist, err := query.CheckTemplateName(req.Name); err != nil {
		errhandler.Handler(c, err, "Error checking template name")
		return
	} else if exist {
		errhandler.Handler(c, fmt.Errorf("template name "+req.Name+" already exists"), "Error checking template name")
		return
	}
	// Process the template data
	raw, err := template.ToRaw(req)
	if err != nil {
		errhandler.Handler(c, err, "Error processing template data")
	}
	// Add the template to the database
	id, err := query.AddTemplate(raw)
	if err != nil {
		errhandler.Handler(c, err, "Error adding template data")
		return
	}
	// Send the response object
	c.JSON(http.StatusOK, gin.H{
		"isSuccess":   true,
		"template_id": id,
	})
}

func DeleteTemplate(c *gin.Context) {
	// Get the template id from the request
	id := c.Param("id")
	// Check if the template id exists
	if exist, err := query.CheckTemplateID(id); err != nil {
		errhandler.Handler(c, err, "Error checking template id")
		return
	} else if !exist {
		errhandler.Handler(c, fmt.Errorf("template id"+id+" does not exist"), "Error checking template id")
		return
	}
	err := query.DeleteTemplate(id)
	if err != nil {
		errhandler.Handler(c, err, "Error deleting template data")
		return
	}
	// Send the response object
	c.JSON(http.StatusOK, gin.H{
		"isSuccess": true,
	})
}

func UpdateTemplate(c *gin.Context) {
	// Get the template id from the request
	id := c.Param("id")
	// Check if the template id exists
	if exist, err := query.CheckTemplateID(id); err != nil {
		errhandler.Handler(c, err, "Error checking template id")
		return
	} else if !exist {
		errhandler.Handler(c, fmt.Errorf("template id"+id+" does not exist"), "Error checking template id")
		return
	}
	// Get the template data from the request
	var req template.Template
	if err := c.ShouldBindJSON(&req); err != nil {
		errhandler.Handler(c, err, "Error binding template data")
		return
	}
	// Check if the name exists
	if exist, err := query.CheckTemplateName(req.Name); err != nil {
		errhandler.Handler(c, err, "Error checking template name")
		return
	} else if exist {
		errhandler.Handler(c, fmt.Errorf("template name "+req.Name+" already exists"), "Error checking template name")
		return
	}
	// Process the template data
	raw, err := template.ToRaw(req)
	if err != nil {
		errhandler.Handler(c, err, "Error processing template data")
		return
	}
	// Update the template in the database
	err = query.UpdateTemplate(id, raw)
	if err != nil {
		errhandler.Handler(c, err, "Error updating template data")
		return
	}
	// Send the response object
	c.JSON(http.StatusOK, gin.H{
		"isSuccess": true,
	})
}

func GetTemplate(c *gin.Context) {
	// Get the template id from the request
	id := c.Param("id")
	// Check if the template id exists
	if exist, err := query.CheckTemplateID(id); err != nil {
		errhandler.Handler(c, err, "Error checking template id")
		return
	} else if !exist {
		errhandler.Handler(c, fmt.Errorf("template id"+id+" does not exist"), "Error checking template id")
		return
	}
	// Load the raw template data
	raw, err := query.LoadRawTemplate(id)
	if err != nil {
		errhandler.Handler(c, err, "Error loading raw template data")
		return
	}
	// Process the raw template data
	template, err := template.ToTemplate(raw)
	if err != nil {
		errhandler.Handler(c, err, "Error processing template data")
		return
	}
	// Send the response object
	c.JSON(http.StatusOK, gin.H{
		"isSuccess": true,
		"template":  template,
	})
}

func GetTemplateList(c *gin.Context) {
	raws, err := query.LoadAllRawTemplate()
	if err != nil {
		errhandler.Handler(c, err, "Error loading raw template data")
		return
	}
	// Process the raw template data
	templates := []template.Template{}
	for _, raw := range raws {
		t, err := template.ToTemplate(raw)
		if err != nil {
			errhandler.Handler(c, err, "Error processing template data")
			return
		}
		templates = append(templates, t)
	}
	// Send the response object
	c.JSON(http.StatusOK, gin.H{
		"isSuccess": true,
		"templates": templates,
	})
}

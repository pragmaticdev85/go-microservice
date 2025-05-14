package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pragmaticdev85/go-microservice/internal/repositories"
	"github.com/pragmaticdev85/go-microservice/internal/services"
)

type ExampleController struct {
	service *services.ExampleService
}

func NewExampleController(service *services.ExampleService) *ExampleController {
	return &ExampleController{service: service}
}

// CreateExample godoc
// @Summary Create a new example
// @Description Create a new example with the input payload
// @Tags examples
// @Accept  json
// @Produce  json
// @Param example body repositories.Example true "Create example"
// @Success 201 {object} repositories.Example
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /examples [post]
func (c *ExampleController) CreateExample(ctx *gin.Context) {
	var example repositories.Example
	if err := ctx.ShouldBindJSON(&example); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdExample, err := c.service.CreateExample(&example)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, createdExample)
}

// GetExampleByID godoc
// @Summary Get example by ID
// @Description Get example by ID
// @Tags examples
// @Accept  json
// @Produce  json
// @Param id path string true "Example ID"
// @Success 200 {object} repositories.Example
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /examples/{id} [get]
func (c *ExampleController) GetExampleByID(ctx *gin.Context) {
	id := ctx.Param("id")

	example, err := c.service.GetExampleByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, example)
}

// Add other controller methods

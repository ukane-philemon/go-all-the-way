package main

import (
	"net/http"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

// recipes contains all the recipes created.
var recipes map[string]*Recipe

func init() {
	recipes = make(map[string]*Recipe, 0)
}

// Chef is information about a recipe owner.
type Chef struct {
	Name              string `json:"name" binding:"required"`
	Country           string `json:"country" binding:"required"`
	YearsOfExperience int    `json:"yearsOfExperience"`
}

// Recipe is information about a recipe.
type Recipe struct {
	Id           string    `json:"id"`
	Name         string    `json:"name" binding:"required"`
	Keywords     []string  `json:"keywords" binding:"required"`
	Ingredients  []string  `json:"ingredients" binding:"required"`
	Instructions []string  `json:"instructions" binding:"required"`
	Chef         *Chef     `json:"chef" binding:"required"`
	PublishedAt  time.Time `json:"publishedAt"`
}

// DeleteRecipeHandler removes an existing recipe.
func DeleteRecipeHandler(c *gin.Context) {
	id := c.Param("recipe-id")

	_, ok := recipes[id]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Recipe not found",
		})
		return
	}

	delete(recipes, id)
	c.JSON(http.StatusOK, gin.H{
		"message": "Recipe deleted",
	})
}

// UpdateRecipeHandler updates an existing recipe.
func UpdateRecipeHandler(c *gin.Context) {
	id := c.Param("recipe-id")

	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	_, ok := recipes[id]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Recipe not found",
		})
		return
	}

	recipe.Id = id
	recipes[id] = &recipe

	c.JSON(http.StatusOK, recipe)
}

// ListRecipesHandler lists all the available recipes in sorted order.
func ListRecipesHandler(c *gin.Context) {
	var recipes []*Recipe
	for _, recipe := range recipes {
		recipes = append(recipes, recipe)
	}

	sort.SliceStable(recipes, func(i, j int) bool {
		return recipes[i].PublishedAt.After(recipes[i].PublishedAt)
	})

	c.JSON(http.StatusOK, recipes)
}

// NewRecipeHandler creates a new recipe.
func NewRecipeHandler(c *gin.Context) {
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	recipe.Id = xid.New().String()
	recipe.PublishedAt = time.Now()
	recipes[recipe.Id] = &recipe
	c.JSON(http.StatusOK, recipe)
}

func main() {
	router := gin.Default()
	router.POST("/recipes", NewRecipeHandler)
	router.GET("/recipes", ListRecipesHandler)
	router.PUT("recipes/:recipe-id", UpdateRecipeHandler)
	router.DELETE("recipes/:recipe-id", DeleteRecipeHandler)
	router.Run()
}

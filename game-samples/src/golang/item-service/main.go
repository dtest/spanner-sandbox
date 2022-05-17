package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	spanner "cloud.google.com/go/spanner"
	c "github.com/dtest/spanner-game-item-service/config"
	"github.com/dtest/spanner-game-item-service/models"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var configuration c.Configurations

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	viper.SetDefault("server.host", "localhost")
	viper.SetDefault("server.port", 8080)

	err := viper.Unmarshal(&configuration)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}

}

// Mutator to create spanner context and client, and set them in gin
func setSpannerConnection() gin.HandlerFunc {
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, configuration.Spanner.URL())

	if err != nil {
		log.Fatal(err)
	}

	return func(c *gin.Context) {
		c.Set("spanner_client", *client)
		c.Set("spanner_context", ctx)
		c.Next()
	}
}

// Helper function to retrieve spanner client and context
func getSpannerConnection(c *gin.Context) (context.Context, spanner.Client) {
	return c.MustGet("spanner_context").(context.Context),
		c.MustGet("spanner_client").(spanner.Client)
}

func createItem(c *gin.Context) {
	var item models.GameItem

	if err := c.BindJSON(&item); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	ctx, client := getSpannerConnection(c)
	err := item.Create(ctx, client)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.IndentedJSON(http.StatusCreated, item.ItemUUID)
}

func getItem(c *gin.Context) {
	var itemUUID = c.Param("id")

	ctx, client := getSpannerConnection(c)

	item, err := models.GetItemByUUID(ctx, client, itemUUID)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "item not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, item)
}

func main() {
	router := gin.Default()
	// TODO: Better configuration of trusted proxy
	router.SetTrustedProxies(nil)

	router.Use(setSpannerConnection())

	router.POST("/items", createItem)
	router.GET("/items/:id", getItem)

	// TODO: Better configuration of host
	router.Run(configuration.Server.URL())
}

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	spanner "cloud.google.com/go/spanner"
	c "github.com/dtest/spanner-game-match-service/config"
	"github.com/dtest/spanner-game-match-service/models"
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

// Creating a game assigns a list of players not currently playing a game
func createGame(c *gin.Context) {
	var game models.Game

	ctx, client := getSpannerConnection(c)
	err := game.CreateGame(ctx, client)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.IndentedJSON(http.StatusCreated, game.GameUUID)
}

func closeGame(c *gin.Context) {
	var game models.Game

	if err := c.BindJSON(&game); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	ctx, client := getSpannerConnection(c)
	err := game.CloseGame(ctx, client)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.IndentedJSON(http.StatusOK, game.Winner)
}

func getOpenGame(c *gin.Context) {
	ctx, client := getSpannerConnection(c)
	game, err := models.GetOpenGame(ctx, client)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.IndentedJSON(http.StatusOK, game)
}

func main() {
	router := gin.Default()
	// TODO: Better configuration of trusted proxy
	router.SetTrustedProxies(nil)

	router.Use(setSpannerConnection())

	router.GET("/games/open", getOpenGame)
	router.POST("/games/create", createGame)
	router.PUT("/games/close", closeGame)

	router.Run(configuration.Server.URL())
}

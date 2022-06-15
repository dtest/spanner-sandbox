package main

import (
	"context"
	"log"
	"net/http"

	spanner "cloud.google.com/go/spanner"
	"github.com/dtest/spanner-game-profile-service/config"
	"github.com/dtest/spanner-game-profile-service/models"

	"github.com/gin-gonic/gin"
)

// Mutator to create spanner context and client, and set them in gin
func setSpannerConnection(c config.Config) gin.HandlerFunc {
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, c.Spanner.DB())

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

// TODO: used by authentication server to generate load. Should not be called by other entities,
//  so restrictions should be implemented
func getPlayerUUIDs(c *gin.Context) {
	ctx, client := getSpannerConnection(c)

	players, err := models.GetPlayerUUIDs(ctx, client)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "No players exist"})
		return
	}

	c.IndentedJSON(http.StatusOK, players)
}

func getPlayerByID(c *gin.Context) {
	var playerUUID = c.Param("id")

	ctx, client := getSpannerConnection(c)

	player, err := models.GetPlayerByUUID(ctx, client, playerUUID)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "player not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, player)
}

func getPlayerStats(c *gin.Context) {
	var player models.Player
	player.PlayerUUID = c.Param("id")

	ctx, client := getSpannerConnection(c)

	err := player.GetPlayerStats(ctx, client)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "player not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, player)
}

func createPlayer(c *gin.Context) {
	var player models.Player

	if err := c.BindJSON(&player); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	ctx, client := getSpannerConnection(c)
	err := player.AddPlayer(ctx, client)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.IndentedJSON(http.StatusCreated, player.PlayerUUID)
}

func main() {
	configuration, _ := config.NewConfig()

	router := gin.Default()
	// TODO: Better configuration of trusted proxy
	router.SetTrustedProxies(nil)

	router.Use(setSpannerConnection(configuration))

	router.POST("/players", createPlayer)
	router.GET("/players", getPlayerUUIDs)
	router.GET("/players/:id", getPlayerByID)
	// router.GET("/player/login", getPlayerByLogin)
	router.GET("/players/:id/stats", getPlayerStats)

	router.Run(configuration.Server.URL())
}

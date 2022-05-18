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

// Update a player balance with a provided amount. Result is a JSON object that contains PlayerUUID and AccountBalance
func updatePlayerBalance(c *gin.Context) {
	var player models.Player
	var ledger models.PlayerLedger

	if err := c.BindJSON(&ledger); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	ctx, client := getSpannerConnection(c)
	err := ledger.UpdateBalance(ctx, client, &player)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	type PlayerBalance struct {
		PlayerUUID, AccountBalance string
	}

	balance := PlayerBalance{PlayerUUID: player.PlayerUUID, AccountBalance: player.Account_balance.FloatString(2)}
	c.IndentedJSON(http.StatusOK, balance)
}

func getPlayer(c *gin.Context) {
	var player models.Player

	ctx, client := getSpannerConnection(c)
	err := player.GetPlayer(ctx, client)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.IndentedJSON(http.StatusOK, player)
}

// func addPlayerItem(c *gin.Context) {

// }

// func getPlayerItem(c *gin.Context) {

// }

func main() {
	router := gin.Default()
	// TODO: Better configuration of trusted proxy
	router.SetTrustedProxies(nil)

	router.Use(setSpannerConnection())

	router.POST("/items", createItem)
	router.GET("/items/:id", getItem)
	router.PUT("/players/balance", updatePlayerBalance) // TODO: leverage profile service instead
	router.GET("/players", getPlayer)
	// router.POST("/players/items", addPlayerItem)
	// router.GET("/players/items", getPlayerItem)

	router.Run(configuration.Server.URL())
}

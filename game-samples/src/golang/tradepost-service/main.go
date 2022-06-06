package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	spanner "cloud.google.com/go/spanner"
	c "github.com/dtest/spanner-game-tradepost-service/config"
	"github.com/dtest/spanner-game-tradepost-service/models"
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

func getPlayerItem(c *gin.Context) {
	ctx, client := getSpannerConnection(c)

	item, err := models.GetRandomPlayerItem(ctx, client)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "item not found"})
		return
	}

	type RandomItem struct {
		PlayerUUID, PlayerItemUUID, Price string
	}

	ri := RandomItem{PlayerUUID: item.PlayerUUID, PlayerItemUUID: item.PlayerItemUUID, Price: item.Price.FloatString(2)}

	c.IndentedJSON(http.StatusOK, ri)
}

// Get a random open order with a random buyer. Used in trade simulation
func getOpenOrder(c *gin.Context) {
	ctx, client := getSpannerConnection(c)

	// Get an order
	order, err := models.GetRandomOpenOrder(ctx, client)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "order not found"})
		return
	}

	// Get a buyer; can't be the same player as the trade order's lister
	buyer, err := models.GetRandomPlayer(ctx, client, order.Lister, order.ListPrice)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "order not found"})
		return
	}

	type RandomOrder struct {
		OrderUUID, BuyerUUID, ListPrice, AccountBalance string
	}

	ro := RandomOrder{OrderUUID: order.OrderUUID, BuyerUUID: buyer.PlayerUUID, ListPrice: order.ListPrice.FloatString(2), AccountBalance: buyer.AccountBalance.FloatString(2)}

	c.IndentedJSON(http.StatusOK, ro)
}

// Create a sell order
// TODO: Enable buy orders
func createOrder(c *gin.Context) {
	var order models.TradeOrder

	if err := c.BindJSON(&order); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	ctx, client := getSpannerConnection(c)
	err := order.Create(ctx, client)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.IndentedJSON(http.StatusCreated, order.OrderUUID)
}

func purchaseOrder(c *gin.Context) {
	var order models.TradeOrder

	if err := c.BindJSON(&order); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	ctx, client := getSpannerConnection(c)
	err := order.Buy(ctx, client)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.IndentedJSON(http.StatusCreated, order.OrderUUID)
}

func main() {
	router := gin.Default()
	// TODO: Better configuration of trusted proxy
	router.SetTrustedProxies(nil)

	router.Use(setSpannerConnection())

	router.GET("/player_items", getPlayerItem)
	router.POST("/sell", createOrder)
	router.PUT("/buy", purchaseOrder)

	router.GET("/open", getOpenOrder)

	router.Run(configuration.Server.URL())
}

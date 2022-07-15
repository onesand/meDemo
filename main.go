package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"math"
	"math/big"
	"meDemo/client"
	"meDemo/model"
	"net/http"
	"os"
	"time"
)

var db = make(map[string]string)

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "fuck")
	})

	//Balance of
	r.GET("/balance/:addr", func(c *gin.Context) {
		address := c.Param("addr")
		account := common.HexToAddress(address)
		userAddress := model.UserAddress{Address: address, CreatedAt: time.Now()}
		client.DB().Create(&userAddress)
		//println("Add to db:" + result.Error.Error())
		balance, err := client.EthClient().BalanceAt(context.Background(), account, nil)
		if err != nil {
			c.String(http.StatusOK, err.Error())
			return
		}
		fbalance := new(big.Float)
		fbalance.SetString(balance.String())
		ethValue := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))
		fmt.Println(ethValue)
		c.String(http.StatusOK, ethValue.String())
	})

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// Get user value
	r.GET("/user/:name", func(c *gin.Context) {
		user := c.Params.ByName("name")
		value, ok := db[user]
		if ok {
			c.JSON(http.StatusOK, gin.H{"user": user, "value": value})
		} else {
			c.JSON(http.StatusOK, gin.H{"user": user, "status": "no value"})
		}
	})

	// Authorized group (uses gin.BasicAuth() middleware)
	// Same than:
	// authorized := r.Group("/")
	// authorized.Use(gin.BasicAuth(gin.Credentials{
	//	  "foo":  "bar",
	//	  "manu": "123",
	//}))
	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
		"foo":  "bar", // user:foo password:bar
		"manu": "123", // user:manu password:123
	}))

	/* example curl for /admin with basicauth header
	   Zm9vOmJhcg== is base64("foo:bar")

		curl -X POST \
	  	http://localhost:8080/admin \
	  	-H 'authorization: Basic Zm9vOmJhcg==' \
	  	-H 'content-type: application/json' \
	  	-d '{"value":"bar"}'
	*/
	authorized.POST("admin", func(c *gin.Context) {
		user := c.MustGet(gin.AuthUserKey).(string)

		// Parse JSON
		var json struct {
			Value string `json:"value" binding:"required"`
		}

		if c.Bind(&json) == nil {
			db[user] = json.Value
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		}
	})

	return r
}

func main() {
	client.SetupConnections()

	r := setupRouter()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	// Listen and Server in 0.0.0.0:8080
	r.Run(":" + port)
}

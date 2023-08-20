package main

import (
	"github.com/gin-gonic/gin"
	"github.com/zhou78yang/service-router/config"
	"github.com/zhou78yang/service-router/strategy"
)

func main() {
	config.Init()
	strategy.Init()

	r := gin.Default()
	r.NoRoute(strategy.Handle)

	if err := r.Run("0.0.0.0:8080"); err != nil {
		panic(err)
	}
}

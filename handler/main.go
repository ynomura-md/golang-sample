package main

import (
    "time"
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/zap"
    "deliver-endpoint/api"
    "go.uber.org/zap"
)

var logger *zap.Logger


func main() {
    logger, _ = zap.NewDevelopment()
    r := gin.Default()
    r.Use(ginzap.Ginzap(logger, time.RFC3339, true))
    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "pong",
        })
    })
    r.GET("/", handler)
    r.Run(":9999")
}


func handler(ctx *gin.Context) {

    user := api.User{"User", 20}

     logger.Info("object sample", zap.Object("userObj", user))


    ctx.JSON(200, gin.H{
        "user": user,
    })
}

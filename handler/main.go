package main

import (
    "time"
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/zap"
    "deliver-endpoint/api"
    "go.uber.org/zap"
    "github.com/satori/go.uuid"
    "go.uber.org/zap/zapcore"
)

//var logger *zap.Logger
var sugarLogger *zap.SugaredLogger

type AccessLog struct {
  Id string
  Time time.Time
  IP string
  UA string
  UID string
  UID_TYPE string
  PLATFORM string
  AD_UNIT_GROUP_ID int
  AD_UNIT_ID int
  EXT map[string]string

}
func (a AccessLog) MarshalLogObject(enc zapcore.ObjectEncoder) error {
    enc.AddString("id", a.Id)
    enc.AddString("time", a.Time.UTC().Format("2006-01-02T15:04:05.999Z"))
    enc.AddString("ip", a.IP)
    enc.AddString("ua", a.UA)
    enc.AddString("uid", a.UID)
    enc.AddString("ut", a.UID_TYPE)
    enc.AddString("pl", a.PLATFORM)
    enc.AddInt("augi", a.AD_UNIT_GROUP_ID)
    enc.AddInt("aui", a.AD_UNIT_ID)
    enc.OpenNamespace("ext")
    for key, value := range a.EXT {
      enc.AddString(key, value)
    }
    return nil
}

func main() {
    logger, _ := zap.NewProduction()
    sugarLogger = logger.Sugar()
    defer sugarLogger.Sync()
    r := gin.Default()
    r.Use(ginzap.Ginzap(logger, time.RFC3339, true))
    r.Use(sampleMiddleware())
    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "pong",
        })
    })
    r.GET("/", handler)
    r.Run(":9999")
}

func sampleMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        sugarLogger.Infow("before logic")
        if uid := c.Query("uid"); len(uid) > 0 {
          c.Set("uid", uid)
          c.Set("uidType", "wowma") // クエリに混ぜる？
        } else {
          id := uuid.NewV1()
          a,_ := id.MarshalText()
          c.Set("uid", string(a))
          c.Set("uidType", "tpo")
        }

        if pl := c.Query("pl"); len(pl) > 0 {
          c.Set("pl", pl)
        } else {
          c.Set("pl", "unknown")
        }

        c.Next()
        sugarLogger.Infow("after logic")
    }
}

type AdReq struct {
  AD_UNIT_GROUP_ID int `form:"gi"`
  AD_UNIT_ID int `form:"ui"`
}


func handler(ctx *gin.Context) {

    var req AdReq

    ctx.Bind(&req)

    user := api.User{"User", 20}

    id := uuid.NewV1()
    a,_ := id.MarshalText()

    logdata := AccessLog{
      string(a),
      time.Now(),
      ctx.ClientIP(),
      ctx.GetHeader("user-agent"),
      ctx.GetString("uid"),
      ctx.GetString("uidType"),
      ctx.GetString("pl"),
      req.AD_UNIT_GROUP_ID,
      req.AD_UNIT_ID,
      map[string]string{"apple": "a", "banana": "1", "lemon": "1"},
    }

    sugarLogger.Infow("record", zap.Object("access",logdata))


    ctx.JSON(200, gin.H{
        "user": user,
    })
}

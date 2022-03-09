package main

import (
	"embed"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	"go-message-pusher/model"
	"go-message-pusher/router"
	"go-message-pusher/util"
	"html/template"
	"log"
	"os"
	"strconv"
	"time"
)

//go:embed view
var fs embed.FS

func main() {
	if os.Getenv("GIN_MODE") != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}

	var err error

	err = util.ReadConfig()
	if err != nil {
		panic(err)
	}
	log.Println("Config loaded.")

	err = model.InitDatabase()
	if err != nil {
		panic(err)
	}
	log.Println("Database initialized.")

	err = model.InitCache(util.EnableRedis)
	if err != nil {
		panic(err)
	}
	log.Println("Cache initialized.")

	server := gin.Default()
	templ := template.Must(template.New("").ParseFS(fs, "view/**/*"))
	server.SetHTMLTemplate(templ)
	router.SetHtmlRouters(server)
	router.SetMessageRouter(server)
	router.SetUserRouter(server)
	log.Println("Router set.")

	s := gocron.NewScheduler(time.UTC)
	_, err = s.Every(100).Minutes().Do(func() {
		err := model.UpdateAllAccessTokens()
		if err != nil {
			log.Println("Failed to update accessTokens.")
			return
		}
		log.Println("AccessTokens updated.")
	})
	if err != nil {
		panic(err)
	}
	log.Println("Scheduler set.")
	s.StartAsync()

	log.Println("Server started.")
	err = server.Run(":" + strconv.Itoa(util.ServerPort))
	if err != nil {
		panic(err)
	}
}

package router

import (
	"github.com/gin-gonic/gin"
	"go-message-pusher/controller"
	"go-message-pusher/middleware"
	"net/http"
)

func SetHtmlRouters(router *gin.Engine) {
	router.GET("/", controller.Index)
	router.GET("/login", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "pages/login.gohtml", gin.H{})
	})
	router.GET("/signup", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "pages/signup.gohtml", gin.H{})
	})
	router.GET("/config", middleware.PageJWTAuth(), controller.Config)
	router.GET("/message", middleware.PageJWTAuth(), controller.Message)
}

func SetMessageRouter(router *gin.Engine) {
	router.GET("/template/:name", middleware.TokenAuth(), controller.PushTemplateMessage)
	router.GET("/plaintext/:name", middleware.TokenAuth(), controller.PushPlainTextMessage)
	router.GET("/textcard/:name", middleware.TokenAuth(), controller.PushTextCardMessage)

	router.POST("/template/:name", middleware.TokenAuth(), controller.PushTemplateMessage)
	router.POST("/plaintext/:name", middleware.TokenAuth(), controller.PushPlainTextMessage)
	router.POST("/textcard/:name", middleware.TokenAuth(), controller.PushTextCardMessage)

	router.GET("/verify/:name", controller.CheckSignature)
}

func SetUserRouter(router *gin.Engine) {
	router.POST("/auth", controller.Auth)
	router.POST("/register", controller.RegisterUser)

	apiGroup := router.Group("/api")
	apiGroup.Use(middleware.ApiJWTAuth())
	{
		apiGroup.POST("/user", controller.UpdateUser)
		apiGroup.POST("/app", controller.UpdateUserApp)
		apiGroup.POST("/corp", controller.UpdateUserCorp)
	}
}

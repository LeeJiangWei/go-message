package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-message-pusher/model"
	"go-message-pusher/util"
	"log"
	"net/http"
)

func Index(ctx *gin.Context) {
	memstats := util.GetMemoryStatus()
	ctx.HTML(http.StatusOK, "pages/index.gohtml", gin.H{
		"Sys":        fmt.Sprintf("%.2f", memstats.Sys),
		"Alloc":      fmt.Sprintf("%.2f", memstats.Alloc),
		"TotalAlloc": fmt.Sprintf("%.2f", memstats.TotalAlloc),
	})
}

func Config(ctx *gin.Context) {
	user := ctx.MustGet("user").(model.User)
	ctx.HTML(http.StatusOK, "pages/config.gohtml", gin.H{
		"User": user,
	})
}

func Message(ctx *gin.Context) {
	u := ctx.MustGet("user").(model.User)
	user, err := model.RetrieveMessagesByUserID(u.ID)
	if err != nil {
		log.Println(err.Error())
		ctx.String(http.StatusInternalServerError, "Unable to query database.")
		return
	}
	ctx.HTML(http.StatusOK, "pages/message.gohtml", gin.H{
		"TemplateMessages":  user.TemplateMessages,
		"PlainTextMessages": user.PlainTextMessages,
		"TextCardMessages":  user.TextCardMessages,
	})
}

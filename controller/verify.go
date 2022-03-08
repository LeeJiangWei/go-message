package controller

import (
	"github.com/gin-gonic/gin"
	"go-message-pusher/model"
	"go-message-pusher/util"
	"net/http"
)

func CheckSignature(ctx *gin.Context) {
	name := ctx.Param("name")

	user, err := model.RetrieveUserCacheByName(name)
	if err != nil {
		ctx.String(http.StatusBadRequest, "User %v does not exist.", name)
		return
	}

	signature := ctx.Query("signature")
	timestamp := ctx.Query("timestamp")
	nonce := ctx.Query("nonce")
	echostr := ctx.Query("echostr")

	token := user.App.VerifyToken

	if util.Verify(signature, timestamp, nonce, token) {
		ctx.String(http.StatusOK, echostr)
	} else {
		ctx.String(http.StatusOK, "")
	}
}

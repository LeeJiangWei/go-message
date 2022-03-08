package middleware

import (
	"github.com/gin-gonic/gin"
	"go-message-pusher/model"
	"log"
	"net/http"
)

// TokenAuth 验证用户自定义 token，用于推送消息
func TokenAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var token string

		if s, exists := ctx.GetQuery("token"); exists {
			token = s
		} else if s, exists = ctx.GetPostForm("token"); exists {
			token = s
		} else {
			ctx.String(http.StatusUnauthorized, "Empty user token.")
			ctx.Abort()
			return
		}

		name := ctx.Param("name")
		user, err := model.RetrieveUserCacheByName(name)
		if err != nil {
			log.Println(err.Error())
			ctx.String(http.StatusBadRequest, `User "%v" does not exist.`, name)
			ctx.Abort()
			return
		}

		if token != user.Token {
			ctx.String(http.StatusUnauthorized, "Incorrect user token.")
			ctx.Abort()
			return
		}

		ctx.Set("user", user)
		ctx.Next()
	}
}

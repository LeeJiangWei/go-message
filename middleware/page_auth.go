package middleware

import (
	"github.com/gin-gonic/gin"
	"go-message-pusher/model"
	"go-message-pusher/util"
	"net/http"
)

// PageJWTAuth 验证后端渲染页面的 JWT token
func PageJWTAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var token string

		if s, exists := ctx.GetQuery("token"); exists {
			token = s
		} else {
			token = ctx.GetHeader("token")
		}

		if token == "" {
			ctx.Redirect(http.StatusFound, "/")
			ctx.Abort()
			return
		}

		claims, err := util.ParseToken(token)
		if err != nil {
			ctx.Redirect(http.StatusFound, "/")
			ctx.Abort()
			return
		}

		user, err := model.RetrieveUserByID(claims.UserID)
		if err != nil {
			ctx.Redirect(http.StatusFound, "/")
			ctx.Abort()
			return
		}

		ctx.Set("user", user)
		ctx.Next()
	}
}

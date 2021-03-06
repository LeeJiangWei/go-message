package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"go-message-pusher/model"
	"go-message-pusher/util"
	"net/http"
)

// ApiJWTAuth 验证 JWT token，用于登录后台
func ApiJWTAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var token string

		if s, exists := ctx.GetQuery("token"); exists {
			token = s
		} else {
			token = ctx.GetHeader("token")
		}

		if token == "" {
			ctx.String(http.StatusUnauthorized, "未携带 JWT token。")
			ctx.Abort()
			return
		}

		claims, err := util.ParseToken(token)
		if err != nil {
			switch err.(*jwt.ValidationError).Errors {
			case jwt.ValidationErrorExpired:
				ctx.String(http.StatusUnauthorized, "JWT token 已过期。")
			default:
				ctx.String(http.StatusUnauthorized, "JWT token 不正确。")
			}
			ctx.Abort()
			return
		}

		user, err := model.RetrieveUserByID(claims.UserID)
		if err != nil {
			ctx.String(http.StatusUnauthorized, "JWT token 不正确。")
			ctx.Abort()
			return
		}

		ctx.Set("user", user)
		ctx.Next()
	}
}

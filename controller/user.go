package controller

import (
	"github.com/gin-gonic/gin"
	"go-message-pusher/model"
	"go-message-pusher/util"
	"log"
	"net/http"
)

func Auth(ctx *gin.Context) {
	name := ctx.PostForm("name")
	password := ctx.PostForm("password")

	user, err := model.RetrieveUserByName(name)
	if err != nil {
		ctx.String(http.StatusBadRequest, `用户 "%v" 不存在。`, name)
		return
	}

	if password != user.Password {
		ctx.String(http.StatusBadRequest, "密码错误。")
		return
	}

	JWT, err := util.GenerateJWT(user.ID)
	if err != nil {
		log.Println(err.Error())
		ctx.String(http.StatusInternalServerError, "Can not generate JWT.")
		return
	}

	ctx.String(http.StatusOK, JWT)
}

func RegisterUser(ctx *gin.Context) {
	name := ctx.PostForm("name")
	password := ctx.PostForm("password")
	token := ctx.PostForm("token")

	user, err := model.CreateUser(name, password, token)
	if err != nil {
		log.Println(err.Error())
		ctx.String(http.StatusBadRequest, `用户名 "%v" 已经存在。`, name)
		return
	}

	_ = model.UpdateUserCache(user)

	JWT, err := util.GenerateJWT(user.ID)
	if err != nil {
		log.Println(err.Error())
		ctx.String(http.StatusInternalServerError, "Can not generate JWT.")
		return
	}

	ctx.String(http.StatusOK, JWT)
}

func UpdateUser(ctx *gin.Context) {
	user := ctx.MustGet("user").(model.User)

	// 删除缓存
	_ = model.DeleteUserCacheByName(user.Name)

	userID := user.ID
	newName := ctx.PostForm("name")
	password := ctx.PostForm("password")
	token := ctx.PostForm("token")

	// 更新数据库
	_, err := model.UpdateUser(userID, newName, password, token)
	if err != nil {
		log.Println(err.Error())
		ctx.String(http.StatusBadRequest, `用户名 "%v" 已经存在。`, newName)
		return
	}
	ctx.String(http.StatusOK, "ok")
}

func UpdateUserApp(ctx *gin.Context) {
	user := ctx.MustGet("user").(model.User)

	// 删除缓存
	_ = model.DeleteUserCacheByName(user.Name)

	userID := user.ID
	appID := ctx.PostForm("appID")
	appSecret := ctx.PostForm("appSecret")
	templateID := ctx.PostForm("templateID")
	receiverID := ctx.PostForm("receiverID")
	verifyToken := ctx.PostForm("verifyToken")

	// 更新数据库
	_, err := model.CreateOrUpdateUserApp(userID, appID, appSecret, templateID, receiverID, verifyToken)
	if err != nil {
		log.Println(err.Error())
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	ctx.String(http.StatusOK, "ok")
}

func UpdateUserCorp(ctx *gin.Context) {
	user := ctx.MustGet("user").(model.User)

	// 删除缓存
	_ = model.DeleteUserCacheByName(user.Name)

	userID := user.ID
	corpID := ctx.PostForm("corpID")
	agentID := ctx.PostForm("agentID")
	agentSecret := ctx.PostForm("agentSecret")
	receiverID := ctx.PostForm("receiverID")
	cardUrl := ctx.PostForm("cardUrl")

	// 更新数据库
	_, err := model.CreateOrUpdateUserCorp(userID, corpID, agentID, agentSecret, receiverID, cardUrl)
	if err != nil {
		log.Println(err.Error())
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	ctx.String(http.StatusOK, "ok")
}

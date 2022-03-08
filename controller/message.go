package controller

import (
	"github.com/gin-gonic/gin"
	"go-message-pusher/model"
	"go-message-pusher/wechat"
	"log"
	"net/http"
	"net/url"
)

func PushTemplateMessage(ctx *gin.Context) {
	user := ctx.MustGet("user").(model.User)

	var from, desc, remark string
	if ctx.Request.Method == http.MethodGet {
		from = ctx.Query("from")
		desc = ctx.Query("desc")
		remark = ctx.Query("remark")
	} else {
		from = ctx.PostForm("from")
		desc = ctx.PostForm("desc")
		remark = ctx.PostForm("remark")
	}

	template, err := wechat.BuildTemplateMessage(user.App.ReceiverID, user.App.TemplateID, from, desc, remark)
	if err != nil {
		log.Println(err.Error())
		ctx.String(http.StatusInternalServerError, "Error occurred when building message.")
		return
	}

	appAccessToken, _ := model.RetrieveAppAccessTokenCache(user.Name)
	res, err := wechat.PushTemplateMessage(appAccessToken, template)

	if err != nil {
		// 未能推送至微信服务器
		log.Println("Error occurred when pushing message: ", err.Error())
		ctx.String(http.StatusInternalServerError, "Error occurred when pushing message.")

		_, dbErr := model.CreateTemplateMessage(user.ID, from, desc, remark, err.Error())
		if dbErr != nil {
			log.Println("DB Error occurred when creating message: ", dbErr.Error())
		}

	} else if res.ErrCode != 0 {
		// 微信服务器返回错误
		log.Println("Error occurred when pushing message: ", res.ErrCode, res.ErrMsg)
		ctx.String(http.StatusInternalServerError, res.ErrMsg)

		_, dbErr := model.CreateTemplateMessage(user.ID, from, desc, remark, res.ErrMsg)
		if dbErr != nil {
			log.Println("DB Error occurred when creating message: ", dbErr.Error())
		}

	} else {
		// 成功
		_, dbErr := model.CreateTemplateMessage(user.ID, from, desc, remark, "success")
		if dbErr != nil {
			log.Println("DB Error occurred when creating message: ", dbErr.Error())
			ctx.String(http.StatusInternalServerError, "Failed to insert message to DB.")
		} else {
			ctx.String(http.StatusOK, "ok")
		}
	}
}

func PushPlainTextMessage(ctx *gin.Context) {
	user := ctx.MustGet("user").(model.User)

	var content string
	if ctx.Request.Method == http.MethodGet {
		content = ctx.Query("content")
	} else {
		content = ctx.PostForm("content")
	}

	plainText, err := wechat.BuildPlainTextMessage(user.Corp.ReceiverID, user.Corp.AgentID, content)
	if err != nil {
		log.Println(err.Error())
		ctx.String(http.StatusInternalServerError, "Error occurred when building message.")
		return
	}

	corpAccessToken, _ := model.RetrieveCorpAccessTokenCache(user.Name)
	res, err := wechat.PushCorpMessage(corpAccessToken, plainText)
	if err != nil {
		// 未能推送至微信服务器
		log.Println(err.Error())
		ctx.String(http.StatusInternalServerError, "Error occurred when pushing message.")

		_, dbErr := model.CreatePlainTextMessage(user.ID, content, err.Error())
		if dbErr != nil {
			log.Println("DB Error occurred when creating message: ", dbErr.Error())
		}

	} else if res.ErrCode != 0 {
		// 微信服务器返回错误
		log.Println(res.ErrCode, res.ErrMsg)
		ctx.String(http.StatusInternalServerError, res.ErrMsg)

		_, dbErr := model.CreatePlainTextMessage(user.ID, content, res.ErrMsg)
		if dbErr != nil {
			log.Println("DB Error occurred when creating message: ", dbErr.Error())
		}

	} else {
		// 成功
		_, dbErr := model.CreatePlainTextMessage(user.ID, content, "success")
		if dbErr != nil {
			log.Println("DB Error occurred when creating message: ", dbErr.Error())
			ctx.String(http.StatusInternalServerError, "Failed to insert message to DB.")
		} else {
			ctx.String(http.StatusOK, "ok")
		}
	}
}

func PushTextCardMessage(ctx *gin.Context) {
	user := ctx.MustGet("user").(model.User)

	var title, desc, u string
	if ctx.Request.Method == http.MethodGet {
		title = ctx.Query("title")
		desc = ctx.Query("desc")
		u = ctx.Query("url")
	} else {
		title = ctx.PostForm("title")
		desc = ctx.PostForm("desc")
		u = ctx.PostForm("url")
	}

	if _, err := url.Parse(u); err != nil || u == "" {
		u = user.Corp.CardUrl
	}

	textCard, err := wechat.BuildTextCardMessage(user.Corp.ReceiverID, user.Corp.AgentID, title, desc, u)
	if err != nil {
		log.Println(err.Error())
		ctx.String(http.StatusInternalServerError, "Error occurred when building message.")
		return
	}

	corpAccessToken, _ := model.RetrieveCorpAccessTokenCache(user.Name)
	res, err := wechat.PushCorpMessage(corpAccessToken, textCard)
	if err != nil {
		// 未能推送至微信服务器
		log.Println(err.Error())
		ctx.String(http.StatusInternalServerError, "Error occurred when pushing message.")

		_, dbErr := model.CreateTextCardMessage(user.ID, title, desc, u, err.Error())
		if dbErr != nil {
			log.Println("DB Error occurred when creating message: ", dbErr.Error())
		}

	} else if res.ErrCode != 0 {
		// 微信服务器返回错误
		log.Println(res.ErrCode, res.ErrMsg)
		ctx.String(http.StatusInternalServerError, res.ErrMsg)

		_, dbErr := model.CreateTextCardMessage(user.ID, title, desc, u, res.ErrMsg)
		if dbErr != nil {
			log.Println("DB Error occurred when creating message: ", dbErr.Error())
		}

	} else {
		// 成功
		_, dbErr := model.CreateTextCardMessage(user.ID, title, desc, u, "success")
		if dbErr != nil {
			log.Println("DB Error occurred when creating message: ", dbErr.Error())
			ctx.String(http.StatusInternalServerError, "Failed to insert message to DB.")
		} else {
			ctx.String(http.StatusOK, "ok")
		}
	}
}

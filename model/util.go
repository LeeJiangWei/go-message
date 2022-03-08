package model

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var db *gorm.DB

func InitDatabase() (err error) {
	db, err = gorm.Open(sqlite.Open("data.db"), &gorm.Config{})
	if err != nil {
		return err
	}

	// enforce foreign key constraints
	if err = db.Exec("PRAGMA foreign_keys = ON", nil).Error; err != nil {
		return
	}

	err = db.AutoMigrate(&User{}, &App{}, &Corp{}, &TemplateMessage{}, &PlainTextMessage{}, &TextCardMessage{})
	if err != nil {
		return
	}

	return
}

/* CRUD of user */

func CreateUser(name, password, token string) (user User, err error) {
	user = User{
		Name:     name,
		Password: password,
		Token:    token,
	}
	err = db.Create(&user).Error
	return
}

func RetrieveUserByID(userID uint) (user User, err error) {
	err = db.Preload("App").Preload("Corp").First(&user, userID).Error
	return
}

func RetrieveUserByName(name string) (user User, err error) {
	err = db.Where("name = ?", name).Preload("App").Preload("Corp").First(&user).Error
	return
}

func RetrieveAllUsers() (users []User, err error) {
	err = db.Joins("App").Joins("Corp").Find(&users).Error
	return
}

func UpdateUser(userID uint, name, password, token string) (user User, err error) {
	user = User{
		Name:     name,
		Password: password,
		Token:    token,
	}
	err = db.Table("users").Where("id = ?", userID).Updates(&user).Error
	return
}

/* CRUD of app */

func CreateOrUpdateUserApp(userID uint, appID, appSecret, templateID, receiverID, verifyToken string) (app App, err error) {
	app = App{
		AppID:       appID,
		AppSecret:   appSecret,
		TemplateID:  templateID,
		ReceiverID:  receiverID,
		VerifyToken: verifyToken,
		UserID:      userID,
	}

	// 如果主键不冲突则插入，否则更新除主键外所有字段
	err = db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&app).Error

	return
}

/* CRUD of corp */

func CreateOrUpdateUserCorp(userID uint, corpID, agentID, agentSecret, receiverID, cardUrl string) (corp Corp, err error) {
	corp = Corp{
		UserID:      userID,
		CorpID:      corpID,
		AgentID:     agentID,
		AgentSecret: agentSecret,
		ReceiverID:  receiverID,
		CardUrl:     cardUrl,
	}

	err = db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&corp).Error

	return
}

/* CRUD of message */

func CreateTemplateMessage(userID uint, from, description, remark, status string) (message TemplateMessage, err error) {
	message = TemplateMessage{
		From:        from,
		Description: description,
		Remark:      remark,
		UserID:      userID,
		Status:      status,
	}

	err = db.Create(&message).Error
	return
}

func CreatePlainTextMessage(userID uint, content, status string) (message PlainTextMessage, err error) {
	message = PlainTextMessage{
		Model:   gorm.Model{},
		Content: content,
		UserID:  userID,
		Status:  status,
	}

	err = db.Create(&message).Error
	return
}

func CreateTextCardMessage(userID uint, title, description, url, status string) (message TextCardMessage, err error) {
	message = TextCardMessage{
		Title:       title,
		Description: description,
		Url:         url,
		UserID:      userID,
		Status:      status,
	}

	err = db.Create(&message).Error
	return
}

/* Retrieve messages */

func RetrieveMessagesByUserID(userID uint) (user User, err error) {
	err = db.Preload("TemplateMessages").Preload("PlainTextMessages").Preload("TextCardMessages").First(&user, userID).Error
	return
}

func RetrieveTemplateMessagesByUserID(userID uint) (user User, err error) {
	err = db.Preload("TemplateMessages").Limit(20).First(&user, userID).Error
	return
}

func RetrieveTextCardMessagesByUserID(userID uint) (user User, err error) {
	err = db.Preload("TextCardMessages").Limit(20).First(&user, userID).Error
	return
}

func RetrievePlainTextMessagesByUserID(userID uint) (user User, err error) {
	err = db.Preload("PlainTextMessages").Limit(20).First(&user, userID).Error
	return
}

package model

import "gorm.io/gorm"

type User struct {
	gorm.Model

	Name     string `gorm:"unique;not null"`
	Password string
	Token    string

	App  App
	Corp Corp

	TemplateMessages  []TemplateMessage
	TextCardMessages  []TextCardMessage
	PlainTextMessages []PlainTextMessage
}

type App struct {
	UserID uint `gorm:"primaryKey"`

	AppID       string
	AppSecret   string
	TemplateID  string
	ReceiverID  string
	VerifyToken string
}

type Corp struct {
	UserID uint `gorm:"primaryKey"`

	CorpID      string
	AgentID     string
	AgentSecret string
	ReceiverID  string
	CardUrl     string
}

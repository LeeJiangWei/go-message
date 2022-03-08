package model

import "gorm.io/gorm"

type TemplateMessage struct {
	gorm.Model

	From        string
	Description string
	Remark      string
	Status      string

	UserID uint
}

type TextCardMessage struct {
	gorm.Model

	Title       string
	Description string
	Url         string
	Status      string

	UserID uint
}

type PlainTextMessage struct {
	gorm.Model

	Content string
	Status  string

	UserID uint
}

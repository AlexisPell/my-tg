package model

import "gorm.io/gorm"

// User представляет структуру пользователя в базе данных
// Поля: Nickname, Email, Password, DeviceFingerPrint

type User struct {
	gorm.Model
	ID                uint   `gorm:"primaryKey;autoIncrement;"`
	Nickname          string `gorm:"unique;not null;column:nickname;type:varchar(255)"`
	Email             string `gorm:"unique;not null;column:email;type:varchar(255)"`
	Password          string `gorm:"not null;column:password;type:varchar(255)"`
	DeviceFingerPrint string `gorm:"unique;not null;column:device_fingerprint;type:varchar(255)"`
}

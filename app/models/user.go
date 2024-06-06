package models

import (
	"dating-mobile-app/app/config"
	"strings"

	"github.com/jinzhu/gorm"
)

var db *gorm.DB

type User struct {
	UserID         int    `json:"user_id"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	Name           string `json:"name"`
	Gender         string `json:"gender"`
	City           string `json:"city"`
	Verified       bool   `json:"verified"`
	LastViewDate   string `json:"last_view_date"`
	LoginExpiredAt string `json:"login_expired_at"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

type UserLog struct {
	UserLogID    int    `json:"user_log_id"`
	LoginUserID  int    `json:"login_user_id"`
	ViewedUserID int    `json:"viewed_user_id"`
	Status       string `json:"status"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type HomeBodyRequest struct {
	Email        string `json:"email"`
	ViewedUserID int    `json:"viewed_user_id"`
	Status       string `json:"status"`
}

func init() {
	config.Connect()
	db = config.GetDB()
}

func (m *User) CreateUser() *User {
	db.NewRecord(m)
	db.Create(&m)
	return m
}

func GetUserByEmail(Email string) (*User, *gorm.DB) {
	var getUser User
	db := db.Where("email=?", Email).Find(&getUser)
	return &getUser, db
}

func (m *UserLog) CreateUserLog() *UserLog {
	db.NewRecord(m)
	db.Create(&m)
	return m
}

func GetUserLogByLoginUserID(UserId int, DateNow string) ([]UserLog, *gorm.DB) {
	var getUserLog []UserLog
	db := db.Where("login_user_id = ? AND created_at BETWEEN ? AND ?", UserId, DateNow+" 00:00:00", DateNow+" 23:59:59").Find(&getUserLog)
	return getUserLog, db
}

func GetUserPartnerUser(UserIdList []string) (*User, *gorm.DB) {
	var getUser User
	db := db.Where("user_id NOT IN (" + strings.Join(UserIdList, ",") + ")").Find(&getUser)
	return &getUser, db
}

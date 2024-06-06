package controllers

import (
	"bytes"
	"dating-mobile-app/app/models"
	"dating-mobile-app/app/utils"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

const MAX_SWIPE = 10

func Register(w http.ResponseWriter, r *http.Request) {
	var newUser = &models.User{}
	utils.ParseBody(r, newUser)

	userDetail, _ := models.GetUserByEmail(newUser.Email)
	if userDetail.UserID > 0 {
		response := map[string]string{
			"message": "email has been registered",
		}
		res, _ := json.Marshal(response)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	password := newUser.Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if len(hashedPassword) != 0 && err == nil {
		now := time.Now().Format("2006-01-02 15:04:05")
		newUser.LastViewDate = now
		newUser.LoginExpiredAt = now
		newUser.CreatedAt = now
		newUser.UpdatedAt = now

		newUser.Password = bytes.NewBuffer(hashedPassword).String()
		m := newUser.CreateUser()
		res, _ := json.Marshal(m)
		w.WriteHeader(http.StatusOK)
		w.Write(res)
		return
	} else {
		response := map[string]string{
			"message": "error while hashing password",
		}
		res, _ := json.Marshal(response)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	var loginUser = &models.User{}
	utils.ParseBody(r, loginUser)

	userDetail, db := models.GetUserByEmail(loginUser.Email)
	if userDetail.UserID < 1 {
		response := map[string]string{
			"message": "user does not exist",
		}
		w.WriteHeader(http.StatusBadRequest)
		res, _ := json.Marshal(response)
		w.Write(res)
		return
	}

	var passwordCheck = bcrypt.CompareHashAndPassword([]byte(userDetail.Password), []byte(loginUser.Password))

	if passwordCheck == nil {
		expiresAt := time.Now().Add(3600 * time.Second).Format("2006-01-02 15:04:05")
		now := time.Now().Format("2006-01-02 15:04:05")

		db.Model(&userDetail).Update("login_expired_at", expiresAt)
		db.Model(&userDetail).Update("updated_at", now)

		res, _ := json.Marshal(userDetail)
		w.WriteHeader(http.StatusOK)
		w.Write(res)
		return
	} else {
		response := map[string]string{
			"message": "incorrect password",
		}
		res, _ := json.Marshal(response)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	var loginUser = &models.User{}
	utils.ParseBody(r, loginUser)

	userDetail, db := models.GetUserByEmail(loginUser.Email)
	if userDetail.UserID < 1 {
		response := map[string]string{
			"message": "user does not exist",
		}
		w.WriteHeader(http.StatusBadRequest)
		res, _ := json.Marshal(response)
		w.Write(res)
		return
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	db.Model(&userDetail).Update("login_expired_at", now)
	db.Model(&userDetail).Update("updated_at", now)

	response := map[string]string{
		"message": "logout success",
	}
	res, _ := json.Marshal(response)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func CheckLogin(w http.ResponseWriter, r *http.Request, homeBodyReq *models.HomeBodyRequest) (bool, *models.User, *gorm.DB) {
	userDetail, db := models.GetUserByEmail(homeBodyReq.Email)

	if userDetail.UserID < 1 {
		return false, nil, nil
	} else {
		loginExpiredAt, _ := time.Parse("2006-01-02 15:04:05", userDetail.LoginExpiredAt)

		timeNowString := time.Now().Format("2006-01-02 15:04:05")
		timeNow, _ := time.Parse("2006-01-02 15:04:05", timeNowString)

		if loginExpiredAt.Before(timeNow) {
			return false, nil, nil
		}
	}

	return true, userDetail, db
}

func Home(w http.ResponseWriter, r *http.Request) {
	var message string
	var partnerUser = &models.User{}

	var homeBodyReq = &models.HomeBodyRequest{}
	utils.ParseBody(r, homeBodyReq)

	isLogin, userDetail, _ := CheckLogin(w, r, homeBodyReq)
	if !isLogin {
		response := map[string]string{
			"message": "login not valid",
		}
		w.WriteHeader(http.StatusBadRequest)
		res, _ := json.Marshal(response)
		w.Write(res)
		return
	}

	dateNowString := time.Now().Format("2006-01-02")
	userLogList, _ := models.GetUserLogByLoginUserID(userDetail.UserID, dateNowString)
	if !userDetail.Verified && len(userLogList) >= MAX_SWIPE {
		message = "you have reach maximum swipe today"
		partnerUser = nil
	} else {
		var UserIdList []string
		for _, val := range userLogList {
			UserIdList = append(UserIdList, strconv.Itoa(val.ViewedUserID))
		}
		UserIdList = append(UserIdList, strconv.Itoa(userDetail.UserID))
		partnerUser, _ = models.GetUserPartnerUser(UserIdList)

		if partnerUser.UserID < 1 {
			message = "there is not new partner anymore"
			partnerUser = nil
		}
	}

	response := map[string]interface{}{
		"message": message,
		"user":    userDetail,
		"partner": partnerUser,
	}

	w.WriteHeader(http.StatusOK)
	res, _ := json.Marshal(response)
	w.Write(res)
}

func Swipe(w http.ResponseWriter, r *http.Request) {
	var message string
	var partnerUser = &models.User{}

	var homeBodyReq = &models.HomeBodyRequest{}
	utils.ParseBody(r, homeBodyReq)

	isLogin, userDetail, _ := CheckLogin(w, r, homeBodyReq)
	if !isLogin {
		response := map[string]string{
			"message": "login not valid",
		}
		w.WriteHeader(http.StatusBadRequest)
		res, _ := json.Marshal(response)
		w.Write(res)
		return
	}

	dateNowString := time.Now().Format("2006-01-02")
	userLogList, _ := models.GetUserLogByLoginUserID(userDetail.UserID, dateNowString)

	if !userDetail.Verified && len(userLogList) >= MAX_SWIPE {
		message = "you have reach maximum swipe today"
		partnerUser = nil
	} else {
		var UserIdList []string
		for _, val := range userLogList {
			UserIdList = append(UserIdList, strconv.Itoa(val.ViewedUserID))
		}
		UserIdList = append(UserIdList, strconv.Itoa(userDetail.UserID))
		UserIdList = append(UserIdList, strconv.Itoa(homeBodyReq.ViewedUserID))
		partnerUser, _ = models.GetUserPartnerUser(UserIdList)

		if partnerUser.UserID < 1 {
			message = "there is not new partner anymore"
			partnerUser = nil
		}
	}

	var newUserLog = &models.UserLog{}
	now := time.Now().Format("2006-01-02 15:04:05")

	newUserLog.LoginUserID = userDetail.UserID
	newUserLog.ViewedUserID = homeBodyReq.ViewedUserID
	newUserLog.Status = homeBodyReq.Status
	newUserLog.CreatedAt = now
	newUserLog.UpdatedAt = now

	newUserLog.CreateUserLog()

	response := map[string]interface{}{
		"message": message,
		"partner": partnerUser,
	}

	w.WriteHeader(http.StatusOK)
	res, _ := json.Marshal(response)
	w.Write(res)
}

func VerifiedUser(w http.ResponseWriter, r *http.Request) {
	var homeBodyReq = &models.HomeBodyRequest{}
	utils.ParseBody(r, homeBodyReq)

	isLogin, userDetail, db := CheckLogin(w, r, homeBodyReq)
	if !isLogin {
		response := map[string]string{
			"message": "login not valid",
		}
		w.WriteHeader(http.StatusBadRequest)
		res, _ := json.Marshal(response)
		w.Write(res)
		return
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	db.Model(&userDetail).Update("verified", 1)
	db.Model(&userDetail).Update("updated_at", now)

	res, _ := json.Marshal(userDetail)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

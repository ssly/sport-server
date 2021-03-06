package user

import (
	"net/http"
)

// User 用户结构
type User struct {
	UUID   string
	Cookie *http.Cookie
}

var userList []User

func init() {
	userList = make([]User, 0)
}

// AddUser 添加用户
func AddUser(uuid string, cookie *http.Cookie) {
	user := User{
		UUID:   uuid,
		Cookie: cookie,
	}
	// 用户不存在于列表，则添加
	if index := HasUserByUUID(uuid); index == -1 {
		userList = append(userList, user)
	} else {
		userList[index].Cookie = cookie
	}
}

// HasUserByUUID 用户是否存在
func HasUserByUUID(uuid string) int {
	index := -1
	for i, v := range userList {
		if v.UUID == uuid {
			index = i
		}
	}
	return index
}

// HasUserByCookie 用户是否存在
func HasUserByCookie(cookie *http.Cookie) int {
	index := -1
	for i, v := range userList {
		if v.Cookie.Name == cookie.Name && v.Cookie.Value == cookie.Value {
			index = i
		}
	}
	return index
}

// GetUserInfoByCookie 获取用户信息
func GetUserInfoByCookie(cookie *http.Cookie) User {
	var user User
	for _, v := range userList {
		if v.Cookie.Name == cookie.Name && v.Cookie.Value == cookie.Value {
			user = v
		}
	}
	return user
}

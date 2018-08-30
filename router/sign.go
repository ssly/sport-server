package router

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	db "sport/db"
	"sport/user"
	"sport/utils"
	"strconv"

	"gopkg.in/mgo.v2/bson"
)

// SignIn 登录接口
func SignIn(w http.ResponseWriter, r *http.Request) {
	// var err error
	defer r.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(r.Body)
	type Body struct {
		EncryptedData string `json:"encryptedData"`
		Iv            string `json:"iv"`
		RawData       string `json:"rawData"`
		Signature     string `json:"signature"`
	}
	var body Body
	json.Unmarshal(bodyBytes, &body)
	if body.Iv == "" {
		w.WriteHeader(400)
		w.Write(utils.FormatResult(1, "Sign in error."))
		return
	}

	session := db.Session()
	defer session.Close()
	userItem := struct {
		Iv string `json:"iv"`
	}{}
	c := session.DB("ly").C("sport_user")
	c.Find(bson.M{"iv": body.Iv}).One(&userItem)
	// 未存在用户
	if userItem.Iv == "" {
		c.Insert(bson.M{
			"iv":            body.Iv,
			"rawData":       body.RawData,
			"encryptedData": body.EncryptedData,
			"signature":     body.Signature,
		})
	}
	// 生成随机token
	token := strconv.Itoa(rand.Intn(10000))

	cookie := &http.Cookie{
		Name:     "S-Access-Token",
		Value:    token,
		HttpOnly: true,
	}
	// 如果登录成功，保存session
	user.AddUser(body.Iv, cookie)

	http.SetCookie(w, cookie)
	result := struct {
		Token string `json:"token"`
	}{
		Token: token,
	}
	w.Write(utils.FormatResult(0, result))
}

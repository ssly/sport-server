package router

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	db "sport/db"
	"sport/user"
	"sport/utils"
	"strconv"

	"gopkg.in/mgo.v2/bson"
)

const wechatURL = "https://api.weixin.qq.com/sns/jscode2session"

var (
	appID     string
	appSecret string
)

func init() {
	// 数据库获取 AppID, appSecret
	session := db.Session()
	defer session.Close()
	result := struct {
		AppID     string `bson:"appId"`
		AppSecret string `bson:"appSecret"`
	}{}
	c := session.DB("ly").C("sport_wechat")
	err := c.Find(nil).One(&result)
	if err != nil {
		fmt.Println("获取小程序参数失败")
	}

	appID = result.AppID
	appSecret = result.AppSecret
}

// SignIn 登录接口
func SignIn(w http.ResponseWriter, r *http.Request) {
	// var err error
	defer r.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(r.Body)
	type Body struct {
		Code    string `json:"code"`
		RawData string `json:"rawData"`
	}
	var body Body
	json.Unmarshal(bodyBytes, &body)
	if body.RawData == "" || body.Code == "" {
		w.WriteHeader(400)
		w.Write(utils.FormatResult(1, "Sign in error."))
		return
	}

	// 获取微信小程序openid
	openID := GetWechatInfo(body.Code)
	if openID == "" {
		w.WriteHeader(400)
		w.Write(utils.FormatResult(1, "openid 获取失败"))
		return
	}

	session := db.Session()
	defer session.Close()
	userItem := struct {
		ID     bson.ObjectId `json:"id" bson:"_id"`
		OpenID string        `bson:"openId"`
	}{}
	c := session.DB("ly").C("sport_user")
	c.Find(bson.M{"openId": openID}).One(&userItem)
	// 未存在用户
	if userItem.OpenID == "" {
		c.Insert(bson.M{
			"openId":  openID,
			"rawData": body.RawData,
		})
		// 查询用户数据的ID
		c.Find(bson.M{"openId": openID}).One(&userItem)
	}
	fmt.Println("登录成功啊", userItem.ID.Hex())

	// 生成随机token
	token := strconv.Itoa(rand.Intn(10000))

	cookie := &http.Cookie{
		Name:     "S-Access-Token",
		Value:    token,
		HttpOnly: true,
	}
	// 如果登录成功，保存session
	user.AddUser(userItem.ID.Hex(), cookie)

	http.SetCookie(w, cookie)
	result := struct {
		Token string `json:"token"`
	}{
		Token: token,
	}
	w.Write(utils.FormatResult(0, result))
}

// GetWechatInfo 获取微信openid
func GetWechatInfo(code string) string {
	params := "?appid=" + appID + "&secret=" + appSecret + "&js_code=" + code + "&grant_type=authorization_code"
	resp, err := http.Get(wechatURL + params)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	body := struct {
		OpenID     string `json:"openid"`
		SessionKey string `json:"session_key"`
	}{}
	json.Unmarshal(bodyBytes, &body)
	return body.OpenID
}

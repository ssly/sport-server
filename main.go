package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	db "sport/db"
	"sport/router"
	"sport/user"
	"sport/utils"
	"strconv"
	"time"

	mw "sport/middlewares"

	"gopkg.in/mgo.v2/bson"
)

var (
	port = "5432"
)

// 获取打卡记录
// param {string} year
// param {string} month
// returns {array} 每天的打卡情况 [0, 1, 1, ...]
func gerRecord(w http.ResponseWriter, r *http.Request) {
	// 获取Cookie
	cookie, err := r.Cookie("S-Access-Token")
	var userItem user.User
	if err == nil {
		userItem = user.GetUserInfoByCookie(cookie)
	}
	if userItem.UUID == "" {
		w.WriteHeader(403)
		w.Write(utils.FormatResult(1, "用户未授权"))
		return
	}

	type Result struct {
		Year    int  `json:"year"`
		Month   int  `json:"month"`
		Date    int  `json:"date"`
		IsPunch bool `json:"isPunch" bson:"isPunch"`
	}
	session := db.Session()
	defer session.Close()
	c := session.DB("ly").C("sport_result_" + userItem.UUID)

	r.ParseForm()
	year, _ := strconv.Atoi(r.Form.Get("year"))
	month, _ := strconv.Atoi(r.Form.Get("month"))

	if r.Form.Get("date") != "" {
		var resultItem Result
		// 只需要查询一天数据
		date, _ := strconv.Atoi(r.Form.Get("date"))
		c.Find(bson.M{
			"year":  year,
			"month": month,
			"date":  date,
		}).One(&resultItem)
		w.Write(utils.FormatResult(0, resultItem))
		return
	}

	dayInCurMonth := utils.GetDayInMonthOf(year, month)
	resultList := make([]Result, dayInCurMonth)

	c.Find(bson.M{
		"year":  year,
		"month": month,
	}).All(&resultList)

	defaultList := make([]Result, dayInCurMonth)
	for i := 0; i < dayInCurMonth; i++ {
		defaultList[i].Year = year
		defaultList[i].Month = month
		defaultList[i].Date = i + 1
		// 查找处理已打卡的日期
		for _, v := range resultList {
			if defaultList[i].Date == v.Date {
				defaultList[i].IsPunch = true
			}
		}
	}

	w.Write(utils.FormatResult(0, defaultList))
}

// 保存打卡记录
func updateRecord(w http.ResponseWriter, r *http.Request) {
	var err error
	// 获取用户信息，建立对应的表
	cookie, err := r.Cookie("S-Access-Token")
	var userItem user.User
	if err == nil {
		userItem = user.GetUserInfoByCookie(cookie)
	}
	if userItem.UUID == "" {
		w.WriteHeader(403)
		w.Write(utils.FormatResult(1, "用户未授权"))
		return
	}

	session := db.Session()
	defer session.Close()

	defer r.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(r.Body)
	type Body struct {
		Year  int `json:"year"`
		Month int `json:"month"`
		Date  int `json:"date"`
	}
	var bodyItem Body

	json.Unmarshal(bodyBytes, &bodyItem)

	// 如果没有传入年月日，则默认取当天的年月日
	if bodyItem.Year == 0 {
		bodyItem.Year = time.Now().Year()
	}
	if bodyItem.Month == 0 {
		bodyItem.Month = int(time.Now().Month())
	}
	if bodyItem.Date == 0 {
		bodyItem.Date = time.Now().Day()
	}

	currentTime := time.Now().UnixNano() / 1000000 // ms时间戳
	queryItem := bson.M{
		"year":  bodyItem.Year,
		"month": bodyItem.Month,
		"date":  bodyItem.Date,
	}

	todayItem := make(map[string]interface{})

	// 返回给前端的数据
	var res struct {
		IsPunch bool `json:"isPunch"`
	}
	c := session.DB("ly").C("sport_result_" + userItem.UUID)
	err = c.Find(queryItem).One(&todayItem)

	if todayItem["isPunch"] == nil {
		// 打卡
		err = c.Insert(bson.M{
			"year":       bodyItem.Year,
			"month":      bodyItem.Month,
			"date":       bodyItem.Date,
			"isPunch":    true,
			"createTime": currentTime,
		})
		if err != nil {
			fmt.Println("Insert error: ", err)
		}
		res.IsPunch = true
	} else {
		// 取消打卡，删除数据
		err = c.Remove(bson.M{
			"year":  bodyItem.Year,
			"month": bodyItem.Month,
			"date":  bodyItem.Date,
		})
		if err != nil {
			fmt.Println("Remove error: ", err)
		}

		res.IsPunch = false
	}

	w.Write(utils.FormatResult(0, res))
}

func startServer() {

	fmt.Println("start server at " + port)
	// 登录接口
	http.Handle("/api/sign-in", mw.LogMiddleware(http.HandlerFunc(router.SignIn)))
	// 获取打卡记录
	http.Handle("/api/sport/get-record", mw.LogMiddleware(http.HandlerFunc(gerRecord)))
	// 保存打卡记录
	http.Handle("/api/sport/update-record", mw.LogMiddleware(http.HandlerFunc(updateRecord)))

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("Listen Error")
	}
}

func main() {
	startServer()
}

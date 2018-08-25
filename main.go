package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	db "sport/db"
	"sport/utils"
	"strconv"
	"time"

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
	session := db.Session()
	defer session.Close()

	r.ParseForm()
	year, _ := strconv.Atoi(r.Form.Get("year"))
	month, _ := strconv.Atoi(r.Form.Get("month"))
	// date, _ := strconv.Atoi(r.Form.Get("date")) // TODO 增加按天过滤

	dayInCurMonth := utils.GetDayInMonthOf(month)

	type Result struct {
		Year    int  `json:"year"`
		Month   int  `json:"month"`
		Date    int  `json:"date"`
		IsPunch bool `json:"isPunch"`
	}

	resultList := make([]Result, dayInCurMonth)
	c := session.DB("ly").C("sport_result")
	c.Find(bson.M{
		"year":  year,
		"month": month,
		// "date": date,
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

	w.Write(formatResult(defaultList))
}

// 保存打卡记录
func updateRecord(w http.ResponseWriter, r *http.Request) {
	var err error
	body, _ := ioutil.ReadAll(r.Body)
	type Body struct {
		Year  int `json:"year"`
		Month int `json:"month"`
		Date  int `json:"date"`
	}
	var bodyItem Body

	json.Unmarshal(body, &bodyItem)

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

	currentTime := time.Now().UnixNano() / 1000000 // 时间戳
	queryItem := bson.M{
		"year":  bodyItem.Year,
		"month": bodyItem.Month,
		"date":  bodyItem.Date,
	}

	session := db.Session()
	defer session.Close()

	todayItem := make(map[string]interface{})

	// 返回给前端的数据
	var res struct {
		IsPunch bool `json:"isPunch"`
	}
	c := session.DB("ly").C("sport_result")
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

	w.Write(formatResult(res))
}

func startServer() {

	fmt.Println("start server at " + port)
	http.HandleFunc("/sport/get-record", gerRecord)       // 获取打卡记录
	http.HandleFunc("/sport/update-record", updateRecord) // 保存打卡记录

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("Listen Error")
	}
}

func formatResult(v interface{}) []byte {
	result := &struct {
		Code byte        `json:"code"`
		Date interface{} `json:"date"`
	}{
		Code: 0,
		Date: v,
	}
	message, err := json.Marshal(result)
	if err != nil {
		fmt.Println("json.Marshal error.")
	}
	return message
}

func main() {
	startServer()
}

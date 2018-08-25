package datebases

import (
	"fmt"

	"gopkg.in/mgo.v2"
)

const (
	db  = "ly"
	url = "mongodb://localhost:27017"
)

var session *mgo.Session

func init() {
	// 初始化数据库
	_session, err := mgo.Dial(url + "/" + db)
	if err != nil {
		fmt.Println("db/mango.go: connect db error.")
	}

	session = _session
}

// Session 获取数据库 session
func Session() *mgo.Session {
	return session.Copy()
}

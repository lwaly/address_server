package mymongo

import (
	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2"
)

var session *mgo.Session

// Conn return mongodb session.
func Conn() *mgo.Session {
	return session.Copy()
}

/*
func Close() {
	session.Close()
}
*/

func init() {
	str1 := beego.AppConfig.String("mongodb::url")
	var err error
	session, err = mgo.Dial(str1)
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
}

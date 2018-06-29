package models

import (
	"AddressServer/models/mymongo"

	"github.com/astaxie/beego"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type AddressRecord struct {
	Id       bson.ObjectId `bson:"_id" json:"Id,omitempty"`
	UId      uint64        `bson:"UId" json:"UId,omitempty"`
	Ip       int64         `bson:"Ip" json:"Ip,omitempty"`
	Country  uint16        `bson:"Country" json:"Country,omitempty"`
	Province uint8         `bson:"Province" json:"Province,omitempty"`
	City     uint8         `bson:"City" json:"City,omitempty"`
	Time     int64         `bson:"Time" json:"Time,omitempty"`
}

type StAddressRecord struct {
	Id       uint64 `bson:"_id" json:"Id,omitempty"`
	Country  uint16 `bson:"Country" json:"Country,omitempty"`
	Province uint8  `bson:"Province" json:"Province,omitempty"`
	City     uint8  `bson:"City" json:"City,omitempty"`
	Time     int64  `bson:"Time" json:"Time,omitempty"`
}

//插入最后登陆记录
func (info *AddressRecord) Insert() (code int, err error) {
	mConn := mymongo.Conn()
	defer mConn.Close()
	info.Id = bson.NewObjectId()
	c := mConn.DB(db).C(address_record)
	err = c.Insert(info)

	if err != nil {
		if mgo.IsDup(err) {
			code = ErrDupRows
		} else {
			code = ErrDatabase
		}
		beego.Error("fail to insert.uid=%d,err=%s", info.Id, err.Error())
		return
	} else {
		code = 0
	}

	return
}

//查询最新登录同城信息
func (info *AddressRecord) FindLastestLogin(field int, PageSzie int32, PageIndex int32) (Recodes []StAddressRecord, code int, err error) {
	mConn := mymongo.Conn()
	defer mConn.Close()

	c := mConn.DB(db).C(address_record)

	switch field {
	case CITY:
		err = c.Find(bson.M{"UId": info.UId, "Country": info.Country, "Province": info.Province, "City": info.City}).Sort("-Time").Skip(int(PageSzie * (PageIndex - 1))).Limit(int(PageSzie)).All(&Recodes)
	case PROVINCE:
		err = c.Find(bson.M{"UId": info.UId, "Country": info.Country, "Province": info.Province}).Sort("-Time").Skip(int(PageSzie * (PageIndex - 1))).Limit(int(PageSzie)).All(&Recodes)
	case COUNTRY:
		err = c.Find(bson.M{"UId": info.UId, "Country": info.Country}).Sort("-Time").Skip(int(PageSzie * (PageIndex - 1))).Limit(int(PageSzie)).All(&Recodes)
	case ALL:
		err = c.Find(bson.M{"UId": info.UId}).Sort("-Time").Skip(int(PageSzie * (PageIndex - 1))).Limit(int(PageSzie)).All(&Recodes)
	}

	if err != nil {
		if err == mgo.ErrNotFound {
			code = ErrNotFound
		} else {
			code = ErrDatabase
		}
		beego.Error("fail to find UserMsg.uid=%d,err=%s", info.UId, err.Error())
		return
	}

	return
}

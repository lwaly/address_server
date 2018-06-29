package models

import (
	"AddressServer/models/mymongo"

	"github.com/astaxie/beego"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type AddressServer struct {
	Id         uint64 `bson:"_id" json:"Id,omitempty"`
	Ip         int64  `bson:"Ip" json:"Ip,omitempty"`
	Country    uint16 `bson:"Country" json:"Country,omitempty"`
	Province   uint8  `bson:"Province" json:"Province,omitempty"`
	City       uint8  `bson:"City" json:"City,omitempty"`
	Status     uint8  `bson:"Status" json:"Status,omitempty"`
	CreateTime int64  `bson:"CreateTime" json:"CreateTimeId,omitempty"`
	UpdateTime int64  `bson:"UpdateTime" json:"UpdateTime,omitempty"`
}

type StLastestLogin struct {
	Id         uint64 `bson:"_id" json:"Id,omitempty"`
	Country    uint16 `bson:"Country" json:"Country,omitempty"`
	Province   uint8  `bson:"Province" json:"Province,omitempty"`
	City       uint8  `bson:"City" json:"City,omitempty"`
	UpdateTime int64  `bson:"UpdateTime" json:"UpdateTime,omitempty"`
}

type StLastestRegister struct {
	Id         uint64 `bson:"_id" json:"Id,omitempty"`
	Country    uint16 `bson:"Country" json:"Country,omitempty"`
	Province   uint8  `bson:"Province" json:"Province,omitempty"`
	City       uint8  `bson:"City" json:"City,omitempty"`
	CreateTime int64  `bson:"CreateTime" json:"CreateTime,omitempty"`
}

//插入最后登陆记录
func (info *AddressServer) Upsert() (code int, err error) {
	mConn := mymongo.Conn()
	defer mConn.Close()

	c := mConn.DB(db).C(address_server)
	err = c.Find(bson.M{"_id": info.Id}).One(nil)
	if err != nil {
		if err == mgo.ErrNotFound {
			err = c.Insert(info)

		} else {
			code = ErrDatabase
			beego.Error("fail to MongoFindByField.err=%s,uid=%d", err.Error(), info.Id)
			return
		}
	} else {
		err = c.Update(bson.M{"_id": info.Id}, bson.M{"$set": bson.M{"Ip": info.Ip, "Country": info.Country, "Province": info.Province, "City": info.City, "UpdateTime": info.UpdateTime, "Status": info.Status}})
	}

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
func (info *AddressServer) FindLastestLogin(field int, PageSzie int32, PageIndex int32) (Recodes []StLastestLogin, code int, err error) {
	mConn := mymongo.Conn()
	defer mConn.Close()

	c := mConn.DB(db).C(address_server)

	switch field {
	case CITY:
		err = c.Find(bson.M{"Country": info.Country, "Province": info.Province, "City": info.City, "Status": info.Status}).Sort("-UpdateTime").Skip(int(PageSzie * (PageIndex - 1))).Limit(int(PageSzie)).All(&Recodes)
	case PROVINCE:
		err = c.Find(bson.M{"Country": info.Country, "Province": info.Province, "Status": info.Status}).Sort("-UpdateTime").Skip(int(PageSzie * (PageIndex - 1))).Limit(int(PageSzie)).All(&Recodes)
	case COUNTRY:
		err = c.Find(bson.M{"Country": info.Country, "Status": info.Status}).Sort("-UpdateTime").Skip(int(PageSzie * (PageIndex - 1))).Limit(int(PageSzie)).All(&Recodes)
	}

	if err != nil {
		if err == mgo.ErrNotFound {
			code = ErrNotFound
		} else {
			code = ErrDatabase
		}
		beego.Error("fail to find UserMsg.uid=%d,err=%s", info.Id, err.Error())
		return
	}

	return
}

//查询最新注册同城信息
func (info *AddressServer) FindLastestRegister(field int, PageSzie int32, PageIndex int32) (Recodes []StLastestRegister, code int, err error) {
	mConn := mymongo.Conn()
	defer mConn.Close()

	c := mConn.DB(db).C(address_server)

	switch field {
	case CITY:
		err = c.Find(bson.M{"Country": info.Country, "Province": info.Province, "City": info.City, "Status": info.Status}).Sort("-CreateTime").Skip(int(PageSzie * (PageIndex - 1))).Limit(int(PageSzie)).All(&Recodes)
	case PROVINCE:
		err = c.Find(bson.M{"Country": info.Country, "Province": info.Province, "Status": info.Status}).Sort("-CreateTime").Skip(int(PageSzie * (PageIndex - 1))).Limit(int(PageSzie)).All(&Recodes)
	case COUNTRY:
		err = c.Find(bson.M{"Country": info.Country, "Status": info.Status}).Sort("-CreateTime").Skip(int(PageSzie * (PageIndex - 1))).Limit(int(PageSzie)).All(&Recodes)
	}

	if err != nil {
		if err == mgo.ErrNotFound {
			code = ErrNotFound
		} else {
			code = ErrDatabase
		}
		beego.Error("fail to find UserMsg.uid=%d,objectid=%d,err=%s", info.Id, err.Error())
		return
	}

	return
}

func (info *AddressServer) FindNearCount(field int) (count int, code int, err error) {
	mConn := mymongo.Conn()
	defer mConn.Close()

	c := mConn.DB(db).C(address_server)

	switch field {
	case CITY:
		count, err = c.Find(bson.M{"Country": info.Country, "Province": info.Province, "City": info.City, "Status": info.Status}).Count()
	case PROVINCE:
		count, err = c.Find(bson.M{"Country": info.Country, "Province": info.Province, "Status": info.Status}).Count()
	case COUNTRY:
		count, err = c.Find(bson.M{"Country": info.Country, "Status": info.Status}).Count()
	}

	if err != nil {
		if err == mgo.ErrNotFound {
			code = ErrNotFound
		} else {
			code = ErrDatabase
		}
		beego.Error("fail to find UserMsg.uid=%d,err=%s", info.Id, err.Error())
		return
	}

	return
}

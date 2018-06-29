package models

import (
	"AddressServer/models/mymongo"

	"github.com/astaxie/beego"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type AddressGps struct {
	Id         uint64  `bson:"_id" json:"Id,omitempty"`
	Location   GeoJson `bson:"location" json:"location"`
	Status     uint8   `bson:"Status" json:"Status,omitempty"`
	CreateTime int64   `bson:"CreateTime" json:"CreateTimeId,omitempty"`
	UpdateTime int64   `bson:"UpdateTime" json:"UpdateTime,omitempty"`
}

type GeoJson struct {
	Type        string     `bson:"type" json:"type"`
	Coordinates [2]float64 `bson:"coordinates" json:"coordinates"`
}

type StAddressLastestLoginGps struct {
	Id         uint64  `bson:"_id" json:"Id,omitempty"`
	Longitude  float64 `bson:"Longitude" json:"Longitude,omitempty"`
	Latitude   float64 `bson:"Latitude" json:"Latitude,omitempty"`
	Dis        float64 `bson:"Dis" json:"Dis,omitempty"`
	UpdateTime int64   `bson:"UpdateTime" json:"UpdateTime,omitempty"`
}

const (
	MaxRet = 30
)

//插入最后登陆记录
func (info *AddressGps) Upsert() (code int, err error) {
	mConn := mymongo.Conn()
	defer mConn.Close()

	c := mConn.DB(db).C(address_gps)
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
		err = c.Update(bson.M{"_id": info.Id}, bson.M{"$set": bson.M{"location": info.Location, "UpdateTime": info.UpdateTime, "Status": info.Status}})
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
func (info *AddressGps) FindLastestLogin(Min int32, Max int32) (count int, Recodes []StAddressLastestLoginGps, code int, err error) {
	mConn := mymongo.Conn()
	defer mConn.Close()

	c := mConn.DB(db)

	result := bson.M{}
	err = c.Run(
		bson.D{
			{"geoNear", address_gps},
			{"spherical", true},
			{"near", bson.M{"type": "Point", "coordinates": [2]float64{info.Location.Coordinates[0], info.Location.Coordinates[1]}}},
			{"minDistance", Min},
			{"maxDistance", Max},
			{"num", MaxRet}},
		&result)

	if err != nil {
		if err == mgo.ErrNotFound {
			code = ErrNotFound
		} else {
			code = ErrDatabase
		}
		beego.Error("fail to find UserMsg.uid=%d,err=%s", info.Id, err.Error())
		return
	}

	valueStats, ok := result["stats"]
	if !ok {
		err = mgo.ErrNotFound
		code = ErrDatabase
		beego.Error("fail to find UserMsg.uid=%d,err=%s", info.Id, err.Error())
		return
	} else {
		tempMap := map[string]interface{}{}
		b, errTemp := bson.Marshal(valueStats)
		if nil != err {
			err = errTemp
			code = ErrDatabase
			beego.Error("fail to find UserMsg.uid=%d,err=%s", info.Id, err.Error())
			return
		}
		err = bson.Unmarshal(b, &tempMap)
		if nil != err {
			code = ErrDatabase
			beego.Error("fail to find UserMsg.uid=%d,err=%s", info.Id, err.Error())
			return
		}
		ff, ok := tempMap["objectsLoaded"]
		if true != ok {
			err = mgo.ErrNotFound
			code = ErrDatabase
			beego.Error("fail to find UserMsg.uid=%d,err=%s", info.Id, err.Error())
			return
		}
		count = ff.(int)
	}

	valueRes, ok := result["results"]
	if ok {
		tempMap := map[string]interface{}{}
		bb, errTemp := bson.Marshal(valueRes)
		if nil != errTemp {
			err = errTemp
			code = ErrDatabase
			beego.Error("fail to find UserMsg.uid=%d,err=%s", info.Id, err.Error())
			return
		}

		err = bson.Unmarshal(bb, &tempMap)
		if nil != err {
			code = ErrDatabase
			beego.Error("fail to find UserMsg.uid=%d,err=%s", info.Id, err.Error())
			return
		}

		for _, value3 := range tempMap {
			var record StAddressLastestLoginGps
			tempMap1 := map[string]interface{}{}
			bb, errTemp := bson.Marshal(value3)
			if nil != errTemp {
				err = errTemp
				code = ErrDatabase
				beego.Error("fail to find UserMsg.uid=%d,err=%s", info.Id, err.Error())
				return
			}

			err = bson.Unmarshal(bb, &tempMap1)

			if nil != err {
				code = ErrDatabase
				beego.Error("fail to find UserMsg.uid=%d,err=%s", info.Id, err.Error())
				return
			}

			//获取距离
			ff, ok := tempMap1["dis"]
			if true != ok {
				err = mgo.ErrNotFound
				code = ErrDatabase
				beego.Error("fail to find UserMsg.uid=%d,err=%s", info.Id, err.Error())
				return
			}
			record.Dis = ff.(float64)

			//获取用户id
			ff, ok = tempMap1["obj"]
			if true != ok {
				err = mgo.ErrNotFound
				code = ErrDatabase
				beego.Error("fail to find UserMsg.uid=%d,err=%s", info.Id, err.Error())
				return
			}
			tempMap2 := map[string]interface{}{}
			bb, errTemp = bson.Marshal(ff)
			if nil != errTemp {
				err = errTemp
				code = ErrDatabase
				beego.Error("fail to find UserMsg.uid=%d,err=%s", info.Id, err.Error())
				return
			}

			err = bson.Unmarshal(bb, &tempMap2)

			if nil != err {
				code = ErrDatabase
				beego.Error("fail to find UserMsg.uid=%d,err=%s", info.Id, err.Error())
				return
			}

			record.Id = uint64(tempMap2["_id"].(int64))

			//获取时间
			record.UpdateTime = tempMap2["UpdateTime"].(int64)
			//获取经纬度坐标
			ff, ok = tempMap2["location"]
			if true != ok {
				err = mgo.ErrNotFound
				code = ErrDatabase
				beego.Error("fail to find UserMsg.uid=%d,err=%s", info.Id, err.Error())
				return
			}
			tempMap2 = map[string]interface{}{}
			bb, errTemp = bson.Marshal(ff)
			if nil != errTemp {
				err = errTemp
				code = ErrDatabase
				beego.Error("fail to find UserMsg.uid=%d,err=%s", info.Id, err.Error())
				return
			}
			tempMap3 := map[string]interface{}{}
			err = bson.Unmarshal(bb, &tempMap3)

			if nil != err {
				code = ErrDatabase
				beego.Error("fail to find UserMsg.uid=%d,err=%s", info.Id, err.Error())
				return
			}

			ff, ok = tempMap3["coordinates"]
			if true != ok {
				err = mgo.ErrNotFound
				code = ErrDatabase
				beego.Error("fail to find UserMsg.uid=%d,err=%s", info.Id, err.Error())
				return
			}
			fff := ff.([]interface{})

			record.Longitude = fff[0].(float64)
			record.Latitude = fff[1].(float64)

			Recodes = append(Recodes, record)
		}
	}
	return
}

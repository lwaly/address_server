package controllers

import "AddressServer/models"

type UploadIpRep struct {
	Id    uint64 `bson:"_id" json:"Id,omitempty"`
	Ip    int64  `bson:"Ip" json:"Ip,omitempty"`
	Token string `json:"Token" valid:"Required"`
}

type UploadIpRsp struct {
	Country  uint16 `bson:"Country" json:"Country,omitempty"`
	Province uint8  `bson:"Province" json:"Province,omitempty"`
	City     uint8  `bson:"City" json:"City,omitempty"`
}

type SearchLastestLoginRep struct {
	Id        uint64 `bson:"_id" json:"Id,omitempty"`
	Ip        int64  `bson:"Ip" json:"Ip,omitempty"`
	Level     uint8  `bson:"Level" json:"Level,omitempty"`
	PageSize  int32  `json:"PageSize" valid:"Required"`
	PageIndex int32  `json:"PageIndex" valid:"Required"`
	Token     string `json:"Token" valid:"Required"`
}

type SearchLastestLoginRsp struct {
	Total int32                   `json:"Total" valid:"Required"`
	Data  []models.StLastestLogin `json:"Data" valid:"Required"`
}

type SearchLastestRegisterRep struct {
	Id        uint64 `bson:"_id" json:"Id,omitempty"`
	Ip        int64  `bson:"Ip" json:"Ip,omitempty"`
	Level     uint8  `bson:"Level" json:"Level,omitempty"`
	PageSize  int32  `json:"PageSize" valid:"Required"`
	PageIndex int32  `json:"PageIndex" valid:"Required"`
	Token     string `json:"Token" valid:"Required"`
}

type SearchLastestRegisterRsp struct {
	Total int32                      `json:"Total" valid:"Required"`
	Data  []models.StLastestRegister `json:"Data" valid:"Required"`
}

type UploadGpsRep struct {
	Id        uint64  `bson:"_id" json:"Id,omitempty"`
	Longitude float64 `bson:"Longitude" json:"Longitude,omitempty"`
	Latitude  float64 `bson:"Latitude" json:"Latitude,omitempty"`
	Token     string  `json:"Token" valid:"Required"`
}

type UploadGpsRsp struct {
}

type SearchGpsRep struct {
	Id        uint64  `bson:"_id" json:"Id,omitempty"`
	Longitude float64 `bson:"Longitude" json:"Longitude,omitempty"`
	Latitude  float64 `bson:"Latitude" json:"Latitude,omitempty"`
	MinDis    int32   `json:"MinDis" valid:"Required"`
	MaxDis    int32   `json:"MaxDis" valid:"Required"`
	Token     string  `json:"Token" valid:"Required"`
}

type SearchGpsRsp struct {
	Total int32                             `json:"Total" valid:"Required"`
	Data  []models.StAddressLastestLoginGps `json:"Data" valid:"Required"`
}

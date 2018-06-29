package controllers

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"AddressServer/models"

	"github.com/astaxie/beego"
)

type AddressServerController struct {
	beego.Controller
}

type Ip2Region struct {
	// db file handler
	dbFileHandler *os.File

	//header block info

	headerSip []int64
	headerPtr []int64
	headerLen int64

	// super block index info
	firstIndexPtr int64
	lastIndexPtr  int64
	totalBlocks   int64

	// for memory mode only
	// the original db binary string

	dbBinStr []byte
	dbFile   string
}

const (
	INDEX_BLOCK_LENGTH  = 12
	TOTAL_HEADER_LENGTH = 8192
)

type IpInfo struct {
	Country  uint16
	Province uint8
	City     uint8
}

var ip2Region Ip2Region

var mapCountry map[string]uint16
var mapCity map[string]map[string]string

func Init() int {
	file, err := os.Open(beego.AppConfig.String("IpFile"))
	if err != nil {
		return -1
	}

	ip2Region = Ip2Region{
		dbFile:        beego.AppConfig.String("IpFile"),
		dbFileHandler: file,
	}

	mapCountry = make(map[string]uint16)
	mapCity = make(map[string]map[string]string)
	file, err = os.Open(beego.AppConfig.String("CountryFile"))
	if err != nil {
		return -1
	}
	buf := bufio.NewReader(file)
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return -1
		}
		line = strings.Replace(line, "\n", "", -1)
		lineSlice := strings.Split(string(line), " ")
		_, ok := mapCountry[lineSlice[1]]
		if ok == true {
			return -1
		}
		countryId, _ := strconv.Atoi(lineSlice[0])
		if 156 == countryId {
			mapCountry[lineSlice[1]] = uint16(countryId)
		}
		mapCountry[lineSlice[1]] = uint16(countryId)
	}

	file, err = os.Open(beego.AppConfig.String("CityFile"))
	if err != nil {
		return -1
	}
	buf = bufio.NewReader(file)
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return -1
		}
		line = strings.Replace(line, "\n", "", -1)
		lineSlice := strings.Split(string(line), " ")
		if 3 == len(lineSlice) {
			iter, ok := mapCity[lineSlice[0]]
			if ok == true {
				_, ok = iter[lineSlice[1]]
				if ok == true {
					return -1
				} else {
					//cityId, _ := strconv.Atoi(lineSlice[2][2:4])
					iter[lineSlice[1]] = lineSlice[2][0:4]
				}
			} else {
				return -1
			}
		} else if 2 == len(lineSlice) {
			_, ok := mapCity[lineSlice[0]]
			if ok == true {
				return -1
			}
			mapTemp := make(map[string]string)
			//特殊直辖市，只有一级
			if lineSlice[0] == "北京" || lineSlice[0] == "上海" || lineSlice[0] == "天津" || lineSlice[0] == "重庆" {
				mapTemp[lineSlice[0]] = lineSlice[1][0:4]
			}
			mapCity[lineSlice[0]] = mapTemp
		} else {
			return -1
		}
	}
	return 0
}

func getIpInfo(cityId int64, line []byte, ipInfo *models.AddressServer) bool {

	lineSlice := strings.Split(string(line), "|")
	length := len(lineSlice)
	//ipInfo.CityId = cityId
	if length < 5 {
		for i := 0; i <= 5-length; i++ {
			lineSlice = append(lineSlice, "")
		}
	}

	var ok bool
	if ipInfo.Country, ok = mapCountry[lineSlice[0]]; true != ok {
		return ok
	}

	if iter, ok := mapCity[lineSlice[2]]; true != ok {
		return ok
	} else {
		var strCity string
		if strCity, ok = iter[lineSlice[3]]; true != ok {
			return ok
		}
		tempId, _ := strconv.Atoi(strCity[0:2])
		ipInfo.Province = uint8(tempId)
		tempId, _ = strconv.Atoi(strCity[2:4])
		ipInfo.City = uint8(tempId)
	}

	return true
}

func (this *Ip2Region) Close() {
	this.dbFileHandler.Close()
}

func getLong(b []byte, offset int64) int64 {

	val := (int64(b[offset]) |
		int64(b[offset+1])<<8 |
		int64(b[offset+2])<<16 |
		int64(b[offset+3])<<24)

	return val

}

func (this *Ip2Region) MemorySearch(ip int64, ipInfo *models.AddressServer) (err error) {
	if this.totalBlocks == 0 {
		this.dbBinStr, err = ioutil.ReadFile(this.dbFile)

		if err != nil {

			return err
		}

		this.firstIndexPtr = getLong(this.dbBinStr, 0)
		this.lastIndexPtr = getLong(this.dbBinStr, 4)
		this.totalBlocks = (this.lastIndexPtr-this.firstIndexPtr)/INDEX_BLOCK_LENGTH + 1
	}

	h := this.totalBlocks
	var dataPtr, l int64
	for l <= h {

		m := (l + h) >> 1
		p := this.firstIndexPtr + m*INDEX_BLOCK_LENGTH
		sip := getLong(this.dbBinStr, p)
		if ip < sip {
			h = m - 1
		} else {
			eip := getLong(this.dbBinStr, p+4)
			if ip > eip {
				l = m + 1
			} else {
				dataPtr = getLong(this.dbBinStr, p+8)
				break
			}
		}
	}
	if dataPtr == 0 {
		return errors.New("not found")
	}

	dataLen := ((dataPtr >> 24) & 0xFF)
	dataPtr = (dataPtr & 0x00FFFFFF)
	getIpInfo(getLong(this.dbBinStr, dataPtr), this.dbBinStr[(dataPtr)+4:dataPtr+dataLen], ipInfo)
	return nil

}

func (info *AddressServerController) UploadIp() {
	form := UploadIpRep{}
	if err := json.Unmarshal(info.Ctx.Input.RequestBody, &form); nil != err {
		beego.Error("Unmarshal:", err)
		info.Data["json"] = NewErrorInfo(-1, err.Error())
		return
	}

	defer info.ServeJSON()

	addressInfo := models.AddressServer{
		Id: form.Id,
		Ip: form.Ip,
	}

	if err := ip2Region.MemorySearch(addressInfo.Ip, &addressInfo); nil != err {
		beego.Error("find:", err)
		info.Data["json"] = NewErrorInfo(-1, err.Error())
		return
	}

	addressInfo.CreateTime = int64(time.Now().UnixNano())
	addressInfo.UpdateTime = addressInfo.CreateTime
	addressInfo.Status = models.VALID

	if code, err := addressInfo.Upsert(); nil != err {
		beego.Error("find:", err)
		info.Data["json"] = NewErrorInfo(code, err.Error())
		return
	}

	addressRecord := models.AddressRecord{
		UId: form.Id,
		Ip:  form.Ip,
	}

	if code, err := addressRecord.Insert(); nil != err {
		beego.Error("find:", err)
		info.Data["json"] = NewErrorInfo(code, err.Error())
		return
	}

	msg := UploadIpRsp{Country: addressInfo.Country, Province: addressInfo.Province, City: addressInfo.City}
	if b, err := json.Marshal(msg); nil == err {
		info.Data["json"] = NewNormalInfo(b)
	} else {
		beego.Error("ErrUnknown:", err)
		info.Data["json"] = NewErrorInfo(ErrCodeUnknown, ErrUnknown)
		return
	}

	return
}

func (info *AddressServerController) SearchLastestRegister() {
	form := SearchLastestRegisterRep{}
	if err := json.Unmarshal(info.Ctx.Input.RequestBody, &form); nil != err {
		beego.Error("Unmarshal:", err)
		info.Data["json"] = NewErrorInfo(-1, err.Error())
		return
	}

	defer info.ServeJSON()

	addressInfo := models.AddressServer{
		Id: form.Id,
		Ip: form.Ip,
	}

	if err := ip2Region.MemorySearch(addressInfo.Ip, &addressInfo); nil != err {
		beego.Error("find:", err)
		info.Data["json"] = NewErrorInfo(-1, err.Error())
		return
	}

	if records, code, err := addressInfo.FindLastestRegister(int(form.Level), form.PageSize, form.PageIndex); nil != err {
		beego.Error("find:", err)
		info.Data["json"] = NewErrorInfo(code, err.Error())
		return
	} else {
		if count, code, err := addressInfo.FindNearCount(int(form.Level)); nil != err {
			beego.Error("find:", err)
			info.Data["json"] = NewErrorInfo(code, err.Error())
			return
		} else {
			retInfo := SearchLastestRegisterRsp{
				Total: int32(count),
				Data:  records,
			}

			if b, err := json.Marshal(retInfo); nil == err {
				info.Data["json"] = NewNormalInfo(b)
			} else {
				beego.Error("ErrUnknown:", err)
				info.Data["json"] = NewErrorInfo(ErrCodeUnknown, ErrUnknown)
				return
			}
		}
	}

	return
}

func (info *AddressServerController) SearchLastestLogin() {
	form := SearchLastestLoginRep{}
	if err := json.Unmarshal(info.Ctx.Input.RequestBody, &form); nil != err {
		beego.Error("Unmarshal:", err)
		info.Data["json"] = NewErrorInfo(-1, err.Error())
		return
	}

	defer info.ServeJSON()

	addressInfo := models.AddressServer{
		Id: form.Id,
		Ip: form.Ip,
	}

	if err := ip2Region.MemorySearch(addressInfo.Ip, &addressInfo); nil != err {
		beego.Error("find:", err)
		info.Data["json"] = NewErrorInfo(-1, err.Error())
		return
	}

	if records, code, err := addressInfo.FindLastestLogin(int(form.Level), form.PageSize, form.PageIndex); nil != err {
		beego.Error("find:", err)
		info.Data["json"] = NewErrorInfo(code, err.Error())
		return
	} else {
		if count, code, err := addressInfo.FindNearCount(int(form.Level)); nil != err {
			beego.Error("find:", err)
			info.Data["json"] = NewErrorInfo(code, err.Error())
			return
		} else {
			retInfo := SearchLastestLoginRsp{
				Total: int32(count),
				Data:  records,
			}

			if b, err := json.Marshal(retInfo); nil == err {
				info.Data["json"] = NewNormalInfo(b)
			} else {
				beego.Error("ErrUnknown:", err)
				info.Data["json"] = NewErrorInfo(ErrCodeUnknown, ErrUnknown)
				return
			}
		}
	}

	return
}

func (info *AddressServerController) UploadGps() {
	form := UploadGpsRep{}
	if err := json.Unmarshal(info.Ctx.Input.RequestBody, &form); nil != err {
		beego.Error("Unmarshal:", err)
		info.Data["json"] = NewErrorInfo(-1, err.Error())
		return
	}

	defer info.ServeJSON()

	addressInfo := models.AddressGps{
		Id:       form.Id,
		Location: models.GeoJson{"Point", [2]float64{form.Longitude, form.Latitude}},
	}

	addressInfo.CreateTime = int64(time.Now().UnixNano())
	addressInfo.UpdateTime = addressInfo.CreateTime
	addressInfo.Status = models.VALID

	if code, err := addressInfo.Upsert(); nil != err {
		beego.Error("find:", err)
		info.Data["json"] = NewErrorInfo(code, err.Error())
		return
	}

	msg := UploadGpsRsp{}
	if b, err := json.Marshal(msg); nil == err {
		info.Data["json"] = NewNormalInfo(b)
	} else {
		beego.Error("ErrUnknown:", err)
		info.Data["json"] = NewErrorInfo(ErrCodeUnknown, ErrUnknown)
		return
	}

	return
}

func (info *AddressServerController) SearchGps() {
	form := SearchGpsRep{}
	if err := json.Unmarshal(info.Ctx.Input.RequestBody, &form); nil != err {
		beego.Error("Unmarshal:", err)
		info.Data["json"] = NewErrorInfo(-1, err.Error())
		return
	}

	defer info.ServeJSON()

	addressInfo := models.AddressGps{
		Id:       form.Id,
		Location: models.GeoJson{"Point", [2]float64{form.Longitude, form.Latitude}},
	}

	if count, records, code, err := addressInfo.FindLastestLogin(form.MinDis, form.MaxDis); nil != err {
		beego.Error("find:", err)
		info.Data["json"] = NewErrorInfo(code, err.Error())
		return
	} else {

		retInfo := SearchGpsRsp{
			Total: int32(count),
			Data:  records,
		}

		if b, err := json.Marshal(retInfo); nil == err {
			info.Data["json"] = NewNormalInfo(b)
		} else {
			beego.Error("ErrUnknown:", err)
			info.Data["json"] = NewErrorInfo(ErrCodeUnknown, ErrUnknown)
			return
		}
	}

	return
}

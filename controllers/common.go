package controllers

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	secret = "test"
)
const (
	ErrInputData    = "数据输入错误"
	ErrDatabase     = "数据库操作错误"
	ErrDupUser      = "用户信息已存在"
	ErrNoUser       = "用户信息不存在"
	ErrPass         = "密码不正确"
	ErrNoUserPass   = "用户信息不存在或密码不正确"
	ErrNoUserChange = "用户信息不存在或数据未改变"
	ErrInvalidUser  = "用户信息不正确"
	ErrOpenFile     = "打开文件出错"
	ErrWriteFile    = "写文件出错"
	ErrSystem       = "操作系统错误"
	ErrExpired      = "验证token过期"
	ErrUnknown      = "未知错误"
)

const (
	ErrCodeInputData    = 1000
	ErrCodeDatabase     = 1001
	ErrCodeDupUser      = 1002
	ErrCodeNoUser       = 1003
	ErrCodePass         = 1004
	ErrCodeNoUserPass   = 1005
	ErrCodeNoUserChange = 1006
	ErrCodeInvalidUser  = 1007
	ErrCodeOpenFile     = 1008
	ErrCodeWriteFile    = 1009
	ErrCodeSystem       = 1010
	ErrCodeExpired      = 1011
	ErrCodeUnknown      = -1
)

type Claims struct {
	//Appid string `json:"Appid"`
	// recommended having
	Userid int `json:"Userid"`
	jwt.StandardClaims
}

func Base64Encode(src []byte) []byte {
	return []byte(base64.StdEncoding.EncodeToString(src))
}
func To_md5(encode string) (decode string) {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(encode))
	cipherStr := md5Ctx.Sum(nil)
	return string(Base64Encode(cipherStr))
}

func Create_token(id int) string {
	expireToken := time.Now().Add(time.Hour * 24).Unix()
	claims := Claims{
		//info.Appid,
		id,
		jwt.StandardClaims{
			ExpiresAt: expireToken,
		},
	}

	// Create the token using your claims
	c_token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Signs the token with a secret.
	signedToken, _ := c_token.SignedString([]byte(secret))

	return signedToken
}

func Token_auth(signedToken string) (int, error) {
	token, err := jwt.ParseWithClaims(signedToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		//fmt.Printf("%v %v", claims.Username, claims.StandardClaims.ExpiresAt)
		//fmt.Println(reflect.TypeOf(claims.StandardClaims.ExpiresAt))
		//return claims.Appid, err
		return claims.Userid, err
	}
	return 0, err
}

// CodeInfo definiton.
type CodeInfo struct {
	Code int             `json:"Code"`
	Info string          `json:"Info"`
	Body json.RawMessage `json:"body"`
}

// NewErrorInfo return a CodeInfo represents error.
func NewErrorInfo(code int, info string) *CodeInfo {
	var msg []byte
	return &CodeInfo{code, info, msg}
}

// NewNormalInfo return a CodeInfo represents OK.
func NewNormalInfo(msg []byte) *CodeInfo {
	return &CodeInfo{0, "ok", msg}
}

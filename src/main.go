package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// 키 발급 만료 제한시간 (unixtimestamp)
const G_EXPIRE_LIMIT = 1

// 액세스 토큰 구조체
type accessToken struct {
	tk1     string
	tk2     string
	expired int64
}

// 인증키 구조체
type authKeyResponse struct {
	rtu_id       string `json:"rtu_id"`
	access_token string `json:"access_token"`
	now          int64  `json:"now"`
	expired      int64  `json:"expired"`
}

// 요청 데이터 구조체
type dataResponse struct {
	cid   string `json:"cid"`
	multi int    `json:"multi"`
	data  string `json:"data"`
}

// 생성된 액세스 토큰 저장 배열
// var keyArr []accessToken

func main() {
	StartServer()
}

func StartServer() {
	router := gin.Default()
	v1 := router.Group("/v1")
	api := v1.Group("/api")

	api.GET("/auth", responseAuthKey)

	api.GET("/lte/data", responseData)

	router.Run(":8080")
}

func responseAuthKey(c *gin.Context) {

	now_unixtime := time.Now().Unix()
	accessToken := accessToken{"01234567890ABCDEF01234567890ABCD", "EF0123456789ABCDEF0123456789ABCD", now_unixtime + G_EXPIRE_LIMIT}

	authKey := authKeyResponse{}
	authKey.rtu_id = ""
	authKey.access_token = accessToken.tk1 + accessToken.tk2
	authKey.now = now_unixtime
	authKey.expired = accessToken.expired

	c.JSON(http.StatusOK, authKey)
}

func responseData(c *gin.Context) {

	//data := dataResponse{}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Ok"})
}

// func main() {

// http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {

// 	log.Print(req.Header)

// 	w.Header().Set("Content-Type", "application/json")

// 	w.Write([]byte("hihi"))

// })

// http.HandleFunc("/v1/api/auth", func(w http.ResponseWriter, req *http.Request) {

// 	req.Header.Get("CID")
// 	req.Header.Get("Authorization")

// 	log.Print(req.Body)

// 	w.Header().Set("Content-Type", "application/json")

// 	w.Write([]byte("hi"))

// 	//resp := authResponse{}
// 	//resp.rtu_id =

// })

// http.HandleFunc("/v1/api/lte/data", func(w http.ResponseWriter, req *http.Request) {

// 	// REQ HEADER
// 	// tk1 32byte
// 	// tk2 32byte

// 	// RESP PARAM
// 	// cid
// 	// multi
// 	// data

// })

// http.ListenAndServe(":60000", nil)

// // err := http.ListenAndServeTLS(":443", "server.crt", "server.key", nil)
// // if err != nil {
// // 	log.Fatal("ListenAndServe : ", err)
// // }
// }

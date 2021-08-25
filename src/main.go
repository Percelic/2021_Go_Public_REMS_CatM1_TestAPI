package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// 키 발급 만료까지 제한시간 (unixtime , second )
const G_EXPIRE_LIMIT = 3 * 60 * 60

// 액세스 토큰 구조체
type accessToken struct {
	tk1     string // 64 byte Access_token에서 앞부분 32 byte
	tk2     string // 64 byte Access_token에서 뒷부분 32 byte
	expired int64  // 만료 시각 (unixtime)
}

// 인증키 구조체
type authKeyResponse struct {
	Rtu_id       string `json:"rtu_id"`
	Access_token string `json:"access_token"`
	Now          int64  `json:"now"`
	Expired      int64  `json:"expired"`
}

// 요청 데이터 구조체
type dataResponse struct {
	Cid   string `json:"cid"`
	Multi int    `json:"multi"`
	Data  string `json:"data"`
}

// 생성된 액세스 토큰 저장 배열
var stored_tokens []accessToken

func main() {
	StartServer()
}

func StartServer() {
	router := gin.Default()
	v1 := router.Group("/v1")
	api := v1.Group("/api")

	api.GET("/auth", responseAuthKey)

	api.POST("/lte/data", responseData)

	router.Run(":8080")
	// HTTPS 실행 시.. ( .crt , .key 인증서 파일이 main.go 경로에 같이 있어야 됨 )
	// router.RunTLS(":8080","server.crt","server.key")
}

func responseAuthKey(c *gin.Context) {

	now_unixtime := time.Now().Unix()

	token := generateToken()
	// accessToken := accessToken{tk1: "1234", tk2: "5678", expired: 10230402045}
	accessToken := accessToken{tk1: token.tk1, tk2: token.tk2, expired: token.expired}

	var authKey authKeyResponse
	authKey.Rtu_id = "00000001"
	authKey.Access_token = accessToken.tk1 + accessToken.tk2
	authKey.Now = now_unixtime
	authKey.Expired = accessToken.expired

	log.Print(c.Request.Header)

	log.Print("Cid : ", c.Request.Header.Get("Cid"))
	log.Print("Authorization : ", c.Request.Header.Get("Authorization"))

	c.JSON(http.StatusOK, authKey)
}

func responseData(c *gin.Context) {

	// 입력 header 확인 프로세스

	log.Print(c.Request.Header)
	log.Print(c.Request.Body)

	log.Print("tk1 : ", c.Request.Header.Get("tk1"))
	log.Print("tk2 : ", c.Request.Header.Get("tk2"))

	//c.IndentedJSON(http.StatusOK, data)
}

// 액세스 토큰 객체 생성
func generateToken() accessToken {

	tk1 := make([]byte, 32)
	tk2 := make([]byte, 32)

	var token accessToken

	for {

		// 0 ~ 9  : 0x30 ~ 0x39
		// A ~ F  : 0x41 ~ 0x46
		for i := 0; i < len(tk1); i++ {
			// 0 ~ 15 사이 무작위 난수 입력
			tk1[i] = byte(rand.Intn(16))
			tk2[i] = byte(rand.Intn(16))

			// HEX to ASCII 변환
			if tk1[i] >= 10 {
				tk1[i] = tk1[i] + 0x37
			} else {
				tk1[i] = tk1[i] + 0x30
			}

			if tk2[i] >= 10 {
				tk2[i] = tk2[i] + 0x37
			} else {
				tk2[i] = tk2[i] + 0x30
			}
		}

		token = accessToken{string(tk1), string(tk2), time.Now().Unix() + G_EXPIRE_LIMIT}

		if checkRepeat(token) {
			break
		}
	}

	stored_tokens = append(stored_tokens, token)

	return token
}

// 생성 토큰 중복 확인
func checkRepeat(token accessToken) bool {

	for _, stored := range stored_tokens {
		if stored.tk1 == token.tk1 {
			return false
		}
	}

	return true
}

package decrypt_test


import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"github.com/google/uuid"
	"sync"
	"crypto/rsa"
	"crypto/rand"
	"encoding/pem"
	"crypto/x509"
	"fmt"
	"encoding/base64"
)

type Test struct {
	count      uint
	privateKey *rsa.PrivateKey
}

type Tests struct {
	mut sync.RWMutex
	m map[string]*Test

}

var tests = Tests {
	m: make(map[string]*Test),
}


func keyToString(key rsa.PublicKey) string {
	pubDer, err := x509.MarshalPKIXPublicKey(&key)
	if err != nil {
		fmt.Println("Failed to get der format for PublicKey.", err)
		return ""
	}

	pubBlk := pem.Block {
		Type:    "PUBLIC KEY",
		Headers: nil,
		Bytes:   pubDer,
	}
	pubPem := string(pem.EncodeToMemory(&pubBlk))

	return pubPem
}

func startTest(c *gin.Context) {
	defer tests.mut.Unlock()
	tests.mut.Lock()

	newTestUid := uuid.New()

	key, _ := rsa.GenerateKey(rand.Reader, 2048)
	err := key.Validate()
	if err != nil {
		fmt.Println("Validation failed.", err)
	}


	tests.m[newTestUid.String()] = &Test{
		count: 0,
		privateKey: key,
	}

	c.JSON(200, gin.H{
		"test_id": newTestUid.String(),
		"public_key": keyToString(key.PublicKey),
	})
}

func endTest(c *gin.Context) {

}

func getData(c *gin.Context) {
	tests.mut.RLock()
	defer tests.mut.RUnlock()

	uuid := c.Param("uuid")

	c.JSON(200, gin.H{
		"test_id": uuid,
		"string": randStringRunes(100),
	})
}

type VerifyData struct {
	Encrypted string `json:"encrypted"`
	Original string `json:"original"`
}

func doVerify(c *gin.Context) {
	tests.mut.Lock()
	defer tests.mut.Unlock()

	uuid := c.Param("uuid")

	var verifyData VerifyData
	c.BindJSON(&verifyData)

	curTest := tests.m[uuid]
	fmt.Println(uuid)
	fmt.Println(curTest.count)

	decoded, _ := base64.StdEncoding.DecodeString(verifyData.Encrypted) // hex.DecodeString(verifyData.Encrypted)
	decrypted, err := rsa.DecryptPKCS1v15(rand.Reader, curTest.privateKey, decoded)
	if err != nil {
		fmt.Println("cant decrypt", err)
		c.Err()
		return
	}

	if verifyData.Original == string(decrypted) {
		c.Status(200)
		curTest.count++
		return
	} else {

	}







}

func GetRoutes() *gin.Engine {
	r := gin.New()

	config := cors.DefaultConfig()

	config.AllowOriginFunc = func(origin string) bool {
		return true
	}

	r.Use(cors.New(config))


	test := r.Group("/test")
	{
		test.GET("/start", startTest)
		test.GET("/end", endTest)
		test.GET("/data/:uuid", getData)
		test.POST("/verify/:uuid", doVerify)
	}


	return r
}
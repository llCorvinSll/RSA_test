package decrypt_test


import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"github.com/google/uuid"
	"sync"
	"crypto/rsa"
	"crypto/rand"
	"encoding/pem"
	"bytes"
	"encoding/asn1"
)

type Test struct {
	count      uint
	privateKey *rsa.PrivateKey
}

type Tests struct {
	mut sync.RWMutex
	m map[string]Test

}

var tests = Tests {
	m: make(map[string]Test),
}


func keyToString(key interface{}) string {
	asn1Bytes, _ := asn1.Marshal(key)

	var pemkey = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: asn1Bytes,
	}

	buf := new(bytes.Buffer)

	_ = pem.Encode(buf, pemkey)

	return  buf.String()
}

func startTest(c *gin.Context) {
	defer tests.mut.Unlock()
	tests.mut.Lock()

	newTestUid := uuid.New()

	key, _ := rsa.GenerateKey(rand.Reader, 2048)


	tests.m[newTestUid.String()] = Test{
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

func doVerify(c *gin.Context) {

}

func GetRoutes() *gin.Engine {
	r := gin.New()

	config := cors.DefaultConfig()
	// config.AllowOrigins = []string{"http://localhost", "http://localhost:88" }

	config.AllowOriginFunc = func(origin string) bool {
		return true
	}
	// config.AllowOrigins == []string{"http://google.com", "http://facebook.com"}

	r.Use(cors.New(config))


	test := r.Group("/test")
	{
		test.GET("/start", startTest)
		test.GET("/end", endTest)
		test.GET("/data/:uuid", getData)
		test.GET("/verify", doVerify)
	}


	return r
}
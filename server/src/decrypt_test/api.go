package decrypt_test


import (
	"github.com/gin-gonic/gin"
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
	mut sync.Mutex
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
		"public_key": keyToString(key.PublicKey),
		"test_id": newTestUid.String(),
	})
}

func endTest(c *gin.Context) {

}

func getData(c *gin.Context) {

}

func doVerify(c *gin.Context) {

}

func GetRoutes() *gin.Engine {
	r := gin.New()

	test := r.Group("/test")
	{
		test.GET("/start", startTest)
		test.GET("/end", endTest)
		test.GET("/data", getData)
		test.GET("/verify", doVerify)
	}

	return r
}
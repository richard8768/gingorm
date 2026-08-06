package util

import (
	"crypto"
	"encoding/hex"
	"errors"
	"fmt"
	"gin_demo/internal/config"
	MRand "math/rand"
	"net/http"
	"os"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func HttpResponse(context *gin.Context, code int, message any, data any) {
	context.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": message,
		"data":    data,
	})
}

func CheckReqBindJson(context *gin.Context, obj any) any {
	if err := context.ShouldBindJSON(obj); err != nil {
		//errs, ok := err.(validator.ValidationErrors)
		var errs validator.ValidationErrors
		ok := errors.As(err, &errs)
		if !ok {
			return err.Error()
		}
		return RemoveTopStruct(errs.Translate(Trans))
	}
	return nil
}

func CheckReqBindQuery(context *gin.Context, obj any) any {
	if err := context.ShouldBindQuery(obj); err != nil {
		//errs, ok := err.(validator.ValidationErrors)
		var errs validator.ValidationErrors
		ok := errors.As(err, &errs)
		if !ok {
			return err.Error()
		}
		return RemoveTopStruct(errs.Translate(Trans))
	}
	return nil
}

func CheckReqBindHeader(context *gin.Context, obj any) any {
	if err := context.ShouldBindHeader(obj); err != nil {
		//errs, ok := err.(validator.ValidationErrors)
		var errs validator.ValidationErrors
		ok := errors.As(err, &errs)
		if !ok {
			return err.Error()
		}
		return RemoveTopStruct(errs.Translate(Trans))
	}
	return nil
}
func CheckReqBind(context *gin.Context, obj any) any {
	if err := context.ShouldBind(obj); err != nil {
		//errs, ok := err.(validator.ValidationErrors)
		var errs validator.ValidationErrors
		ok := errors.As(err, &errs)
		if !ok {
			return err.Error()
		}
		return RemoveTopStruct(errs.Translate(Trans))
	}
	return nil
}

func TimeoutMiddleware(duration time.Duration) gin.HandlerFunc {
	return timeout.New(
		timeout.WithTimeout(duration),
		timeout.WithResponse(TimeoutResponse),
	)
}

func TimeoutResponse(c *gin.Context) {
	c.JSON(http.StatusGatewayTimeout, gin.H{
		"code":    http.StatusGatewayTimeout,
		"message": "Gateway Timeout",
		"data":    nil,
	})
}
func GetUserId(context *gin.Context) (uint, error) {
	var boolean bool
	var userIds any
	userId, err := config.RedisClient.Get("UserId").Result()
	if err != nil {
		userIds, boolean = context.Get("UserId")
		if !boolean {
			return 0, errors.New("意外的错误")
		}
		userId = userIds.(string)
	}
	num, err := strconv.ParseUint(userId, 10, 64)
	if err != nil {
		return 0, errors.New("意外的类型错误")
	}
	return uint(num), nil
}

var num int64

const (
	Continuity = "20060102150405"
)

// Generate 生成24位订单号
// 前面17位代表时间精确到毫秒，中间3位代表进程id，最后4位代表序号
func GenerateOrderNo(t time.Time) string {
	s := t.Format(Continuity)
	m := t.UnixNano()/1e6 - t.UnixNano()/1e9*1e3
	ms := sup(m, 3)
	p := os.Getpid() % 1000
	ps := sup(int64(p), 3)
	i := atomic.AddInt64(&num, 1)
	r := i % 10000
	rs := sup(r, 4)
	n := fmt.Sprintf("%s%s%s%s", s, ms, ps, rs)
	return n
}

// 对长度不足n的数字前面补0
func sup(i int64, n int) string {
	m := fmt.Sprintf("%d", i)
	for len(m) < n {
		m = fmt.Sprintf("0%s", m)
	}
	return m
}

func GenRandStrings(r *MRand.Rand, n int, randtype string) string {
	var str string
	if randtype == "number" {
		str = "0123456789"
	} else if randtype == "password" {
		str = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ~!@#$%^&*()-+_=,."
	} else {
		str = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	}
	bytes := []byte(str)
	var result []byte
	//r := MRand.New(MRand.NewSource(time.Now().UnixNano()))
	lenth := len(bytes)
	for i := 0; i < n; i++ {
		result = append(result, bytes[r.Intn(lenth)])
	}
	return string(result)
}
func GenMd5Signature(str string) string {
	omd5 := crypto.MD5.New()
	omd5.Write([]byte(str))
	md5str := hex.EncodeToString(omd5.Sum(nil))
	return md5str
}
func GenSha256Signature(str string) string {
	osha256 := crypto.SHA256.New()
	osha256.Write([]byte(str))
	re := osha256.Sum(nil)
	sha256str := hex.EncodeToString(re)
	return sha256str
}

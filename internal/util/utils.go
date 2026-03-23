package util

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"gin_demo/internal/config"
	"gin_demo/internal/dto"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-module/carbon"
)

func HttpResponse(context *gin.Context, code int, message any, data any) {
	context.JSON(http.StatusOK, &dto.HttpResponse{
		Code:    code,
		Message: message,
		Data:    data,
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

func GenFileName() string {
	dateStr := carbon.Now().Format("YmdHis")
	data := []byte(dateStr)
	h := md5.New()
	h.Write(data)
	sum := h.Sum(nil)
	md5Str := dateStr + "_" + hex.EncodeToString(sum)
	return md5Str
}

func GetImageFileType(fileType string) string {
	fileType = strings.ToLower(fileType)
	switch fileType {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/gif":
		return ".gif"
	}

	return ""
}

var num int64

const (
	Normal     = "2006-01-02 15:04:05"
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

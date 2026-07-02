package util

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"gin_demo/internal/config"
	"gin_demo/internal/dto"
	MRand "math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-module/carbon"
	"github.com/xuri/excelize/v2"
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
	//https://blog.csdn.net/qq_26086231/article/details/135589839
	switch fileType {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/gif":
		return ".gif"
	case "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":
		return ".xlsx"
	case "application/x-7z-compressed":
		return ".7z"
	case "application/x-rar-compressed":
		return ".rar"
	case "application/zip":
		return ".zip"
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

func ExportToExcel(context *gin.Context, titleList []string, data [][]string, fileName string) {
	// 生成一个新的文件
	file := excelize.NewFile()
	defer file.Close()
	// 添加sheet页
	sheetName := "Sheet1"

	file.SetSheetName("Sheet1", sheetName)
	sheetID, _ := file.GetSheetIndex(sheetName)
	file.SetActiveSheet(sheetID)
	currentSheet := file.GetSheetName(sheetID)
	// 插入表头
	for colIdx, header := range titleList {
		cell, _ := excelize.CoordinatesToCellName(colIdx+1, 1)
		file.SetCellValue(currentSheet, cell, header)
	}
	// 插入内容
	for rowIdx, row := range data {
		for colIdx, value := range row {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, rowIdx+2)
			file.SetCellValue(currentSheet, cell, value)
		}
	}
	// 设置 HTTP 响应的头信息
	context.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	context.Header("Content-Disposition", "attachment; filename="+fileName)
	// 将 Excel 文件写入 HTTP 响应
	if err := file.Write(context.Writer); err != nil {
		HttpResponse(context, 500, "failed", nil)
		return
	}
}

func CheckReqBind(context *gin.Context, obj any) any {
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

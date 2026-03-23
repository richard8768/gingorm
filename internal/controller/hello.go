package controller

import (
	contexts "context"
	"fmt"
	"gin_demo/internal/config"
	"gin_demo/internal/dto"
	"gin_demo/internal/util"
	"mime/multipart"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/dustin/go-humanize"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/golang-module/carbon"
)

// HelloWorld test
// @Summary HelloWorld
// @Schemes
// @Description HelloWorld
// @Tags HelloWorldHandler
// @Accept json
// @Produce json
// @Success 200 {string} hello world!
// @Router / [get]
func HelloWorldHandler(context *gin.Context) {
	util.HttpResponse(context, 200, "success", "hello world!")
}

// Get test
// @Summary Get tests
// @Schemes
// @Description Get testss
// @Tags GetHandler
// @Accept json
// @Produce json
// @Success 200 {string} 15
// @Router /get [get]
func GetHandler(context *gin.Context) {
	fmt.Println(config.AppConfigs)
	times := carbon.Now().TimestampMilli()
	fmt.Println("times:", times)
	id := "15"
	util.HttpResponse(context, 200, "success", id)
}

// upload single file
// @Summary upload single file
// @Schemes
// @Description upload single file
// @Tags UploadHandler
// @Accept json
// @Produce json
// @Param body body dto.FileRequest true "请求body"
// @Success 200 {object} dto.FileResponse
// @Router /upload [post]
func UploadHandler(context *gin.Context) {
	req := dto.FileRequest{}
	if err := context.ShouldBind(&req); err != nil {
		util.HttpResponse(context, 500, err.Error(), nil)
		return
	}

	//file, err := context.FormFile("file")
	//if err != nil {
	//	util.HttpResponse(context, 500, err.Error(), nil)
	//	return
	//}

	form, err := context.MultipartForm()
	if err != nil {
		util.HttpResponse(context, 500, err.Error(), nil)
		return
	}
	files := form.File["file"]
	if files == nil {
		util.HttpResponse(context, 500, "file is empty", nil)
		return
	}
	file := files[0]

	fileName := util.GenFileName()
	fileType := util.GetImageFileType(file.Header["Content-Type"][0])
	saveFilePath := config.GetLocalUploadPath() + fileName + fileType
	context.SaveUploadedFile(file, saveFilePath)

	var rsp dto.FileResponse
	rsp = dto.FileResponse{FileName: saveFilePath}
	util.HttpResponse(context, 200, "success", rsp)
}

func AliOssHandler(context *gin.Context) {
	localFileName := "20260319213817_5f3082c299b92ea01cc04da7bf7708f7.jpg"
	localFilepath := config.GetLocalUploadPath() + localFileName
	fmt.Println(localFilepath)

	securityToken := config.AppConfigs.AliOss.SecurityToken
	endpoint := config.AppConfigs.AliOss.EndPoint
	if len(securityToken) == 0 {
		securityToken = ""
	}
	var ossCfg *oss.Config
	if len(endpoint) > 0 {
		ossCfg = oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewStaticCredentialsProvider(config.AppConfigs.AliOss.AccessKeyId, config.AppConfigs.AliOss.AccessKeySecret, securityToken)).
			WithRegion(config.AppConfigs.AliOss.RegionId).
			WithEndpoint(endpoint).
			WithUseCName(true).
			WithConnectTimeout(10 * time.Second).
			WithReadWriteTimeout(30 * time.Second).
			WithRetryMaxAttempts(5).
			WithDisableSSL(true) // 设置不使用HTTPS协议。默认使用HTTPS
	} else {
		ossCfg = oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewStaticCredentialsProvider(config.AppConfigs.AliOss.AccessKeyId, config.AppConfigs.AliOss.AccessKeySecret, securityToken)).
			WithRegion(config.AppConfigs.AliOss.RegionId).
			WithConnectTimeout(10 * time.Second).
			WithReadWriteTimeout(30 * time.Second).
			WithRetryMaxAttempts(5).
			WithDisableSSL(true) // 设置不使用HTTPS协议。默认使用HTTPS
	}

	client := oss.NewClient(ossCfg)

	putRequest := &oss.PutObjectRequest{
		Bucket: oss.Ptr(config.AppConfigs.AliOss.BucketName),
		Key:    oss.Ptr(localFileName),
	}
	result, err := client.PutObjectFromFile(contexts.TODO(), putRequest, localFilepath)
	if err != nil {
		util.HttpResponse(context, 500, fmt.Sprintf("failed to put object from file %v", err), nil)
	}
	fmt.Printf("Status: %#v\n", result.Status)
	fmt.Printf("RequestId: %#v\n", result.ResultCommon.Headers.Get("X-Oss-Request-Id"))
	fmt.Printf("ETag: %#v\n", *result.ETag)

	util.HttpResponse(context, 200, "success", result)
}

func init() {
	validate, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		validate.RegisterStructValidation(FileUploadValidation, dto.FileRequest{})
	}
}

func FileUploadValidation(sl validator.StructLevel) {
	form := sl.Current().Interface().(dto.FileRequest)
	formType := reflect.TypeOf(form)

	for i := 0; i < formType.NumField(); i++ {
		field := formType.Field(i)
		if field.Type != reflect.TypeOf(&multipart.FileHeader{}) {
			continue
		}

		fileTag := field.Tag.Get("form")
		if fileTag == "" {
			continue
		}

		fileField := reflect.ValueOf(form).Field(i)
		file := fileField.Interface().(*multipart.FileHeader)
		if file == nil {
			continue
		}
		if file.Size != 0 {
			fileSizeTag := field.Tag.Get("fileSize")
			if !CheckFileSize(file, fileSizeTag) {
				sl.ReportError(file, field.Name, "fileSize", "CheckFileSize", fmt.Sprintf("请上传%s大小内的文件", fileSizeTag))
			}

			fileSuffixTag := field.Tag.Get("fileSuffix")
			if !CheckFileSuffix(file, fileSuffixTag) {
				sl.ReportError(file, field.Name, "fileSuffix", "CheckFileSuffix", fmt.Sprintf("请上传扩展名为%s的文件", fileSuffixTag))
			}
		}

	}
}
func CheckFileSize(file *multipart.FileHeader, fileSize string) bool {
	maxSize, err := humanize.ParseBytes(fileSize)
	if err != nil {
		maxSize = 10 * 1024 * 1024 // 10M
	}

	return uint64(file.Size) <= maxSize
}

func CheckFileSuffix(file *multipart.FileHeader, allowedSuffixes string) bool {
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowed := strings.Split(allowedSuffixes, "|")
	for _, suffix := range allowed {
		if ext == "."+suffix {
			return true
		}
	}

	return false
}

package controller

import (
	contexts "context"
	"fmt"
	"gin_demo/internal/config"
	"gin_demo/internal/util"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/gin-gonic/gin"
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

func AliOssHandler(context *gin.Context) {
	localFileName := "20260319213817_5f3082c299b92ea01cc04da7bf7708f7.jpg"
	localFilepath, _ := config.GetLocalUploadPath()
	localFilepath = localFilepath + localFileName
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

package config

import (
	"errors"
	"fmt"
	"slices"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type AllConfig struct {
	*Server
	*DataBase
	*Redis
	*Log
	*Jwt
	*AliOss
}

type Server struct {
	Mode            string `mapstructure:"mode"`
	Port            int    `mapstructure:"port"`
	Name            string `mapstructure:"name"`
	Version         string `mapstructure:"version"`
	Level           string `mapstructure:"level"`
	LocalUploadPath string `mapstructure:"local_upload_path"`
}

type DataBase struct {
	Driver       string
	Host         string
	Port         string
	UserName     string
	Password     string
	DBName       string `mapstructure:"db_name"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	MaxLifeTime  int    `mapstructure:"max_life_time"`
	LogLevel     int    `mapstructure:"log_level"`
	Prefix       string
	Config       string
}

type Redis struct {
	Host     string
	Port     string
	Password string
	DataBase int `mapstructure:"data_base"`
}

type Log struct {
	Level    string
	FilePath string
}

type JwtBaseOption struct {
	Secret string `mapstructure:"secret"`
	TTL    int    `mapstructure:"ttl"`
}

type Jwt struct {
	UserList []string      `mapstructure:"userlist"`
	Admin    JwtBaseOption `mapstructure:"admin"`
	User     JwtBaseOption `mapstructure:"user"`
}

type AliOss struct {
	RegionId        string `mapstructure:"region_id"`
	EndPoint        string `mapstructure:"endpoint"`
	AccessKeyId     string `mapstructure:"access_key_id"`
	AccessKeySecret string `mapstructure:"access_key_secret"`
	SecurityToken   string `mapstructure:"security_token"`
	BucketName      string `mapstructure:"bucket_name"`
}

var envPtr = pflag.String("env", "dev", "Environment: dev test or prod")

var AppConfigs *AllConfig // 全局Config

func InitLoadConfig() (*AllConfig, error) {
	// 使用pflag库来读取命令行参数，用于指定环境，默认为"dev"
	pflag.Parse()

	config := viper.New()
	// 设置读取路径
	config.AddConfigPath("./config")
	// 设置读取文件名字
	config.SetConfigName(fmt.Sprintf("application-%s", *envPtr))
	// 设置读取文件类型
	config.SetConfigType("yaml")
	// 读取文件载体
	var configData *AllConfig
	// 读取配置文件
	err := config.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Use Viper ReadInConfig Fatal error config err:%s \n", err))
	}
	//fmt.Println(config.AllKeys())
	//fmt.Println(config.AllSettings())
	//fmt.Println(config.Get("database.prefix"))
	//fmt.Println(config.Get("jwt.user.secret"), config.Get("jwt.user.ttl"))
	//fmt.Println(&configData)
	// 查找对应配置文件
	err = config.Unmarshal(&configData)
	if err != nil {
		panic(fmt.Errorf("read config file to struct err: %s\n", err))
	}
	AppConfigs = configData
	// 打印配置文件信息
	//fmt.Printf("配置文件信息：%+v", configData)
	return configData, nil
}

func GetHttpPort() int {
	return AppConfigs.Server.Port
}

func GetRunMode() string {
	mode := AppConfigs.Server.Level
	definedMode := [3]string{"debug", "release", "test"}
	posIndex := -1
	length := len(definedMode)
	for index := 0; index < length; index++ {
		if mode == definedMode[index] {
			posIndex = index
		}
	}
	if posIndex == -1 {
		mode = "debug"
	}
	return mode
}

func GetJwtCfg(userType string) (*JwtBaseOption, error) {
	userTypeList := AppConfigs.Jwt.UserList
	index := slices.Index(userTypeList, userType)
	if index == -1 {
		return nil, errors.New("意外的用户类型")
	}
	if userType == "admin" {
		return &AppConfigs.Jwt.Admin, nil
	} else {
		return &AppConfigs.Jwt.User, nil
	}
}

func GetLogPath() string {
	return AppConfigs.Log.FilePath
}

func GetLocalUploadPath() string {
	return AppConfigs.Server.LocalUploadPath
}

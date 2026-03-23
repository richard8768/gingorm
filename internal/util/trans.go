package util

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
)

var Trans ut.Translator

// InitTrans 修改 Gin 表单验证的翻译器
func InitTrans(locale string) (err error) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 注册自定义验证器
		RegisterCustomValidators()

		// 注册一个获取 json 的 tag 的自定义方法
		v.RegisterTagNameFunc(func(field reflect.StructField) string {
			name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
		zhT := zh.New()
		enT := en.New()
		// 第一个参数是备用，后面的参数是默认的
		uni := ut.New(enT, zhT, zhT)
		trans, ok := uni.GetTranslator(locale)
		Trans = trans
		if !ok {
			return fmt.Errorf("uni.GetTranslator %v", locale)
		}
		switch locale {
		case "zh":
			err := zhTranslations.RegisterDefaultTranslations(v, Trans)
			if err != nil {
				return err
			}
		default:
			err := enTranslations.RegisterDefaultTranslations(v, Trans)
			if err != nil {
				return err
			}
		}
		return
	}
	return
}
func RemoveTopStruct(fields map[string]string) map[string]string {
	rsp := map[string]string{}
	for field, err := range fields {
		key := field[:strings.Index(field, ".")+1]
		finalKey := strings.Replace(field, key, "", 1)
		finalErr := strings.ReplaceAll(err, key, "")
		rsp[finalKey] = finalErr
		//rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}

package util

import (
	"fmt"
	"gin_demo/internal/dto"
	"mime/multipart"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/go-playground/validator/v10"
)

func FileUploadValidation(sl validator.StructLevel) {
	form := sl.Current().Interface().(dto.UserSingleFileUploadRequest)
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

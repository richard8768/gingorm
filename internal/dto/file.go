package dto

import "mime/multipart"

type UserAvatarUploadRequest struct {
	//File *multipart.FileHeader `form:"file" binding:"required"`
	File *multipart.FileHeader `form:"file" binding:"required" fileSize:"1M" fileSuffix:"jpg|png|gif" msg:"请上传1M大小内的文件"`
}
type UserSingleFileUploadRequest struct {
	//File *multipart.FileHeader `form:"file" binding:"required"`
	File *multipart.FileHeader `form:"file" binding:"required" fileSize:"3M" fileSuffix:"jpg|png|gif|zip|rar|7z" msg:"请上传3M大小内的文件"`
}
type UserSingleFileUploadResponse struct {
	FileName string `json:"file_name"`
	FilePath string `json:"file_path"`
}

type UserSingleFileDownloadRequest struct {
	ID uint `form:"id"    binding:"required,number,gt=0"`
}
type UserAddressUploadRequest struct {
	File *multipart.FileHeader `form:"file" binding:"required" fileSize:"5" fileSuffix:"xlsx" msg:"请上传5M大小内的文件"`
}

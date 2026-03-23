package dto

import "mime/multipart"

type FileRequest struct {
	//File *multipart.FileHeader `form:"file" binding:"required"`
	File *multipart.FileHeader `form:"file" binding:"required" fileSize:"1M" fileSuffix:"jpg|png|gif" msg:"请上传1M大小内的图片"`
}
type FileResponse struct {
	FileName string `json:"file_path"`
}

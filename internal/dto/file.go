package dto

import "mime/multipart"

type UserAvatarUploadRequest struct {
	//File *multipart.FileHeader `form:"file" binding:"required"`
	File *multipart.FileHeader `form:"file" binding:"required" fileSize:"1M" fileSuffix:"jpg|png|gif" msg:"请上传1M大小内的文件"`
}
type UserNormalFileUploadRequest struct {
	//File *multipart.FileHeader `form:"file" binding:"required"`
	File *multipart.FileHeader `form:"file" binding:"required" fileSize:"3M" fileSuffix:"jpg|png|gif|zip|rar|7z" msg:"请上传3M大小内的文件"`
}
type UserNormalFileUploadResponse struct {
	FileName string `json:"file_name"`
	FilePath string `json:"file_path"`
}

type UserNormalFileDownloadRequest struct {
	ID uint `form:"id"    binding:"required,number,gt=0"`
}
type UserAddressUploadRequest struct {
	File *multipart.FileHeader `form:"file" binding:"required" fileSize:"4" fileSuffix:"xlsx" msg:"请上传4M大小内的文件"`
}

type UserLargeFileUploadInitRequest struct {
	FileName string `form:"file_name" binding:"required"`
}
type UserLargeFileUploadInitResponse struct {
	UploadId     string `json:"upload_id" `
	Signature    string `json:"signature" `
	UploadStatus string `json:"upload_status"`
}
type UserLargeFileUploadRequest struct {
	UploadId    string `header:"upload_id" binding:"required"`
	ChunkIndex  int    `header:"chunk_index" binding:"required,number,gt=0,ltefield=TotalChunks"`
	ChunkMd5    string `header:"chunk_md5" binding:"required"`
	TotalChunks int    `header:"total_chunks" binding:"required,number,gt=0"`
}

type UserLargeFileUploadResponse struct {
	UploadId          string `json:"upload_id" `
	ChunkIndex        int    `json:"chunk_index" `
	ChunkMd5          string `json:"chunk_md5" `
	ChunkUploadStatus string `json:"chunk_upload_status"`
}

type UserChunkUploadIdRequest struct {
	UploadId string `form:"upload_id" binding:"required"`
}

// 文件上传状态记录
type ChunkUploadInfo struct {
	ChunkIndex int  `json:"chunk_index"`
	Uploaded   bool `json:"uploaded"`
}

type UserChunkUploadList struct {
	UploadId         string             `json:"upload_id" `
	UploadChunksList []*ChunkUploadInfo `json:"upload_chunk_list"`
}

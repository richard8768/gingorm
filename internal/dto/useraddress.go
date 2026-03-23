package dto

type UserAddressCreateRequest struct {
	Address       string `json:"address" binding:"required,min=10,max=200"`
	Tel           string `json:"tel" binding:"required,mobile"`
	ConsigneeName string `json:"consignee_name"    binding:"required,min=2,max=100"`
	Post          string `json:"post"   binding:"required,number"`
	ProvinceID    int    `json:"province_id" binding:"required,number"`
	CityID        int    `json:"city_id"    binding:"required,number"`
	AreaID        int    `json:"area_id"   binding:"required,number"`
}

type UserAddressRequest struct {
	ID uint `json:"id"    binding:"required,number,gt=0"`
}

type UserAddressGetRequest struct {
	ID uint `form:"id"    binding:"required,number,gt=0"`
}

type UserAddressUpdateRequest struct {
	ID uint `json:"id"    binding:"required,number,gt=0"`
	UserAddressCreateRequest
}

type UserAddressDeleteRequest struct {
	ID []uint `json:"id"    binding:"required,gt=0"`
}

type UserAddressDeleteResponse struct {
	ID           uint  `json:"id"`
	DeleteStatus int64 `json:"delete_status"`
}

type UserAddressSearchRequest struct {
	Keyword  string `form:"keyword" json:"keyword"  binding:"max=20"`
	Page     int    `form:"page" json:"page"  binding:"number,gt=0"`
	PageSize int    `form:"page_size" json:"page_size"  binding:"number,gt=0,oneof=5 10 15 20 25 30 35 40 45 50 100"`
}
type UserAddressResponse struct {
	ID            uint   `json:"id"`
	Address       string `json:"address"`
	Tel           string `json:"tel"`
	ConsigneeName string `json:"consignee_name"`
	Post          string `json:"post"`
	ProvinceID    int    `json:"province_id"`
	CityID        int    `json:"city_id"`
	AreaID        int    `json:"area_id"`
}

type UserAddressListResponse struct {
	Total           int64                 `json:"total"`
	Page            int                   `json:"current_page"`
	PageSize        int                   `json:"page_size"`
	UserAddressList []UserAddressResponse `json:"user_address_list"`
}

package endpoints

type GetUserResponse struct {
	Uid string `json:"uid"`
}

type LoginResponse struct {
	Id          int32  `json:"id"`
	Username    string `json:"username"`
	Avatar      string `json:"avatar"`
	RoleID      int32  `json:"roleID"`
	RecentTime  string `json:"recentTime"`
	CreatedTime string `json:"createdTime"`
}

type SendCodeResponse struct {
	CodeID string `json:"codeID"`
}

type UserListResponse struct {
	Count int             `json:"count"`
	Data  []*UserResponse `json:"data"`
}

type UserResponse struct {
	Id          int32  `json:"id"`
	Username    string `json:"username"`
	Avatar      string `json:"avatar"`
	RoleID      int32  `json:"roleID"`
	RecentTime  string `json:"recentTime"`
	CreatedTime string `json:"createdTime"`
}

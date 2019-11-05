package endpoints

type LoginResponse struct {
	Id          int32  `json:"id"`
	Username    string `json:"username"`
	Avatar      string `json:"avatar"`
	RoleID      int32  `json:"roleID"`
	RecentTime  string `json:"recentTime"`
	CreatedTime string `json:"createdTime"`
}

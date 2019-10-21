package endpoints

type GetUserRequest struct {
	Uid string `json:"uid"`
}

type LoginRequest struct {
	Username string `json:"username" validator:"required||string=[6|10]"`
	Password string `json:"password" validator:"required||string=[6|10]"`
}

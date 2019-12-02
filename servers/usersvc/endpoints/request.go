package endpoints

type GetUserRequest struct {
	Uid string `json:"uid"`
}

type LoginRequest struct {
	Username string `json:"username" validator:"required||string=[6|10]"`
	Password string `json:"password" validator:"required||string=[6|10]"`
}

type RegisterRequest struct {
	Username string `json:"username" validator:"required||string=[6|10]"`
	Password string `json:"password" validator:"required||string=[6|10]"`
	CodeID   int32  `json:"codeID" validator:"required||len=6"`
}

type UserListRequest struct {
	Page int32 `json:"page" validator:"required||number=0|_"`
	Size int32 `json:"size" validator:"required||number=0|_"`
}

type LogoutRequest struct {
	SID string
}

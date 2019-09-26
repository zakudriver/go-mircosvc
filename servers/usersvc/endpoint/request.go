package endpoint

type GetUserRequest struct {
	UID string `json:"s"`
}

type LoginRequest struct {
	Username string `json:"username" validator:"required||string=[6|10]"`
	Password string `json:"password" validator:"required||string=[6|10]"`
}

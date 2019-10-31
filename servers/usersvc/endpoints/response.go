package endpoints

type GetUserResponse struct {
	Name string `json:"name"`
	Err  error  `json:"err"`
}

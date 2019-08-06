package entity


type Order struct {
	Id     string `json:"orderId"`
	Source string `json:"source"`
	IsPay  int    `json:"isPay"`
}

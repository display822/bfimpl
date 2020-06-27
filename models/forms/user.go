package forms

type ReqUser struct {
	Name     string `json:"name" valid:"Required"`
	Email    string `json:"email" valid:"Required"`
	Wx       string `json:"wx" valid:"Required"`
	Phone    string `json:"phone" valid:"Required"`
	UserType int    `json:"useType" valid:"Required"`
}

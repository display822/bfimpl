package forms

type RedisStringReq struct {
	Key   string `form:"key" valid:"Required" json:"key"`
	Value string `form:"value" valid:"Required" json:"value"`
}

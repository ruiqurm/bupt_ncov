package buptncov

type LoginStruct struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type GeneralResponse struct {
	Error   int                    `json:"e"`
	Message string                 `json:"m"`
	Data    map[string]interface{} `json:"d"`
}

const MAX_RETRY_TIMES = 3

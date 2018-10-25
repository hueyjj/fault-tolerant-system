package response

type Response struct {
	Replaced string `json:"replaced,omitempty"`
	Msg      string `json:"msg,omitempty"`
	Value    string `json:"value,omitempty"`
	Result   string `json:"result,omitempty"`
	IsExists string `json:"isExists,omitempty"`
}

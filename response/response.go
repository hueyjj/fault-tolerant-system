package response

type Response struct {
	Replaced bool   `json:"replaced,omitempty"`
	Msg      string `json:"msg,omitempty"`
	Value    int    `json:"value,omitempty"`
	Result   string `json:"result,omitempty"`
	IsExists bool   `json:"isExists,omitempty"`
}

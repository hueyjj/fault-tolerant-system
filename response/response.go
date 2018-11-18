package response

type Response struct {
	Replaced *bool   `json:"replaced,omitempty"`
	Msg      string  `json:"msg,omitempty"`
	Value    string  `json:"value,omitempty"`
	Result   string  `json:"result,omitempty"`
	IsExists *bool   `json:"isExists,omitempty"`
	Payload  Payload `json:"payload,omitempty"`
}

type Payload struct {
	VectorClocks map[string]int    `json:"vectorclocks"`
	LookupIPs    map[string]string `json:"lookupips"`
}

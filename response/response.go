package response

import "bitbucket.org/cmps128gofour/homework4/vectorclock"

type Response struct {
	Replaced *bool   `json:"replaced,omitempty"`
	Msg      string  `json:"msg,omitempty"`
	Owner    *int    `json:"owner,omitempty"`
	Value    string  `json:"value,omitempty"`
	Result   string  `json:"result,omitempty"`
	IsExists *bool   `json:"isExists,omitempty"`
	Payload  Payload `json:"payload,omitempty"`
}

type ViewResponse struct {
	Msg     string         `json:"msg,omitempty"`
	Result  string         `json:"result,omitempty"`
	IPTable map[string]int `json:"iptable,omitempty"`
}

type ShardResponse struct {
	ID       *int   `json:"id,omitempty"`
	Result   string `json:"result,omitempty"`
	ShardIDs string `json:"shard_ids,omitempty"`
	Members  string `json:"members,omitempty"`
	Count    *int   `json:"Count,omitempty"`
	Msg      string `json:"msg,omitempty"`
}

type IPTableResponse struct {
	View string `json:"view",omitempty`
}

type Payload struct {
	VectorClocks map[string]vectorclock.Unit `json:"vectorclocks,omitempty"`
	IPTable      map[string]int              `json:"iptable,omitempty"`
}

type Update struct {
	VectorClocks map[string]vectorclock.Unit `json:"vectorclocks,omitempty"`
	KVS          map[string]string           `json:"kvs,omitempty"`
}

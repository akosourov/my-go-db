package server

type RequestBody struct {
	ValueStr      string            `json:"value_str"`
	ValueInt      int               `json:"value_int"`
	ItemText      string            `json:"item_text"`
	ItemTextArray []string          `json:"item_text_array"`
	ItemTextDict  map[string]string `json:"item_text_dict"`
	TTL           int               `json:"ttl"`
}

type ResponseBody struct {
	Success       bool        `json:"success"`
	Message       string      `json:"message,omitempty"`
	Value         interface{} `json:"value,omitempty"`
	ValueStr      string      `json:"value_str,omitempty"`
	ValueInt      int         `json:"valuet_int,omitempty"`
	ValueSliceInt []int       `json:"value_slice_int,omitempty"`
}

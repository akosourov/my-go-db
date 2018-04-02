package server

type Payload struct {
	ItemText      string            `json:"item_text"`
	ItemTextArray []string          `json:"item_text_array"`
	ItemTextDict  map[string]string `json:"item_text_dict"`
	TTL           int64             `json:"ttl"`
}

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Value   interface{} `json:"value,omitempty"`
}

package server

type RequestBody struct {
	String        string            `json:"string,omitempty"`
	Int           int               `json:"int,omitempty"`
	StringList    []string          `json:"string_list,omitempty"`
	IntList       []int             `json:"int_list,omitempty"`
	StringDict    map[string]string `json:"string_dict,omitempty"`
	IntDict       map[string]int    `json:"int_dict,omitempty"`

	TTL           int               `json:"ttl,omitempty"`
}

type ResponseBody struct {
	Success       bool              `json:"success"`
	Message       string            `json:"message,omitempty"`

	String        string            `json:"string,omitempty"`
	Int           int               `json:"int,omitempty"`
	StringList    []string          `json:"string_list,omitempty"`
	IntList       []int             `json:"int_list,omitempty"`
	StringDict    map[string]string `json:"string_dict,omitempty"`
	IntDict       map[string]int    `json:"int_dict,omitempty"`

	Keys          []string          `json:"keys,omitempty"`
}

package pkg

type UserRequest struct {
	UserId              string `json:"user_id"`
	Ip                  string `json:"ip"`
	Method              string `json:"method"`
	Status              string `json:"status"`
	Error               string `json:"error,omitempty"`
	RequestId           string `json:"request_id"`
	URL                 string `json:"url"`
	ResponseContentType string `json:"response_content_type,omitempty"`
	ResponseStatus      string `json:"response_status,omitempty"`
	BodySize            int    `json:"body_size"`
	InitTime            int64  `json:"init_time"`
	UpdateTime          int64  `json:"update_time,omitempty"`
	ResponseSize        int    `json:"response_size,omitempty"`
} // Re-listed

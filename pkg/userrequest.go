package pkg

import (
	"time"
)

type UserRequest struct {
	UserId              string    `json:"user_id"`
	Ip                  string    `json:"ip"`
	Method              string    `json:"method"`
	Status              string    `json:"status"`
	Error               string    `json:"error,omitempty"`
	RequestId           string    `json:"request_id"`
	URL                 string    `json:"url"`
	ResponseContentType string    `json:"response_content_type,omitempty"`
	ResponseStatus      string    `json:"response_status,omitempty"`
	BodySize            int64     `json:"body_size"`
	InitTime            time.Time `json:"init_time"`             // Unnecessary on influx
	UpdateTime          time.Time `json:"update_time,omitempty"` // Unnecessary on infux
	ResponseSize        int64     `json:"response_size,omitempty"`
} // Re-listed

func (u UserRequest) ToMapI() map[string]interface{} {
	return map[string]interface{}{
		"user_id":               u.UserId,
		"ip":                    u.Ip,
		"method":                u.Method,
		"status":                u.Status,
		"error":                 u.Error,
		"request_id":            u.RequestId,
		"url":                   u.URL,
		"response_content_type": u.ResponseContentType,
		"response_status":       u.ResponseStatus,
		"body_size":             u.BodySize,
		"init_time":             u.InitTime,
		"update_time":           u.UpdateTime,
		"response_size":         u.ResponseSize,
	}
}

func UserRequestFromMapI(m map[string]interface{}) (u UserRequest) {
	// or mapstructure.Decode(resultmap, &u)
	u.UserId, _ = m["user_id"].(string)
	u.Ip, _ = m["ip"].(string)
	u.Method, _ = m["method"].(string)
	u.Status, _ = m["status"].(string)
	u.Error, _ = m["error"].(string)
	u.RequestId, _ = m["request_id"].(string)
	u.URL, _ = m["url"].(string)
	u.ResponseContentType, _ = m["response_content_type"].(string)
	u.ResponseStatus, _ = m["response_status"].(string)
	u.BodySize, _ = m["body_size"].(int64)
	u.InitTime, _ = m["init_time"].(time.Time)
	u.UpdateTime, _ = m["update_time"].(time.Time)
	u.ResponseSize, _ = m["response_size"].(int64)
	// This way is more fast
	return
}

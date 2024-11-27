package perrors

import "strings"

type Response struct {
	Code    Code   `json:"code"`
	Message string `json:"message"`
}

type Code string

func (c Code) Error() string {
	return string(c)
}

func (c Code) MakeJSON(message ...string) map[string]interface{} {
	body := map[string]interface{}{
		"code": string(c),
	}

	if len(message) != 0 {
		body["message"] = strings.Join(message, "")
	}

	return body
}

func (c Code) WithJSON(detail interface{}, message ...string) map[string]interface{} {
	body := map[string]interface{}{
		"code":   string(c),
		"detail": detail,
	}

	if len(message) != 0 {
		body["message"] = strings.Join(message, "")
	}

	return body
}

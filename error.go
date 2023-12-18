package mistral

type ErrorContext struct {
	EnumValues []string `json:"enum_values"`
}

type ErrorDetail []struct {
	Loc     []any        `json:"loc"`
	Msg     string       `json:"msg"`
	Type    string       `json:"type"`
	Context ErrorContext `json:"ctx"`
}

type ErrorMessage struct {
	Detail ErrorDetail `json:"detail"`
}

type ErrorResponse struct {
	Code    string  `json:"code"`
	Message string  `json:"message"`
	Object  string  `json:"object"`
	Param   *string `json:"param,omitempty"`
	Type    string  `json:"type"`
}

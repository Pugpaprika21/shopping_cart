package dto

type Respone struct {
	Message    string      `json:"message,omitempty"`
	Statusbool bool        `json:"statusbool,omitempty"`
	Data       interface{} `json:"data,omitempty"`
}

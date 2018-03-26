package simpleapm

type PostHttp struct {
	*Post
	Data PostData `json:"data"`
}

type PostData struct {
	Path     string  `json:"path"`
	Method   string  `json:"method"`
	Duration float64 `json:"duration"`
	Query    string  `json:"query"`
	HttpCode int     `json:"httpCode"`
}

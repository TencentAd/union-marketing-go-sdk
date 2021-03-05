package http_tools

// HttpClient
type HttpConfig struct {
	BasePath      string            `json:"basePath,omitempty"`
	Host          string            `json:"host,omitempty"`
	Scheme        string            `json:"scheme,omitempty"`
	DefaultHeader map[string]string `json:"defaultHeader,omitempty"`
	UserAgent     string            `json:"userAgent,omitempty"`
}

func (c *HttpConfig) AddDefaultHeader(key string, value string) {
	if c.DefaultHeader == nil {
		c.DefaultHeader = make( map[string]string)
	}
	c.DefaultHeader[key] = value
}
package Agent

type ContextItem struct {
	MessageID int64  `json:"MessageID"`
	Align     string `json:"align,omitempty"`
	Variant   string `json:"variant,omitempty"`
	Context   string `json:"context,omitempty"`
	Name      string `json:"name,omitempty"`
	Time      string `json:"time,omitempty"`
}

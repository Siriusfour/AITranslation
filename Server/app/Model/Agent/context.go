package Agent

type ContextItem struct {
	MessageID int64  `json:"MessageID"`
	Align     string `json:"align,omitempty"`
	Variant   string `json:"variant,omitempty"`
	Time      string `json:"time,omitempty"`

	Type string `json:"type"` // human / ai / system / tool

	Data ContextData `json:"data"` //
}

type ContextData struct {
	Content          string                 `json:"content"`                     // ✅ 必须
	AdditionalKwargs map[string]interface{} `json:"additional_kwargs,omitempty"` // 可选
	ResponseMetadata map[string]interface{} `json:"response_metadata,omitempty"` // 可选
}

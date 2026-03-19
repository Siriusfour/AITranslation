package Agent

type ContextItem struct {
	MessageID int64      `json:"MessageID"`
	Align     string     `json:"align,omitempty"`
	Variant   string     `json:"variant,omitempty"`
	Time      string     `json:"time,omitempty"`

	Type string         `json:"type"` // human / ai / system / tool

	Data ContextData    `json:"data"` //
}

type ContextData struct {
	Content           string                 `json:"content"`                     // ✅ 必须
	AdditionalKwargs  map[string]interface{} `json:"additional_kwargs,omitempty"` // 可选
	ResponseMetadata  map[string]interface{} `json:"response_metadata,omitempty"` // 可选
}

[
{ "type": "system", "data": { "content": "你是助手" } },
{ "type": "human", "data": { "content": "内毒素检测需要什么试剂" } },
{ "type": "ai", "data": { "content": "常见有鲎试剂、标准品..." } }
]
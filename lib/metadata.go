package lib

import "encoding/json"

type Metadata struct {
	Type       string
	Provider   string
	attributes map[string]interface{} `json:"attributes"`
}

func (m *Metadata) SetType(typeStr string) {
	m.Type = typeStr
}

func (m *Metadata) SetProvider(provider string) {
	m.Provider = provider
}

func (m *Metadata) SetAttr(name string, value interface{}) {
	m.attributes[name] = value
}

func (m *Metadata) Attr(name string) (interface{}, bool) {
	value, ok := m.attributes[name]

	return value, ok
}

func (m Metadata) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":       m.Type,
		"provider":   m.Provider,
		"attributes": m.attributes,
	})
}

func NewMetadata() *Metadata {
	m := new(Metadata)
	m.attributes = make(map[string]interface{})

	return m
}

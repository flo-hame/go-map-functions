package converter

type FieldMapping struct {
	Field        string              `json:"field,omitempty"`
	Type         *string             `json:"type,omitempty"`
	FixValue     *string             `json:"fix_value,omitempty"`
	ValueMapping []FieldValueMapping `json:"value_mapping,omitempty"`
}

type FieldValueMapping struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

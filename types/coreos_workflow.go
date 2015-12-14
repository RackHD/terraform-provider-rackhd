package types

type CoreosWorkflow struct {
	Name    string        `json:"name"`
	Options CoreosOptions `json:"options"`
}

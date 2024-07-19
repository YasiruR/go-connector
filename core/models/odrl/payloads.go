package odrl

type Action string
type Assigner string
type Assignee string
type Target string

type Permission struct {
	Action      Action       `json:"odrl:action"`
	Constraints []Constraint `json:"odrl:constraint"`
}

type Constraint struct {
	LeftOperand  string `json:"odrl:leftOperand"`
	Operator     string `json:"odrl:operator"`
	RightOperand struct {
		Value string `json:"@value"`
		Type  string `json:"@type"`
	} `json:"odrl:rightOperand"`
}

type Duty struct {
	Action Action `json:"odrl:action"`
}

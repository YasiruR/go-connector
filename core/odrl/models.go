package odrl

type Policy struct {
	Id       string
	Assigner string
}

type Permission struct {
	Action      string
	Constraints []Constraint
}

type Constraint struct {
	LeftOperand  string
	Operator     string
	RightOperand string
}

type Duty struct {
	Action string
}

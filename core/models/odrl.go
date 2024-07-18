package models

type ODRLPolicy struct {
	Id       string
	Assigner string
}

type ODRLPermission struct {
	Action      string
	Constraints []ODRLConstraint
}

type ODRLConstraint struct {
	LeftOperand  string
	Operator     string
	RightOperand string
}

type ODRLDuty struct {
	Action string
}

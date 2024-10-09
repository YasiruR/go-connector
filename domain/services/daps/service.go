package daps

type Service interface {
	Metadata()
	Register(jwt string)
}

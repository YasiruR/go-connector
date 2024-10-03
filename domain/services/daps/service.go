package daps

type Service interface {
	Metadata()
	RequestToken()
}

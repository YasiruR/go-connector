package clearing_house

type Service interface {
	LogMessage()
	QueryByProcessID()
}

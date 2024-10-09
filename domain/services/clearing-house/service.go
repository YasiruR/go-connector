package clearing_house

type Service interface {
	CreateProcess()
	LogMessage(processId, msg string)
	QueryMessage(msgId string) string
	QueryMessages(processId string) []string
}

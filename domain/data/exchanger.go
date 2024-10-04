package data

type Exchanger interface {
	NewToken(participantId, datasetId string) string
	Push(datasetId, host, db, usr, pw string) error
	Pull(datasetId, endpoint, token string)
}

package models


type OutboxMessage struct {
	ID      int64  `json:"id"`
	Topic   string `json:"topic"`
	Key     []byte `json:"key"`
	Message []byte `json:"message"`
	Sent    string   `json:"sent"`
}

package queue

type QueueInterface interface {
	SendID(id string)
	RecieveID() (bool, string)
}

type QueueClient struct {
	QueueInterface
}

func NewQueueClient(queue QueueInterface) *QueueClient {
	return &QueueClient{queue}
}

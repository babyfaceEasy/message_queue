package models

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

// Message the model to hold message info
type Message struct {
	ID          int
	QueueID     int
	Message     string
	Status      string
	AvailableAt *time.Time
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}

// CreateMessage this adds a new message to the given queue
func (m *Message) CreateMessage(queueName string) (Message, error) {
	// get queue ID from name
	//err := new(Queue).GetQueueByName(queueName)
	queueDets := Queue{}
	err := queueDets.GetQueueByName(queueName)
	if err != nil {
		return Message{}, err
	}

	// queue doesn't exist
	if (Queue{}) == queueDets {
		return Message{}, errors.New("Queue : queue doesn't exist")
	}

	db, err := gorm.Open("mysql", "root:root@/message_queue?charset=utf8&parseTime=True&loc=Local")
	defer db.Close()
	if err != nil {
		return Message{}, err
	}

	// save details
	m.QueueID = queueDets.ID
	possible := db.NewRecord(&m)
	if possible {
		db.Create(&m)
	}

	return *m, nil
}

// GetQueue returns the details of the queue this message belongs to
func (m Message) GetQueue() (Queue, error) {
	db, err := gorm.Open("mysql", "root:root@/message_queue?charset=utf8&parseTime=True&loc=Local")
	defer db.Close()
	if err != nil {
		return Queue{}, err
	}

	queue := new(Queue)
	db.Model(&m).Related(&queue)

	return *queue, nil
}

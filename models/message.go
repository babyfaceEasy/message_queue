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
	Status      MessageStatus `sql:"not null;type:ENUM('created', 'in_transit', 'requeued', 'processed', 'queued')"`
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
	m.Status = created
	err = db.Create(&m).Error
	if err != nil {
		return Message{}, err
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
	if err := db.Model(&m).Related(&queue).Error; err != nil {
		return Queue{}, err
	}

	return *queue, nil
}

// GetMessageByID returns a given message if the message exist else returns empty object
func (m *Message) GetMessageByID(id int64) error {
	db, err := gorm.Open("mysql", "root:root@/message_queue?charset=utf8&parseTime=True&loc=Local")
	defer db.Close()
	if err != nil {
		return err
	}

	if err := db.First(&m, id).Error; err != nil {
		return err
	}

	return nil
}

// UpdateStatus moves the status of the message model in the right direction
func (m *Message) UpdateStatus() error {
	var newStatus MessageStatus = m.Status
	if m.Status == inTransit {
		newStatus = queued
	} else if m.Status == queued {
		newStatus = processed
	}
	db, err := gorm.Open("mysql", "root:root@/message_queue?charset=utf8&parseTime=True&loc=Local")
	defer db.Close()
	if err != nil {
		return err
	}

	if err := db.Model(&m).Update("status", newStatus).Error; err != nil {
		return err
	}

	return nil
}

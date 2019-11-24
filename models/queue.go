package models

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// Queue the struct to hold queue information
type Queue struct {
	ID        int
	Name      string
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

// CreateQueue adds a queue to the database
func (q Queue) CreateQueue() (Queue, error) {
	if q.Name == "" {
		err := errors.New("Name must be specified")
		return Queue{}, err
	}
	connString := fmt.Sprintf(
		"%s:%s@/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	db, err := gorm.Open(os.Getenv("DIALECT"), connString)
	defer db.Close()
	if err != nil {
		return Queue{}, err
	}
	// check if queue already exist
	exists := Queue{}
	db.Where("name = ?", q.Name).First(&exists)
	if (Queue{}) != exists {
		return Queue{}, errors.New("Queue already exist")
	}
	possible := db.NewRecord(&q)
	if possible {
		db.Create(&q)
	}
	return q, nil
}

// GetQueueByID returns queue detail of the given id
func (q *Queue) GetQueueByID(id int) error {
	db, err := gorm.Open("mysql", "root:root@/message_queue?charset=utf8&parseTime=True&loc=Local")
	defer db.Close()
	if err != nil {
		return err
	}

	db.First(&q, id)

	return nil
}

// GetQueueByName returns queue details of the given queue_name
func (q *Queue) GetQueueByName(queueName string) error {
	db, err := gorm.Open("mysql", "root:root@/message_queue?charset=utf8&parseTime=True&loc=Local")
	defer db.Close()
	if err != nil {
		return err
	}
	db.Where("name = ? ", queueName).First(&q)

	return nil
}

// GetMessages returns the messages attached to this queue
func (q Queue) GetMessages() ([]Message, error) {
	db, err := gorm.Open("mysql", "root:root@/message_queue?charset=utf8&parseTime=True&loc=Local")
	defer db.Close()
	if err != nil {
		return []Message{}, err
	}

	var messages []Message
	db.Model(&q).Related(&messages)

	return messages, nil
}

// GetMessage returns the oldest message inside the queue
func (q Queue) GetMessage() (Message, error) {
	db, err := gorm.Open("mysql", "root:root@/message_queue?charset=utf8&parseTime=True&loc=Local")
	defer db.Close()
	if err != nil {
		return Message{}, err
	}

	var messages []Message
	db.Model(&q).Related(&messages)
	mostRecent := len(messages) - 1

	return messages[mostRecent], nil
}

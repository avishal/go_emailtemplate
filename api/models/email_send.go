package models

import (
	// "log"
	// "context"
	"fmt"
	"time"

	// "go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/bson/primitive"
	// "encoding/json"
	// "github.com/gin-gonic/gin"
	// "net/http"
	// "io/ioutil"
	// "tradewindsnew/graph/model"
	// "math/rand"
	"gorm.io/gorm"
)

type EmailSend struct {
	Id            uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Client_id     string    `gorm:"size:255;not null;" json:"client_id"`
	Template_id   string    `gorm:"size:255;not null;" json:"template_id"`
	Is_sent       string    `gorm:"size:255;" json:"is_sent"`
	Status        string    `gorm:"size:255;not null;" json:"status"`
	Sent_datetime time.Time `sql:"default:null" json:"sent_datetime"`
	Send_from     string    `gorm:"size:255;not null;" json:"send_from"`
	Send_to       string    `gorm:"size:255;not null;" json:"send_to"`
	Cc            string    `gorm:"size:255;not null;" json:"cc"`
	Subject       string    `gorm:"size:255;not null;" json:"subject"`
	Body          string    `gorm:"size:255;not null;" json:"body"`
	Receiver_type bool      `gorm:"size:10;not null;" sql:"default:0" json:"receiver_type"`
	Receiver_id   string    `gorm:"size:25;not null;" sql:"default:no" json:"receiver_id"`
	Error_message string    `gorm:"size:25;null;" sql:"default:no" json:"error_message"`
	Created_dt    time.Time `sql:"default:CURRENT_TIMESTAMP" json:"created_dt"`
}

type EmailObj struct {
	Client_id      uint32   `json:"client_id"`
	Template_id    uint32   `json:"template_id"`
	Template_title string   `json:"template_title"`
	Receiver_id    string   `json:"receiver_id"`
	Receiver_type  string   `json:"receiver_type"`
	Variables      []string `json:"variables"`
}

func (u *EmailSend) Prepare() {
	// u.Username = html.EscapeString(strings.TrimSpace(u.Username))
	// u.Email = html.EscapeString(strings.TrimSpace(u.Email))
	u.Created_dt = time.Now()
	u.Sent_datetime = time.Now()
	// u.Deleted_dt = nil
}

func (u *EmailSend) FindAllEmails(db *gorm.DB) (*[]EmailSend, error) {
	var err error
	ets := []EmailSend{}
	err = db.Debug().Model(&EmailSend{}).Limit(100).Find(&ets).Error
	if err != nil {
		return &[]EmailSend{}, err
	}
	return &ets, err
}

func (u *EmailSend) SaveEmail(db *gorm.DB) (*EmailSend, error) {

	var err error
	fmt.Println(u)
	err = db.Debug().Create(&u).Error
	if err != nil {
		return &EmailSend{}, err
	}
	return u, nil
}

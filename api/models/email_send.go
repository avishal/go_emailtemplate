package models

import (
	// "log"
	// "context"
	"fmt"
	"time"
	// "github.com/go-sql-driver/mysql"
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
	Template_id   uint32    `gorm:"size:255;not null;" json:"template_id"`
	Is_sent       uint32    `json:"is_sent"`
	Status        string    `gorm:"size:255;not null;" json:"status"`
	Sent_datetime time.Time `sql:"default:null" json:"sent_datetime"`
	Send_from     string    `gorm:"size:255;not null;" json:"send_from"`
	Send_to       string    `gorm:"size:255;not null;" json:"send_to"`
	Cc            string    `gorm:"size:255;not null;" json:"cc"`
	Subject       string    `gorm:"size:255;not null;" json:"subject"`
	Body          string    `gorm:"size:255;not null;" json:"body"`
	Receiver_type string      `gorm:"size:10;not null;" sql:"default:0" json:"receiver_type"`
	Receiver_id   uint32    `gorm:"size:25;not null;" sql:"default:no" json:"receiver_id"`
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
	err = db.Debug().Model(&EmailSend{}).Find(&ets).Error
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

func (u *EmailSend) SaveUpdateEmail(db *gorm.DB, id uint32) (*EmailSend, error) {

	var err error
	fmt.Println(u)
	// err = db.Debug().Create(&u).Error
	err = db.Debug().Model(&EmailSend{}).Where("id = ?", id).Updates(&u).Error
	if err != nil {
		return &EmailSend{}, err
	}
	return u, nil
}

func (u *EmailSend) FindEmailByCol(db *gorm.DB, field string, value uint32) (*[]EmailSend, error) {
	var err error
	ets := []EmailSend{}
	err = db.Debug().Model(&EmailSend{}).Where(field, value).Find(&ets).Error
	if err != nil {
		return &[]EmailSend{}, err
	}
	return &ets, err
}

func (u *EmailSend) FindEmailByStringCol(db *gorm.DB, field string, value string) (*[]EmailSend, error) {
	var err error
	ets := []EmailSend{}
	err = db.Debug().Model(&EmailSend{}).Where(field, value).Find(&ets).Error
	if err != nil {
		return &[]EmailSend{}, err
	}
	return &ets, err
}


func (u *EmailSend) FindEmailByIdType(db *gorm.DB, id uint32, receiver_type string) (*[]EmailSend, error) {
	var err error
	ets := []EmailSend{}
	err = db.Debug().Model(&EmailSend{}).Where("receiver_id = ?", id).Where("receiver_type = ?", receiver_type).Find(&ets).Error
	if err != nil {
		return &[]EmailSend{}, err
	}
	return &ets, err
}

func (u *EmailSend) DeleteEmail(db *gorm.DB, uid uint32) (int64, error) {
	db.Debug().Model(EmailSend{}).Where("id = ?", uid).Take(&u).Delete(&EmailSend{})
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}


func (u *EmailSend) FindAllUnsentEmails(db *gorm.DB) (*[]EmailSend, error) {
	var err error
	ets := []EmailSend{}
	err = db.Debug().Model(&EmailSend{}).Where("status = ?", "new").Where("is_sent = ?", 0).Find(&ets).Error
	if err != nil {
		return &[]EmailSend{}, err
	}
	return &ets, err
}
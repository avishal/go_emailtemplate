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

type EmailTemplate struct {
	Id uint32 `gorm:"primary_key;auto_increment" json:"id"`
	Client_id string `gorm:"size:255;not null;" json:"client_id"`
	Template_title    string `gorm:"size:255;not null;" json:"template_title"`
	Template_description string `gorm:"size:255;" json:"template_description"`
	Template_type string `gorm:"size:255;not null;" json:"template_type"`
	Send_from string `gorm:"size:255;not null;" json:"send_from"`
	Send_to string `gorm:"size:255;not null;" json:"send_to"`
	Cc string `gorm:"size:255;not null;" json:"cc"`
	Subject string `gorm:"size:255;not null;" json:"subject"`
	Body string `gorm:"size:255;not null;" json:"body"`
	Has_attachments bool `gorm:"size:10;not null;" sql:"default:0" json:"has_attachments"`
	Attachments string `gorm:"size:25;not null;" sql:"default:no" json:"attachments"`
	Created_dt time.Time `sql:"default:CURRENT_TIMESTAMP" json:"created_dt"`
	Updated_dt time.Time `sql:"default:CURRENT_TIMESTAMP" json:"updated_dt"`
	Deleted_dt time.Time `gorm:"default:null" sql:"default:null" json:"deleted_dt"`
}

type EmailTemplateData struct {
	// Template_data EmailTemplate
	Id uint32
	Client_id string
	Template_title string
	Variables []string
}

func (u *EmailTemplate) Prepare() {
	// u.Username = html.EscapeString(strings.TrimSpace(u.Username))
	// u.Email = html.EscapeString(strings.TrimSpace(u.Email))
	u.Created_dt = time.Now()
	u.Updated_dt = time.Now()
	// u.Deleted_dt = nil
}

func (u *EmailTemplate) FindAllTemplates(db *gorm.DB) (*[]EmailTemplate, error) {
	var err error
	ets := []EmailTemplate{}
	// err = db.Debug().Model(&EmailTemplate{}).Limit(100).Find(&ets).Error
	err = db.Debug().Model(&EmailTemplate{}).Find(&ets).Error
	if err != nil {
		return &[]EmailTemplate{}, err
	}
	return &ets, err
}

func (u *EmailTemplate) SaveTemplate(db *gorm.DB) (*EmailTemplate, error) {

	var err error
	fmt.Println(u);
	err = db.Debug().Create(&u).Error
	if err != nil {
		return &EmailTemplate{}, err
	}
	return u, nil
}

func (u *EmailTemplate) UpdateTemplate(db *gorm.DB, id uint32) (*EmailTemplate, error) {

	var err error
	fmt.Println(u);
	err = db.Debug().Model(&EmailTemplate{}).Where("id = ?", id).Updates(&u).Error
	if err != nil {
		return &EmailTemplate{}, err
	}
	return u, nil
}


func (u *EmailTemplate) FindTemplateByID(db *gorm.DB, uid uint32) (*EmailTemplate, error) {
	var err error
	err = db.Debug().Model(EmailTemplate{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &EmailTemplate{}, err
	}
	// if gorm.IsRecordNotFoundError(err) {
	// 	return &User{}, errors.New("User Not Found")
	// }
	return u, err
}

func (u *EmailTemplate) DeleteTemplateByID(db *gorm.DB, uid uint32) (int64, error) {
	db.Debug().Model(EmailTemplate{}).Where("id = ?", uid).Take(&u).Delete(&EmailTemplate{})
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

func (p *EmailTemplate) FindTemplateByClientID(db *gorm.DB,clientid string) (*[]EmailTemplate, error) {
	var err error
	lookups := []EmailTemplate{}
	
	err = db.Debug().Model(&EmailTemplate{}).Where("client_id = ?", clientid).Find(&lookups).Error
	if err != nil {
		return &[]EmailTemplate{}, err
	}

	return &lookups, nil
}


func (p *EmailTemplate) FindTemplateByClientIDTemplateTitle(db *gorm.DB,clientid string,title string) (*[]EmailTemplate, error) {
	var err error
	templates := []EmailTemplate{}
	
	err = db.Debug().Model(&EmailTemplate{}).Where("client_id = ?", clientid).Where("template_title = ?", title).Find(&templates).Error
	if err != nil {
		return &[]EmailTemplate{}, err
	}

	return &templates, nil
}
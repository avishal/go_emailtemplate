package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"

	// "log"
	"net/http"
	"reflect"
	"regexp"
	"unsafe"

	// "tradewindsnew/api/fileupload"

	// "github.com/joho/godotenv"
	// "golang.org/x/crypto/bcrypt"

	// "email-template/api/auth"
	"email-template/api/models"
	// "email-template/api/security"
	// "email-template/api/utils/formaterror"

	"github.com/gin-gonic/gin"
	// "time"
	"bytes"
	"net/smtp"
	"text/template"

	"github.com/jasonlvhit/gocron"
)

func (server *Server) CreateEmail(c *gin.Context) {

	//clear previous error if any
	errList = map[string]string{}
	// c.MultipartForm()
	// c.PostForm.
	// for key, value := range c.Request.PostForm {
	// 	log.Printf("%v = %v \n", key, value)
	// }

	body, err := ioutil.ReadAll(c.Request.Body)

	if err != nil {
		errList["Invalid_body"] = "Unable to get request"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	var email map[string]interface{}
	err = json.Unmarshal(body, &email)

	// err = json.Unmarshal(email["data"], &varData)
	varData := email["data"].(map[string]interface{})

	if err != nil {
		errList["Unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	// fmt.Println("------>>>")
	// templateData, err := etmp.FindTemplateByClientIDTemplateTitle(server.DB,  uint32(email["client_id"].(float64)), email["template_title"] )
	etmp := models.EmailTemplate{}
	templateData, err := etmp.FindTemplateByClientIDTemplateTitle(server.DB, (email["client_id"]).(string), (email["template_title"]).(string))
	// fmt.Println("-----1")
	// fmt.Println(templateData)
	if len(*templateData) > 0 {
		lookinto := []string{"Send_from", "Send_to", "Cc", "Subject", "Body"}
		// lookinto := []string{"Send_from"};

		s := reflect.ValueOf(&(*templateData)[0]).Elem()
		// fmt.Println(s.FieldByName("Send_from").String())
		for i := 0; i < len(lookinto); i++ {
			// fmt.Println("-----0")
			for key, result := range varData {
				// r, _ := regexp.Compile("{[\\w ]+}")
				r := regexp.MustCompile("[$1][{1]" + key + "[}+]")
				// fmt.Println("key: ",key)
				// fmt.Println(lookinto[i],s.FieldByName(lookinto[i]).String())
				// fmt.Println("-----1")
				// fmt.Println(result)
				// fmt.Println("-----2")
				fitem := s.FieldByName(lookinto[i])
				replaceItem := r.ReplaceAllString(fitem.String(), result.(string))

				rf := reflect.NewAt(fitem.Type(), unsafe.Pointer(fitem.UnsafeAddr())).Elem()
				ri := reflect.ValueOf(&replaceItem).Elem()
				rf.Set(ri)
				// fmt.Println(key,lookinto[i], s.FieldByName(lookinto[i]).String());
				// fmt.Println("-----3")
			}
		}

		newemail := models.EmailSend{}
		newemail.Send_from = (*templateData)[0].Send_from
		newemail.Send_to = (*templateData)[0].Send_to
		newemail.Cc = (*templateData)[0].Cc
		newemail.Subject = (*templateData)[0].Subject
		newemail.Body = (*templateData)[0].Body
		newemail.Client_id = (*templateData)[0].Client_id
		newemail.Template_id = (*templateData)[0].Id
		newemail.Receiver_type = email["receiver_type"].(string)
		newemail.Receiver_id = uint32(email["receiver_id"].(float64))
		newemail.Status = "new"
		newemail.Error_message = ""

		newemail.Is_sent = 0
		newemail.Prepare()
		fmt.Println("-----2")
		emailCreated, err := newemail.SaveEmail(server.DB)
		if err != nil {
			// formattedError := formaterror.FormatError(err.Error())
			// errList = formattedError
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": http.StatusInternalServerError,
				"error":  err.Error(),
			})
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"status":   http.StatusCreated,
			"response": emailCreated,
		})

		// fmt.Println(newemail)

	} else {
		c.JSON(http.StatusCreated, gin.H{
			"status":   http.StatusCreated,
			"response": "Template not found",
		})
	}
	// email.Prepare()

	/*userCreated, err := user.SaveEmail(server.DB)
	if err != nil {
		// formattedError := formaterror.FormatError(err.Error())
		// errList = formattedError
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status":   http.StatusCreated,
		"response": userCreated,
	})*/
	/*errorMessages := user.Validate("")
	if len(errorMessages) > 0 {
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}*/
}

func (server *Server) GetEmail(c *gin.Context) {

	//clear previous error if any
	errList = map[string]string{}

	emailsend := models.EmailSend{}
	field := c.Param("field")
	value := c.Param("value")

	if field == "id" {
		tempId, err := strconv.ParseUint(value, 10, 32)
		emails, err := emailsend.FindEmailByCol(server.DB, field, uint32(tempId))

		if err != nil {
			errList["No_emailsend"] = "No email Found"
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": http.StatusInternalServerError,
				"error":  errList,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":   http.StatusOK,
			"response": emails,
		})

	} else if field == "client_id" {
		// tempId, err := strconv.ParseUint(value, 10, 32)
		emails, err := emailsend.FindEmailByStringCol(server.DB, field, value)
		if err != nil {
			errList["No_emailsend"] = "No email Found"
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": http.StatusInternalServerError,
				"error":  errList,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":   http.StatusOK,
			"response": emails,
		})

	} else if field == "receiver_id" {
		tempId, err := strconv.ParseUint(value, 10, 32)
		emails, err := emailsend.FindEmailByCol(server.DB, field, uint32(tempId))
		if err != nil {
			errList["No_emailsend"] = "No email Found"
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": http.StatusInternalServerError,
				"error":  errList,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":   http.StatusOK,
			"response": emails,
		})

	} else if field == "receiver_type" {
		emails, err := emailsend.FindEmailByStringCol(server.DB, field, value)
		if err != nil {
			errList["No_emailsend"] = "No email Found"
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": http.StatusInternalServerError,
				"error":  errList,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":   http.StatusOK,
			"response": emails,
		})

	} else if field == "status" {

		emails, err := emailsend.FindEmailByStringCol(server.DB, field, value)
		if err != nil {
			errList["No_emailsend"] = "No email Found"
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": http.StatusInternalServerError,
				"error":  errList,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":   http.StatusOK,
			"response": emails,
		})
	}

}

func (server *Server) GetAllEmail(c *gin.Context) {

	//clear previous error if any
	errList = map[string]string{}

	template := models.EmailSend{}

	templates, err := template.FindAllEmails(server.DB)
	if err != nil {
		errList["No_template"] = "No emails Found"
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  errList,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": templates,
	})
}

func (server *Server) GetEmailByIdType(c *gin.Context) {

	//clear previous error if any
	errList = map[string]string{}

	template := models.EmailSend{}
	id := c.Param("id")
	templId, err := strconv.ParseUint(id, 10, 32)
	receiver_type := c.Param("receiver_type")
	templates, err := template.FindEmailByIdType(server.DB, uint32(templId), receiver_type)
	if err != nil {
		errList["No_template"] = "No email Found"
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  errList,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": templates,
	})
}

func (server *Server) GetUnSentSellerMails(c *gin.Context) {

	//clear previous error if any
	errList = map[string]string{}

	template := models.EmailSend{}

	templates, err := template.FindAllUnsentSellerEmails(server.DB)
	if err != nil {
		errList["No_template"] = "No email Found"
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  errList,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": templates,
	})
}

func (server *Server) DeleteEmail(c *gin.Context) {

	//clear previous error if any
	errList = map[string]string{}

	template := models.EmailSend{}
	id := c.Param("id")

	templId, err := strconv.ParseUint(id, 10, 32)
	_, err = template.DeleteEmail(server.DB, uint32(templId))

	if err != nil {
		errList["No_template"] = "No email Found"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": "Email Deleted",
	})

}

func (server *Server) UpdateEmail(c *gin.Context) {

	//clear previous error if any
	errList = map[string]string{}
	//template := models.EmailSend{}
	id := c.Param("id")

	tempId, err := strconv.ParseUint(id, 10, 32)
	body, err := ioutil.ReadAll(c.Request.Body)

	if err != nil {
		errList["Invalid_body"] = "Unable to get request"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	var email map[string]interface{}
	err = json.Unmarshal(body, &email)

	// err = json.Unmarshal(email["data"], &varData)
	varData := email["data"].(map[string]interface{})

	if err != nil {
		errList["Unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	// fmt.Println("------>>>")
	// templateData, err := etmp.FindTemplateByClientIDTemplateTitle(server.DB,  uint32(email["client_id"].(float64)), email["template_title"] )
	etmp := models.EmailTemplate{}
	templateData, err := etmp.FindTemplateByClientIDTemplateTitle(server.DB, (email["client_id"]).(string), (email["template_title"]).(string))
	// fmt.Println("-----1")
	// fmt.Println(templateData)
	if len(*templateData) > 0 {
		lookinto := []string{"Send_from", "Send_to", "Cc", "Subject", "Body"}
		// lookinto := []string{"Send_from"};

		s := reflect.ValueOf(&(*templateData)[0]).Elem()
		// fmt.Println(s.FieldByName("Send_from").String())
		for i := 0; i < len(lookinto); i++ {
			// fmt.Println("-----0")
			for key, result := range varData {
				// r, _ := regexp.Compile("{[\\w ]+}")
				r := regexp.MustCompile("[$1][{1]" + key + "[}+]")
				// fmt.Println("key: ",key)
				// fmt.Println(lookinto[i],s.FieldByName(lookinto[i]).String())
				// fmt.Println("-----1")
				// fmt.Println(result)
				// fmt.Println("-----2")
				fitem := s.FieldByName(lookinto[i])
				replaceItem := r.ReplaceAllString(fitem.String(), result.(string))

				rf := reflect.NewAt(fitem.Type(), unsafe.Pointer(fitem.UnsafeAddr())).Elem()
				ri := reflect.ValueOf(&replaceItem).Elem()
				rf.Set(ri)
				// fmt.Println(key,lookinto[i], s.FieldByName(lookinto[i]).String());
				// fmt.Println("-----3")
			}
		}

		newemail := models.EmailSend{}
		newemail.Send_from = (*templateData)[0].Send_from
		newemail.Send_to = (*templateData)[0].Send_to
		newemail.Cc = (*templateData)[0].Cc
		newemail.Subject = (*templateData)[0].Subject
		newemail.Body = (*templateData)[0].Body
		newemail.Client_id = (*templateData)[0].Client_id
		newemail.Template_id = (*templateData)[0].Id
		newemail.Receiver_type = email["receiver_type"].(string)
		newemail.Receiver_id = uint32(email["receiver_id"].(float64))
		newemail.Status = "new"
		newemail.Error_message = ""

		newemail.Is_sent = 0
		newemail.Prepare()
		fmt.Println("-----2")
		emailCreated, err := newemail.SaveUpdateEmail(server.DB, uint32(tempId))
		if err != nil {
			// formattedError := formaterror.FormatError(err.Error())
			// errList = formattedError
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": http.StatusInternalServerError,
				"error":  err.Error(),
			})
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"status":   http.StatusCreated,
			"response": emailCreated,
		})

		// fmt.Println(newemail)

	} else {
		c.JSON(http.StatusCreated, gin.H{
			"status":   http.StatusCreated,
			"response": "Template not found",
		})
	}
	// email.Prepare()

	/*userCreated, err := user.SaveEmail(server.DB)
	if err != nil {
		// formattedError := formaterror.FormatError(err.Error())
		// errList = formattedError
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status":   http.StatusCreated,
		"response": userCreated,
	})*/
	/*errorMessages := user.Validate("")
	if len(errorMessages) > 0 {
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}*/
}

func task() {
	fmt.Println("I am running task.")
}

func (server *Server) StartScheduler(c *gin.Context) {
	fmt.Println("-----------------start scheduler-----------------")
	gocron.Every(10).Seconds().Do(server.GetNewMails)
	gocron.Start()
}

func (server *Server) StopScheduler(c *gin.Context) {
	fmt.Println("-----------------stop scheduler-----------------")
	gocron.Clear()
}

func (server *Server) GetNewMails() {
	email := models.EmailSend{}
	allEmails, err := email.FindAllUnsentEmails(server.DB)

	if err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{
		// 	"status": http.StatusInternalServerError,
		// 	"error":  err.Error(),
		// })
		// return
		fmt.Println("mail send err")
	}

	for i := 0; i < len(*allEmails); i++ {
		SendMail((*allEmails)[i])
		(*allEmails)[i].Is_sent = 1
		(*allEmails)[i].Status = "email_sent"
		_, err := (*allEmails)[i].SaveUpdateEmail(server.DB, (*allEmails)[i].Id)
		if err != nil {
			fmt.Println("err: update mail")
		}
	}
}

func SendMail(email_data models.EmailSend) {
	fmt.Println("-----------------Sending E-Mail----------------")
	fmt.Println(email_data.Subject)
	// Sender data.
	from := "a5fba4294ecf28"
	password := "8efe5a5e14b346"

	// Receiver email address.
	to := []string{
		email_data.Send_to,
	}

	// smtp server configuration.
	smtpHost := "smtp.mailtrap.io"
	smtpPort := "2525"

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	t, _ := template.ParseFiles("template.html")

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: "+email_data.Subject+"\n%s\n\n", mimeHeaders)))

	t.Execute(&body, struct {
		Body string
	}{
		Body: email_data.Body,
	})

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, email_data.Send_from, to, body.Bytes())
	// fmt.Println(body);
	// message := "This is a test email message."
	// err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message.Bytes())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Email Sent!")
}

func (server *Server) TestMail(c *gin.Context) {
	fmt.Println("-----------------Test Mail----------------")
	// Sender data.
	from := "a5fba4294ecf28"
	password := "8efe5a5e14b346"

	// Receiver email address.
	to := []string{
		"thiagu@gmx.com",
	}

	// smtp server configuration.
	smtpHost := "smtp.mailtrap.io"
	smtpPort := "2525"

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	t, _ := template.ParseFiles("template.html")

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: This is a test subject \n%s\n\n", mimeHeaders)))

	t.Execute(&body, struct {
		Name    string
		Message string
	}{
		Name:    "Puneet Singh",
		Message: "This is a test message in a HTML template",
	})

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, "testtest@test.com", to, body.Bytes())
	// fmt.Println(body);
	// message := "This is a test email message."
	// err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message.Bytes())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Email Sent!")
}

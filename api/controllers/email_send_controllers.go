package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	// "strconv"

	// "log"
	"net/http"
	"regexp"
	"reflect"
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
	templateData, err := etmp.FindTemplateByClientIDTemplateTitle(server.DB,  (email["client_id"]).(string), (email["template_title"]).(string) )
	// fmt.Println("-----1")
	// fmt.Println(templateData)
	if len(*templateData) > 0 {
		lookinto := []string{"Send_from", "Send_to", "Cc", "Subject", "Body"};
		// lookinto := []string{"Send_from"};
		
		s := reflect.ValueOf(&(*templateData)[0]).Elem()
		// fmt.Println(s.FieldByName("Send_from").String())
		for i:=0; i< len(lookinto); i++ {
			// fmt.Println("-----0")
			for key, result := range varData {
				// r, _ := regexp.Compile("{[\\w ]+}")
				r := regexp.MustCompile("[$1][{1]"+key+"[}+]")
				// fmt.Println("key: ",key)
				// fmt.Println(lookinto[i],s.FieldByName(lookinto[i]).String())
				// fmt.Println("-----1")
				// fmt.Println(result)
				// fmt.Println("-----2")
				fitem := s.FieldByName(lookinto[i])
				replaceItem := r.ReplaceAllString(fitem.String(), result.(string))

				rf := reflect.NewAt(fitem.Type(), unsafe.Pointer(fitem.UnsafeAddr())).Elem()
				ri := reflect.ValueOf(&replaceItem).Elem() 
				rf.Set(ri);
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
		newemail.Error_message =""
		
		newemail.Is_sent= 0
		newemail.Prepare();
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

	}else {
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

func (server *Server) GetEmailSends(c *gin.Context) {

	//clear previous error if any
	errList = map[string]string{}

	template := models.EmailSend{}

	templates, err := template.FindAllEmails(server.DB)
	if err != nil {
		errList["No_template"] = "No templates Found"
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

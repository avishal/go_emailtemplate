package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	// "strconv"

	// "log"
	"net/http"

	// "tradewindsnew/api/fileupload"

	// "github.com/joho/godotenv"
	// "golang.org/x/crypto/bcrypt"

	// "email-template/api/auth"
	"email-template/api/models"
	// "email-template/api/security"
	// "email-template/api/utils/formaterror"

	"github.com/gin-gonic/gin"
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

	// email := models.EmailSend{}
	var email map[string]interface{}
	err = json.Unmarshal(body, &email)
	
	// err = json.Unmarshal(email["data"], &varData)
	varData := email["data"].(map[string]interface{})
	for key, result := range varData {
		fmt.Println(key)
		fmt.Println("-----")
		fmt.Println(result)
	}

	if err != nil {
		errList["Unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	// templateData, err := etmp.FindTemplateByClientIDTemplateTitle(server.DB,  uint32(email["client_id"].(float64)), email["template_title"] )
	/*etmp := models.EmailTemplate{}
	templateData, err := etmp.FindTemplateByClientIDTemplateTitle(server.DB,  email["client_id"], email["template_title"] )
	if( len(templateData) > 0)
	{
		lookinto := []string{"Send_from", "Send_to", "Cc", "Subject", "Body"};
		r, _ := regexp.Compile("{[\\w ]+}")
		s := reflect.ValueOf(&templateData[0]).Elem()
		for i:=0; i< len(lookinto); i++ {
			r.ReplaceAllString(s.Elem().FieldByName(lookinto[i]).String(), email[data])
		}

	}else {
		c.JSON(http.StatusCreated, gin.H{
			"status":   http.StatusCreated,
			"response": "Template not found",
		})
	}*/
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

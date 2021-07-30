package controllers


import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	// "log"
	"net/http"
	"strconv"
	"regexp"
	"reflect"

	// "tradewindsnew/api/fileupload"

	// "github.com/joho/godotenv"
	// "golang.org/x/crypto/bcrypt"

	// "email-template/api/auth"
	"email-template/api/models"
	// "email-template/api/security"
	// "email-template/api/utils/formaterror"

	"github.com/gin-gonic/gin"
)

func (server *Server) CreateEmailTemplate(c *gin.Context) {

	//clear previous error if any
	errList = map[string]string{}
	// c.MultipartForm()
	// c.PostForm.
	// for key, value := range c.Request.PostForm {
	// 	log.Printf("%v = %v \n", key, value)
	// }

	
	body, err := ioutil.ReadAll(c.Request.Body)
	fmt.Println(string(body))
	if err != nil {
		errList["Invalid_body"] = "Unable to get request"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	user := models.EmailTemplate{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		errList["Unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	user.Prepare()
	/*errorMessages := user.Validate("")
	if len(errorMessages) > 0 {
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}*/
	userCreated, err := user.SaveTemplate(server.DB)
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
	})
}

func (server *Server) UpdateEmailTemplate(c *gin.Context) {

	//clear previous error if any
	errList = map[string]string{}
	// c.MultipartForm()
	// c.PostForm.
	// for key, value := range c.Request.PostForm {
	// 	log.Printf("%v = %v \n", key, value)
	// }
	templateId := c.Param("id")

	templId, err := strconv.ParseUint(templateId, 10, 32)
	
	body, err := ioutil.ReadAll(c.Request.Body)
	fmt.Println(string(body))
	if err != nil {
		errList["Invalid_body"] = "Unable to get request"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	template := models.EmailTemplate{}
	err = json.Unmarshal(body, &template)
	if err != nil {
		errList["Unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	template.Prepare()
	/*errorMessages := template.Validate("")
	if len(errorMessages) > 0 {
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}*/
	templateCreated, err := template.UpdateTemplate(server.DB, uint32(templId))
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
		"response": templateCreated,
	})
}

func (server *Server) GetEmailTemplates(c *gin.Context) {

	//clear previous error if any
	errList = map[string]string{}

	template := models.EmailTemplate{}

	templates, err := template.FindAllTemplates(server.DB)
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

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func (server *Server) GetEmailTemplateById(c *gin.Context) {

	//clear previous error if any
	errList = map[string]string{}

	templateId := c.Param("id")

	templId, err := strconv.ParseUint(templateId, 10, 32)
	if err != nil {
		errList["Invalid_request"] = "Invalid Request"
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  errList,
		})
		return
	}
	user := models.EmailTemplate{}
	lookinto := []string{"Send_from", "Send_to", "Cc", "Subject", "Body"};
	templateGotten, err := user.FindTemplateByID(server.DB, uint32(templId))
	
	r, _ := regexp.Compile("{[\\w ]+}")
	// match, _ := regexp.MatchString("{[\\w ]+}", templateGotten.Body)

	// fmt.Println("-------------");
	// fmt.Println(match);
	// fmt.Println("-------------");
	s := reflect.ValueOf(&templateGotten).Elem()
    // println(s.String())
    // fmt.Println(r)

	// fmt.Println(r);
	variablesData := make([]string, 0);
	for i:=0; i< len(lookinto); i++ {
		// fmt.Println(s.Elem().FieldByName(lookinto[i]).String())
		patternItems := r.FindAllString(s.Elem().FieldByName(lookinto[i]).String(),-1)
		for j:=0; j< len(patternItems); j++ {
			itm := patternItems[j][1 : len(patternItems[j])-1];
			if(!contains(variablesData, itm)){
				variablesData = append(variablesData, itm)
			}
		}
			
	}

	et := models.EmailTemplateData{}
	et.Id = templateGotten.Id
	et.Client_id = templateGotten.Client_id
	et.Template_title = templateGotten.Template_title
	et.Variables = variablesData
	// for i :=0 ; i < reflect.ValueOf(&templateGotten).Elem().Len(); i++ {
	// 	fmt.Println(s.Elem().Field(i).String());
	// }
	// fmt.Println(variablesData)
	// templateGotten.variables =
	if err != nil {
		errList["No_template"] = "No Template Found"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": et,
	})
}


func (server *Server) GetEmailTemplateByClientId(c *gin.Context) {

	//clear previous error if any
	errList = map[string]string{}

	clientID := c.Param("client_id")
	user := models.EmailTemplate{}
	// lookinto := []string{"Send_from", "Send_to", "Cc", "Subject", "Body"};
	templateGotten, err := user.FindTemplateByClientID(server.DB, clientID)
	// fmt.Println("=---------");
	/*fmt.Println((*templateGotten)[0]);
	allets := make([]models.EmailTemplateData, 0)

	r, _ := regexp.Compile("{[\\w ]+}")
	for j:=0; j< len(*templateGotten); j++ {
		x := (*templateGotten)[j]
		s := reflect.ValueOf(&x).Elem()
		variablesData := make([]string, 0);
		for i:=0; i< len(lookinto); i++ {
			// fmt.Println(s.Elem().FieldByName(lookinto[i]).String())
			patternItems := r.FindAllString(s.Elem().FieldByName(lookinto[i]).String(),-1)
			for j:=0; j< len(patternItems); j++ {
				itm := patternItems[j][1 : len(patternItems[j])-1];
				if(!contains(variablesData, itm)){
					variablesData = append(variablesData, itm)
				}
			}
				
		}

		et := models.EmailTemplateData{}
		et.Id = x.Id
		et.Client_id = x.Client_id
		et.Template_title = x.Template_title
		et.Variables = variablesData
		allets = append(allets, et)
	}*/
	if err != nil {
		errList["No_template"] = "No Template Found"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": templateGotten,
	})
}

func (server *Server) GetEmailTemplateByClientIdTemplateTitle(c *gin.Context) {

	//clear previous error if any
	errList = map[string]string{}

	clientID := c.Param("client_id")
	title := c.Param("title")
	user := models.EmailTemplate{}
	// lookinto := []string{"Send_from", "Send_to", "Cc", "Subject", "Body"};
	templateGotten, err := user.FindTemplateByClientIDTemplateTitle(server.DB, clientID, title)
	// fmt.Println("=---------");
	/*fmt.Println((*templateGotten)[0]);
	allets := make([]models.EmailTemplateData, 0)

	r, _ := regexp.Compile("{[\\w ]+}")
	for j:=0; j< len(*templateGotten); j++ {
		x := (*templateGotten)[j]
		s := reflect.ValueOf(&x).Elem()
		variablesData := make([]string, 0);
		for i:=0; i< len(lookinto); i++ {
			// fmt.Println(s.Elem().FieldByName(lookinto[i]).String())
			patternItems := r.FindAllString(s.Elem().FieldByName(lookinto[i]).String(),-1)
			for j:=0; j< len(patternItems); j++ {
				itm := patternItems[j][1 : len(patternItems[j])-1];
				if(!contains(variablesData, itm)){
					variablesData = append(variablesData, itm)
				}
			}
				
		}

		et := models.EmailTemplateData{}
		et.Id = x.Id
		et.Client_id = x.Client_id
		et.Template_title = x.Template_title
		et.Variables = variablesData
		allets = append(allets, et)
	}*/
	if err != nil {
		errList["No_template"] = "No Template Found"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": templateGotten,
	})
}

func (server *Server) DeleteEmailTemplate(c *gin.Context) {

	//clear previous error if any
	errList = map[string]string{}

	templateId := c.Param("id")

	templId, err := strconv.ParseUint(templateId, 10, 32)
	if err != nil {
		errList["Invalid_request"] = "Invalid Request"
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  errList,
		})
		return
	}
	user := models.EmailTemplate{}
	_, err = user.DeleteTemplateByID(server.DB, uint32(templId))
	
	if err != nil {
		errList["No_template"] = "No Template Found"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": "Template Deleted",
	})
}
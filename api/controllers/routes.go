package controllers

import (
	"net/http"
	// "email-template/api/middlewares"
	// _ "email-template/docs"

	"github.com/gin-gonic/gin"
	// swaggerFiles "github.com/swaggo/files"
	// ginSwagger "github.com/swaggo/gin-swagger"
	// "github.com/swaggo/gin-swagger/swaggerFiles"
)

func (s *Server) initializeRoutes() {

	v1 := s.Router.Group("/api/v1")
	{
		v1.GET("/", HealthCheck)

		// url := ginSwagger.URL("http://localhost:8888/api/v1/swagger/doc.json") // The url pointing to API definition
		// v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

		v1.GET("/get_templates", s.GetEmailTemplates)
		v1.GET("/get_template/:id", s.GetEmailTemplateById)
		v1.GET("/get_template/client/:client_id", s.GetEmailTemplateByClientId)
		v1.GET("/get_templates/:client_id/:title", s.GetEmailTemplateByClientIdTemplateTitle)

		v1.POST("/email_template/create", s.CreateEmailTemplate)
		v1.POST("/email_template/update/:id", s.UpdateEmailTemplate)
		v1.POST("/email_template/delete/:id", s.DeleteEmailTemplate)

		v1.POST("/send_email", s.CreateEmail)
		v1.GET("/get_email/:field/:value", s.GetEmail)
		v1.GET("/get_email_by_id_type/:id/:receiver_type", s.GetEmailByIdType)
		v1.GET("/get_all_emails", s.GetAllEmail)
		v1.POST("/update_email/:id", s.UpdateEmail)
		v1.POST("/delete_email/:id", s.DeleteEmail)
		
		v1.GET("/start-scheduler", s.StartScheduler)
		v1.GET("/stop-scheduler", s.StopScheduler)
		v1.GET("/test-mail", s.TestMail)
		// v1.POST("/get_email/client/:client_id", s.CreateEmail)
		// v1.POST("/get_email/receiver/:receiver_id", s.CreateEmail)
		// v1.POST("/get_email/:receiver_type", s.CreateEmail)
		// v1.POST("/get_email/status/:status", s.CreateEmail)
		//Customers routes
		// swagger:route POST /api/v1/customers CreateCustomers
		// v1.POST("/create_customer/", s.CreateCustomer)
		// // swagger:route GET /api/v1/customers GetCustomers
		// v1.GET("/get_customers/", s.GetCustomers)
		// // swagger:route GET /api/v1/customers/:id FindCustomer
		// v1.GET("/get_customer/:id/", s.GetCustomer)
		// // swagger:route GET /api/v1/get_active_customers GetActiveCustomers
		// v1.GET("/get_active_customers/", s.GetActiveCustomers)
		// // swagger:route GET /api/v1/search_customers/ SearchCustomers
		// v1.GET("/search_customers/", s.SearchCustomer)
		// // swagger:route PUT /api/v1/customers/:id UpdateCustomer
		// v1.PUT("/update_customer/:id/", s.UpdateCustomer)
		// // swagger:route DELETE /api/v1/customers/:id DeleteCustomer
		// v1.DELETE("/delete_customer/:id/", middlewares.TokenAuthMiddleware(), s.DeleteCustomer)
		// // v1.GET("/user_posts/:id", s.GetUserPosts)

		// //Sellers routes
		// // swagger:route POST /api/v1/sellers CreateSellers
		// v1.POST("/create_seller/", s.CreateSeller)
		// // swagger:route GET /api/v1/sellers GetSellers
		// v1.GET("/get_sellers/", s.GetSellers)
		// // swagger:route GET /api/v1/sellers/:id FindSeller
		// v1.GET("/get_seller/:id/", s.GetSeller)
		// // swagger:route GET /api/v1/get_active_sellers GetActiveSellers
		// v1.GET("/get_active_sellers/", s.GetActiveSellers)
		// // swagger:route GET /api/v1/search_sellers/ SearchSellers
		// v1.GET("/search_sellers/", s.SearchSeller)
		// // swagger:route PUT /api/v1/sellers/:id UpdateSeller
		// v1.PUT("/update_seller/:id/", s.UpdateSeller)
		// // swagger:route DELETE /api/v1/sellers/:id DeleteSeller
		// v1.DELETE("/delete_seller/:id/", s.DeleteSeller)

		// //Contents routes
		// // swagger:route POST /api/v1/sellers CreateContents
		// v1.POST("/create_content/", s.CreateContent)
		// // swagger:route GET /api/v1/contents GetContents
		// v1.GET("/get_contents/", s.GetContents)
		// // swagger:route GET /api/v1/get_active_contents GetActiveContents
		// v1.GET("/get_active_contents/", s.GetActiveContents)
		// // swagger:route GET /api/v1/contents FindContent
		// v1.GET("/get_content/:id/", s.GetContent)
		// // swagger:route GET /api/v1/search_contents/ SearchContents
		// v1.GET("/search_contents/", s.SearchContent)
		// // swagger:route PUT /api/v1/contents/:id UpdateSeller
		// v1.PUT("/update_content/:id/", s.UpdateContent)
		// // v1.PUT("/contents/:id", middlewares.TokenAuthMiddleware(), s.UpdateContent)
		// // swagger:route DELETE /api/v1/contents/:id DeleteContent
		// v1.DELETE("/delete_content/:id/", s.DeleteContent)

		// //Lookups routes
		// // swagger:route POST /api/v1/lookups CreateLookups
		// v1.POST("/create_lookup/", s.CreateLookup)
		// // swagger:route GET /api/v1/lookups GetLookups
		// v1.GET("/get_lookups/", s.GetLookups)
		// // swagger:route GET /api/v1/lookups FindLookup
		// v1.GET("/get_lookups/:id/", s.GetLookup)
		// // swagger:route GET /api/v1/get_active_lookups GetActiveLookups
		// v1.GET("/get_active_lookups/", s.GetActiveLookups)
		// // swagger:route GET /api/v1/search_lookups/ SearchLookup
		// v1.GET("/search_lookups/", s.SearchLookup)
		// // v1.PUT("/lookups/:id", middlewares.TokenAuthMiddleware(), s.UpdateLookup)
		// // swagger:route PUT /api/v1/lookups/:id UpdateLookup
		// v1.PUT("/update_lookup/:id/", s.UpdateLookup)
		// // swagger:route DELETE /api/v1/lookups/:id DeleteLookup
		// v1.DELETE("/delete_lookup/:id/", s.DeleteLookup)

		// // Routes for splash
		// v1.POST("/splashes", s.CreateSplash)
		// // swagger:route GET /api/v1/lookups GetLookups
		// v1.GET("/splashes", s.GetSplashes)

		// v1.POST("/users/", s.CreateUser)

		// //Users routes
		// v1.POST("/users", s.CreateUser)
		// // The user of the app have no business getting all the users.
		// v1.GET("/users", s.GetUsers)
		// // v1.GET("/users/:id", s.GetUser)
		// v1.PUT("/users/:id", middlewares.TokenAuthMiddleware(), s.UpdateUser)
		// v1.PUT("/avatar/users/:id", s.UpdateAvatar)
		// v1.DELETE("/users/:id", middlewares.TokenAuthMiddleware(), s.DeleteUser)

		// v1.GET("/search_audit/", s.SearchAudit)

	}
}


// HealthCheck godoc
// @Summary Show the status of server.
// @Description get the status of server.
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router / [get]
func HealthCheck(c *gin.Context) {
	res := map[string]interface{}{
		"data": "Server is up and running",
	}

	c.JSON(http.StatusOK, res)
}

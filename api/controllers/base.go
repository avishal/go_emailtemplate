package controllers

import (
	"fmt"
	"log"
	"net/http"

	"email-template/api/middlewares"
	"email-template/api/models"

	"github.com/gin-gonic/gin"
	// "github.com/jinzhu/gorm"
	// _ "github.com/jinzhu/gorm/dialects/sqlite" //mysql database driver
	"gorm.io/driver/mysql"
	// "gorm.io/driver/postgres"
	// "gorm.io/driver/sqlite"
	// "gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	// "context"
	// "time"
)

type Server struct {
	DB     *gorm.DB
	Router *gin.Engine
}

var errList = make(map[string]string)

func (server *Server) Initialize(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName, DbSchema string) {

	var err error
	fmt.Println(Dbdriver)

	// If you are using mysql, i added support for you here(dont forgot to edit the .env file)
	if Dbdriver == "mysql" {
		// dsns := "tradewinds_support_dev:synergy2021pwd@tcp(agency-cloud-non-prod-cluster.cluster-cha1o3asnyvz.us-east-2.rds.amazonaws.com:3306)/tradewinds_dev?charset=utf8&parseTime=True&loc=Local"
		DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", DbUser, DbPassword, DbHost, DbPort, DbName)
		fmt.Println(DBURL)
		server.DB, err = gorm.Open(mysql.Open(DBURL), &gorm.Config{})
		if err != nil {
			fmt.Printf("Cannot connect to %s database", Dbdriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database", Dbdriver)
		}
	} /*else if Dbdriver == "postgres" {
		DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
		server.DB, err = gorm.Open(postgres.Open(DBURL), &gorm.Config{})
		if err != nil {
			fmt.Printf("Cannot connect to %s database", Dbdriver)
			log.Fatal("This is the error connecting to postgres:", err)
		} else {
			fmt.Printf("We are connected to the %s database", Dbdriver)
		}
	} else if Dbdriver == "sqlserver" {
		DBURLs := fmt.Sprintf("sqlserver://%s:%s@%s:%s/?database=%s&schema=%s", DbUser, DbPassword, DbHost, DbPort, DbName, DbSchema)
		fmt.Println(DBURLs)
		// DBURL := "sqlserver://wmarket_dev:WMPlace@2021!@wmarketplace-dev.database.windows.net:1433/?database=wmarketplace-dev&schema=Dev"
		server.DB, err = gorm.Open(sqlserver.Open(DBURLs), &gorm.Config{})
		if err != nil {
			fmt.Printf("Cannot connect to %s database", Dbdriver)
			log.Fatal("This is the error connecting to postgres:", err)
		} else {
			fmt.Printf("We are connected to the %s database", Dbdriver)
		}
	} else {
		server.DB, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
		if err != nil {
			panic("Failed to connect to database!")
		}

		// fmt.Println("Unknown Driver")
	}*/

	// server.DB.DropTableIfExists(&models.Audit{})
	// server.DB.Migrator().DropTable(&models.EmailTemplate{})
	// server.DB.Migrator().DropTable(&models.Customer{})
	// server.DB.Migrator().DropTable(&models.Audit{})
	// server.DB.Migrator().DropTable(&models.User{})
	//database migration
	// message := server.DB.Migrator().CreateTable(&models.EmailTemplate{})
	message := server.DB.Migrator().CreateTable(&models.EmailSend{})
	fmt.Println(message)
	// message = server.DB.Migrator().CreateTable(&models.Customer{})
	// message = server.DB.Migrator().CreateTable(&models.Content{})
	// message = server.DB.Migrator().CreateTable(&models.Lookup{})
	// message = server.DB.Migrator().CreateTable(&models.Splash{})
	// message := server.DB.Migrator().CreateTable(&models.Audit{})
	// fmt.Println(message)
	server.DB.Debug().AutoMigrate(
		// &models.Customer{},
		// &models.Seller{},
		// &models.Content{},
		// &models.Lookup{},
		// &models.Splash{},
		// &models.Post{},
		// &models.ResetPassword{},
		// &models.Like{},
		// &models.Comment{},
		// &models.User{},
		// &models.Audit{},
	)

	server.Router = gin.Default()
	server.Router.Use(middlewares.CORSMiddleware())

	server.initializeRoutes()

}


func (server *Server) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, server.Router))
}
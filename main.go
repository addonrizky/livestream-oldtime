package main

import (
	"fmt"
	"time"
	//"github.com/addonrizky/sagaracrud/validator"
	"gopkg.in/go-playground/validator.v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	//"github.com/asumsi/livestream/modelgeneral"
	roomHttpDeliver "github.com/asumsi/livestream/room/delivery/http"
	roomRepo "github.com/asumsi/livestream/room/repository"
	roomUsecase "github.com/asumsi/livestream/room/usecase"
	"github.com/asumsi/livestream/utility"

	userHttpDeliver "github.com/asumsi/livestream/auth/user/delivery/http"
	userRepo "github.com/asumsi/livestream/auth/user/repository"
	userUsecase "github.com/asumsi/livestream/auth/user/usecase"

	googleHttpDeliver "github.com/asumsi/livestream/auth/google/delivery/http"
	googleRepo "github.com/asumsi/livestream/auth/google/repository"
	googleUsecase "github.com/asumsi/livestream/auth/google/usecase"

	streamHttpDeliver "github.com/asumsi/livestream/stream/delivery/http"
	//streamRepo "github.com/asumsi/livestream/stream/repository"
	streamUsecase "github.com/asumsi/livestream/stream/usecase"

	"github.com/labstack/echo/v4/middleware"

	"github.com/labstack/echo/v4"
	"github.com/swaggo/echo-swagger"
	_ "github.com/asumsi/livestream/docs"

	//_ "github.com/swaggo/echo-swagger/example/docs"
)

func init() {
	if utility.GetConfigBool(`debug`) {
		fmt.Println("Service RUN on DEBUG mode")
	}

	//validator.Init()
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		fmt.Println("wah error nih")
		return err
	}
	return nil
}


// @title Livestream API
// @version 1.0
// @description Livestream API enable you to live stream.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host api.inlive.app
// @BasePath /v1
func main() {
	e := echo.New()

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	e.Validator = &CustomValidator{validator: validator.New()}
	e.Use(middleware.Logger())

	dbHost := utility.GetConfigString(`db_host`)
	dbUser := utility.GetConfigString(`db_user`)
	dbPass := utility.GetConfigString(`db_pass`)
	dbName := utility.GetConfigString(`db_name`)
	connection := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", dbUser, dbPass, dbHost, dbName)
	db, err := gorm.Open(postgres.Open(connection), &gorm.Config{})
	if err != nil && utility.GetConfigBool("debug") {
		fmt.Println(err)
	}

	// migration db
	/*db.AutoMigrate(modelgeneral.Roles{})
	db.AutoMigrate(modelgeneral.Rooms{})
	db.AutoMigrate(modelgeneral.Streams{})
	db.AutoMigrate(modelgeneral.Subcriptions{})
	db.AutoMigrate(modelgeneral.Users{})
	*/

	timeoutContext := time.Duration(utility.GetConfigInt("context_timeout")) * time.Second

	// module user
	userRepo := userRepo.NewUserRepository(db)
	userUC := userUsecase.NewUserUsecase(userRepo, timeoutContext)
	userHttpDeliver.NewUserHttpHandler(e, userUC)

	// module auth google
	googleRepo := googleRepo.NewGoogleRepository(db)
	googleUC := googleUsecase.NewGoogleUsecase(googleRepo, timeoutContext)
	googleHttpDeliver.NewGoogleHttpHandler(e, googleUC, userUC)

	// module room
	roomRepo := roomRepo.NewRoomRepository(db)
	ar := roomUsecase.NewRoomUsecase(roomRepo, timeoutContext)
	roomHttpDeliver.NewRoomHttpHandler(e, ar)

	// module stream
	as := streamUsecase.NewStreamUsecase(timeoutContext)
	streamHttpDeliver.NewStreamHttpHandler(e, as)

	e.Start(utility.GetConfigString("server_address"))

}

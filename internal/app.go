package internal

import (
	"database/sql"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/joho/godotenv"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/ptflp/godecoder"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"petstore/internal/responder"
	_userController "petstore/internal/user/controller"
	_userMiddleware "petstore/internal/user/controller/middleware"
	_userRepo "petstore/internal/user/repository"
	_userUsecase "petstore/internal/user/usecase"

	_petController "petstore/internal/pet/controller"
	_petRepo "petstore/internal/pet/repository"
	_petUsecase "petstore/internal/pet/usecase"

	_orderController "petstore/internal/order/controller"
	_orderRepo "petstore/internal/order/repository"
	_orderUsecase "petstore/internal/order/usecase"

	"time"

	_ "github.com/lib/pq"
	_ "petstore/internal/docs"
)

//	@title			PetStore
//	@version		1.0
//	@description	This is implementation of PetStore API
//
// @securitydefinitions.apikey ApiKeyAuth
// @in Login
// @name Authorization
//
//	@host		localhost:8080
//	@BasePath	/

func RunApp() {
	err := godotenv.Load()

	if err != nil {
		panic(err)
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)

	db := initDB()
	resp := responder.NewResponder(godecoder.NewDecoder(), zap.NewExample())
	tokenAuth := jwtauth.New("HS256", []byte("secret"), nil, jwt.WithAcceptableSkew(30*time.Second))

	userRepo := _userRepo.NewUserRepository(db)
	auRepo := _userRepo.NewAuthRepository(db)
	userUsecase := _userUsecase.NewUserUsecase(userRepo, auRepo, tokenAuth)

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	r.Group(func(r chi.Router) {
		_userController.NewUserController(r, resp, userUsecase)
	})

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(_userMiddleware.Authenticator(resp, userUsecase))

		petRepo := _petRepo.NewPetRepository(db)
		categoryRepo := _petRepo.NewCategoryRepository(db)
		tagRepo := _petRepo.NewTagRepository(db)

		petUsecase := _petUsecase.NewPetUsecase(petRepo, categoryRepo, tagRepo)
		_petController.NewPetController(r, resp, petUsecase)
	})

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(_userMiddleware.Authenticator(resp, userUsecase))

		orderRepo := _orderRepo.NewOrderRepository(db)
		orderUsecase := _orderUsecase.NewOrderUsecase(orderRepo)

		_orderController.NewOrderController(r, resp, orderUsecase)
	})

	log.Println("server starting...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("failed run server", err)
	}
}

func initDB() *sql.DB {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	pgInfo := fmt.Sprintf("host = %s port = %s "+
		"user = %s password = %s dbname = %s sslmode = disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", pgInfo)

	if err != nil {
		log.Panicln("failed to connect to database", err)
	}

	if err := db.Ping(); err != nil {
		log.Panicln("failed to ping database", err)
	}

	return db
}

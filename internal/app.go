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
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"petstore/internal/responder"
	_userController "petstore/internal/user/controller"
	_userRepo "petstore/internal/user/repository"
	_userUsecase "petstore/internal/user/usecase"
	"time"

	_ "github.com/lib/pq"
)

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

	r.Group(func(r chi.Router) {
		userRepo := _userRepo.NewUserRepository(db)
		auRepo := _userRepo.NewAuthRepository(db)
		userUsecase := _userUsecase.NewUserUsecase(userRepo, auRepo, tokenAuth)
		_userController.NewUserController(r, resp, userUsecase)
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

package main

import (
	"api_blog/controller"
	"api_blog/exception"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
)

func main() {
	sentry.Init(sentry.ClientOptions{
		Dsn: os.Getenv("SENTRY_DSN"),
	})
	// database section
	_ = godotenv.Load()
	db, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))
	exception.PanicIfErr(err)
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxIdleTime(5 * time.Second)
	db.SetConnMaxLifetime(60 * time.Second)
	defer db.Close()

	// router section
	router := httprouter.New()
	router.PanicHandler = exception.ErrorHandler

	router.GET("/", controller.HomeIndex)

	postController := controller.NewPostController(db)
	router.GET("/api/posts", postController.Index)
	router.POST("/api/posts", postController.Store)
	router.GET("/api/posts/:postSlug", postController.Show)
	router.PUT("/api/posts/:postSlug", postController.Update)
	router.DELETE("/api/posts/:postSlug", postController.Delete)

	toolController := controller.NewToolController(db)
	router.GET("/api/tools", toolController.Index)

	fmt.Println("listening on http://localhost:8080")
	err = http.ListenAndServe(":8080", router)
	exception.PanicIfErr(err)
}

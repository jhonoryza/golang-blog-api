package main

import (
	"api_blog/controller"
	"api_blog/exception"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"time"
)

func main() {
	router := httprouter.New()
	db, err := sql.Open("pgx", "postgres://postgres:postgres@192.168.18.106:5432/blog?sslmode=disable")
	exception.PanicIfErr(err)
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxIdleTime(5 * time.Second)
	db.SetConnMaxLifetime(60 * time.Second)
	defer db.Close()

	router.PanicHandler = exception.ErrorHandler

	router.GET("/", controller.HomeIndex)

	postController := controller.NewPostController(db)
	router.GET("/api/posts", postController.Index)
	router.GET("/api/posts/:postSlug", postController.Show)

	fmt.Println("listening on http://localhost:8080")
	err = http.ListenAndServe(":8080", router)
	exception.PanicIfErr(err)
}

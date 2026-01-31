package main

import (
	"api_blog/controller"
	"api_blog/exception"
	"api_blog/middleware"
	"api_blog/wilayah"
	"database/sql"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
)

//go:embed public/*
var embeddedFiles embed.FS

func main() {
	isImport := flag.Bool("import", false, "Run wilayah data import")
	flag.Parse()

	if *isImport {
		wilayah.RunImport()
		return
	}

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

	router.GET("/health", controller.HomeIndex)

	authController := controller.NewAuthController(db)
	router.POST("/api/login", authController.Login)

	postController := controller.NewPostController(db)
	router.GET("/api/posts", postController.Index)
	router.GET("/api/posts/:postSlug", postController.Show)
	router.POST("/api/posts", middleware.AuthMiddleware(db, postController.Store))
	router.PUT("/api/posts/:postSlug", middleware.AuthMiddleware(db, postController.Update))
	router.DELETE("/api/posts/:postSlug", middleware.AuthMiddleware(db, postController.Delete))

	toolController := controller.NewToolController(db)
	router.GET("/api/tools", toolController.Index)
	router.POST("/api/tools", middleware.AuthMiddleware(db, toolController.Store))
	router.PUT("/api/tools/:toolId", middleware.AuthMiddleware(db, toolController.Update))
	router.DELETE("/api/tools/:toolId", middleware.AuthMiddleware(db, toolController.Delete))

	prayController := controller.NewPrayController(db)
	router.POST("/api/prayers", prayController.Index)
	router.GET("/api/hijri/calendar", prayController.HijriCalendar)

	timeZoneController := controller.NewTimezoneController(db)
	router.GET("/api/timezones", timeZoneController.Index)

	profileController := controller.NewProfileController(db)
	router.GET("/api/profile", profileController.Show)

	productController := controller.NewProductController(db)
	router.GET("/api/products", productController.Index)
	router.GET("/api/products/:productId", productController.Show)

	serviceController := controller.NewServiceController(db)
	router.GET("/api/services", serviceController.Index)

	quranController := controller.NewQuranController(db)
	router.GET("/api/surah", quranController.Index)
	router.GET("/api/ayah", quranController.Show)

	wilayahController := controller.NewWilayahController(db)
	router.GET("/api/provinces", wilayahController.Provinces)
	router.GET("/api/cities", wilayahController.Cities)
	router.GET("/api/districts", wilayahController.Districts)
	router.GET("/api/subdistricts", wilayahController.SubDistricts)

	sub, _ := fs.Sub(embeddedFiles, "public")
    fileServer := http.FileServer(http.FS(sub))

    router.Handler("GET", "/public/*filepath",
        http.StripPrefix("/public/", fileServer),
    )

    router.GET("/doc", serveEmbed("public/doc.html"))
    router.GET("/doc/swagger", serveEmbed("public/swagger.html"))
    router.GET("/doc/redoc", serveEmbed("public/redoc.html"))
    router.GET("/doc/stoplight", serveEmbed("public/stoplight.html"))
    router.GET("/doc/scalar", serveEmbed("public/scalar.html"))

	fmt.Println("listening on http://localhost:8080")
	err = http.ListenAndServe(":8080", router)
	exception.PanicIfErr(err)
}

func serveEmbed(path string) httprouter.Handle {
    return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
        data, err := embeddedFiles.ReadFile(path)
        if err != nil {
            http.NotFound(w, r)
            return
        }
        w.Header().Set("Content-Type", "text/html")
        w.Write(data)
    }
}

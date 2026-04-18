package main

import (
	"api_blog/cache"
	"api_blog/controller"
	"api_blog/exception"
	"api_blog/middleware"
	"api_blog/wilayah"
	"context"
	"crypto/tls"
	"database/sql"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	"github.com/redis/go-redis/v9"
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
	dbURL := os.Getenv("DATABASE_URL")
	fmt.Printf("[DATABASE] Connecting to: %s\n", dbURL)

	db, err := sql.Open("pgx", dbURL)
	exception.PanicIfErr(err)
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxIdleTime(5 * time.Second)
	db.SetConnMaxLifetime(60 * time.Second)
	defer db.Close()

	// test database connection
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		fmt.Printf("[DATABASE] Connection failed: %v\n", err)
	} else {
		fmt.Printf("[DATABASE] Connected successfully\n")
	}

	// redis section
	redisAddr := os.Getenv("REDIS_ADDR")
	redisUsername := os.Getenv("REDIS_USERNAME")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	fmt.Printf("[REDIS] Connecting to: %s\n", redisAddr)
	fmt.Printf("[REDIS] Username: %s\n", redisUsername)

	redisClient := redis.NewClient(&redis.Options{
		Addr:            redisAddr,
		Username:        redisUsername,
		Password:        redisPassword,
		DB:              0,
		DialTimeout:     15 * time.Second,
		ReadTimeout:     15 * time.Second,
		WriteTimeout:    15 * time.Second,
		PoolSize:        10,
		MinIdleConns:    2,
		DisableIdentity: true,
		TLSConfig:       &tls.Config{MinVersion: tls.VersionTLS12},
	})
	cacheManager := cache.NewCacheManager(redisClient, 24*time.Hour)
	defer redisClient.Close()

	// test redis connection
	ctx, cancel = context.WithTimeout(context.Background(), 15*time.Second)
	pingResult, err := redisClient.Ping(ctx).Result()
	cancel()
	if err != nil {
		fmt.Printf("[REDIS] Connection failed: %v\n", err)
	} else {
		fmt.Printf("[REDIS] Connected successfully: %s\n", pingResult)
	}

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

	quranController := controller.NewQuranController(db, cacheManager)
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
	err = http.ListenAndServe(":8080", GlobalLogger(router))
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

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func GlobalLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &responseWriter{w, http.StatusOK}

		next.ServeHTTP(rw, r)

		log.Printf("[REQUEST] %s %s -> %d (%v)",
			r.Method,
			r.URL.Path,
			rw.statusCode,
			time.Since(start),
		)
	})
}
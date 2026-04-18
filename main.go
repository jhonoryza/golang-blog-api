package main

import (
	"api_blog/delivery/http/controller"
	"api_blog/delivery/http/middleware"
	"api_blog/infrastructure/exception"
	"api_blog/infrastructure/cache"
	"api_blog/infrastructure/database"
	"api_blog/infrastructure/persistence"
	"api_blog/usecase/auth"
	"api_blog/usecase/post"
	"api_blog/usecase/product"
	"api_blog/usecase/quran"
	"api_blog/usecase/service"
	"api_blog/usecase/tool"
	"api_blog/usecase/user"
	wilayahUsecase "api_blog/usecase/wilayah"
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

var embeddedFiles embed.FS

func main() {
	isImport := flag.Bool("import", false, "Run wilayah data import")
	flag.Parse()

	if *isImport {
		runWilayahImport()
		return
	}

	sentry.Init(sentry.ClientOptions{
		Dsn: os.Getenv("SENTRY_DSN"),
	})
	_ = godotenv.Load()

	db := initDatabase()
	defer db.Close()

	redisClient := initRedis()
	cacheManager := cache.NewCacheManager(redisClient, 24*time.Hour)
	defer redisClient.Close()

	repos := initRepositories(db)
	ucs := initUsecases(repos)
	ctrls := initControllers(ucs, cacheManager, db)

	router := httprouter.New()
	router.PanicHandler = exception.ErrorHandler
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[REQUEST] %s %s", r.Method, r.URL.Path)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"code":404,"message":"Not Found","data":null}`))
	})

	setupRoutes(router, ctrls, ucs)

	sub, _ := fs.Sub(embeddedFiles, "public")
	router.Handler("GET", "/public/*filepath", http.StripPrefix("/public/", http.FileServer(http.FS(sub))))

	router.GET("/doc", serveEmbed("public/doc.html"))
	router.GET("/doc/swagger", serveEmbed("public/swagger.html"))
	router.GET("/doc/redoc", serveEmbed("public/redoc.html"))
	router.GET("/doc/stoplight", serveEmbed("public/stoplight.html"))
	router.GET("/doc/scalar", serveEmbed("public/scalar.html"))

	router.GET("/", controller.HomeIndex)
	router.GET("/health", controller.HomeIndex)

	fmt.Println("listening on http://localhost:8080")
	err := http.ListenAndServe(":8080", GlobalLogger(router))
	exception.PanicIfErr(err)
}

func initDatabase() *sql.DB {
	dbURL := os.Getenv("DATABASE_URL")
	fmt.Printf("[DATABASE] Connecting to: %s\n", dbURL)

	db, err := database.NewPostgresConnection(dbURL, database.Config{
		MaxIdleConns:    10,
		MaxOpenConns:    100,
		ConnMaxIdleTime: 5 * time.Second,
		ConnMaxLifetime: 60 * time.Second,
	})
	if err != nil {
		fmt.Printf("[DATABASE] Connection failed: %v\n", err)
	}
	return db
}

func initRedis() *redis.Client {
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

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	pingResult, err := redisClient.Ping(ctx).Result()
	cancel()
	if err != nil {
		fmt.Printf("[REDIS] Connection failed: %v\n", err)
	} else {
		fmt.Printf("[REDIS] Connected successfully: %s\n", pingResult)
	}

	return redisClient
}

type Repositories struct {
	Post    *persistence.PostRepository
	Tool    *persistence.ToolRepository
	User    *persistence.UserRepository
	Quran   *persistence.QuranRepository
	Wilayah *persistence.WilayahRepository
}

func initRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		Post:    persistence.NewPostRepository(db),
		Tool:    persistence.NewToolRepository(db),
		User:    persistence.NewUserRepository(db),
		Quran:   persistence.NewQuranRepository(db),
		Wilayah: persistence.NewWilayahRepository(db),
	}
}

type Usecases struct {
	Post    *post.PostUsecase
	Tool    *tool.ToolUsecase
	User    *user.UserUsecase
	Auth    *auth.AuthUsecase
	Quran   *quran.QuranUsecase
	Wilayah *wilayahUsecase.WilayahUsecase
	Product *product.ProductUsecase
	Service *service.ServiceUsecase
}

func initUsecases(repos *Repositories) *Usecases {
	return &Usecases{
		Post:    post.NewPostUsecase(repos.Post),
		Tool:    tool.NewToolUsecase(repos.Tool),
		User:    user.NewUserUsecase(repos.User),
		Auth:    auth.NewAuthUsecase(repos.User),
		Quran:   quran.NewQuranUsecase(repos.Quran),
		Wilayah: wilayahUsecase.NewWilayahUsecase(repos.Wilayah),
		Product: product.NewProductUsecase(),
		Service: service.NewServiceUsecase(),
	}
}

type Controllers struct {
	Auth     *controller.AuthController
	Post     *controller.PostController
	Tool     *controller.ToolController
	Pray     *controller.PrayController
	Profile  *controller.ProfileController
	Product  *controller.ProductController
	Service  *controller.ServiceController
	Quran    *controller.QuranController
	Wilayah  *controller.WilayahController
	Timezone *controller.TimezoneController
}

func initControllers(ucs *Usecases, cacheMgr *cache.CacheManager, db *sql.DB) *Controllers {
	return &Controllers{
		Auth:     controller.NewAuthController(ucs.Auth, db),
		Post:     controller.NewPostController(ucs.Post),
		Tool:     controller.NewToolController(ucs.Tool, db),
		Pray:     controller.NewPrayController(),
		Profile:  controller.NewProfileController(),
		Product:  controller.NewProductController(ucs.Product),
		Service:  controller.NewServiceController(),
		Quran:    controller.NewQuranController(ucs.Quran, cacheMgr),
		Wilayah:  controller.NewWilayahController(ucs.Wilayah, db),
		Timezone: controller.NewTimezoneController(db),
	}
}

func setupRoutes(router *httprouter.Router, ctrls *Controllers, ucs *Usecases) {
	router.POST("/api/login", ctrls.Auth.Login)

	router.GET("/api/posts", ctrls.Post.Index)
	router.GET("/api/posts/:postSlug", ctrls.Post.Show)
	router.POST("/api/posts", middleware.AuthMiddleware(ucs.User, ctrls.Post.Store))
	router.PUT("/api/posts/:postSlug", middleware.AuthMiddleware(ucs.User, ctrls.Post.Update))
	router.DELETE("/api/posts/:postSlug", middleware.AuthMiddleware(ucs.User, ctrls.Post.Delete))

	router.GET("/api/tools", ctrls.Tool.Index)
	router.POST("/api/tools", middleware.AuthMiddleware(ucs.User, ctrls.Tool.Store))
	router.PUT("/api/tools/:toolId", middleware.AuthMiddleware(ucs.User, ctrls.Tool.Update))
	router.DELETE("/api/tools/:toolId", middleware.AuthMiddleware(ucs.User, ctrls.Tool.Delete))

	router.POST("/api/prayers", ctrls.Pray.Index)
	router.GET("/api/hijri/calendar", ctrls.Pray.HijriCalendar)

	router.GET("/api/timezones", ctrls.Timezone.Index)

	router.GET("/api/profile", ctrls.Profile.Show)

	router.GET("/api/products", ctrls.Product.Index)
	router.GET("/api/products/:productId", ctrls.Product.Show)

	router.GET("/api/services", ctrls.Service.Index)

	router.GET("/api/surah", ctrls.Quran.Index)
	router.GET("/api/ayah", ctrls.Quran.Show)

	router.GET("/api/provinces", ctrls.Wilayah.Provinces)
	router.GET("/api/cities", ctrls.Wilayah.Cities)
	router.GET("/api/districts", ctrls.Wilayah.Districts)
	router.GET("/api/subdistricts", ctrls.Wilayah.SubDistricts)
}

func runWilayahImport() {
	fmt.Println("Wilayah import not yet migrated to new structure")
}

var fileServer = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
})

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
		log.Printf("[REQUEST] %s %s -> %d (%v)", r.Method, r.URL.Path, rw.statusCode, time.Since(start))
	})
}
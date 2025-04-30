package controller

import (
	"api_blog/exception"
	"api_blog/response"
	"database/sql"
	"net/http"
	"os"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type ProductController struct {
	DB *sql.DB
}

func NewProductController(db *sql.DB) *ProductController {
	return &ProductController{
		DB: db,
	}
}

type Product struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Image        string    `json:"image"`
	Category     string    `json:"category"`
	Features     []string  `json:"features"`
	IsOpenSource bool      `json:"is_opensource"`
	Pricing      []Pricing `json:"pricing"`
	TechStack    []string  `json:"techStack"`
	Docs         string    `json:"docs"`
	PolicyURL    string    `json:"policy_url,omitempty"`
}

type Pricing struct {
	Plan  string `json:"plan"`
	Price string `json:"price"`
}

func (c *ProductController) Index(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	products := getProducts()
	resp := response.ApiResponse{
		Code:    http.StatusOK,
		Message: "OK",
		Data:    products,
	}

	resp.ToJson(w)
}

func (c *ProductController) Show(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	productId, err := strconv.Atoi(p.ByName("productId"))
	if err != nil {
		exception.BadRequestError(w, r, "BAD REQUEST")
		return
	}
	products := getProducts()
	var product Product
	for _, p := range products {
		if p.ID == productId {
			product = p
		}
	}
	if product.ID == 0 {
		exception.BadRequestError(w, r, "BAD REQUEST")
		return
	}

	resp := response.ApiResponse{
		Code:    http.StatusOK,
		Message: "OK",
		Data:    product,
	}

	resp.ToJson(w)
}

func getProducts() []Product {
	var baseImageUrl = os.Getenv("IMAGE_BASE_URL") + "/blog/image"
	products := []Product{
		{
			ID:          1,
			Name:        "Tinkering",
			Description: "A rapid prototyping tool for Laravel developers. easily run PHP code, run PHP, run a Laravel Eloquent query on your database.",
			Image:       baseImageUrl + "/tinkering.png",
			Category:    "Desktop apps",
			Features: []string{
				"Intuitive user interfaces",
				"Customizable php version path",
				"Customizable laravel path",
				"Running artisan commands",
				"Vertical or horizontal mode",
				"Supports mac, windows and linux",
			},
			IsOpenSource: true,
			Pricing: []Pricing{
				{"Pro", "free"},
			},
			TechStack: []string{"Rust", "Vue"},
			Docs:      "https://dev.to/jhonoryza/tinkering-app-42f0",
		},
		{
			ID:          2,
			Name:        "Reprox",
			Description: "A tunneling for http and tcp, easily expose your private network to the public",
			Image:       baseImageUrl + "/reprox.png",
			Category:    "CLI apps",
			Features: []string{
				"Self hosted",
				"Customizable with your own domain",
				"Tunneling http port",
				"Tunneling tcp port",
				"Supports mac, windows and linux",
			},
			IsOpenSource: true,
			Pricing: []Pricing{
				{"Pro", "free"},
			},
			TechStack: []string{"Golang"},
			Docs:      "https://github.com/jhonoryza/reprox",
		},
		{
			ID:          3,
			Name:        "Logdesk",
			Description: "A debugger tools for Laravel developers, improves your debugging experience by sending laravel logs to this app.",
			Image:       baseImageUrl + "/logdesk2.png",
			Category:    "Desktop apps",
			Features: []string{
				"Logging laravel log in realtime",
				"Toggle on top",
				"Supports mac, windows and linux",
			},
			IsOpenSource: true,
			Pricing: []Pricing{
				{"Pro", "free"},
			},
			TechStack: []string{"Node.js", "Vue", "Electron"},
			Docs:      "https://github.com/jhonoryza/logdesk",
		},
		{
			ID:          4,
			Name:        "Alquran",
			Description: "A quran apps, no internet needed, offline mode, no ads, free, qibla direction and qrscanner.",
			Image:       baseImageUrl + "/quran.png",
			Category:    "Android apps",
			Features: []string{
				"Dark mode",
				"Indonesian translation",
				"Qibla direction",
				"QR Scanner",
				"offline mode",
			},
			IsOpenSource: true,
			Pricing: []Pricing{
				{"Pro", "free"},
			},
			TechStack: []string{"Dart", "Flutter"},
			Docs:      "https://play.google.com/store/apps/details?id=com.labkita.baca_alquran",
			PolicyURL: "https://labkita.my.id/alquran",
		},
	}
	return products
}

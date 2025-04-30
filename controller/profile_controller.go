package controller

import (
	"api_blog/response"
	"database/sql"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
)

type ProfileController struct {
	DB *sql.DB
}

func NewProfileController(db *sql.DB) *ProfileController {
	return &ProfileController{
		DB: db,
	}
}

type Project struct {
	ID          int      `json:"id"`
	Year        int      `json:"year"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Image       string   `json:"image"`
	State       string   `json:"state"`
	Stacks      []string `json:"stacks"`
	Link        string   `json:"link"`
}

var baseImageUrl = os.Getenv("IMAGE_BASE_URL") + "/blog/image"

var projects = []Project{
	{
		ID:          1,
		Year:        2024,
		Title:       "Temubisnis Web Platform",
		Description: "Provide local or umkm product information",
		Image:       baseImageUrl + "/temubisnis.png",
		State:       "production",
		Stacks: []string{
			"PHP 8.2", "Laravel", "Livewire", "Volt", "Filament", "Alpinejs", "Tailwind", "Docker", "Minio", "MySQL", "Redis",
		},
		Link: "https://temubisnis.id",
	},
	{
		ID:          2,
		Year:        2024,
		Title:       "Halal Product Web Platform",
		Description: "Provide halal product information",
		Image:       baseImageUrl + "/ehh.png",
		State:       "development",
		Stacks: []string{
			"PHP 8.2", "Laravel", "Inertia", "Vue", "Filament", "Alpinejs", "Tailwind", "Docker", "Minio", "MySQL", "Redis",
		},
		Link: "https://ehhdev.com",
	},
	{
		ID:          3,
		Year:        2019,
		Title:       "Service Aja Web Platform",
		Description: "Provide company internal ticket tracking",
		Image:       baseImageUrl + "/serviceaja.png",
		State:       "production",
		Stacks: []string{
			"PHP 7.2", "Laravel", "HTML", "CSS", "Javascript", "MySQL", "Redis", "Nginx",
		},
		Link: "https://www.serviceaja.com",
	},
	{
		ID:          4,
		Year:        2020,
		Title:       "RS Mata SMEC Web Platform",
		Description: "Provide hospital information & booking system",
		Image:       baseImageUrl + "/rssmec.png",
		State:       "production",
		Stacks: []string{
			"PHP 7.2", "Laravel", "HTML", "CSS", "Javascript", "MySQL", "Redis", "Nginx",
		},
		Link: "https://rsmatasmec.com",
	},
	{
		ID:          5,
		Year:        2020,
		Title:       "MOST Web Platform",
		Description: "Provide halal product information",
		Image:       baseImageUrl + "/most.png",
		State:       "production",
		Stacks: []string{
			"PHP 7.4", "Laravel", "Nuxt", "Livewire", "Alpinejs", "Docker", "Minio", "MySQL", "Redis",
		},
		Link: "https://most.co.id/",
	},
	{
		ID:          6,
		Year:        2021,
		Title:       "IHWG Web Platform",
		Description: "Provide organization information",
		Image:       baseImageUrl + "/ihwg.png",
		State:       "production",
		Stacks: []string{
			"PHP 7.4", "Laravel", "Nuxt", "Livewire", "Alpinejs", "Docker", "Minio", "MySQL", "Redis",
		},
		Link: "https://ihwg.or.id/",
	},
	{
		ID:          7,
		Year:        2021,
		Title:       "Mandiri Sekuritas Web Platform",
		Description: "Provide company information",
		Image:       baseImageUrl + "/mandiri.png",
		State:       "production",
		Stacks: []string{
			"PHP 7.4", "Laravel", "Nuxt", "Livewire", "Alpinejs", "Docker", "Minio", "MySQL", "Redis",
		},
		Link: "https://www.mandirisekuritas.co.id/",
	},
	{
		ID:          8,
		Year:        2021,
		Title:       "SCTV Web Platform",
		Description: "Provide tv program information",
		Image:       baseImageUrl + "/sctv.png",
		State:       "production",
		Stacks: []string{
			"PHP 8.0", "Laravel", "Nuxt", "Livewire", "Alpinejs", "Docker", "Minio", "MySQL", "Redis",
		},
		Link: "https://sctv.co.id/",
	},
	{
		ID:          9,
		Year:        2021,
		Title:       "Suitmedia Web Platform",
		Description: "Provide company information",
		Image:       baseImageUrl + "/oldsuit.png",
		State:       "production",
		Stacks: []string{
			"PHP 8.0", "Laravel", "Nuxt", "Livewire", "Alpinejs", "Docker", "Minio", "MySQL", "Redis",
		},
		Link: "https://archive.suitmedia.com/",
	},
	{
		ID:          10,
		Year:        2022,
		Title:       "Eiger Care Web Platform",
		Description: "Provide integration of multiple marketplace, provide order, warehouse and pos management",
		Image:       baseImageUrl + "/care.png",
		State:       "production",
		Stacks: []string{
			"PHP 8.0", "Laravel", "Nuxt", "Livewire", "Alpinejs", "Docker", "Minio", "MySQL", "Redis",
		},
		Link: "https://care.eigerindo.co.id/",
	},
	{
		ID:          11,
		Year:        2022,
		Title:       "Eiger Adventure Web Platform",
		Description: "Provide loyalty member system like vouchers, point, discount.",
		Image:       baseImageUrl + "/careadv.png",
		State:       "production",
		Stacks: []string{
			"PHP 8.0", "Laravel", "Nuxt", "Livewire", "Alpinejs", "Docker", "Minio", "MySQL", "Redis",
		},
		Link: "https://care.eigerindo.co.id/",
	},
	{
		ID:          12,
		Year:        2023,
		Title:       "KSEI Web Platform",
		Description: "Provide company & stock information.",
		Image:       baseImageUrl + "/ksei.png",
		State:       "production",
		Stacks: []string{
			"PHP 8.2", "Laravel", "Nuxt", "Livewire", "Alpinejs", "Docker", "Minio", "MySQL", "Redis",
		},
		Link: "#",
	},
	{
		ID:          13,
		Year:        2022,
		Title:       "Danone Tracenility Web Platform",
		Description: "Provide tracking material products & internal task management.",
		Image:       baseImageUrl + "/tracey.png",
		State:       "production",
		Stacks: []string{
			"PHP 7.4", "Laravel", "Nuxt", "Livewire", "Alpinejs", "Docker", "Minio", "MySQL", "Redis",
		},
		Link: "#",
	},
	{
		ID:          14,
		Year:        2024,
		Title:       "GPI-G Web Platform",
		Description: "Provide company information.",
		Image:       baseImageUrl + "/gpi.png",
		State:       "production",
		Stacks: []string{
			"PHP 8.2", "Laravel", "Nuxt", "Livewire", "Alpinejs", "Docker", "Minio", "MySQL", "Redis",
		},
		Link: "https://www.gpi-g.com/",
	},
	{
		ID:          15,
		Year:        2024,
		Title:       "Teman Muslim Web Platform",
		Description: "Provide muslim task information, like quran, shalat, qiblat etc.",
		Image:       baseImageUrl + "/temanmuslim.png",
		State:       "production",
		Stacks: []string{
			"PHP 8.2", "Laravel", "Nuxt", "Livewire", "Alpinejs", "Docker", "Minio", "MySQL", "Redis",
		},
		Link: "https://temanmuslim.com/",
	},
	{
		ID:          16,
		Year:        2023,
		Title:       "Alifmart E-commerce Platform",
		Description: "Provide maintenance online shopping system.",
		Image:       baseImageUrl + "/alifmart.png",
		State:       "production",
		Stacks: []string{
			"PHP 7.4", "Laravel", "Nuxt", "Livewire", "Alpinejs", "Docker", "Minio", "MySQL", "Redis",
		},
		Link: "https://alifmart.online/",
	},
	{
		ID:          17,
		Year:        2024,
		Title:       "Solitaire Game Platform",
		Description: "Provide backend api for telegram games.",
		Image:       baseImageUrl + "/solitaire.jpg",
		State:       "production",
		Stacks: []string{
			"NestJs", "Redis", "AWS Gamelift",
		},
		Link: "#",
	},
	{
		ID:          18,
		Year:        2024,
		Title:       "Uno Game Platform",
		Description: "Provide backend api for telegram games.",
		Image:       baseImageUrl + "/uno.jpg",
		State:       "production",
		Stacks: []string{
			"NestJs", "Redis", "AWS Gamelift",
		},
		Link: "#",
	},
	{
		ID:          19,
		Year:        2024,
		Title:       "General Service Game Platform",
		Description: "Provide backend general api for telegram games.",
		Image:       baseImageUrl + "/uno.jpg",
		State:       "production",
		Stacks: []string{
			"NestJs", "Redis", "AWS Gamelift",
		},
		Link: "#",
	},
	{
		ID:          20,
		Year:        2024,
		Title:       "Waifuclash Game Platform",
		Description: "Maintenance backend api for telegram games.",
		Image:       baseImageUrl + "/waifuclash.jpg",
		State:       "production",
		Stacks: []string{
			"NestJs", "PostgreSQL", "Redis", "AWS Gamelift",
		},
		Link: "#",
	},
	{
		ID:          21,
		Year:        2024,
		Title:       "Kittysolitaire Service Game Platform",
		Description: "Maintenance backend api for telegram games.",
		Image:       baseImageUrl + "/kittysolitaire.jpg",
		State:       "production",
		Stacks: []string{
			"NestJs", "PostgreSQL", "Redis", "AWS Gamelift",
		},
		Link: "#",
	},
	{
		ID:          22,
		Year:        2024,
		Title:       "Morinaga parenthing Web Platform",
		Description: "Maintenance backend api.",
		Image:       baseImageUrl + "/morigro.png",
		State:       "production",
		Stacks: []string{
			"NestJs", "PostgreSQL", "Redis", "AWS Gamelift",
		},
		Link: "#",
	},
	{
		ID:          23,
		Year:        2025,
		Title:       "Treasury Gamification Web Platform",
		Description: "Provide prize drawing system.",
		Image:       baseImageUrl + "/treasury.png",
		State:       "production",
		Stacks: []string{
			"PHP 8.2",
			"Laravel",
			"PostgreSQL",
			"Redis",
			"Nginx",
		},
		Link: "#",
	},
}

func (c *ProfileController) Show(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	exp := strconv.Itoa(time.Now().Year() - 2019)
	var profile = map[string]any{
		"id":    1,
		"name":  "Fajar",
		"role":  "Software Engineer",
		"image": baseImageUrl + "/fajarsp.jpg",
		"skills": []string{
			"AWS", "Linode", "Oracle", "Alibaba", "GCP", "Docker",
			"PHP", "Laravel", "Go", "Nestjs", "Nodejs", "Javascript",
			"Flutter", "Rust", "Html", "CSS",
			"PostgreSQL", "SQLite", "MySQL", "MariaDB",
		},
		"bio": "Fajar is a full-stack developer with over " + exp + " years of experience in building scalable web applications and has a strong background in cloud technologies.",
		"experience": []string{
			"System Information", "Ticket tracking", "Booking system", "Loyalty system",
			"Stock management", "E-commerce platform", "etc",
		},
		"education": "Bachelor Degree from Physics Department, Padjadjaran University",
		"socialMedia": map[string]string{
			"github":   "https://github.com/jhonoryza",
			"youtube":  "https://www.youtube.com/@labkita",
			"linkedin": "https://www.linkedin.com/in/fajar-sidik-priatna-8b31a788",
		},
		"projectsCount": len(projects),
		"projects":      projects,
	}

	resp := response.ApiResponse{
		Code:    http.StatusOK,
		Message: "OK",
		Data:    profile,
	}

	resp.ToJson(w)
}

package wilayah

import (
	"api_blog/exception"
	"bufio"
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/schollz/progressbar/v3"
)

func RunImport() {
	_ = godotenv.Load()
	db, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))
	exception.PanicIfErr(err)
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxIdleTime(5 * time.Second)
	db.SetConnMaxLifetime(60 * time.Second)
	defer db.Close()

	// download data.csv disini https://kodewilayah.id
	file, err := os.Open("./wilayah/data.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	total := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			total++
		}
	}

	// Step 2: proses import ulang
	file.Seek(0, 0)
	reader := csv.NewReader(bufio.NewReader(file))
	reader.Comma = ','

	bar := progressbar.Default(int64(total))

	// Counter kategori
	var (
		provinceCount    int
		cityCount        int
		districtCount    int
		subdistrictCount int
	)

	for {
		record, err := reader.Read()
		if err != nil {
			break
		}

		if len(record) < 2 {
			continue
		}

		id := strings.TrimSpace(record[0])
		name := strings.TrimSpace(record[1])
		parts := strings.Split(id, ".")

		switch len(parts) {
		case 1:
			insertProvince(db, id, name)
			provinceCount++
		case 2:
			insertCity(db, id, name, parts[0])
			cityCount++
		case 3:
			insertDistrict(db, id, name, fmt.Sprintf("%s.%s", parts[0], parts[1]))
			districtCount++
		case 4:
			insertSubdistrict(db, id, name, fmt.Sprintf("%s.%s.%s", parts[0], parts[1], parts[2]))
			subdistrictCount++
		}

		bar.Add(1)
	}

	fmt.Println("\n✅ Data import selesai.")
	fmt.Printf("📌 Total Provinsi    : %d\n", provinceCount)
	fmt.Printf("📌 Total Kota/Kab    : %d\n", cityCount)
	fmt.Printf("📌 Total Kecamatan   : %d\n", districtCount)
	fmt.Printf("📌 Total Kelurahan   : %d\n", subdistrictCount)
}

func insertProvince(db *sql.DB, id, name string) {
	_, err := db.Exec(`
        INSERT INTO provinces (id, name)
        VALUES ($1, $2)
        ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name
    `, id, name)
	if err != nil {
		log.Printf("Gagal insert/update province %s: %v", name, err)
	}
}

func insertCity(db *sql.DB, id, name, provinceID string) {
	_, err := db.Exec(`
        INSERT INTO cities (id, name, province_id)
        VALUES ($1, $2, $3)
        ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name, province_id = EXCLUDED.province_id
    `, id, name, provinceID)
	if err != nil {
		log.Printf("Gagal insert/update city %s: %v", name, err)
	}
}

func insertDistrict(db *sql.DB, id, name, cityID string) {
	_, err := db.Exec(`
        INSERT INTO districts (id, name, city_id)
        VALUES ($1, $2, $3)
        ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name, city_id = EXCLUDED.city_id
    `, id, name, cityID)
	if err != nil {
		log.Printf("Gagal insert/update district %s: %v", name, err)
	}
}

func insertSubdistrict(db *sql.DB, id, name, districtID string) {
	_, err := db.Exec(`
        INSERT INTO subdistricts (id, name, district_id)
        VALUES ($1, $2, $3)
        ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name, district_id = EXCLUDED.district_id
    `, id, name, districtID)
	if err != nil {
		log.Printf("Gagal insert/update subdistrict %s: %v", name, err)
	}
}

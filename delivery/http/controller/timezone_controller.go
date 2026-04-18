package controller

import (
	"api_blog/delivery/http/response"
	"database/sql"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"github.com/julienschmidt/httprouter"
)

var timezoneAlias = map[string]string{
	"+VERSION":              "", // Skip
	"Africa/Asmera":         "Africa/Asmara",
	"Africa/Timbuktu":       "Africa/Bamako",
	"America/Atka":          "America/Adak",
	"America/Coral_Harbour": "America/Atikokan",
	"America/Ensenada":      "America/Tijuana",
	"America/Fort_Wayne":    "America/Indiana/Indianapolis",
	"America/Godthab":       "America/Nuuk",
	"America/Indianapolis":  "America/Indiana/Indianapolis",
	"America/Jujuy":         "America/Argentina/Jujuy",
	"America/Knox_IN":       "America/Indiana/Knox",
	"America/Louisville":    "America/Kentucky/Louisville",
	"America/Mendoza":       "America/Argentina/Mendoza",
	"America/Porto_Acre":    "America/Rio_Branco",
	"America/Rosario":       "America/Argentina/Cordoba",
	"America/Santa_Isabel":  "America/Tijuana",
	"America/Shiprock":      "America/Denver",
	"Asia/Ashkhabad":        "Asia/Ashgabat",
	"Asia/Calcutta":         "Asia/Kolkata",
	"Asia/Chongqing":        "Asia/Shanghai",
	"Asia/Chungking":        "Asia/Shanghai",
	"Asia/Dacca":            "Asia/Dhaka",
	"Asia/Istanbul":         "Europe/Istanbul",
	"Asia/Katmandu":         "Asia/Kathmandu",
	"Asia/Saigon":           "Asia/Ho_Chi_Minh",
	"Asia/Tel_Aviv":         "Asia/Jerusalem",
	"Asia/Thimbu":           "Asia/Thimphu",
	"Asia/Ujung_Pandang":    "Asia/Makassar",
	"Australia/ACT":         "Australia/Sydney",
	"Australia/Canberra":    "Australia/Sydney",
	"Australia/LHI":         "Australia/Lord_Howe",
	"Australia/NSW":         "Australia/Sydney",
	"Australia/Queensland":  "Australia/Brisbane",
	"Australia/South":       "Australia/Adelaide",
	"Australia/Tasmania":    "Australia/Hobart",
	"Australia/Victoria":    "Australia/Melbourne",
	"Australia/West":        "Australia/Perth",
	"Europe/Andorra":        "Europe/Madrid",
	"Europe/Athens":         "Europe/Istanbul",
	"Europe/Berlin":         "Europe/Budapest",
	"Europe/Brussels":       "Europe/Paris",
	"Europe/Bucharest":      "Europe/Istanbul",
	"Europe/Copenhagen":     "Europe/Oslo",
	"Europe/Helsinki":       "Europe/Stockholm",
	"Europe/Istanbul":       "Europe/Istanbul",
	"Europe/Lisbon":         "Europe/London",
	"Europe/Luxembourg":     "Europe/Paris",
	"Europe/Madrid":         "Europe/Paris",
	"Europe/Minsk":          "Europe/Moscow",
	"Europe/Moscow":         "Europe/Moscow",
	"Europe/Oslo":           "Europe/Oslo",
	"Europe/Paris":          "Europe/Paris",
	"Europe/Prague":         "Europe/Prague",
	"Europe/Rome":           "Europe/Rome",
	"Europe/Sofia":          "Europe/Sofia",
	"Europe/Stockholm":      "Europe/Stockholm",
	"Europe/Vienna":         "Europe/Vienna",
	"Europe/Warsaw":         "Europe/Warsaw",
	"Europe/Zagreb":         "Europe/Zagreb",
	"Indian/Chagos":         "Indian/Indian_Ocean",
	"Indian/Christmas":      "Indian/Christmas",
	"Indian/Cocos":          "Indian/Cocos",
	"Indian/Maldives":       "Indian/Maldives",
	"Indian/Reunion":        "Indian/Reunion",
	"Pacific/Apia":          "Pacific/Apia",
	"Pacific/Auckland":      "Pacific/Auckland",
	"Pacific/Chatham":       "Pacific/Chatham",
	"Pacific/Easter":        "Pacific/Easter",
	"Pacific/Enderbury":     "Pacific/Enderbury",
	"Pacific/Fakaofo":       "Pacific/Fakaofo",
	"Pacific/Guam":          "Pacific/Guam",
	"Pacific/Honolulu":      "Pacific/Honolulu",
	"Pacific/Johnston":      "Pacific/Johnston",
	"Pacific/Kiritimati":    "Pacific/Kiritimati",
	"Pacific/Kosrae":        "Pacific/Kosrae",
	"Pacific/Nauru":         "Pacific/Nauru",
	"Pacific/Niue":          "Pacific/Niue",
	"Pacific/Pago_Pago":     "Pacific/Pago_Pago",
	"Pacific/Palau":         "Pacific/Palau",
	"Pacific/Pitcairn":      "Pacific/Pitcairn",
	"Pacific/Ponape":        "Pacific/Ponape",
	"Pacific/Port_Moresby":  "Pacific/Port_Moresby",
	"Pacific/Rarotonga":     "Pacific/Rarotonga",
	"Pacific/Tarawa":        "Pacific/Tarawa",
	"Pacific/Tongatapu":     "Pacific/Tongatapu",
	"Pacific/Truk":          "Pacific/Truk",
	"Pacific/Wake":          "Pacific/Wake",
	"Pacific/Wallis":        "Pacific/Wallis",
}

type TimezoneController struct {
	DB *sql.DB
}

func NewTimezoneController(db *sql.DB) *TimezoneController {
	return &TimezoneController{
		DB: db,
	}
}

func (c *TimezoneController) Index(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	timezones, err := listTimezones("/usr/share/zoneinfo")
	if err != nil {
		panic(err)
	}

	// for _, tz := range timezones {
	// 	fmt.Println(tz)
	// }

	resp := response.ApiResponse{
		Code:    http.StatusOK,
		Message: "OK",
		Data:    timezones,
	}

	resp.ToJson(w)
}

func listTimezones(zoneinfoPath string) ([]string, error) {
	var timezones []string

	resolvedPath, err := resolveSymlinkRecursive(zoneinfoPath)
	if err != nil {
		return nil, err
	}

	err = filepath.WalkDir(resolvedPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if strings.HasPrefix(path, resolvedPath+"/posix") ||
			strings.HasPrefix(path, resolvedPath+"/right") ||
			strings.Contains(path, "leap-seconds.list") ||
			strings.Contains(path, "zone.tab") ||
			strings.Contains(path, "zone1970.tab") ||
			strings.Contains(path, "tzdata.zi") ||
			strings.Contains(path, "localtime") ||
			strings.Contains(path, "Factory") ||
			strings.Contains(path, "+VERSION") {
			return nil
		}

		relPath := strings.TrimPrefix(path, resolvedPath+"/")

		if alias, ok := timezoneAlias[relPath]; ok && alias != "" {
			if slices.Contains(timezones, alias) {
				return nil
			}
			timezones = append(timezones, relPath)
			return nil
		}

		timezones = append(timezones, relPath)
		return nil
	})

	if err != nil {
		return nil, err
	}

	sort.Strings(timezones)
	return timezones, nil
}

func resolveSymlinkRecursive(path string) (string, error) {
	for {
		info, err := os.Lstat(path)
		if err != nil {
			return "", err
		}

		// Kalau bukan symlink, selesai
		if info.Mode()&os.ModeSymlink == 0 {
			return path, nil
		}

		// Kalau symlink, resolve targetnya
		target, err := os.Readlink(path)
		if err != nil {
			return "", err
		}

		// Kalau target relative, gabung dengan parent
		if !filepath.IsAbs(target) {
			path = filepath.Join(filepath.Dir(path), target)
		} else {
			path = target
		}
	}
}

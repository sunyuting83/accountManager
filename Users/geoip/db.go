package geoip

import (
	"colaAPI/Users/utils"
	"os"
	"strings"

	"github.com/oschwald/geoip2-golang"
)

func LoadGeoFile() (*geoip2.Reader, error) {
	CurrentPath, _ := utils.GetCurrentPath()
	dbPath := strings.Join([]string{CurrentPath, "GeoIP"}, "/")
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		os.MkdirAll(dbPath, 0755)
	}
	dbFile := strings.Join([]string{dbPath, "GeoLite2-City.mmdb"}, "/")
	// fmt.Println(dbFile)
	db, err := geoip2.Open(dbFile)
	if err != nil {
		return nil, err
	}
	return db, nil
}

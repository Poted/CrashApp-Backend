package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type GeoIPResponse struct {
	IP          string  `json:"ip"`
	CountryCode string  `json:"country_code"`
	CountryName string  `json:"country_name"`
	RegionCode  string  `json:"region_code"`
	RegionName  string  `json:"region_name"`
	City        string  `json:"city"`
	ZipCode     string  `json:"zip_code"`
	TimeZone    string  `json:"time_zone"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	MetroCode   int     `json:"metro_code"`
}

func getIP(c *fiber.Ctx) error {

	geoIP := GetIP(c)

	return c.JSON(map[string]interface{}{
		"ip":           geoIP.IP,
		"country_code": geoIP.CountryCode,
		"country_name": geoIP.CountryName,
		"region_code":  geoIP.RegionCode,
		"region_name":  geoIP.RegionName,
		"city":         geoIP.City,
		"zip_code":     geoIP.ZipCode,
		"time_zone":    geoIP.TimeZone,
		"latitude":     geoIP.Latitude,
		"longitude":    geoIP.Longitude,
		"metro_code":   geoIP.MetroCode,
	})
}

func GetIP(c *fiber.Ctx) *GeoIPResponse {

	ip := c.IP()

	resp, err := http.Get("https://freegeoip.app/json/" + ip)
	if err != nil {
		fmt.Println("Error fetching GeoIP data:", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
	}

	var geoIP GeoIPResponse
	err = json.Unmarshal(body, &geoIP)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
	}

	return &geoIP
}

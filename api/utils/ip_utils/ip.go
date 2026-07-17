package ip_utils

import (
	"encoding/json"
	"net"
	"strings"

	"github.com/biter777/countries"
	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"
)

type countryAPIResponse struct {
	Country     string `json:"country"`
	CountryCode string `json:"countryCode"`
}

func GetCountryFromIP(ip string) (string, error) {

	client := resty.New()
	resp, err := client.R().
		EnableTrace().
		Get("https://extreme-ip-lookup.com/json/")

	if err != nil {
		return "", err
	}

	countryResponse := countryAPIResponse{}

	err = json.Unmarshal(resp.Body(), &countryResponse)
	if err != nil {
		return "", err
	}

	return countryResponse.Country, nil
}

// cloudflare automatically sends us the country code header which we can use to get the name of the country
func GetCountryFromEchoContext(c echo.Context) (string, error) {

	countryCode := c.Request().Header.Get("CF-IPCountry")
	if countryCode == "" {
		return "", nil
	}

	// use a go package to get the name of the country from the country code
	country := countries.ByName(countryCode)
	return country.Info().Name, nil
}

func GetUserIpWithHeaders(xForwardedfor string, remoteAddress string) string {
	remoteAddr, _, _ := net.SplitHostPort(remoteAddress)
	if xForwardedfor != "" {
		forwardedForArray := strings.Split(xForwardedfor, ",")
		if len(forwardedForArray) < 1 {
			return remoteAddr
		}
		return forwardedForArray[0]
	}
	return remoteAddr
}

func GetUserIPEcho(c echo.Context) string {
	remoteAddr, _, _ := net.SplitHostPort(c.Request().RemoteAddr)
	forwardedFor := c.Request().Header.Get("X-Forwarded-For")
	if forwardedFor != "" {
		forwardedForArray := strings.Split(forwardedFor, ",")
		if len(forwardedForArray) < 1 {
			return remoteAddr
		}
		return forwardedForArray[0]
	}
	return remoteAddr
}

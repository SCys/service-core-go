package core

import (
	"net"

	"github.com/oschwald/geoip2-golang"
)

var GeoLiteReader *geoip2.Reader

type GeoIP struct {
	// basic info
	CountryISOCode string `json:"country_iso_code"`
	Country        string `json:"country"`
	Province       string `json:"province"`
	City           string `json:"city"`

	// extra info
	ISP      string  `json:"isp"`
	Lat      float64 `json:"lat"`
	Lng      float64 `json:"lng"`
	TimeZone string  `json:"time_zone"`
}

func IPGeo(ip net.IP) *GeoIP {
	if GeoLiteReader == nil {
		return nil
	}

	record, err := GeoLiteReader.City(ip)
	if err != nil {
		E("lookup ip failed:%s", err, ip)
		return nil
	}

	lang := "en"
	if record.Country.IsoCode == "CN" {
		lang = "zh-CN"
	}

	return &GeoIP{
		Country:  record.Country.Names[lang],
		Province: record.Subdivisions[0].Names[lang],
		City:     record.City.Names[lang],

		ISP:      "unknown",
		Lat:      record.Location.Latitude,
		Lng:      record.Location.Longitude,
		TimeZone: record.Location.TimeZone,
	}
}

// IPGeoInit 初始化
func IPGeoInit(dbPath string) {
	var err error

	GeoLiteReader, err = geoip2.Open(dbPath)
	if err != nil {
		E("open geolite failed:%s", err, dbPath)
	}
}

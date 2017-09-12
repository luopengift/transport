package utils

type EsIP struct {
	Timezone      string  `json:"timezone"`
	Ip            string  `json:"ip"`
	Latitude      float64 `json:"latitude"`
	ContinentCode string  `json:"continent_code"`
	CityName      string  `json:"city_name"`
	CountryName   string  `json:"country_name"`
	CountryCode2  string  `json:"country_code2"`
	DmaCode       int     `json:"dma_code"`
	CountryCode3  string  `json:"country_code3"`
	RegionName    string  `json:"region_name"`
	Location      struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"location"`
	PostalCode string  `json:"postal_code"`
	RegionCode string  `json:"region_code"`
	Longitude  float64 `json:"longitude"`
}

func GeoToEsIP(geoip *GeoIP) *EsIP {
	esip := EsIP{}
	esip.Timezone = geoip.Location.TimeZone
	esip.Ip = geoip.IP
	esip.Latitude = geoip.Location.Latitude
	esip.ContinentCode = geoip.Continent.Code
	esip.CityName = geoip.City.Names.EN
	esip.CountryName = geoip.Country.Names.EN
	esip.CountryCode2 = geoip.Country.IsoCode
	esip.DmaCode = geoip.Location.MetroCode
	esip.CountryCode3 = geoip.Country.IsoCode
	if len(geoip.Subdivisions) > 0 {
		esip.RegionName = geoip.Subdivisions[0].Names.EN
		esip.RegionCode = geoip.Subdivisions[0].IsoCode
	}
	esip.Location.Lon = geoip.Location.Longitude
	esip.Location.Lat = geoip.Location.Latitude
	esip.PostalCode = geoip.Postal.Code
	esip.Longitude = geoip.Location.Longitude
	return &esip
}

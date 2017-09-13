package utils

type EsIP struct {
	Timezone      string  `json:"timezone,omitempty"`
	Ip            string  `json:"ip,omitempty"`
	Latitude      float64 `json:"latitude,omitempty"`
	ContinentCode string  `json:"continent_code,omitempty"`
	CityName      string  `json:"city_name,omitempty"`
	CountryName   string  `json:"country_name,omitempty"`
	CountryCode2  string  `json:"country_code2,omitempty"`
	DmaCode       int     `json:"dma_code,omitempty"`
	CountryCode3  string  `json:"country_code3,omitempty"`
	RegionName    string  `json:"region_name,omitempty"`
	Location      struct {
		Lon float64 `json:"lon,omitempty"`
		Lat float64 `json:"lat,omitempty"`
	} `json:"location,omitempty"`
	PostalCode string  `json:"postal_code,omitempty"`
	RegionCode string  `json:"region_code,omitempty"`
	Longitude  float64 `json:"longitude,omitempty"`
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

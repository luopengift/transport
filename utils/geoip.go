package utils

import (
	"github.com/oschwald/maxminddb-golang"
    "github.com/luopengift/golibs/logger"
    "net"
)

type ID = uint64
type Names struct {
	ES   string `json:"es"`
	JA   string `json:"ja"`
	PTBR string `json:"pt-BR"`
	EN   string `json:"en"`
	FR   string `json:"fr"`
	RU   string `json:"ru"`
	ZHCN string `json:"zh-CN"`
	DE   string `json:"de"`
}

type City struct {
	GeonameId ID    `json:"geoname_id"`
	Names     Names `json:"names"`
}

type Continent struct {
	GeonameId ID     `json:"geoname_id"`
	Names     Names  `json:"names"`
	Code      string `json:"code"`
}

type Country struct {
	GeonameId ID     `json:"geoname_id"`
	Names     Names  `json:"names"`
	IsoCode   string `json:"iso_code"`
}

type Location struct {
	Longitude      float64 `json:"longitude"`
	TimeZone       string  `json:"time_zone"`
	AccuracyRadius int     `json:"accuracy_radius"`
	Latitude       float64 `json:"latitude"`
	MetroCode      int     `json:"metro_code"`
}

type Postal struct {
	Code string `json:"code"`
}

type RegisteredCountry struct {
	GeonameId ID    `json:"geoname_id"`
	Names     Names `json:"names"`
}

type Subdivision struct {
	GeonameId ID     `json:"geoname_id"`
	Names     Names  `json:"names"`
	IsoCode   string `json:"iso_code"`
}

type GeoIP struct {
	IP                string            `json:"ip"`
	City              City              `json:"city"`
	Continent         Continent         `json:"continent"`
	Country           Country           `json:"country"`
	Location          Location          `json:"location"`
	Postal            Postal            `json:"postal"`
	RegisteredCountry RegisteredCountry `json:"registered_country"`
	Subdivisions      []Subdivision     `json:"subdivisions"`
}

var (
    GeoDB = "GeoLite2-City.mmdb"
)

type Client struct {
	*maxminddb.Reader
}

func NewClient(db string) (*Client, error) {
    var err error
	c := new(Client)
	c.Reader, err = maxminddb.Open(db)
    if err != nil {
        logger.Error("GeoIP db open error: %v", err)
        return nil, err
    }
	return c, nil
}

func (c *Client) Close() error {
	return c.Reader.Close()
}

func (c *Client) Search(ip string) (*GeoIP, error) {
	record := map[string]interface{}{}
	err := c.Reader.Lookup(net.ParseIP(ip), &record)
	if err != nil {
		return nil, err
	}
	geoip := &GeoIP{IP: ip}
	err = Format(record, geoip)
	return geoip, err
}

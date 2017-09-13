package utils

import (
	"encoding/json"
	"fmt"
	"testing"
)

func Test_GeoIP(t *testing.T) {
	c, _ := NewClient(GeoDB)
	defer c.Close()
	rest, err := c.Search("54.245.168.2")
	fmt.Println(rest, err)
	rest, err = c.Search("191.158.113.240")
	esip := GeoToEsIP(rest)
	fmt.Println(fmt.Sprintf("%#v, %v", esip, err))
	b, _ := json.Marshal(esip)
	fmt.Println(string(b))
}

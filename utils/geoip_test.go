package utils

import (
	"fmt"
	"testing"
)

func Test_GeoIP(t *testing.T) {
	c, _ := NewClient(GeoDB)
	defer c.Close()
	rest, err := c.Search("54.245.168.2")
	fmt.Println(rest, err)
}

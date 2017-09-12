package codec

import (
	"encoding/json"
	"fmt"
	"github.com/luopengift/gohttp"
	"github.com/luopengift/golibs/logger"
	"github.com/luopengift/transport"
	"github.com/luopengift/transport/utils"
	"strconv"
	"strings"
)

type KVHandler struct {
	Keys   [][]string             `json:"keys"`
	Split  string                 `json:"split"`
	Ignore string                 `json:"ignore"`
	GeoIP  string                 `json:"geoip"` //tell program which keys format to geoip, eg: "ip => geoip"
	IpDB   string                 `json:"ipdb"`
	Tags   map[string]interface{} `json:"tags"`

	geomap map[string]string
}

func (kv *KVHandler) Init(config transport.Configer) error {
	kv.Ignore = "-"
	kv.GeoIP = ""
	kv.IpDB = utils.GeoDB
	kv.geomap = map[string]string{}
	err := config.Parse(kv)
	if err != nil {
		return err
	}
	if kv.GeoIP != "" && strings.Count(kv.GeoIP, "=>") == 1 {
		geo := strings.Split(kv.GeoIP, "=>") //key is need to geoip, value is return key
		kv.geomap[strings.TrimSpace(geo[0])] = strings.TrimSpace(geo[1])
		utils.GeoIPClient, err = utils.NewClient(kv.IpDB)
	}
	return err
}

func (kv *KVHandler) Handle(in, out []byte) (int, error) {
	o := make(map[string]interface{})
	for index, value := range strings.Split(string(in), kv.Split) {
		if index == len(kv.Keys) {
			logger.Warn("index<%d> out of len(kv.Keys)<%d>: %s", index, len(kv.Keys), string(in))
			return 0, fmt.Errorf("index out of range")
		}
		key := kv.Keys[index][0]
		valueType := kv.Keys[index][1]
		if key == kv.Ignore || valueType == kv.Ignore {
			continue
		}
		switch valueType {
		case "string", "geoip":
			o[key] = value
		case "int":
			if v, err := strconv.Atoi(value); err == nil {
				o[key] = v
			}
		case "int64":
			if v, err := strconv.ParseInt(value, 10, 64); err == nil {
				o[key] = v
			}
		case "float64":
			if v, err := strconv.ParseFloat(value, 64); err == nil {
				o[key] = v
			}
		case "json":
			v := map[string]interface{}{}
			if err := json.Unmarshal(gohttp.StringToBytes(value), &v); err == nil {
				o[key] = v
			}
		default:
			return 0, fmt.Errorf("type<%v> is unknown", valueType)
		}
	}

	if kv.GeoIP != "" {
		for key, value := range kv.geomap {
			geoip, err := utils.GeoIPClient.Search(o[key].(string))
			if err == nil {
				o[value] = *utils.GeoToEsIP(geoip)
				continue
			}
			logger.Warn("GeoIP is fail:%v", o[key])
		}
	}

	for key, value := range kv.Tags {
		o[key] = value
	}

	b, err := gohttp.ToBytes(o)
	if err != nil {
		return 0, err
	}
	n := copy(out, b)
	return n, nil
}

func (kv *KVHandler) Version() string {
	return "0.0.3"
}

func init() {
	transport.RegistHandler("kv", new(KVHandler))
}

package elasticsearch

import (
	"encoding/json"
	"github.com/luopengift/gohttp"
	"github.com/luopengift/golibs/logger"
)

const (
	enter = '\n'
)

//action_and_meta_data\n
//optional_source\n

type Meta struct {
	Index string `json:"_index"`
	Type  string `json:"_type"`
	Id    string `json:"_id,omitempty"`
}

func NewMeta(_index, _type, _id string) *Meta {
	return &Meta{_index, _type, _id}
}

type Bulk struct {
	Index  *Meta  `json:"index,omitempty"`
	Update *Meta  `json:"update,omitempty"`
	Source []byte `json:"-"`
}

func NewBulkIndex(_index, _type, _id string, source []byte) *Bulk {
	return &Bulk{
		Index:  NewMeta(_index, _type, _id),
		Source: source,
	}
}
func NewBulkUpdate(_index, _type, _id string, source []byte) *Bulk {
	return &Bulk{
		Update: NewMeta(_index, _type, _id),
		Source: source,
	}
}

//构建/_bulk接口所需的数据格式
func (b *Bulk) Bytes() ([]byte, error) {
	action, err := json.Marshal(b)
	if err != nil {
		return nil, err
	}
	p := make([]byte, 0, 1000)
	p = append(p, action...)
	p = append(p, enter)
	p = append(p, b.Source...)
	p = append(p, enter)
	return p, nil
}

// A BulkIndexer is used to index documents in ElasticSearch
type BulkIndexer interface {
	// Index documents
	Index(body []byte) (err error, retry bool)
	// Check if a flush is needed
	CheckFlush(count int, length int) bool
}

// A HttpBulkIndexer uses the HTTP REST Bulk Api of ElasticSearch
// in order to index documents
type HttpBulk struct {
	// Protocol (http or https).
	Protocol string
	// Host name and port number (default to "localhost:9200").
	Addrs string
	// Path (default to "")
	Path string
	// Maximum number of documents.
	MaxCount int
	// Internal HTTP Client.
	Client *gohttp.Client
	// Optional username for HTTP authentication
	username string
	// Optional password for HTTP authentication
	password string
}

func NewHttpBulk(protocol, addrs, path string, maxCount int,
	username string, password string) *HttpBulk {

	h := &HttpBulk{
		Protocol: protocol,
		Addrs:   addrs,
		Path:     path,
		MaxCount: maxCount,
		username: username,
		password: password,
	}

	h.Client = gohttp.NewClient().Url(h.Protocol+"://"+h.Addrs).Path("/_bulk").Header("Accept", "application/json")
	return h
}

func (h *HttpBulk) Index(p []byte) error {
    resp,err := h.Client.Body(p).Post()
	if err != nil {
		logger.Error("<bulk post error>%#v",err)
		return err
	}
    if resp.Code() != 200 {
        logger.Warn("resp%#v",resp.String())
        return nil    
    }
    logger.Debug("Response is %v",resp.String())
    return nil
}



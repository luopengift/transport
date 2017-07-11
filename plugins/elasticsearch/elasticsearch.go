//package elasticsearch
package main

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
	Source string `json:"-"`
}

func NewBulkIndex(_index, _type, _id, source string) *Bulk {
	return &Bulk{
		Index:  NewMeta(_index, _type, _id),
		Source: source,
	}
}
func NewBulkUpdate(_index, _type, _id, source string) *Bulk {
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
	source := []byte(b.Source)
	p := make([]byte, 0, 1000)
	p = append(p, action...)
	p = append(p, enter)
	p = append(p, source...)
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
	client *gohttp.Client
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

	h.client = gohttp.NewClient().URL(h.Protocol+"://"+h.Addrs).Path(path).Path("/_bulk").Header("Accept", "application/json")
	return h
}

func (h *HttpBulk) Index(p []byte) error {
	resp,err := h.client.Body(p).Post()
	if err != nil {
		logger.Error("<bulk post error>resp:%v,err:%v",resp.String(),err)
		return err
	}
	return nil
}


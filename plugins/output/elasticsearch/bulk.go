package elasticsearch

import (
    "encoding/json"
    "bytes"
    "github.com/luopengift/gohttp"
    "fmt"
)

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
    var b bytes.Buffer
    b.Write(action)
    b.Write('\n')
    b.Write(b.Source)
    b.Write('\n')
    return b.Bytes(), nil
}

func Send(addr string, p []byte) error {
    resp, err := gohttp.Url("http://"+addr).Path("/_bulk").Header("Accept", "application/json").Body(p).Post()
    if err != nil {
        return err
    }
    if resp.Code() != 200 {
        return fmt.Errorf(resp.String())
    }
    return nil

}


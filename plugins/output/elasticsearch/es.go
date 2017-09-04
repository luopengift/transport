package elasticsearch

import (
	"github.com/luopengift/golibs/file"
	"github.com/luopengift/golibs/logger"
	"github.com/luopengift/transport"

    "context"
    es "gopkg.in/olivere/elastic.v5"


)

type EsOutput struct {
	Addrs   []string `json:"addrs"` //es addrs
	Index   string   `json:"index"` //es index
	Type    string   `json:"type"`  //es type
	Timeout int      `json:"time"`  //Pool timeout
	Batch   int      `json:"batch"` //多少条数据提交一次

    buffer chan []byte
    ctx context.Context
    client  *es.Client
}

func NewEsOutput() *EsOutput {
	return new(EsOutput)
}

func (out *EsOutput) Init(config transport.Configer) error {
	out.Timeout = 5
	out.Batch = 1
	err := config.Parse(out)
	if err != nil {
		return err
	}

    out.buffer = make(chan []byte, out.Batch * 2)

    // 连接es 
    out.ctx = context.Background()
    out.client, err = es.NewClient(es.SetURL("http://"+out.Addrs[0]), es.SetSniff(false))
    if err != nil {
        return err
    }

    // 检查index是否存在，如果不存在则创建index 
    exists, err := out.client.IndexExists(file.TimeRule.Handle(out.Index)).Do(out.ctx)
    if err != nil {
        return err
    }
    if !exists {
        _, err := out.client.CreateIndex(file.TimeRule.Handle(out.Index)).Do(out.ctx)
        if err != nil {
            return err
        }
    }
    return nil

}

func (out *EsOutput) Write(p []byte) (int, error) {
	out.buffer <- p
	return len(p), nil
}

func (out *EsOutput) Start() error {
    for {
        bulkRequest := out.client.Bulk()
        for tmp := out.Batch; tmp>0; tmp-- {
            req := es.NewBulkIndexRequest().Index(out.Index).Type(out.Type).Doc(string(<-out.buffer))
            bulkRequest.Add(req)
        }
        bulkResponse, err := bulkRequest.Do(out.ctx)
        if err != nil {
            logger.Error("bulkResponse error %v",err)
        }
        indexed := bulkResponse.Indexed() 
        logger.Info("导入了",len(indexed),"条数据") 
    }
    return nil
}

func (out *EsOutput) Close() error {
	return nil
}

func init() {
	transport.RegistOutputer("elasticsearch", NewEsOutput())
}

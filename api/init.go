package api

import (
	"github.com/luopengift/gohttp"
)

type RootHandler struct {
	gohttp.HttpHandler
}

func (ctx *RootHandler) GET() {
	ctx.Output("root")
}

type StatsHandler struct {
	gohttp.HttpHandler
}

func (ctx *StatsHandler) GET() {
	/*input_stat := []string{}
	  for _,input := range t.Inputs {
	      input_stat = append(input_stat, fmt.Sprintf("%v:%v",input.Name,input.Cnt))
	  }
	  codec_stat := []string{}
	  for _,codec := range t.Codecs {
	      codec_stat = append(codec_stat, fmt.Sprintf("%v:%v",codec.Name,codec.Cnt))
	  }
	  output_stat := []string{}
	  for _,output := range t.Outputs {
	      output_stat = append(output_stat, fmt.Sprintf("%v:%v",output.Name,output.Cnt))
	  }
	  msg := fmt.Sprinf("stat=> inputs:%s|codecs:%s|outputs:%s", strings.Join(input_stat,","),strings.Join(codec_stat,","),strings.Join(output_stat,","))

	*/
	msg := "status"

	ctx.Output(msg)
}

type StoreHandler struct {
	gohttp.HttpHandler
}

func init() {
	app := gohttp.Init()
	app.Route("^/$", &RootHandler{})
	app.Route("^/stats$", &StatsHandler{})
	go app.Run(":38888")
}

package transport

var ExampleConfig = `
{
    "runtime":{
        "DEBUG":false,
        "MAXPROCS":4,
        "VERSION":"0.0.1"
    },
    "input":{
        "type":"file",
        "path":"a.log"
    },
    "handle":{
        "type":"null"
    },
    "output":{
        "type":"file",
        "path":"test.test"
    }
}`

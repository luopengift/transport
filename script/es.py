#!/usr/bin/env python
#-*-coding:utf8-*-
import elasticsearch
import json
from datetime import datetime
from elasticsearch import Elasticsearch
from concurrent.futures import ThreadPoolExecutor
from concurrent.futures import as_completed


'''i
es = elasticsearch.Elasticsearch(hosts=
        [{'host':'10.10.10.102','port':9200},{'host':'10.10.10.103','port':9200}])

re = es.indices.get(index='zhizi-log-*')
res = es.search(index="zhizi-log-*", body={"query":{"bool":{"must":[{"term":{"module":"honeybee.redis.fallback"}}],"must_not":[],"should":[]}},"from":0,"size":10,"sort":[],"aggs":{}})
print res
'''

def log_w(text):
    logfile = '/tmp/request_id.log'
    tt = str(text) + "\n"
    f = open(logfile,'a+')
    f.write(tt)
    f.close() 

url = "http://10.10.10.102:9200/"
index_name  = "zhizi-log-2017.07.20"

es = Elasticsearch(url,timeout=120)

query = """
{
  "query": {
    "bool": {
      "must": [
        {
          "term": {
            "module": "honeybee.redis.fallback"
          }
        }
      ],
      "must_not": [],
      "should": []
    }
  },
  "sort": [],
  "aggs": {}
}
"""


return_fields = [
    '_scroll_id',
    'hits.hits._source.request_id',
    'hits.hits._source.cost',
    'hits.hits._source.module',
    #hits.hits._source.uid',
    #hits.hits._source.request_time',
    #hits.hits._source.timestamp',
]

Dict ={}

fallbacklist = []


def query_requestid(request_id):
    request_cycle = list()
    res = es.search(index="zhizi-log-*", body={"query":{"bool":{"must":[{"term":{"request_id":request_id}}],"must_not":[],"should":[]}},"from":0,"size":50,"sort":[],"aggs":{}})
    for data in res['hits']['hits']:
        value = data["_source"]
        if value["module"] not in ["honeybee.rpc.zhizicore","core.total_cost","honeybee.infolist","core.invoke_pipeline"]:
            request_cycle.append((value['module'],value['cost']))
    sort = sorted(request_cycle, key=lambda x:x[1])
    return request_id, sort[-1][0],sort[-1][1]




def main():
    response = es.search(
            index=index_name,
            body=query,
            scroll="2m"
        )
    n = 0
    scrollId=response.get("_scroll_id","")  # 获取scrollID
    results = response.get("hits",{})
    while results != {}:
        response= es.scroll(scroll_id=scrollId, scroll= "2m", filter_path=return_fields)
        scrollId = response.get("_scroll_id","")
        results = response.get("hits",{})
        if results != {}:
            for data in results.get("hits"):
                value = data["_source"]
                fallbacklist.append(value["request_id"])
    print "fallbacklist length is:","\t",len(fallbacklist)


    max_works = 10
    with ThreadPoolExecutor(max_workers=max_works) as executor:
        futures = [executor.submit(query_requestid,request_id) for request_id in fallbacklist]    
        for future in as_completed(futures):
            res = future.result()
            print res[0],"\t",res[1],"\t",res[2]
if __name__ == "__main__":
    main()

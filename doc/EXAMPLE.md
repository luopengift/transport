## 
将Nginx里的日志收集，分析，展示

1. 先放图
![nginx](https://github.com/luopengift/transport/blob/master/Image/png/nginx_log_analysis.png)
![nginx](https://github.com/luopengift/transport/blob/master/Image/png/worldmap.png.png)

2. 需要的工具
  * Nginx
  * Transport
  * Elasticsearch
  * Grafana
  
3. Nginx配置
```
log_format main '$time_iso8601|$http_x_forwarded_for|$remote_addr|$upstream_addr|$server_addr|$hostname|$http_host|$server_name|$http_referer|$status|$body_bytes_sent|$upstream_response_time|$request_time|$request_method|$https|$scheme|$request_uri|$http_user_agent|$args|$request_body';
```

4. Transport配置
file-kv-elasticsearch [详情](https://github.com/luopengift/transport/blob/master/test/nginx-kv-es.json)

5. Elasticsearch配置
添加template解析geo_point类型, [详情](https://github.com/luopengift/transport/blob/master/test/elasticsearch/geoip.template.json)

6. Grafana配置
  * ./grafana-cli plugins install grafana-worldmap-panel
  * ./grafana-cli plugins install grafana-piechart-pane
  * ttps://github.com/luopengift/transport/blob/master/test/grafana/nginx_es.template.json





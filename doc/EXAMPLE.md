## 
将Nginx里的日志收集，分析，展示

1. 先放图
![nginx](https://github.com/luopengift/transport/blob/master/Image/png/nginx_log_analysis.png)
![nginx](https://github.com/luopengift/transport/blob/master/Image/png/worldmap.png)

2. 需要的工具
  * Nginx
  * Transport
  * Elasticsearch 5.x
  * Grafana
  
3. [Nginx](https://www.nginx.com/resources/wiki/)配置
```
log_format main '$time_iso8601|$http_x_forwarded_for|$remote_addr|$upstream_addr|$server_addr|$hostname|$http_host|$server_name|$http_referer|$status|$body_bytes_sent|$upstream_response_time|$request_time|$request_method|$https|$scheme|$request_uri|$http_user_agent|$args|$request_body';
```

4. [Transport](https://github.com/luopengift/transport)配置
file-kv-elasticsearch [详情](https://github.com/luopengift/transport/blob/master/test/nginx-kv-es.json)

5. [Elasticsearch](https://www.elastic.co/guide/en/elasticsearch/reference/current/index.html)配置
添加template解析geo_point类型, [详情](https://github.com/luopengift/transport/blob/master/test/elasticsearch/geoip.template.json)

6. [Grafana](http://docs.grafana.org/)配置
  * [plugin] ./grafana-cli plugins install grafana-worldmap-panel
  * [plugin] ./grafana-cli plugins install grafana-piechart-pane
  * [Dashboard](https://github.com/luopengift/transport/blob/master/test/grafana/nginx_es.template.json)
  * NOTE: Dashboard由vincentliu提供的[Elasticserch Nginx Logs](https://grafana.com/dashboards/2292)改进而来





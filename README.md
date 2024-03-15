# Stage-2024

## Config

does not work

```json
{
  "clickhouseSettings": "",
  "connector.class": "com.clickhouse.kafka.connect.ClickHouseSinkConnector",
  "consumer.override.max.partition.fetch.bytes": "5242880",
  "consumer.override.max.poll.records": "5000",
  "database": "default",
  "errors.retry.timeout": "60",
  "exactlyOnce": "false",
  "hostname": "clickhouse",
  "name": "clickhouse-testing",
  "password": "root",
  "port": "8443",
  "ssl": "false",
  "tasks.max": "1",
  "topics": "bolt-test",
  "username": "root",
  "value.converter": "org.apache.kafka.connect.json.JsonConverter",
  "value.converter.schemas.enable": "false"
}
```

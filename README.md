# Stage-2024

## Mongo

[Mongo Express](http://localhost:8081/)
admin
pass

## Connector config

database should already exist

```json
{
  "connector.class": "com.mongodb.kafka.connect.MongoSinkConnector",
  "name": "mongodb-sink-connector-test",
  "topics.regex": "donkey-locations",
  "database": "test",
  "connection.url": "mongodb://mongo:27017/",
  "connection.password": "example",
  "connection.username": "root",
  "errors.deadletterqueue.context.headers.enable": "false",
  "timeseries.timefield.auto.convert": "false",
  "value.converter": "io.confluent.connect.protobuf.ProtobufConverter",
  "value.converter.schema.registry.url": "http://redpanda-0:8081",
  "key.converter": "io.confluent.connect.protobuf.ProtobufConverter",
  "key.converter.schema.registry.url": "http://redpanda-0:8081"
}
```

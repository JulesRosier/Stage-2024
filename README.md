# Stage-2024

## Mongo

[Mongo Express](http://localhost:8081/)
admin
pass

## Connector config

database should already exist

```json
{
  "connection.password": "example",
  "connection.uri": "mongodb://root:example@mongo:27017/",
  "connection.url": "mongodb://mongo:27017/",
  "connection.username": "root",
  "connector.class": "com.mongodb.kafka.connect.MongoSinkConnector",
  "database": "test",
  "errors.deadletterqueue.context.headers.enable": "false",
  "key.converter": "io.confluent.connect.protobuf.ProtobufConverter",
  "key.converter.schema.registry.url": "http://redpanda-0:8081",
  "name": "mongodb-sink-connector-juuc",
  "timeseries.timefield.auto.convert": "false",
  "topics.regex": "donkey-locations",
  "value.converter": "io.confluent.connect.protobuf.ProtobufConverter",
  "value.converter.schema.registry.url": "http://redpanda-0:8081"
}
```

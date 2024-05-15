# from pyspark.sql import SparkSession
# from pprint import pprint 
# from datetime import datetime, date
# import pandas as pd
# from pyspark.sql import Row

# spark = SparkSession.builder.getOrCreate()


# df = spark.createDataFrame([
#     Row(a=1, b=2., c='string1', d=date(2000, 1, 1), e=datetime(2000, 1, 1, 12, 0)),
#     Row(a=2, b=3., c='string2', d=date(2000, 2, 1), e=datetime(2000, 1, 2, 12, 0)),
#     Row(a=4, b=5., c='string3', d=date(2000, 3, 1), e=datetime(2000, 1, 3, 12, 0))
# ])
# pprint(df)


from pyspark.sql import SparkSession
import os
os.environ['PYSPARK_SUBMIT_ARGS'] = '--packages org.apache.spark:spark-streaming-kafka-0-10_2.12:3.2.0,org.apache.spark:spark-sql-kafka-0-10_2.12:3.2.0 pyspark-shell'

from pyspark.sql.protobuf.functions import from_protobuf, to_protobuf

KAFKA_BOOTSTRAP_SERVERS = "localhost:19092"
KAFKA_TOPIC = "bike_picked_up"

spark = SparkSession.builder.appName("read_test_straeam").getOrCreate()

# Reduce logging
spark.sparkContext.setLogLevel("WARN")

df = spark.readStream.format("kafka") \
    .option("kafka.bootstrap.servers", KAFKA_BOOTSTRAP_SERVERS) \
    .option("subscribe", KAFKA_TOPIC) \
    .option("startingOffsets", "earliest") \
    .load()

schema_registry_options = {
  "schema.registry.subject" : "bike_picked_up-value",
  "schema.registry.address" : "https://localhost:18081/"
}
# df.selectExpr("CAST(key AS STRING)", "CAST(value AS STRING)") \
#     .writeStream \
#     .format("console") \
#     .outputMode("append") \
#     .start() \
#     .awaitTermination()
# output = df \
#   .select(from_protobuf("value", "AppEvent", "../proto/bikes/bike_picked_up.proto").alias("event")) \
#   .where('event.name == "alice"') \
#   .select(to_protobuf("event", "AppEvent", "../proto/bikes/bike_picked_up.proto").alias("event")) 

# stream1_proto_select = df.select(from_protobuf("value", "bike_picked_up", options=schema_registry_options).alias("event"))
# stream1_proto_select.printSchema()
proto_events_df = (
  df
    .select(
      from_protobuf("proto_bytes", "bike_picked_up",options = schema_registry_options)
        .alias("proto_event")
    )
)
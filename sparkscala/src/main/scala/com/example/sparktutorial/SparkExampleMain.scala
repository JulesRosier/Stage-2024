package com.example.sparktutorial

import org.apache.kafka.clients.consumer.ConsumerRecord
import org.apache.kafka.common.serialization.StringDeserializer
import org.apache.spark.streaming.kafka010._
import org.apache.spark.streaming.kafka010.LocationStrategies.PreferConsistent
import org.apache.spark.streaming.kafka010.ConsumerStrategies.Subscribe

object SparkExampleMain extends App {
  println("Hello world")

//  val spark: SparkSession = SparkSession.builder()
//    .master("local[*]")
//    .appName("SparkByExamples.com")
//    .getOrCreate()
//
//  val df = spark
//    .readStream
//    .format("kafka")
//    .option("kafka.bootstrap.servers", "localhost:19092")
//    .option("subscirbe", "bike_picked_up")
//    .option("includeHeaders", "true")
//    .load()
//  df.selectExpr("CAST(key AS STRING)", "CAST(value AS STRING)")
//    .as[(String, String)]

//  df.show()
val kafkaParams = Map[String, Object](
  "bootstrap.servers" -> "localhost:9092,anotherhost:9092",
  "key.deserializer" -> classOf[StringDeserializer],
  "value.deserializer" -> classOf[StringDeserializer],
  "group.id" -> "use_a_separate_group_id_for_each_stream",
  "auto.offset.reset" -> "latest",
  "enable.auto.commit" -> (false: java.lang.Boolean)
)

  val topics = Array("topicA", "topicB")
  val stream = KafkaUtils.createDirectStream[String, String](
    streamingContext,
    PreferConsistent,
    Subscribe[String, String](topics, kafkaParams)
  )

  stream.map(record => (record.key, record.value))
}
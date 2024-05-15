scalaVersion := "2.12.18"

name := "sparktutorial"
organization := "com.example"
version := "1.0"

libraryDependencies ++= Seq(
  "org.apache.spark" %% "spark-core" % "3.5.0",
  "org.apache.spark" %% "spark-sql" % "3.5.0"
)

fork := true
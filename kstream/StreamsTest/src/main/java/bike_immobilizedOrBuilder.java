// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: bikes/bike_immobilized.proto

// Protobuf Java Version: 4.26.1
public interface bike_immobilizedOrBuilder extends
    // @@protoc_insertion_point(interface_extends:bike_immobilized)
    com.google.protobuf.MessageOrBuilder {

  /**
   * <code>.google.protobuf.Timestamp time_stamp = 1 [json_name = "timeStamp"];</code>
   * @return Whether the timeStamp field is set.
   */
  boolean hasTimeStamp();
  /**
   * <code>.google.protobuf.Timestamp time_stamp = 1 [json_name = "timeStamp"];</code>
   * @return The timeStamp.
   */
  com.google.protobuf.Timestamp getTimeStamp();
  /**
   * <code>.google.protobuf.Timestamp time_stamp = 1 [json_name = "timeStamp"];</code>
   */
  com.google.protobuf.TimestampOrBuilder getTimeStampOrBuilder();

  /**
   * <code>.location_bike bike = 2 [json_name = "bike"];</code>
   * @return Whether the bike field is set.
   */
  boolean hasBike();
  /**
   * <code>.location_bike bike = 2 [json_name = "bike"];</code>
   * @return The bike.
   */
  location_bike getBike();
  /**
   * <code>.location_bike bike = 2 [json_name = "bike"];</code>
   */
  location_bikeOrBuilder getBikeOrBuilder();
}

// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: stations/station_occupation_decreased.proto
// Protobuf Java Version: 4.26.1

public final class StationOccupationDecreasedProto {
  private StationOccupationDecreasedProto() {}
  static {
    com.google.protobuf.RuntimeVersion.validateProtobufGencodeVersion(
      com.google.protobuf.RuntimeVersion.RuntimeDomain.PUBLIC,
      /* major= */ 4,
      /* minor= */ 26,
      /* patch= */ 1,
      /* suffix= */ "",
      StationOccupationDecreasedProto.class.getName());
  }
  public static void registerAllExtensions(
      com.google.protobuf.ExtensionRegistryLite registry) {
  }

  public static void registerAllExtensions(
      com.google.protobuf.ExtensionRegistry registry) {
    registerAllExtensions(
        (com.google.protobuf.ExtensionRegistryLite) registry);
  }
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_station_occupation_decreased_descriptor;
  static final 
    com.google.protobuf.GeneratedMessage.FieldAccessorTable
      internal_static_station_occupation_decreased_fieldAccessorTable;

  public static com.google.protobuf.Descriptors.FileDescriptor
      getDescriptor() {
    return descriptor;
  }
  private static  com.google.protobuf.Descriptors.FileDescriptor
      descriptor;
  static {
    java.lang.String[] descriptorData = {
      "\n+stations/station_occupation_decreased." +
      "proto\032\037google/protobuf/timestamp.proto\032%" +
      "stations/station_identification.proto\"\230\002" +
      "\n\034station_occupation_decreased\0229\n\ntime_s" +
      "tamp\030\001 \001(\0132\032.google.protobuf.TimestampR\t" +
      "timeStamp\0221\n\007station\030\002 \001(\0132\027.station_ide" +
      "ntificationR\007station\022)\n\020amount_decreased" +
      "\030\003 \001(\005R\017amountDecreased\022<\n\032current_avail" +
      "able_capacity\030\004 \001(\005R\030currentAvailableCap" +
      "acity\022!\n\014max_capacity\030\005 \001(\005R\013maxCapacity" +
      "BDB\037StationOccupationDecreasedProtoP\001Z\037s" +
      "tage2024/pkg/protogen/stationsb\006proto3"
    };
    descriptor = com.google.protobuf.Descriptors.FileDescriptor
      .internalBuildGeneratedFileFrom(descriptorData,
        new com.google.protobuf.Descriptors.FileDescriptor[] {
          com.google.protobuf.TimestampProto.getDescriptor(),
          StationIdentificationProto.getDescriptor(),
        });
    internal_static_station_occupation_decreased_descriptor =
      getDescriptor().getMessageTypes().get(0);
    internal_static_station_occupation_decreased_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessage.FieldAccessorTable(
        internal_static_station_occupation_decreased_descriptor,
        new java.lang.String[] { "TimeStamp", "Station", "AmountDecreased", "CurrentAvailableCapacity", "MaxCapacity", });
    descriptor.resolveAllFeaturesImmutable();
    com.google.protobuf.TimestampProto.getDescriptor();
    StationIdentificationProto.getDescriptor();
  }

  // @@protoc_insertion_point(outer_class_scope)
}

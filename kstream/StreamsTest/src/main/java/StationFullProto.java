// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: stations/station_full.proto
// Protobuf Java Version: 4.26.1

public final class StationFullProto {
  private StationFullProto() {}
  static {
    com.google.protobuf.RuntimeVersion.validateProtobufGencodeVersion(
      com.google.protobuf.RuntimeVersion.RuntimeDomain.PUBLIC,
      /* major= */ 4,
      /* minor= */ 26,
      /* patch= */ 1,
      /* suffix= */ "",
      StationFullProto.class.getName());
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
    internal_static_station_full_descriptor;
  static final 
    com.google.protobuf.GeneratedMessage.FieldAccessorTable
      internal_static_station_full_fieldAccessorTable;

  public static com.google.protobuf.Descriptors.FileDescriptor
      getDescriptor() {
    return descriptor;
  }
  private static  com.google.protobuf.Descriptors.FileDescriptor
      descriptor;
  static {
    java.lang.String[] descriptorData = {
      "\n\033stations/station_full.proto\032\037google/pr" +
      "otobuf/timestamp.proto\032%stations/station" +
      "_identification.proto\"\237\001\n\014station_full\0229" +
      "\n\ntime_stamp\030\001 \001(\0132\032.google.protobuf.Tim" +
      "estampR\ttimeStamp\0221\n\007station\030\002 \001(\0132\027.sta" +
      "tion_identificationR\007station\022!\n\014max_capa" +
      "city\030\003 \001(\005R\013maxCapacityB5B\020StationFullPr" +
      "otoP\001Z\037stage2024/pkg/protogen/stationsb\006" +
      "proto3"
    };
    descriptor = com.google.protobuf.Descriptors.FileDescriptor
      .internalBuildGeneratedFileFrom(descriptorData,
        new com.google.protobuf.Descriptors.FileDescriptor[] {
          com.google.protobuf.TimestampProto.getDescriptor(),
          StationIdentificationProto.getDescriptor(),
        });
    internal_static_station_full_descriptor =
      getDescriptor().getMessageTypes().get(0);
    internal_static_station_full_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessage.FieldAccessorTable(
        internal_static_station_full_descriptor,
        new java.lang.String[] { "TimeStamp", "Station", "MaxCapacity", });
    descriptor.resolveAllFeaturesImmutable();
    com.google.protobuf.TimestampProto.getDescriptor();
    StationIdentificationProto.getDescriptor();
  }

  // @@protoc_insertion_point(outer_class_scope)
}

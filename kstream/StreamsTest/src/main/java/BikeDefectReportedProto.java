// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: bikes/bike_defect_reported.proto
// Protobuf Java Version: 4.26.1

public final class BikeDefectReportedProto {
  private BikeDefectReportedProto() {}
  static {
    com.google.protobuf.RuntimeVersion.validateProtobufGencodeVersion(
      com.google.protobuf.RuntimeVersion.RuntimeDomain.PUBLIC,
      /* major= */ 4,
      /* minor= */ 26,
      /* patch= */ 1,
      /* suffix= */ "",
      BikeDefectReportedProto.class.getName());
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
    internal_static_bike_defect_reported_descriptor;
  static final 
    com.google.protobuf.GeneratedMessage.FieldAccessorTable
      internal_static_bike_defect_reported_fieldAccessorTable;
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_defect_bike_descriptor;
  static final 
    com.google.protobuf.GeneratedMessage.FieldAccessorTable
      internal_static_defect_bike_fieldAccessorTable;

  public static com.google.protobuf.Descriptors.FileDescriptor
      getDescriptor() {
    return descriptor;
  }
  private static  com.google.protobuf.Descriptors.FileDescriptor
      descriptor;
  static {
    java.lang.String[] descriptorData = {
      "\n bikes/bike_defect_reported.proto\032\037goog" +
      "le/protobuf/timestamp.proto\032\037bikes/bike_" +
      "identification.proto\032\037users/user_identif" +
      "ication.proto\032\025common/location.proto\"\306\001\n" +
      "\024bike_defect_reported\0229\n\ntime_stamp\030\001 \001(" +
      "\0132\032.google.protobuf.TimestampR\ttimeStamp" +
      "\022 \n\004bike\030\002 \001(\0132\014.defect_bikeR\004bike\022(\n\004us" +
      "er\030\003 \001(\0132\024.user_identificationR\004user\022\'\n\017" +
      "reported_defect\030\004 \001(\tR\016reportedDefect\"\246\001" +
      "\n\013defect_bike\022(\n\004bike\030\001 \001(\0132\024.bike_ident" +
      "ificationR\004bike\022%\n\010location\030\002 \001(\0132\t.loca" +
      "tionR\010location\022\037\n\013is_electric\030\003 \001(\010R\nisE" +
      "lectric\022%\n\016is_immobilized\030\004 \001(\010R\risImmob" +
      "ilizedB9B\027BikeDefectReportedProtoP\001Z\034sta" +
      "ge2024/pkg/protogen/bikesb\006proto3"
    };
    descriptor = com.google.protobuf.Descriptors.FileDescriptor
      .internalBuildGeneratedFileFrom(descriptorData,
        new com.google.protobuf.Descriptors.FileDescriptor[] {
          com.google.protobuf.TimestampProto.getDescriptor(),
          BikeIdentificationProto.getDescriptor(),
          UserIdentificationProto.getDescriptor(),
          LocationProto.getDescriptor(),
        });
    internal_static_bike_defect_reported_descriptor =
      getDescriptor().getMessageTypes().get(0);
    internal_static_bike_defect_reported_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessage.FieldAccessorTable(
        internal_static_bike_defect_reported_descriptor,
        new java.lang.String[] { "TimeStamp", "Bike", "User", "ReportedDefect", });
    internal_static_defect_bike_descriptor =
      getDescriptor().getMessageTypes().get(1);
    internal_static_defect_bike_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessage.FieldAccessorTable(
        internal_static_defect_bike_descriptor,
        new java.lang.String[] { "Bike", "Location", "IsElectric", "IsImmobilized", });
    descriptor.resolveAllFeaturesImmutable();
    com.google.protobuf.TimestampProto.getDescriptor();
    BikeIdentificationProto.getDescriptor();
    UserIdentificationProto.getDescriptor();
    LocationProto.getDescriptor();
  }

  // @@protoc_insertion_point(outer_class_scope)
}

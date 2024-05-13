// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: stations/station_deprecated.proto

// Protobuf Java Version: 4.26.1
/**
 * Protobuf type {@code deprecated_station}
 */
public final class deprecated_station extends
    com.google.protobuf.GeneratedMessage implements
    // @@protoc_insertion_point(message_implements:deprecated_station)
    deprecated_stationOrBuilder {
private static final long serialVersionUID = 0L;
  static {
    com.google.protobuf.RuntimeVersion.validateProtobufGencodeVersion(
      com.google.protobuf.RuntimeVersion.RuntimeDomain.PUBLIC,
      /* major= */ 4,
      /* minor= */ 26,
      /* patch= */ 1,
      /* suffix= */ "",
      deprecated_station.class.getName());
  }
  // Use deprecated_station.newBuilder() to construct.
  private deprecated_station(com.google.protobuf.GeneratedMessage.Builder<?> builder) {
    super(builder);
  }
  private deprecated_station() {
  }

  public static final com.google.protobuf.Descriptors.Descriptor
      getDescriptor() {
    return StationDeprecatedProto.internal_static_deprecated_station_descriptor;
  }

  @java.lang.Override
  protected com.google.protobuf.GeneratedMessage.FieldAccessorTable
      internalGetFieldAccessorTable() {
    return StationDeprecatedProto.internal_static_deprecated_station_fieldAccessorTable
        .ensureFieldAccessorsInitialized(
            deprecated_station.class, deprecated_station.Builder.class);
  }

  private int bitField0_;
  public static final int STATION_FIELD_NUMBER = 1;
  private station_identification station_;
  /**
   * <code>.station_identification station = 1 [json_name = "station"];</code>
   * @return Whether the station field is set.
   */
  @java.lang.Override
  public boolean hasStation() {
    return ((bitField0_ & 0x00000001) != 0);
  }
  /**
   * <code>.station_identification station = 1 [json_name = "station"];</code>
   * @return The station.
   */
  @java.lang.Override
  public station_identification getStation() {
    return station_ == null ? station_identification.getDefaultInstance() : station_;
  }
  /**
   * <code>.station_identification station = 1 [json_name = "station"];</code>
   */
  @java.lang.Override
  public station_identificationOrBuilder getStationOrBuilder() {
    return station_ == null ? station_identification.getDefaultInstance() : station_;
  }

  public static final int IS_ACTIVE_FIELD_NUMBER = 2;
  private boolean isActive_ = false;
  /**
   * <code>bool is_active = 2 [json_name = "isActive"];</code>
   * @return The isActive.
   */
  @java.lang.Override
  public boolean getIsActive() {
    return isActive_;
  }

  private byte memoizedIsInitialized = -1;
  @java.lang.Override
  public final boolean isInitialized() {
    byte isInitialized = memoizedIsInitialized;
    if (isInitialized == 1) return true;
    if (isInitialized == 0) return false;

    memoizedIsInitialized = 1;
    return true;
  }

  @java.lang.Override
  public void writeTo(com.google.protobuf.CodedOutputStream output)
                      throws java.io.IOException {
    if (((bitField0_ & 0x00000001) != 0)) {
      output.writeMessage(1, getStation());
    }
    if (isActive_ != false) {
      output.writeBool(2, isActive_);
    }
    getUnknownFields().writeTo(output);
  }

  @java.lang.Override
  public int getSerializedSize() {
    int size = memoizedSize;
    if (size != -1) return size;

    size = 0;
    if (((bitField0_ & 0x00000001) != 0)) {
      size += com.google.protobuf.CodedOutputStream
        .computeMessageSize(1, getStation());
    }
    if (isActive_ != false) {
      size += com.google.protobuf.CodedOutputStream
        .computeBoolSize(2, isActive_);
    }
    size += getUnknownFields().getSerializedSize();
    memoizedSize = size;
    return size;
  }

  @java.lang.Override
  public boolean equals(final java.lang.Object obj) {
    if (obj == this) {
     return true;
    }
    if (!(obj instanceof deprecated_station)) {
      return super.equals(obj);
    }
    deprecated_station other = (deprecated_station) obj;

    if (hasStation() != other.hasStation()) return false;
    if (hasStation()) {
      if (!getStation()
          .equals(other.getStation())) return false;
    }
    if (getIsActive()
        != other.getIsActive()) return false;
    if (!getUnknownFields().equals(other.getUnknownFields())) return false;
    return true;
  }

  @java.lang.Override
  public int hashCode() {
    if (memoizedHashCode != 0) {
      return memoizedHashCode;
    }
    int hash = 41;
    hash = (19 * hash) + getDescriptor().hashCode();
    if (hasStation()) {
      hash = (37 * hash) + STATION_FIELD_NUMBER;
      hash = (53 * hash) + getStation().hashCode();
    }
    hash = (37 * hash) + IS_ACTIVE_FIELD_NUMBER;
    hash = (53 * hash) + com.google.protobuf.Internal.hashBoolean(
        getIsActive());
    hash = (29 * hash) + getUnknownFields().hashCode();
    memoizedHashCode = hash;
    return hash;
  }

  public static deprecated_station parseFrom(
      java.nio.ByteBuffer data)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data);
  }
  public static deprecated_station parseFrom(
      java.nio.ByteBuffer data,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data, extensionRegistry);
  }
  public static deprecated_station parseFrom(
      com.google.protobuf.ByteString data)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data);
  }
  public static deprecated_station parseFrom(
      com.google.protobuf.ByteString data,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data, extensionRegistry);
  }
  public static deprecated_station parseFrom(byte[] data)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data);
  }
  public static deprecated_station parseFrom(
      byte[] data,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data, extensionRegistry);
  }
  public static deprecated_station parseFrom(java.io.InputStream input)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessage
        .parseWithIOException(PARSER, input);
  }
  public static deprecated_station parseFrom(
      java.io.InputStream input,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessage
        .parseWithIOException(PARSER, input, extensionRegistry);
  }

  public static deprecated_station parseDelimitedFrom(java.io.InputStream input)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessage
        .parseDelimitedWithIOException(PARSER, input);
  }

  public static deprecated_station parseDelimitedFrom(
      java.io.InputStream input,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessage
        .parseDelimitedWithIOException(PARSER, input, extensionRegistry);
  }
  public static deprecated_station parseFrom(
      com.google.protobuf.CodedInputStream input)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessage
        .parseWithIOException(PARSER, input);
  }
  public static deprecated_station parseFrom(
      com.google.protobuf.CodedInputStream input,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessage
        .parseWithIOException(PARSER, input, extensionRegistry);
  }

  @java.lang.Override
  public Builder newBuilderForType() { return newBuilder(); }
  public static Builder newBuilder() {
    return DEFAULT_INSTANCE.toBuilder();
  }
  public static Builder newBuilder(deprecated_station prototype) {
    return DEFAULT_INSTANCE.toBuilder().mergeFrom(prototype);
  }
  @java.lang.Override
  public Builder toBuilder() {
    return this == DEFAULT_INSTANCE
        ? new Builder() : new Builder().mergeFrom(this);
  }

  @java.lang.Override
  protected Builder newBuilderForType(
      com.google.protobuf.GeneratedMessage.BuilderParent parent) {
    Builder builder = new Builder(parent);
    return builder;
  }
  /**
   * Protobuf type {@code deprecated_station}
   */
  public static final class Builder extends
      com.google.protobuf.GeneratedMessage.Builder<Builder> implements
      // @@protoc_insertion_point(builder_implements:deprecated_station)
      deprecated_stationOrBuilder {
    public static final com.google.protobuf.Descriptors.Descriptor
        getDescriptor() {
      return StationDeprecatedProto.internal_static_deprecated_station_descriptor;
    }

    @java.lang.Override
    protected com.google.protobuf.GeneratedMessage.FieldAccessorTable
        internalGetFieldAccessorTable() {
      return StationDeprecatedProto.internal_static_deprecated_station_fieldAccessorTable
          .ensureFieldAccessorsInitialized(
              deprecated_station.class, deprecated_station.Builder.class);
    }

    // Construct using deprecated_station.newBuilder()
    private Builder() {
      maybeForceBuilderInitialization();
    }

    private Builder(
        com.google.protobuf.GeneratedMessage.BuilderParent parent) {
      super(parent);
      maybeForceBuilderInitialization();
    }
    private void maybeForceBuilderInitialization() {
      if (com.google.protobuf.GeneratedMessage
              .alwaysUseFieldBuilders) {
        getStationFieldBuilder();
      }
    }
    @java.lang.Override
    public Builder clear() {
      super.clear();
      bitField0_ = 0;
      station_ = null;
      if (stationBuilder_ != null) {
        stationBuilder_.dispose();
        stationBuilder_ = null;
      }
      isActive_ = false;
      return this;
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.Descriptor
        getDescriptorForType() {
      return StationDeprecatedProto.internal_static_deprecated_station_descriptor;
    }

    @java.lang.Override
    public deprecated_station getDefaultInstanceForType() {
      return deprecated_station.getDefaultInstance();
    }

    @java.lang.Override
    public deprecated_station build() {
      deprecated_station result = buildPartial();
      if (!result.isInitialized()) {
        throw newUninitializedMessageException(result);
      }
      return result;
    }

    @java.lang.Override
    public deprecated_station buildPartial() {
      deprecated_station result = new deprecated_station(this);
      if (bitField0_ != 0) { buildPartial0(result); }
      onBuilt();
      return result;
    }

    private void buildPartial0(deprecated_station result) {
      int from_bitField0_ = bitField0_;
      int to_bitField0_ = 0;
      if (((from_bitField0_ & 0x00000001) != 0)) {
        result.station_ = stationBuilder_ == null
            ? station_
            : stationBuilder_.build();
        to_bitField0_ |= 0x00000001;
      }
      if (((from_bitField0_ & 0x00000002) != 0)) {
        result.isActive_ = isActive_;
      }
      result.bitField0_ |= to_bitField0_;
    }

    @java.lang.Override
    public Builder mergeFrom(com.google.protobuf.Message other) {
      if (other instanceof deprecated_station) {
        return mergeFrom((deprecated_station)other);
      } else {
        super.mergeFrom(other);
        return this;
      }
    }

    public Builder mergeFrom(deprecated_station other) {
      if (other == deprecated_station.getDefaultInstance()) return this;
      if (other.hasStation()) {
        mergeStation(other.getStation());
      }
      if (other.getIsActive() != false) {
        setIsActive(other.getIsActive());
      }
      this.mergeUnknownFields(other.getUnknownFields());
      onChanged();
      return this;
    }

    @java.lang.Override
    public final boolean isInitialized() {
      return true;
    }

    @java.lang.Override
    public Builder mergeFrom(
        com.google.protobuf.CodedInputStream input,
        com.google.protobuf.ExtensionRegistryLite extensionRegistry)
        throws java.io.IOException {
      if (extensionRegistry == null) {
        throw new java.lang.NullPointerException();
      }
      try {
        boolean done = false;
        while (!done) {
          int tag = input.readTag();
          switch (tag) {
            case 0:
              done = true;
              break;
            case 10: {
              input.readMessage(
                  getStationFieldBuilder().getBuilder(),
                  extensionRegistry);
              bitField0_ |= 0x00000001;
              break;
            } // case 10
            case 16: {
              isActive_ = input.readBool();
              bitField0_ |= 0x00000002;
              break;
            } // case 16
            default: {
              if (!super.parseUnknownField(input, extensionRegistry, tag)) {
                done = true; // was an endgroup tag
              }
              break;
            } // default:
          } // switch (tag)
        } // while (!done)
      } catch (com.google.protobuf.InvalidProtocolBufferException e) {
        throw e.unwrapIOException();
      } finally {
        onChanged();
      } // finally
      return this;
    }
    private int bitField0_;

    private station_identification station_;
    private com.google.protobuf.SingleFieldBuilder<
        station_identification, station_identification.Builder, station_identificationOrBuilder> stationBuilder_;
    /**
     * <code>.station_identification station = 1 [json_name = "station"];</code>
     * @return Whether the station field is set.
     */
    public boolean hasStation() {
      return ((bitField0_ & 0x00000001) != 0);
    }
    /**
     * <code>.station_identification station = 1 [json_name = "station"];</code>
     * @return The station.
     */
    public station_identification getStation() {
      if (stationBuilder_ == null) {
        return station_ == null ? station_identification.getDefaultInstance() : station_;
      } else {
        return stationBuilder_.getMessage();
      }
    }
    /**
     * <code>.station_identification station = 1 [json_name = "station"];</code>
     */
    public Builder setStation(station_identification value) {
      if (stationBuilder_ == null) {
        if (value == null) {
          throw new NullPointerException();
        }
        station_ = value;
      } else {
        stationBuilder_.setMessage(value);
      }
      bitField0_ |= 0x00000001;
      onChanged();
      return this;
    }
    /**
     * <code>.station_identification station = 1 [json_name = "station"];</code>
     */
    public Builder setStation(
        station_identification.Builder builderForValue) {
      if (stationBuilder_ == null) {
        station_ = builderForValue.build();
      } else {
        stationBuilder_.setMessage(builderForValue.build());
      }
      bitField0_ |= 0x00000001;
      onChanged();
      return this;
    }
    /**
     * <code>.station_identification station = 1 [json_name = "station"];</code>
     */
    public Builder mergeStation(station_identification value) {
      if (stationBuilder_ == null) {
        if (((bitField0_ & 0x00000001) != 0) &&
          station_ != null &&
          station_ != station_identification.getDefaultInstance()) {
          getStationBuilder().mergeFrom(value);
        } else {
          station_ = value;
        }
      } else {
        stationBuilder_.mergeFrom(value);
      }
      if (station_ != null) {
        bitField0_ |= 0x00000001;
        onChanged();
      }
      return this;
    }
    /**
     * <code>.station_identification station = 1 [json_name = "station"];</code>
     */
    public Builder clearStation() {
      bitField0_ = (bitField0_ & ~0x00000001);
      station_ = null;
      if (stationBuilder_ != null) {
        stationBuilder_.dispose();
        stationBuilder_ = null;
      }
      onChanged();
      return this;
    }
    /**
     * <code>.station_identification station = 1 [json_name = "station"];</code>
     */
    public station_identification.Builder getStationBuilder() {
      bitField0_ |= 0x00000001;
      onChanged();
      return getStationFieldBuilder().getBuilder();
    }
    /**
     * <code>.station_identification station = 1 [json_name = "station"];</code>
     */
    public station_identificationOrBuilder getStationOrBuilder() {
      if (stationBuilder_ != null) {
        return stationBuilder_.getMessageOrBuilder();
      } else {
        return station_ == null ?
            station_identification.getDefaultInstance() : station_;
      }
    }
    /**
     * <code>.station_identification station = 1 [json_name = "station"];</code>
     */
    private com.google.protobuf.SingleFieldBuilder<
        station_identification, station_identification.Builder, station_identificationOrBuilder> 
        getStationFieldBuilder() {
      if (stationBuilder_ == null) {
        stationBuilder_ = new com.google.protobuf.SingleFieldBuilder<
            station_identification, station_identification.Builder, station_identificationOrBuilder>(
                getStation(),
                getParentForChildren(),
                isClean());
        station_ = null;
      }
      return stationBuilder_;
    }

    private boolean isActive_ ;
    /**
     * <code>bool is_active = 2 [json_name = "isActive"];</code>
     * @return The isActive.
     */
    @java.lang.Override
    public boolean getIsActive() {
      return isActive_;
    }
    /**
     * <code>bool is_active = 2 [json_name = "isActive"];</code>
     * @param value The isActive to set.
     * @return This builder for chaining.
     */
    public Builder setIsActive(boolean value) {

      isActive_ = value;
      bitField0_ |= 0x00000002;
      onChanged();
      return this;
    }
    /**
     * <code>bool is_active = 2 [json_name = "isActive"];</code>
     * @return This builder for chaining.
     */
    public Builder clearIsActive() {
      bitField0_ = (bitField0_ & ~0x00000002);
      isActive_ = false;
      onChanged();
      return this;
    }

    // @@protoc_insertion_point(builder_scope:deprecated_station)
  }

  // @@protoc_insertion_point(class_scope:deprecated_station)
  private static final deprecated_station DEFAULT_INSTANCE;
  static {
    DEFAULT_INSTANCE = new deprecated_station();
  }

  public static deprecated_station getDefaultInstance() {
    return DEFAULT_INSTANCE;
  }

  private static final com.google.protobuf.Parser<deprecated_station>
      PARSER = new com.google.protobuf.AbstractParser<deprecated_station>() {
    @java.lang.Override
    public deprecated_station parsePartialFrom(
        com.google.protobuf.CodedInputStream input,
        com.google.protobuf.ExtensionRegistryLite extensionRegistry)
        throws com.google.protobuf.InvalidProtocolBufferException {
      Builder builder = newBuilder();
      try {
        builder.mergeFrom(input, extensionRegistry);
      } catch (com.google.protobuf.InvalidProtocolBufferException e) {
        throw e.setUnfinishedMessage(builder.buildPartial());
      } catch (com.google.protobuf.UninitializedMessageException e) {
        throw e.asInvalidProtocolBufferException().setUnfinishedMessage(builder.buildPartial());
      } catch (java.io.IOException e) {
        throw new com.google.protobuf.InvalidProtocolBufferException(e)
            .setUnfinishedMessage(builder.buildPartial());
      }
      return builder.buildPartial();
    }
  };

  public static com.google.protobuf.Parser<deprecated_station> parser() {
    return PARSER;
  }

  @java.lang.Override
  public com.google.protobuf.Parser<deprecated_station> getParserForType() {
    return PARSER;
  }

  @java.lang.Override
  public deprecated_station getDefaultInstanceForType() {
    return DEFAULT_INSTANCE;
  }

}


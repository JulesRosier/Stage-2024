// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: stations/station_full.proto

// Protobuf Java Version: 4.26.1
/**
 * Protobuf type {@code station_full}
 */
public final class station_full extends
    com.google.protobuf.GeneratedMessage implements
    // @@protoc_insertion_point(message_implements:station_full)
    station_fullOrBuilder {
private static final long serialVersionUID = 0L;
  static {
    com.google.protobuf.RuntimeVersion.validateProtobufGencodeVersion(
      com.google.protobuf.RuntimeVersion.RuntimeDomain.PUBLIC,
      /* major= */ 4,
      /* minor= */ 26,
      /* patch= */ 1,
      /* suffix= */ "",
      station_full.class.getName());
  }
  // Use station_full.newBuilder() to construct.
  private station_full(com.google.protobuf.GeneratedMessage.Builder<?> builder) {
    super(builder);
  }
  private station_full() {
  }

  public static final com.google.protobuf.Descriptors.Descriptor
      getDescriptor() {
    return StationFullProto.internal_static_station_full_descriptor;
  }

  @java.lang.Override
  protected com.google.protobuf.GeneratedMessage.FieldAccessorTable
      internalGetFieldAccessorTable() {
    return StationFullProto.internal_static_station_full_fieldAccessorTable
        .ensureFieldAccessorsInitialized(
            station_full.class, station_full.Builder.class);
  }

  private int bitField0_;
  public static final int TIME_STAMP_FIELD_NUMBER = 1;
  private com.google.protobuf.Timestamp timeStamp_;
  /**
   * <code>.google.protobuf.Timestamp time_stamp = 1 [json_name = "timeStamp"];</code>
   * @return Whether the timeStamp field is set.
   */
  @java.lang.Override
  public boolean hasTimeStamp() {
    return ((bitField0_ & 0x00000001) != 0);
  }
  /**
   * <code>.google.protobuf.Timestamp time_stamp = 1 [json_name = "timeStamp"];</code>
   * @return The timeStamp.
   */
  @java.lang.Override
  public com.google.protobuf.Timestamp getTimeStamp() {
    return timeStamp_ == null ? com.google.protobuf.Timestamp.getDefaultInstance() : timeStamp_;
  }
  /**
   * <code>.google.protobuf.Timestamp time_stamp = 1 [json_name = "timeStamp"];</code>
   */
  @java.lang.Override
  public com.google.protobuf.TimestampOrBuilder getTimeStampOrBuilder() {
    return timeStamp_ == null ? com.google.protobuf.Timestamp.getDefaultInstance() : timeStamp_;
  }

  public static final int STATION_FIELD_NUMBER = 2;
  private station_identification station_;
  /**
   * <code>.station_identification station = 2 [json_name = "station"];</code>
   * @return Whether the station field is set.
   */
  @java.lang.Override
  public boolean hasStation() {
    return ((bitField0_ & 0x00000002) != 0);
  }
  /**
   * <code>.station_identification station = 2 [json_name = "station"];</code>
   * @return The station.
   */
  @java.lang.Override
  public station_identification getStation() {
    return station_ == null ? station_identification.getDefaultInstance() : station_;
  }
  /**
   * <code>.station_identification station = 2 [json_name = "station"];</code>
   */
  @java.lang.Override
  public station_identificationOrBuilder getStationOrBuilder() {
    return station_ == null ? station_identification.getDefaultInstance() : station_;
  }

  public static final int MAX_CAPACITY_FIELD_NUMBER = 3;
  private int maxCapacity_ = 0;
  /**
   * <code>int32 max_capacity = 3 [json_name = "maxCapacity"];</code>
   * @return The maxCapacity.
   */
  @java.lang.Override
  public int getMaxCapacity() {
    return maxCapacity_;
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
      output.writeMessage(1, getTimeStamp());
    }
    if (((bitField0_ & 0x00000002) != 0)) {
      output.writeMessage(2, getStation());
    }
    if (maxCapacity_ != 0) {
      output.writeInt32(3, maxCapacity_);
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
        .computeMessageSize(1, getTimeStamp());
    }
    if (((bitField0_ & 0x00000002) != 0)) {
      size += com.google.protobuf.CodedOutputStream
        .computeMessageSize(2, getStation());
    }
    if (maxCapacity_ != 0) {
      size += com.google.protobuf.CodedOutputStream
        .computeInt32Size(3, maxCapacity_);
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
    if (!(obj instanceof station_full)) {
      return super.equals(obj);
    }
    station_full other = (station_full) obj;

    if (hasTimeStamp() != other.hasTimeStamp()) return false;
    if (hasTimeStamp()) {
      if (!getTimeStamp()
          .equals(other.getTimeStamp())) return false;
    }
    if (hasStation() != other.hasStation()) return false;
    if (hasStation()) {
      if (!getStation()
          .equals(other.getStation())) return false;
    }
    if (getMaxCapacity()
        != other.getMaxCapacity()) return false;
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
    if (hasTimeStamp()) {
      hash = (37 * hash) + TIME_STAMP_FIELD_NUMBER;
      hash = (53 * hash) + getTimeStamp().hashCode();
    }
    if (hasStation()) {
      hash = (37 * hash) + STATION_FIELD_NUMBER;
      hash = (53 * hash) + getStation().hashCode();
    }
    hash = (37 * hash) + MAX_CAPACITY_FIELD_NUMBER;
    hash = (53 * hash) + getMaxCapacity();
    hash = (29 * hash) + getUnknownFields().hashCode();
    memoizedHashCode = hash;
    return hash;
  }

  public static station_full parseFrom(
      java.nio.ByteBuffer data)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data);
  }
  public static station_full parseFrom(
      java.nio.ByteBuffer data,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data, extensionRegistry);
  }
  public static station_full parseFrom(
      com.google.protobuf.ByteString data)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data);
  }
  public static station_full parseFrom(
      com.google.protobuf.ByteString data,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data, extensionRegistry);
  }
  public static station_full parseFrom(byte[] data)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data);
  }
  public static station_full parseFrom(
      byte[] data,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data, extensionRegistry);
  }
  public static station_full parseFrom(java.io.InputStream input)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessage
        .parseWithIOException(PARSER, input);
  }
  public static station_full parseFrom(
      java.io.InputStream input,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessage
        .parseWithIOException(PARSER, input, extensionRegistry);
  }

  public static station_full parseDelimitedFrom(java.io.InputStream input)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessage
        .parseDelimitedWithIOException(PARSER, input);
  }

  public static station_full parseDelimitedFrom(
      java.io.InputStream input,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessage
        .parseDelimitedWithIOException(PARSER, input, extensionRegistry);
  }
  public static station_full parseFrom(
      com.google.protobuf.CodedInputStream input)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessage
        .parseWithIOException(PARSER, input);
  }
  public static station_full parseFrom(
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
  public static Builder newBuilder(station_full prototype) {
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
   * Protobuf type {@code station_full}
   */
  public static final class Builder extends
      com.google.protobuf.GeneratedMessage.Builder<Builder> implements
      // @@protoc_insertion_point(builder_implements:station_full)
      station_fullOrBuilder {
    public static final com.google.protobuf.Descriptors.Descriptor
        getDescriptor() {
      return StationFullProto.internal_static_station_full_descriptor;
    }

    @java.lang.Override
    protected com.google.protobuf.GeneratedMessage.FieldAccessorTable
        internalGetFieldAccessorTable() {
      return StationFullProto.internal_static_station_full_fieldAccessorTable
          .ensureFieldAccessorsInitialized(
              station_full.class, station_full.Builder.class);
    }

    // Construct using station_full.newBuilder()
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
        getTimeStampFieldBuilder();
        getStationFieldBuilder();
      }
    }
    @java.lang.Override
    public Builder clear() {
      super.clear();
      bitField0_ = 0;
      timeStamp_ = null;
      if (timeStampBuilder_ != null) {
        timeStampBuilder_.dispose();
        timeStampBuilder_ = null;
      }
      station_ = null;
      if (stationBuilder_ != null) {
        stationBuilder_.dispose();
        stationBuilder_ = null;
      }
      maxCapacity_ = 0;
      return this;
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.Descriptor
        getDescriptorForType() {
      return StationFullProto.internal_static_station_full_descriptor;
    }

    @java.lang.Override
    public station_full getDefaultInstanceForType() {
      return station_full.getDefaultInstance();
    }

    @java.lang.Override
    public station_full build() {
      station_full result = buildPartial();
      if (!result.isInitialized()) {
        throw newUninitializedMessageException(result);
      }
      return result;
    }

    @java.lang.Override
    public station_full buildPartial() {
      station_full result = new station_full(this);
      if (bitField0_ != 0) { buildPartial0(result); }
      onBuilt();
      return result;
    }

    private void buildPartial0(station_full result) {
      int from_bitField0_ = bitField0_;
      int to_bitField0_ = 0;
      if (((from_bitField0_ & 0x00000001) != 0)) {
        result.timeStamp_ = timeStampBuilder_ == null
            ? timeStamp_
            : timeStampBuilder_.build();
        to_bitField0_ |= 0x00000001;
      }
      if (((from_bitField0_ & 0x00000002) != 0)) {
        result.station_ = stationBuilder_ == null
            ? station_
            : stationBuilder_.build();
        to_bitField0_ |= 0x00000002;
      }
      if (((from_bitField0_ & 0x00000004) != 0)) {
        result.maxCapacity_ = maxCapacity_;
      }
      result.bitField0_ |= to_bitField0_;
    }

    @java.lang.Override
    public Builder mergeFrom(com.google.protobuf.Message other) {
      if (other instanceof station_full) {
        return mergeFrom((station_full)other);
      } else {
        super.mergeFrom(other);
        return this;
      }
    }

    public Builder mergeFrom(station_full other) {
      if (other == station_full.getDefaultInstance()) return this;
      if (other.hasTimeStamp()) {
        mergeTimeStamp(other.getTimeStamp());
      }
      if (other.hasStation()) {
        mergeStation(other.getStation());
      }
      if (other.getMaxCapacity() != 0) {
        setMaxCapacity(other.getMaxCapacity());
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
                  getTimeStampFieldBuilder().getBuilder(),
                  extensionRegistry);
              bitField0_ |= 0x00000001;
              break;
            } // case 10
            case 18: {
              input.readMessage(
                  getStationFieldBuilder().getBuilder(),
                  extensionRegistry);
              bitField0_ |= 0x00000002;
              break;
            } // case 18
            case 24: {
              maxCapacity_ = input.readInt32();
              bitField0_ |= 0x00000004;
              break;
            } // case 24
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

    private com.google.protobuf.Timestamp timeStamp_;
    private com.google.protobuf.SingleFieldBuilder<
        com.google.protobuf.Timestamp, com.google.protobuf.Timestamp.Builder, com.google.protobuf.TimestampOrBuilder> timeStampBuilder_;
    /**
     * <code>.google.protobuf.Timestamp time_stamp = 1 [json_name = "timeStamp"];</code>
     * @return Whether the timeStamp field is set.
     */
    public boolean hasTimeStamp() {
      return ((bitField0_ & 0x00000001) != 0);
    }
    /**
     * <code>.google.protobuf.Timestamp time_stamp = 1 [json_name = "timeStamp"];</code>
     * @return The timeStamp.
     */
    public com.google.protobuf.Timestamp getTimeStamp() {
      if (timeStampBuilder_ == null) {
        return timeStamp_ == null ? com.google.protobuf.Timestamp.getDefaultInstance() : timeStamp_;
      } else {
        return timeStampBuilder_.getMessage();
      }
    }
    /**
     * <code>.google.protobuf.Timestamp time_stamp = 1 [json_name = "timeStamp"];</code>
     */
    public Builder setTimeStamp(com.google.protobuf.Timestamp value) {
      if (timeStampBuilder_ == null) {
        if (value == null) {
          throw new NullPointerException();
        }
        timeStamp_ = value;
      } else {
        timeStampBuilder_.setMessage(value);
      }
      bitField0_ |= 0x00000001;
      onChanged();
      return this;
    }
    /**
     * <code>.google.protobuf.Timestamp time_stamp = 1 [json_name = "timeStamp"];</code>
     */
    public Builder setTimeStamp(
        com.google.protobuf.Timestamp.Builder builderForValue) {
      if (timeStampBuilder_ == null) {
        timeStamp_ = builderForValue.build();
      } else {
        timeStampBuilder_.setMessage(builderForValue.build());
      }
      bitField0_ |= 0x00000001;
      onChanged();
      return this;
    }
    /**
     * <code>.google.protobuf.Timestamp time_stamp = 1 [json_name = "timeStamp"];</code>
     */
    public Builder mergeTimeStamp(com.google.protobuf.Timestamp value) {
      if (timeStampBuilder_ == null) {
        if (((bitField0_ & 0x00000001) != 0) &&
          timeStamp_ != null &&
          timeStamp_ != com.google.protobuf.Timestamp.getDefaultInstance()) {
          getTimeStampBuilder().mergeFrom(value);
        } else {
          timeStamp_ = value;
        }
      } else {
        timeStampBuilder_.mergeFrom(value);
      }
      if (timeStamp_ != null) {
        bitField0_ |= 0x00000001;
        onChanged();
      }
      return this;
    }
    /**
     * <code>.google.protobuf.Timestamp time_stamp = 1 [json_name = "timeStamp"];</code>
     */
    public Builder clearTimeStamp() {
      bitField0_ = (bitField0_ & ~0x00000001);
      timeStamp_ = null;
      if (timeStampBuilder_ != null) {
        timeStampBuilder_.dispose();
        timeStampBuilder_ = null;
      }
      onChanged();
      return this;
    }
    /**
     * <code>.google.protobuf.Timestamp time_stamp = 1 [json_name = "timeStamp"];</code>
     */
    public com.google.protobuf.Timestamp.Builder getTimeStampBuilder() {
      bitField0_ |= 0x00000001;
      onChanged();
      return getTimeStampFieldBuilder().getBuilder();
    }
    /**
     * <code>.google.protobuf.Timestamp time_stamp = 1 [json_name = "timeStamp"];</code>
     */
    public com.google.protobuf.TimestampOrBuilder getTimeStampOrBuilder() {
      if (timeStampBuilder_ != null) {
        return timeStampBuilder_.getMessageOrBuilder();
      } else {
        return timeStamp_ == null ?
            com.google.protobuf.Timestamp.getDefaultInstance() : timeStamp_;
      }
    }
    /**
     * <code>.google.protobuf.Timestamp time_stamp = 1 [json_name = "timeStamp"];</code>
     */
    private com.google.protobuf.SingleFieldBuilder<
        com.google.protobuf.Timestamp, com.google.protobuf.Timestamp.Builder, com.google.protobuf.TimestampOrBuilder> 
        getTimeStampFieldBuilder() {
      if (timeStampBuilder_ == null) {
        timeStampBuilder_ = new com.google.protobuf.SingleFieldBuilder<
            com.google.protobuf.Timestamp, com.google.protobuf.Timestamp.Builder, com.google.protobuf.TimestampOrBuilder>(
                getTimeStamp(),
                getParentForChildren(),
                isClean());
        timeStamp_ = null;
      }
      return timeStampBuilder_;
    }

    private station_identification station_;
    private com.google.protobuf.SingleFieldBuilder<
        station_identification, station_identification.Builder, station_identificationOrBuilder> stationBuilder_;
    /**
     * <code>.station_identification station = 2 [json_name = "station"];</code>
     * @return Whether the station field is set.
     */
    public boolean hasStation() {
      return ((bitField0_ & 0x00000002) != 0);
    }
    /**
     * <code>.station_identification station = 2 [json_name = "station"];</code>
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
     * <code>.station_identification station = 2 [json_name = "station"];</code>
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
      bitField0_ |= 0x00000002;
      onChanged();
      return this;
    }
    /**
     * <code>.station_identification station = 2 [json_name = "station"];</code>
     */
    public Builder setStation(
        station_identification.Builder builderForValue) {
      if (stationBuilder_ == null) {
        station_ = builderForValue.build();
      } else {
        stationBuilder_.setMessage(builderForValue.build());
      }
      bitField0_ |= 0x00000002;
      onChanged();
      return this;
    }
    /**
     * <code>.station_identification station = 2 [json_name = "station"];</code>
     */
    public Builder mergeStation(station_identification value) {
      if (stationBuilder_ == null) {
        if (((bitField0_ & 0x00000002) != 0) &&
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
        bitField0_ |= 0x00000002;
        onChanged();
      }
      return this;
    }
    /**
     * <code>.station_identification station = 2 [json_name = "station"];</code>
     */
    public Builder clearStation() {
      bitField0_ = (bitField0_ & ~0x00000002);
      station_ = null;
      if (stationBuilder_ != null) {
        stationBuilder_.dispose();
        stationBuilder_ = null;
      }
      onChanged();
      return this;
    }
    /**
     * <code>.station_identification station = 2 [json_name = "station"];</code>
     */
    public station_identification.Builder getStationBuilder() {
      bitField0_ |= 0x00000002;
      onChanged();
      return getStationFieldBuilder().getBuilder();
    }
    /**
     * <code>.station_identification station = 2 [json_name = "station"];</code>
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
     * <code>.station_identification station = 2 [json_name = "station"];</code>
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

    private int maxCapacity_ ;
    /**
     * <code>int32 max_capacity = 3 [json_name = "maxCapacity"];</code>
     * @return The maxCapacity.
     */
    @java.lang.Override
    public int getMaxCapacity() {
      return maxCapacity_;
    }
    /**
     * <code>int32 max_capacity = 3 [json_name = "maxCapacity"];</code>
     * @param value The maxCapacity to set.
     * @return This builder for chaining.
     */
    public Builder setMaxCapacity(int value) {

      maxCapacity_ = value;
      bitField0_ |= 0x00000004;
      onChanged();
      return this;
    }
    /**
     * <code>int32 max_capacity = 3 [json_name = "maxCapacity"];</code>
     * @return This builder for chaining.
     */
    public Builder clearMaxCapacity() {
      bitField0_ = (bitField0_ & ~0x00000004);
      maxCapacity_ = 0;
      onChanged();
      return this;
    }

    // @@protoc_insertion_point(builder_scope:station_full)
  }

  // @@protoc_insertion_point(class_scope:station_full)
  private static final station_full DEFAULT_INSTANCE;
  static {
    DEFAULT_INSTANCE = new station_full();
  }

  public static station_full getDefaultInstance() {
    return DEFAULT_INSTANCE;
  }

  private static final com.google.protobuf.Parser<station_full>
      PARSER = new com.google.protobuf.AbstractParser<station_full>() {
    @java.lang.Override
    public station_full parsePartialFrom(
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

  public static com.google.protobuf.Parser<station_full> parser() {
    return PARSER;
  }

  @java.lang.Override
  public com.google.protobuf.Parser<station_full> getParserForType() {
    return PARSER;
  }

  @java.lang.Override
  public station_full getDefaultInstanceForType() {
    return DEFAULT_INSTANCE;
  }

}


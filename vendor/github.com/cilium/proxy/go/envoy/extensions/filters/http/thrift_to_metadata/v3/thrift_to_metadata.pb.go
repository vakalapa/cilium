// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.23.0
// 	protoc        v5.29.2
// source: envoy/extensions/filters/http/thrift_to_metadata/v3/thrift_to_metadata.proto

package thrift_to_metadatav3

import (
	v3 "github.com/cilium/proxy/go/envoy/extensions/filters/network/thrift_proxy/v3"
	_ "github.com/cncf/xds/go/udpa/annotations"
	_ "github.com/cncf/xds/go/xds/annotations/v3"
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	structpb "google.golang.org/protobuf/types/known/structpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type Field int32

const (
	// The Thrift method name, string value.
	Field_METHOD_NAME Field = 0
	// The Thrift protocol name, string value. Values are "binary", "binary/non-strict", and "compact", with "(auto)" suffix if
	// :ref:`protocol <envoy_v3_api_field_extensions.filters.http.thrift_to_metadata.v3.ThriftToMetadata.protocol>`
	// is set to :ref:`AUTO_PROTOCOL<envoy_v3_api_enum_value_extensions.filters.network.thrift_proxy.v3.ProtocolType.AUTO_PROTOCOL>`
	Field_PROTOCOL Field = 1
	// The Thrift transport name, string value. Values are "framed", "header", and "unframed", with "(auto)" suffix if
	// :ref:`transport <envoy_v3_api_field_extensions.filters.http.thrift_to_metadata.v3.ThriftToMetadata.transport>`
	// is set to :ref:`AUTO_TRANSPORT<envoy_v3_api_enum_value_extensions.filters.network.thrift_proxy.v3.TransportType.AUTO_TRANSPORT>`
	Field_TRANSPORT Field = 2
	// The Thrift message type, singed 16-bit integer value.
	Field_HEADER_FLAGS Field = 3
	// The Thrift sequence ID, singed 32-bit integer value.
	Field_SEQUENCE_ID Field = 4
	// The Thrift message type, string value. Values in request are "call" and "oneway", and in response are "reply" and "exception".
	Field_MESSAGE_TYPE Field = 5
	// The Thrift reply type, string value. This is only valid for response rules. Values are "success" and "error".
	Field_REPLY_TYPE Field = 6
)

// Enum value maps for Field.
var (
	Field_name = map[int32]string{
		0: "METHOD_NAME",
		1: "PROTOCOL",
		2: "TRANSPORT",
		3: "HEADER_FLAGS",
		4: "SEQUENCE_ID",
		5: "MESSAGE_TYPE",
		6: "REPLY_TYPE",
	}
	Field_value = map[string]int32{
		"METHOD_NAME":  0,
		"PROTOCOL":     1,
		"TRANSPORT":    2,
		"HEADER_FLAGS": 3,
		"SEQUENCE_ID":  4,
		"MESSAGE_TYPE": 5,
		"REPLY_TYPE":   6,
	}
)

func (x Field) Enum() *Field {
	p := new(Field)
	*p = x
	return p
}

func (x Field) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Field) Descriptor() protoreflect.EnumDescriptor {
	return file_envoy_extensions_filters_http_thrift_to_metadata_v3_thrift_to_metadata_proto_enumTypes[0].Descriptor()
}

func (Field) Type() protoreflect.EnumType {
	return &file_envoy_extensions_filters_http_thrift_to_metadata_v3_thrift_to_metadata_proto_enumTypes[0]
}

func (x Field) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Field.Descriptor instead.
func (Field) EnumDescriptor() ([]byte, []int) {
	return file_envoy_extensions_filters_http_thrift_to_metadata_v3_thrift_to_metadata_proto_rawDescGZIP(), []int{0}
}

type KeyValuePair struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The namespace — if this is empty, the filter's namespace will be used.
	MetadataNamespace string `protobuf:"bytes,1,opt,name=metadata_namespace,json=metadataNamespace,proto3" json:"metadata_namespace,omitempty"`
	// The key to use within the namespace.
	Key string `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	// When used for on_present case, if value is non-empty it'll be used instead
	// of the field value.
	//
	// When used for on_missing case, a non-empty value must be provided.
	Value *structpb.Value `protobuf:"bytes,3,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *KeyValuePair) Reset() {
	*x = KeyValuePair{}
	if protoimpl.UnsafeEnabled {
		mi := &file_envoy_extensions_filters_http_thrift_to_metadata_v3_thrift_to_metadata_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *KeyValuePair) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KeyValuePair) ProtoMessage() {}

func (x *KeyValuePair) ProtoReflect() protoreflect.Message {
	mi := &file_envoy_extensions_filters_http_thrift_to_metadata_v3_thrift_to_metadata_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KeyValuePair.ProtoReflect.Descriptor instead.
func (*KeyValuePair) Descriptor() ([]byte, []int) {
	return file_envoy_extensions_filters_http_thrift_to_metadata_v3_thrift_to_metadata_proto_rawDescGZIP(), []int{0}
}

func (x *KeyValuePair) GetMetadataNamespace() string {
	if x != nil {
		return x.MetadataNamespace
	}
	return ""
}

func (x *KeyValuePair) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *KeyValuePair) GetValue() *structpb.Value {
	if x != nil {
		return x.Value
	}
	return nil
}

type FieldSelector struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// field name to log
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// field id to match
	Id int32 `protobuf:"varint,2,opt,name=id,proto3" json:"id,omitempty"`
	// next node of the field selector
	Child *FieldSelector `protobuf:"bytes,3,opt,name=child,proto3" json:"child,omitempty"`
}

func (x *FieldSelector) Reset() {
	*x = FieldSelector{}
	if protoimpl.UnsafeEnabled {
		mi := &file_envoy_extensions_filters_http_thrift_to_metadata_v3_thrift_to_metadata_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FieldSelector) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FieldSelector) ProtoMessage() {}

func (x *FieldSelector) ProtoReflect() protoreflect.Message {
	mi := &file_envoy_extensions_filters_http_thrift_to_metadata_v3_thrift_to_metadata_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FieldSelector.ProtoReflect.Descriptor instead.
func (*FieldSelector) Descriptor() ([]byte, []int) {
	return file_envoy_extensions_filters_http_thrift_to_metadata_v3_thrift_to_metadata_proto_rawDescGZIP(), []int{1}
}

func (x *FieldSelector) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *FieldSelector) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *FieldSelector) GetChild() *FieldSelector {
	if x != nil {
		return x.Child
	}
	return nil
}

// [#next-free-field: 6]
type Rule struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The field to match on. If set, takes precedence over field_selector.
	Field Field `protobuf:"varint,1,opt,name=field,proto3,enum=envoy.extensions.filters.http.thrift_to_metadata.v3.Field" json:"field,omitempty"`
	// Specifies that a match will be performed on the value of a field in the thrift body.
	// If set, the whole http body will be buffered to extract the field value, which
	// may have performance implications.
	//
	// It's a thrift over http version of
	// :ref:`field_selector<envoy_v3_api_field_extensions.filters.network.thrift_proxy.filters.payload_to_metadata.v3.PayloadToMetadata.Rule.field_selector>`.
	//
	// See also `payload-to-metadata <https://www.envoyproxy.io/docs/envoy/latest/configuration/other_protocols/thrift_filters/payload_to_metadata_filter>`_
	// for more reference.
	//
	// Example:
	//
	// .. code-block:: yaml
	//
	//	method_name: foo
	//	field_selector:
	//	  name: info
	//	  id: 2
	//	  child:
	//	    name: version
	//	    id: 1
	//
	// The above yaml will match on value of “info.version“ in the below thrift schema as input of
	// :ref:`on_present<envoy_v3_api_field_extensions.filters.http.thrift_to_metadata.v3.Rule.on_present>` or
	// :ref:`on_missing<envoy_v3_api_field_extensions.filters.http.thrift_to_metadata.v3.Rule.on_missing>`
	// while we are processing “foo“ method. This rule won't be applied to “bar“ method.
	//
	// .. code-block:: thrift
	//
	//	struct Info {
	//	  1: required string version;
	//	}
	//	service Server {
	//	  bool foo(1: i32 id, 2: Info info);
	//	  bool bar(1: i32 id, 2: Info info);
	//	}
	FieldSelector *FieldSelector `protobuf:"bytes,2,opt,name=field_selector,json=fieldSelector,proto3" json:"field_selector,omitempty"`
	// If specified, :ref:`field_selector<envoy_v3_api_field_extensions.filters.http.thrift_to_metadata.v3.Rule.field_selector>`
	// will be used to extract the field value *only* on the thrift message with method name.
	MethodName string `protobuf:"bytes,3,opt,name=method_name,json=methodName,proto3" json:"method_name,omitempty"`
	// The key-value pair to set in the *filter metadata* if the field is present
	// in *thrift metadata*.
	//
	// If the value in the KeyValuePair is non-empty, it'll be used instead
	// of field value.
	OnPresent *KeyValuePair `protobuf:"bytes,4,opt,name=on_present,json=onPresent,proto3" json:"on_present,omitempty"`
	// The key-value pair to set in the *filter metadata* if the field is missing
	// in *thrift metadata*.
	//
	// The value in the KeyValuePair must be set, since it'll be used in lieu
	// of the missing field value.
	OnMissing *KeyValuePair `protobuf:"bytes,5,opt,name=on_missing,json=onMissing,proto3" json:"on_missing,omitempty"`
}

func (x *Rule) Reset() {
	*x = Rule{}
	if protoimpl.UnsafeEnabled {
		mi := &file_envoy_extensions_filters_http_thrift_to_metadata_v3_thrift_to_metadata_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Rule) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Rule) ProtoMessage() {}

func (x *Rule) ProtoReflect() protoreflect.Message {
	mi := &file_envoy_extensions_filters_http_thrift_to_metadata_v3_thrift_to_metadata_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Rule.ProtoReflect.Descriptor instead.
func (*Rule) Descriptor() ([]byte, []int) {
	return file_envoy_extensions_filters_http_thrift_to_metadata_v3_thrift_to_metadata_proto_rawDescGZIP(), []int{2}
}

func (x *Rule) GetField() Field {
	if x != nil {
		return x.Field
	}
	return Field_METHOD_NAME
}

func (x *Rule) GetFieldSelector() *FieldSelector {
	if x != nil {
		return x.FieldSelector
	}
	return nil
}

func (x *Rule) GetMethodName() string {
	if x != nil {
		return x.MethodName
	}
	return ""
}

func (x *Rule) GetOnPresent() *KeyValuePair {
	if x != nil {
		return x.OnPresent
	}
	return nil
}

func (x *Rule) GetOnMissing() *KeyValuePair {
	if x != nil {
		return x.OnMissing
	}
	return nil
}

// The configuration for transforming thrift metadata into filter metadata.
//
// [#next-free-field: 7]
type ThriftToMetadata struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The list of rules to apply to http request body to extract thrift metadata.
	RequestRules []*Rule `protobuf:"bytes,1,rep,name=request_rules,json=requestRules,proto3" json:"request_rules,omitempty"`
	// The list of rules to apply to http response body to extract thrift metadata.
	ResponseRules []*Rule `protobuf:"bytes,2,rep,name=response_rules,json=responseRules,proto3" json:"response_rules,omitempty"`
	// Supplies the type of transport that the Thrift proxy should use. Defaults to
	// :ref:`AUTO_TRANSPORT<envoy_v3_api_enum_value_extensions.filters.network.thrift_proxy.v3.TransportType.AUTO_TRANSPORT>`.
	Transport v3.TransportType `protobuf:"varint,3,opt,name=transport,proto3,enum=envoy.extensions.filters.network.thrift_proxy.v3.TransportType" json:"transport,omitempty"`
	// Supplies the type of protocol that the Thrift proxy should use. Defaults to
	// :ref:`AUTO_PROTOCOL<envoy_v3_api_enum_value_extensions.filters.network.thrift_proxy.v3.ProtocolType.AUTO_PROTOCOL>`.
	// Note that :ref:`LAX_BINARY<envoy_v3_api_enum_value_extensions.filters.network.thrift_proxy.v3.ProtocolType.LAX_BINARY>`
	// is not distinguished by :ref:`AUTO_PROTOCOL<envoy_v3_api_enum_value_extensions.filters.network.thrift_proxy.v3.ProtocolType.AUTO_PROTOCOL>`,
	// which is the same with :ref:`thrift_proxy network filter <envoy_v3_api_msg_extensions.filters.network.thrift_proxy.v3.ThriftProxy>`.
	// Note that :ref:`TWITTER<envoy_v3_api_enum_value_extensions.filters.network.thrift_proxy.v3.ProtocolType.TWITTER>` is
	// not supported due to deprecation in envoy.
	Protocol v3.ProtocolType `protobuf:"varint,4,opt,name=protocol,proto3,enum=envoy.extensions.filters.network.thrift_proxy.v3.ProtocolType" json:"protocol,omitempty"`
	// Allowed content-type for thrift payload to filter metadata transformation.
	// Default to “{"application/x-thrift"}“.
	//
	// Set “allow_empty_content_type“ if empty/missing content-type header
	// is allowed.
	AllowContentTypes []string `protobuf:"bytes,5,rep,name=allow_content_types,json=allowContentTypes,proto3" json:"allow_content_types,omitempty"`
	// Allowed empty content-type for thrift payload to filter metadata transformation.
	// Default to false.
	AllowEmptyContentType bool `protobuf:"varint,6,opt,name=allow_empty_content_type,json=allowEmptyContentType,proto3" json:"allow_empty_content_type,omitempty"`
}

func (x *ThriftToMetadata) Reset() {
	*x = ThriftToMetadata{}
	if protoimpl.UnsafeEnabled {
		mi := &file_envoy_extensions_filters_http_thrift_to_metadata_v3_thrift_to_metadata_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ThriftToMetadata) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ThriftToMetadata) ProtoMessage() {}

func (x *ThriftToMetadata) ProtoReflect() protoreflect.Message {
	mi := &file_envoy_extensions_filters_http_thrift_to_metadata_v3_thrift_to_metadata_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ThriftToMetadata.ProtoReflect.Descriptor instead.
func (*ThriftToMetadata) Descriptor() ([]byte, []int) {
	return file_envoy_extensions_filters_http_thrift_to_metadata_v3_thrift_to_metadata_proto_rawDescGZIP(), []int{3}
}

func (x *ThriftToMetadata) GetRequestRules() []*Rule {
	if x != nil {
		return x.RequestRules
	}
	return nil
}

func (x *ThriftToMetadata) GetResponseRules() []*Rule {
	if x != nil {
		return x.ResponseRules
	}
	return nil
}

func (x *ThriftToMetadata) GetTransport() v3.TransportType {
	if x != nil {
		return x.Transport
	}
	return v3.TransportType_AUTO_TRANSPORT
}

func (x *ThriftToMetadata) GetProtocol() v3.ProtocolType {
	if x != nil {
		return x.Protocol
	}
	return v3.ProtocolType_AUTO_PROTOCOL
}

func (x *ThriftToMetadata) GetAllowContentTypes() []string {
	if x != nil {
		return x.AllowContentTypes
	}
	return nil
}

func (x *ThriftToMetadata) GetAllowEmptyContentType() bool {
	if x != nil {
		return x.AllowEmptyContentType
	}
	return false
}

// Thrift to metadata configuration on a per-route basis, which overrides the global configuration for
// request rules and responses rules.
type ThriftToMetadataPerRoute struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The list of rules to apply to http request body to extract thrift metadata.
	RequestRules []*Rule `protobuf:"bytes,1,rep,name=request_rules,json=requestRules,proto3" json:"request_rules,omitempty"`
	// The list of rules to apply to http response body to extract thrift metadata.
	ResponseRules []*Rule `protobuf:"bytes,2,rep,name=response_rules,json=responseRules,proto3" json:"response_rules,omitempty"`
}

func (x *ThriftToMetadataPerRoute) Reset() {
	*x = ThriftToMetadataPerRoute{}
	if protoimpl.UnsafeEnabled {
		mi := &file_envoy_extensions_filters_http_thrift_to_metadata_v3_thrift_to_metadata_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ThriftToMetadataPerRoute) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ThriftToMetadataPerRoute) ProtoMessage() {}

func (x *ThriftToMetadataPerRoute) ProtoReflect() protoreflect.Message {
	mi := &file_envoy_extensions_filters_http_thrift_to_metadata_v3_thrift_to_metadata_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ThriftToMetadataPerRoute.ProtoReflect.Descriptor instead.
func (*ThriftToMetadataPerRoute) Descriptor() ([]byte, []int) {
	return file_envoy_extensions_filters_http_thrift_to_metadata_v3_thrift_to_metadata_proto_rawDescGZIP(), []int{4}
}

func (x *ThriftToMetadataPerRoute) GetRequestRules() []*Rule {
	if x != nil {
		return x.RequestRules
	}
	return nil
}

func (x *ThriftToMetadataPerRoute) GetResponseRules() []*Rule {
	if x != nil {
		return x.ResponseRules
	}
	return nil
}

var File_envoy_extensions_filters_http_thrift_to_metadata_v3_thrift_to_metadata_proto protoreflect.FileDescriptor

var file_envoy_extensions_filters_http_thrift_to_metadata_v3_thrift_to_metadata_proto_rawDesc = []byte{
	0x0a, 0x4c, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x2f, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f,
	0x6e, 0x73, 0x2f, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x73, 0x2f, 0x68, 0x74, 0x74, 0x70, 0x2f,
	0x74, 0x68, 0x72, 0x69, 0x66, 0x74, 0x5f, 0x74, 0x6f, 0x5f, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61,
	0x74, 0x61, 0x2f, 0x76, 0x33, 0x2f, 0x74, 0x68, 0x72, 0x69, 0x66, 0x74, 0x5f, 0x74, 0x6f, 0x5f,
	0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x33,
	0x65, 0x6e, 0x76, 0x6f, 0x79, 0x2e, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x73,
	0x2e, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x73, 0x2e, 0x68, 0x74, 0x74, 0x70, 0x2e, 0x74, 0x68,
	0x72, 0x69, 0x66, 0x74, 0x5f, 0x74, 0x6f, 0x5f, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61,
	0x2e, 0x76, 0x33, 0x1a, 0x43, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x2f, 0x65, 0x78, 0x74, 0x65, 0x6e,
	0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x73, 0x2f, 0x6e, 0x65,
	0x74, 0x77, 0x6f, 0x72, 0x6b, 0x2f, 0x74, 0x68, 0x72, 0x69, 0x66, 0x74, 0x5f, 0x70, 0x72, 0x6f,
	0x78, 0x79, 0x2f, 0x76, 0x33, 0x2f, 0x74, 0x68, 0x72, 0x69, 0x66, 0x74, 0x5f, 0x70, 0x72, 0x6f,
	0x78, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x78, 0x64, 0x73, 0x2f, 0x61, 0x6e, 0x6e, 0x6f,
	0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x76, 0x33, 0x2f, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1d, 0x75, 0x64, 0x70, 0x61, 0x2f, 0x61, 0x6e,
	0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65,
	0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x86, 0x01, 0x0a, 0x0c, 0x4b, 0x65, 0x79, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x50, 0x61, 0x69, 0x72,
	0x12, 0x2d, 0x0a, 0x12, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x6e, 0x61, 0x6d,
	0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x11, 0x6d, 0x65,
	0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x4e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x12,
	0x19, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42,
	0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x2c, 0x0a, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x56, 0x61, 0x6c, 0x75,
	0x65, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0xb6, 0x01, 0x0a, 0x0d, 0x46, 0x69, 0x65,
	0x6c, 0x64, 0x53, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x12, 0x1b, 0x0a, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x10,
	0x01, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x24, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x05, 0x42, 0x14, 0xfa, 0x42, 0x11, 0x1a, 0x0f, 0x18, 0xff, 0xff, 0x01, 0x28, 0x80,
	0x80, 0xfe, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x01, 0x52, 0x02, 0x69, 0x64, 0x12, 0x58, 0x0a,
	0x05, 0x63, 0x68, 0x69, 0x6c, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x42, 0x2e, 0x65,
	0x6e, 0x76, 0x6f, 0x79, 0x2e, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2e,
	0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x73, 0x2e, 0x68, 0x74, 0x74, 0x70, 0x2e, 0x74, 0x68, 0x72,
	0x69, 0x66, 0x74, 0x5f, 0x74, 0x6f, 0x5f, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2e,
	0x76, 0x33, 0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x53, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72,
	0x52, 0x05, 0x63, 0x68, 0x69, 0x6c, 0x64, 0x3a, 0x08, 0xd2, 0xc6, 0xa4, 0xe1, 0x06, 0x02, 0x08,
	0x01, 0x22, 0xbc, 0x03, 0x0a, 0x04, 0x52, 0x75, 0x6c, 0x65, 0x12, 0x50, 0x0a, 0x05, 0x66, 0x69,
	0x65, 0x6c, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x3a, 0x2e, 0x65, 0x6e, 0x76, 0x6f,
	0x79, 0x2e, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x66, 0x69, 0x6c,
	0x74, 0x65, 0x72, 0x73, 0x2e, 0x68, 0x74, 0x74, 0x70, 0x2e, 0x74, 0x68, 0x72, 0x69, 0x66, 0x74,
	0x5f, 0x74, 0x6f, 0x5f, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x76, 0x33, 0x2e,
	0x46, 0x69, 0x65, 0x6c, 0x64, 0x52, 0x05, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x12, 0x73, 0x0a, 0x0e,
	0x66, 0x69, 0x65, 0x6c, 0x64, 0x5f, 0x73, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x42, 0x2e, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x2e, 0x65, 0x78, 0x74,
	0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x73, 0x2e,
	0x68, 0x74, 0x74, 0x70, 0x2e, 0x74, 0x68, 0x72, 0x69, 0x66, 0x74, 0x5f, 0x74, 0x6f, 0x5f, 0x6d,
	0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x76, 0x33, 0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64,
	0x53, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x42, 0x08, 0xd2, 0xc6, 0xa4, 0xe1, 0x06, 0x02,
	0x08, 0x01, 0x52, 0x0d, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x53, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f,
	0x72, 0x12, 0x29, 0x0a, 0x0b, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x5f, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x42, 0x08, 0xd2, 0xc6, 0xa4, 0xe1, 0x06, 0x02, 0x08, 0x01,
	0x52, 0x0a, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x60, 0x0a, 0x0a,
	0x6f, 0x6e, 0x5f, 0x70, 0x72, 0x65, 0x73, 0x65, 0x6e, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x41, 0x2e, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x2e, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69,
	0x6f, 0x6e, 0x73, 0x2e, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x73, 0x2e, 0x68, 0x74, 0x74, 0x70,
	0x2e, 0x74, 0x68, 0x72, 0x69, 0x66, 0x74, 0x5f, 0x74, 0x6f, 0x5f, 0x6d, 0x65, 0x74, 0x61, 0x64,
	0x61, 0x74, 0x61, 0x2e, 0x76, 0x33, 0x2e, 0x4b, 0x65, 0x79, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x50,
	0x61, 0x69, 0x72, 0x52, 0x09, 0x6f, 0x6e, 0x50, 0x72, 0x65, 0x73, 0x65, 0x6e, 0x74, 0x12, 0x60,
	0x0a, 0x0a, 0x6f, 0x6e, 0x5f, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6e, 0x67, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x41, 0x2e, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x2e, 0x65, 0x78, 0x74, 0x65, 0x6e,
	0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x73, 0x2e, 0x68, 0x74,
	0x74, 0x70, 0x2e, 0x74, 0x68, 0x72, 0x69, 0x66, 0x74, 0x5f, 0x74, 0x6f, 0x5f, 0x6d, 0x65, 0x74,
	0x61, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x76, 0x33, 0x2e, 0x4b, 0x65, 0x79, 0x56, 0x61, 0x6c, 0x75,
	0x65, 0x50, 0x61, 0x69, 0x72, 0x52, 0x09, 0x6f, 0x6e, 0x4d, 0x69, 0x73, 0x73, 0x69, 0x6e, 0x67,
	0x22, 0x9a, 0x04, 0x0a, 0x10, 0x54, 0x68, 0x72, 0x69, 0x66, 0x74, 0x54, 0x6f, 0x4d, 0x65, 0x74,
	0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x5e, 0x0a, 0x0d, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x5f, 0x72, 0x75, 0x6c, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x39, 0x2e, 0x65,
	0x6e, 0x76, 0x6f, 0x79, 0x2e, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2e,
	0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x73, 0x2e, 0x68, 0x74, 0x74, 0x70, 0x2e, 0x74, 0x68, 0x72,
	0x69, 0x66, 0x74, 0x5f, 0x74, 0x6f, 0x5f, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2e,
	0x76, 0x33, 0x2e, 0x52, 0x75, 0x6c, 0x65, 0x52, 0x0c, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x52, 0x75, 0x6c, 0x65, 0x73, 0x12, 0x60, 0x0a, 0x0e, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x5f, 0x72, 0x75, 0x6c, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x39, 0x2e,
	0x65, 0x6e, 0x76, 0x6f, 0x79, 0x2e, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x73,
	0x2e, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x73, 0x2e, 0x68, 0x74, 0x74, 0x70, 0x2e, 0x74, 0x68,
	0x72, 0x69, 0x66, 0x74, 0x5f, 0x74, 0x6f, 0x5f, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61,
	0x2e, 0x76, 0x33, 0x2e, 0x52, 0x75, 0x6c, 0x65, 0x52, 0x0d, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x52, 0x75, 0x6c, 0x65, 0x73, 0x12, 0x67, 0x0a, 0x09, 0x74, 0x72, 0x61, 0x6e, 0x73,
	0x70, 0x6f, 0x72, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x3f, 0x2e, 0x65, 0x6e, 0x76,
	0x6f, 0x79, 0x2e, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x66, 0x69,
	0x6c, 0x74, 0x65, 0x72, 0x73, 0x2e, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x2e, 0x74, 0x68,
	0x72, 0x69, 0x66, 0x74, 0x5f, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x2e, 0x76, 0x33, 0x2e, 0x54, 0x72,
	0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x54, 0x79, 0x70, 0x65, 0x42, 0x08, 0xfa, 0x42, 0x05,
	0x82, 0x01, 0x02, 0x10, 0x01, 0x52, 0x09, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74,
	0x12, 0x64, 0x0a, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x3e, 0x2e, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x2e, 0x65, 0x78, 0x74, 0x65, 0x6e,
	0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x73, 0x2e, 0x6e, 0x65,
	0x74, 0x77, 0x6f, 0x72, 0x6b, 0x2e, 0x74, 0x68, 0x72, 0x69, 0x66, 0x74, 0x5f, 0x70, 0x72, 0x6f,
	0x78, 0x79, 0x2e, 0x76, 0x33, 0x2e, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x54, 0x79,
	0x70, 0x65, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x82, 0x01, 0x02, 0x10, 0x01, 0x52, 0x08, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x12, 0x3c, 0x0a, 0x13, 0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x5f,
	0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x18, 0x05, 0x20,
	0x03, 0x28, 0x09, 0x42, 0x0c, 0xfa, 0x42, 0x09, 0x92, 0x01, 0x06, 0x22, 0x04, 0x72, 0x02, 0x10,
	0x01, 0x52, 0x11, 0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x54,
	0x79, 0x70, 0x65, 0x73, 0x12, 0x37, 0x0a, 0x18, 0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x5f, 0x65, 0x6d,
	0x70, 0x74, 0x79, 0x5f, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x5f, 0x74, 0x79, 0x70, 0x65,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x08, 0x52, 0x15, 0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x45, 0x6d, 0x70,
	0x74, 0x79, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x22, 0xe6, 0x01,
	0x0a, 0x18, 0x54, 0x68, 0x72, 0x69, 0x66, 0x74, 0x54, 0x6f, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61,
	0x74, 0x61, 0x50, 0x65, 0x72, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x12, 0x5e, 0x0a, 0x0d, 0x72, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x5f, 0x72, 0x75, 0x6c, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x39, 0x2e, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x2e, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73,
	0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x73, 0x2e, 0x68, 0x74, 0x74,
	0x70, 0x2e, 0x74, 0x68, 0x72, 0x69, 0x66, 0x74, 0x5f, 0x74, 0x6f, 0x5f, 0x6d, 0x65, 0x74, 0x61,
	0x64, 0x61, 0x74, 0x61, 0x2e, 0x76, 0x33, 0x2e, 0x52, 0x75, 0x6c, 0x65, 0x52, 0x0c, 0x72, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x52, 0x75, 0x6c, 0x65, 0x73, 0x12, 0x60, 0x0a, 0x0e, 0x72, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x5f, 0x72, 0x75, 0x6c, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x39, 0x2e, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x2e, 0x65, 0x78, 0x74, 0x65, 0x6e,
	0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x73, 0x2e, 0x68, 0x74,
	0x74, 0x70, 0x2e, 0x74, 0x68, 0x72, 0x69, 0x66, 0x74, 0x5f, 0x74, 0x6f, 0x5f, 0x6d, 0x65, 0x74,
	0x61, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x76, 0x33, 0x2e, 0x52, 0x75, 0x6c, 0x65, 0x52, 0x0d, 0x72,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x52, 0x75, 0x6c, 0x65, 0x73, 0x3a, 0x08, 0xd2, 0xc6,
	0xa4, 0xe1, 0x06, 0x02, 0x08, 0x01, 0x2a, 0x7a, 0x0a, 0x05, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x12,
	0x0f, 0x0a, 0x0b, 0x4d, 0x45, 0x54, 0x48, 0x4f, 0x44, 0x5f, 0x4e, 0x41, 0x4d, 0x45, 0x10, 0x00,
	0x12, 0x0c, 0x0a, 0x08, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4f, 0x4c, 0x10, 0x01, 0x12, 0x0d,
	0x0a, 0x09, 0x54, 0x52, 0x41, 0x4e, 0x53, 0x50, 0x4f, 0x52, 0x54, 0x10, 0x02, 0x12, 0x10, 0x0a,
	0x0c, 0x48, 0x45, 0x41, 0x44, 0x45, 0x52, 0x5f, 0x46, 0x4c, 0x41, 0x47, 0x53, 0x10, 0x03, 0x12,
	0x0f, 0x0a, 0x0b, 0x53, 0x45, 0x51, 0x55, 0x45, 0x4e, 0x43, 0x45, 0x5f, 0x49, 0x44, 0x10, 0x04,
	0x12, 0x10, 0x0a, 0x0c, 0x4d, 0x45, 0x53, 0x53, 0x41, 0x47, 0x45, 0x5f, 0x54, 0x59, 0x50, 0x45,
	0x10, 0x05, 0x12, 0x0e, 0x0a, 0x0a, 0x52, 0x45, 0x50, 0x4c, 0x59, 0x5f, 0x54, 0x59, 0x50, 0x45,
	0x10, 0x06, 0x42, 0xd5, 0x01, 0x0a, 0x41, 0x69, 0x6f, 0x2e, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x70,
	0x72, 0x6f, 0x78, 0x79, 0x2e, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x2e, 0x65, 0x78, 0x74, 0x65, 0x6e,
	0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x73, 0x2e, 0x68, 0x74,
	0x74, 0x70, 0x2e, 0x74, 0x68, 0x72, 0x69, 0x66, 0x74, 0x5f, 0x74, 0x6f, 0x5f, 0x6d, 0x65, 0x74,
	0x61, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x76, 0x33, 0x42, 0x15, 0x54, 0x68, 0x72, 0x69, 0x66, 0x74,
	0x54, 0x6f, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50,
	0x01, 0x5a, 0x6f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x65, 0x6e,
	0x76, 0x6f, 0x79, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x2f, 0x67, 0x6f, 0x2d, 0x63, 0x6f, 0x6e, 0x74,
	0x72, 0x6f, 0x6c, 0x2d, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2f, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x2f,
	0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x66, 0x69, 0x6c, 0x74, 0x65,
	0x72, 0x73, 0x2f, 0x68, 0x74, 0x74, 0x70, 0x2f, 0x74, 0x68, 0x72, 0x69, 0x66, 0x74, 0x5f, 0x74,
	0x6f, 0x5f, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2f, 0x76, 0x33, 0x3b, 0x74, 0x68,
	0x72, 0x69, 0x66, 0x74, 0x5f, 0x74, 0x6f, 0x5f, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61,
	0x76, 0x33, 0xba, 0x80, 0xc8, 0xd1, 0x06, 0x02, 0x10, 0x02, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_envoy_extensions_filters_http_thrift_to_metadata_v3_thrift_to_metadata_proto_rawDescOnce sync.Once
	file_envoy_extensions_filters_http_thrift_to_metadata_v3_thrift_to_metadata_proto_rawDescData = file_envoy_extensions_filters_http_thrift_to_metadata_v3_thrift_to_metadata_proto_rawDesc
)

func file_envoy_extensions_filters_http_thrift_to_metadata_v3_thrift_to_metadata_proto_rawDescGZIP() []byte {
	file_envoy_extensions_filters_http_thrift_to_metadata_v3_thrift_to_metadata_proto_rawDescOnce.Do(func() {
		file_envoy_extensions_filters_http_thrift_to_metadata_v3_thrift_to_metadata_proto_rawDescData = protoimpl.X.CompressGZIP(file_envoy_extensions_filters_http_thrift_to_metadata_v3_thrift_to_metadata_proto_rawDescData)
	})
	return file_envoy_extensions_filters_http_thrift_to_metadata_v3_thrift_to_metadata_proto_rawDescData
}

var file_envoy_extensions_filters_http_thrift_to_metadata_v3_thrift_to_metadata_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_envoy_extensions_filters_http_thrift_to_metadata_v3_thrift_to_metadata_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_envoy_extensions_filters_http_thrift_to_metadata_v3_thrift_to_metadata_proto_goTypes = []interface{}{
	(Field)(0),                       // 0: envoy.extensions.filters.http.thrift_to_metadata.v3.Field
	(*KeyValuePair)(nil),             // 1: envoy.extensions.filters.http.thrift_to_metadata.v3.KeyValuePair
	(*FieldSelector)(nil),            // 2: envoy.extensions.filters.http.thrift_to_metadata.v3.FieldSelector
	(*Rule)(nil),                     // 3: envoy.extensions.filters.http.thrift_to_metadata.v3.Rule
	(*ThriftToMetadata)(nil),         // 4: envoy.extensions.filters.http.thrift_to_metadata.v3.ThriftToMetadata
	(*ThriftToMetadataPerRoute)(nil), // 5: envoy.extensions.filters.http.thrift_to_metadata.v3.ThriftToMetadataPerRoute
	(*structpb.Value)(nil),           // 6: google.protobuf.Value
	(v3.TransportType)(0),            // 7: envoy.extensions.filters.network.thrift_proxy.v3.TransportType
	(v3.ProtocolType)(0),             // 8: envoy.extensions.filters.network.thrift_proxy.v3.ProtocolType
}
var file_envoy_extensions_filters_http_thrift_to_metadata_v3_thrift_to_metadata_proto_depIdxs = []int32{
	6,  // 0: envoy.extensions.filters.http.thrift_to_metadata.v3.KeyValuePair.value:type_name -> google.protobuf.Value
	2,  // 1: envoy.extensions.filters.http.thrift_to_metadata.v3.FieldSelector.child:type_name -> envoy.extensions.filters.http.thrift_to_metadata.v3.FieldSelector
	0,  // 2: envoy.extensions.filters.http.thrift_to_metadata.v3.Rule.field:type_name -> envoy.extensions.filters.http.thrift_to_metadata.v3.Field
	2,  // 3: envoy.extensions.filters.http.thrift_to_metadata.v3.Rule.field_selector:type_name -> envoy.extensions.filters.http.thrift_to_metadata.v3.FieldSelector
	1,  // 4: envoy.extensions.filters.http.thrift_to_metadata.v3.Rule.on_present:type_name -> envoy.extensions.filters.http.thrift_to_metadata.v3.KeyValuePair
	1,  // 5: envoy.extensions.filters.http.thrift_to_metadata.v3.Rule.on_missing:type_name -> envoy.extensions.filters.http.thrift_to_metadata.v3.KeyValuePair
	3,  // 6: envoy.extensions.filters.http.thrift_to_metadata.v3.ThriftToMetadata.request_rules:type_name -> envoy.extensions.filters.http.thrift_to_metadata.v3.Rule
	3,  // 7: envoy.extensions.filters.http.thrift_to_metadata.v3.ThriftToMetadata.response_rules:type_name -> envoy.extensions.filters.http.thrift_to_metadata.v3.Rule
	7,  // 8: envoy.extensions.filters.http.thrift_to_metadata.v3.ThriftToMetadata.transport:type_name -> envoy.extensions.filters.network.thrift_proxy.v3.TransportType
	8,  // 9: envoy.extensions.filters.http.thrift_to_metadata.v3.ThriftToMetadata.protocol:type_name -> envoy.extensions.filters.network.thrift_proxy.v3.ProtocolType
	3,  // 10: envoy.extensions.filters.http.thrift_to_metadata.v3.ThriftToMetadataPerRoute.request_rules:type_name -> envoy.extensions.filters.http.thrift_to_metadata.v3.Rule
	3,  // 11: envoy.extensions.filters.http.thrift_to_metadata.v3.ThriftToMetadataPerRoute.response_rules:type_name -> envoy.extensions.filters.http.thrift_to_metadata.v3.Rule
	12, // [12:12] is the sub-list for method output_type
	12, // [12:12] is the sub-list for method input_type
	12, // [12:12] is the sub-list for extension type_name
	12, // [12:12] is the sub-list for extension extendee
	0,  // [0:12] is the sub-list for field type_name
}

func init() { file_envoy_extensions_filters_http_thrift_to_metadata_v3_thrift_to_metadata_proto_init() }
func file_envoy_extensions_filters_http_thrift_to_metadata_v3_thrift_to_metadata_proto_init() {
	if File_envoy_extensions_filters_http_thrift_to_metadata_v3_thrift_to_metadata_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_envoy_extensions_filters_http_thrift_to_metadata_v3_thrift_to_metadata_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*KeyValuePair); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_envoy_extensions_filters_http_thrift_to_metadata_v3_thrift_to_metadata_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FieldSelector); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_envoy_extensions_filters_http_thrift_to_metadata_v3_thrift_to_metadata_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Rule); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_envoy_extensions_filters_http_thrift_to_metadata_v3_thrift_to_metadata_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ThriftToMetadata); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_envoy_extensions_filters_http_thrift_to_metadata_v3_thrift_to_metadata_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ThriftToMetadataPerRoute); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_envoy_extensions_filters_http_thrift_to_metadata_v3_thrift_to_metadata_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_envoy_extensions_filters_http_thrift_to_metadata_v3_thrift_to_metadata_proto_goTypes,
		DependencyIndexes: file_envoy_extensions_filters_http_thrift_to_metadata_v3_thrift_to_metadata_proto_depIdxs,
		EnumInfos:         file_envoy_extensions_filters_http_thrift_to_metadata_v3_thrift_to_metadata_proto_enumTypes,
		MessageInfos:      file_envoy_extensions_filters_http_thrift_to_metadata_v3_thrift_to_metadata_proto_msgTypes,
	}.Build()
	File_envoy_extensions_filters_http_thrift_to_metadata_v3_thrift_to_metadata_proto = out.File
	file_envoy_extensions_filters_http_thrift_to_metadata_v3_thrift_to_metadata_proto_rawDesc = nil
	file_envoy_extensions_filters_http_thrift_to_metadata_v3_thrift_to_metadata_proto_goTypes = nil
	file_envoy_extensions_filters_http_thrift_to_metadata_v3_thrift_to_metadata_proto_depIdxs = nil
}

// Code generated by protoc-gen-go. DO NOT EDIT.
// source: go.chromium.org/luci/buildbucket/proto/build.proto

/*
Package buildbucketpb is a generated protocol buffer package.

It is generated from these files:
	go.chromium.org/luci/buildbucket/proto/build.proto
	go.chromium.org/luci/buildbucket/proto/common.proto
	go.chromium.org/luci/buildbucket/proto/step.proto

It has these top-level messages:
	Build
	BuildInfra
	Builder
	GerritChange
	GitilesCommit
	StringPair
	Step
*/
package buildbucketpb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "github.com/golang/protobuf/ptypes/timestamp"
import google_protobuf1 "github.com/golang/protobuf/ptypes/struct"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// A single build, identified by an int64 id.
// Belongs to a builder.
//
// RPC: see Builds service for build creation and retrieval.
// Some Build fields are marked as excluded from responses by default.
// Use build_fields request field to specify that a field must be included.
//
// BigQuery: this message also defines schema of a BigQuery table of completed builds.
// A BigQuery row is inserted soon after build ends, i.e. a row represents a state of
// a build at completion time and does not change after that.
// All fields are included.
type Build struct {
	// Identifier of the build, unique per LUCI deployment.
	// IDs are monotonically decreasing.
	Id int64 `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
	// Required. The builder this build belongs to.
	//
	// Tuple (builder.project, builder.bucket) defines build ACL
	// which may change after build has ended.
	Builder *Builder_ID `protobuf:"bytes,2,opt,name=builder" json:"builder,omitempty"`
	// Human-oriented identifier of the build with the following properties:
	// - unique within the builder
	// - a monotonically increasing number
	// - mostly contiguous
	// - much shorter than id
	//
	// Caution: populated (positive number) iff build numbers were enabled
	// in the builder configuration at the time of build creation.
	//
	// Caution: Build numbers are not guaranteed to be contiguous.
	// There may be gaps during outages.
	//
	// Caution: Build numbers, while monotonically increasing, do not
	// necessarily reflect source-code order. For example, force builds
	// or rebuilds can allocate new, higher, numbers, but build an older-
	// than-HEAD version of the source.
	Number int32 `protobuf:"varint,3,opt,name=number" json:"number,omitempty"`
	// Verified identity which created this build.
	// Derived by the server from OAuth 2.0 token and verified against Google
	// OAuth servers at the build creation time.
	CreatedBy string `protobuf:"bytes,4,opt,name=created_by,json=createdBy" json:"created_by,omitempty"`
	// URL of a human-oriented build page.
	// Always populated.
	ViewUrl string `protobuf:"bytes,5,opt,name=view_url,json=viewUrl" json:"view_url,omitempty"`
	// When the build was created.
	CreateTime *google_protobuf.Timestamp `protobuf:"bytes,6,opt,name=create_time,json=createTime" json:"create_time,omitempty"`
	// When the build started.
	StartTime *google_protobuf.Timestamp `protobuf:"bytes,7,opt,name=start_time,json=startTime" json:"start_time,omitempty"`
	// When the build ended.
	EndTime *google_protobuf.Timestamp `protobuf:"bytes,8,opt,name=end_time,json=endTime" json:"end_time,omitempty"`
	// When the build was most recently updated.
	//
	// RPC: can be > end_time if, e.g. new tags were attached to a completed
	// build.
	UpdateTime *google_protobuf.Timestamp `protobuf:"bytes,9,opt,name=update_time,json=updateTime" json:"update_time,omitempty"`
	// Status of the build.
	// Must be specified, i.e. not STATUS_UNSPECIFIED.
	//
	// RPC: Responses have most current status.
	//
	// BigQuery: Final status of the build. Cannot be SCHEDULED or STARTED.
	Status Status `protobuf:"varint,12,opt,name=status,enum=buildbucket.v2.Status" json:"status,omitempty"`
	// Input to the build script / recipe.
	Input *Build_Input `protobuf:"bytes,15,opt,name=input" json:"input,omitempty"`
	// Output of the build script / recipe.
	// SHOULD depend only on input field and NOT other fields.
	//
	// RPC: By default, this field is excluded from responses.
	// Updated while the build is running and finalized when the build ends.
	Output *Build_Output `protobuf:"bytes,16,opt,name=output" json:"output,omitempty"`
	// Current list of build steps.
	// Updated as build runs.
	//
	// RPC: By default, this field is excluded from responses.
	Steps []*Step `protobuf:"bytes,17,rep,name=steps" json:"steps,omitempty"`
	// Build infrastructure used by the build.
	//
	// RPC: By default, this field is excluded from responses.
	Infra *BuildInfra `protobuf:"bytes,18,opt,name=infra" json:"infra,omitempty"`
	// Arbitrary annotations for the build.
	// One key may have multiple values, which is why this is not a map<string,string>.
	// Indexed by the server, see also BuildFilter.tags.
	Tags []*StringPair `protobuf:"bytes,19,rep,name=tags" json:"tags,omitempty"`
}

func (m *Build) Reset()                    { *m = Build{} }
func (m *Build) String() string            { return proto.CompactTextString(m) }
func (*Build) ProtoMessage()               {}
func (*Build) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Build) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Build) GetBuilder() *Builder_ID {
	if m != nil {
		return m.Builder
	}
	return nil
}

func (m *Build) GetNumber() int32 {
	if m != nil {
		return m.Number
	}
	return 0
}

func (m *Build) GetCreatedBy() string {
	if m != nil {
		return m.CreatedBy
	}
	return ""
}

func (m *Build) GetViewUrl() string {
	if m != nil {
		return m.ViewUrl
	}
	return ""
}

func (m *Build) GetCreateTime() *google_protobuf.Timestamp {
	if m != nil {
		return m.CreateTime
	}
	return nil
}

func (m *Build) GetStartTime() *google_protobuf.Timestamp {
	if m != nil {
		return m.StartTime
	}
	return nil
}

func (m *Build) GetEndTime() *google_protobuf.Timestamp {
	if m != nil {
		return m.EndTime
	}
	return nil
}

func (m *Build) GetUpdateTime() *google_protobuf.Timestamp {
	if m != nil {
		return m.UpdateTime
	}
	return nil
}

func (m *Build) GetStatus() Status {
	if m != nil {
		return m.Status
	}
	return Status_STATUS_UNSPECIFIED
}

func (m *Build) GetInput() *Build_Input {
	if m != nil {
		return m.Input
	}
	return nil
}

func (m *Build) GetOutput() *Build_Output {
	if m != nil {
		return m.Output
	}
	return nil
}

func (m *Build) GetSteps() []*Step {
	if m != nil {
		return m.Steps
	}
	return nil
}

func (m *Build) GetInfra() *BuildInfra {
	if m != nil {
		return m.Infra
	}
	return nil
}

func (m *Build) GetTags() []*StringPair {
	if m != nil {
		return m.Tags
	}
	return nil
}

// Defines what to build/test.
type Build_Input struct {
	// Arbitrary JSON object. Available at build run time.
	//
	// RPC: By default, this field is excluded from responses.
	//
	// V1 equivalent: corresponds to "properties" key in "parameters_json".
	Properties *google_protobuf1.Struct `protobuf:"bytes,1,opt,name=properties" json:"properties,omitempty"`
	// Gitiles commits to run against.
	// Usually present in CI builds, set by LUCI Scheduler.
	// Usually describe different repositories.
	// If not present, the build may checkout "refs/heads/master".
	// NOT a blamelist.
	//
	// V1 equivalent: supersedes "revision" property and "buildset"
	// tag that starts with "commit/gitiles/".
	GitilesCommits []*GitilesCommit `protobuf:"bytes,2,rep,name=gitiles_commits,json=gitilesCommits" json:"gitiles_commits,omitempty"`
	// Gerrit patchsets to run against.
	// Usually present in tryjobs, set by CQ, Gerrit, git-cl-try.
	// Applied on top of the corresponding commit in gitiles_commits
	// or tip of the tree if the commit is not specified.
	//
	// V1 equivalent: supersedes patch_* properties and "buildset"
	// tag that starts with "patch/gerrit/".
	GerritChanges []*GerritChange `protobuf:"bytes,3,rep,name=gerrit_changes,json=gerritChanges" json:"gerrit_changes,omitempty"`
	// If true, the build does not affect prod. In recipe land, runtime.is_experimental will
	// return true and recipes should not make prod-visible side effects.
	// By default, experimental builds are not surfaced in RPCs, PubSub
	// notifications (unless it is callback), and reported in metrics / BigQuery tables
	// under a different name.
	// See also include_experimental fields to in request messages.
	Experimental bool `protobuf:"varint,5,opt,name=experimental" json:"experimental,omitempty"`
}

func (m *Build_Input) Reset()                    { *m = Build_Input{} }
func (m *Build_Input) String() string            { return proto.CompactTextString(m) }
func (*Build_Input) ProtoMessage()               {}
func (*Build_Input) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 0} }

func (m *Build_Input) GetProperties() *google_protobuf1.Struct {
	if m != nil {
		return m.Properties
	}
	return nil
}

func (m *Build_Input) GetGitilesCommits() []*GitilesCommit {
	if m != nil {
		return m.GitilesCommits
	}
	return nil
}

func (m *Build_Input) GetGerritChanges() []*GerritChange {
	if m != nil {
		return m.GerritChanges
	}
	return nil
}

func (m *Build_Input) GetExperimental() bool {
	if m != nil {
		return m.Experimental
	}
	return false
}

// Output of the build script / recipe.
type Build_Output struct {
	// Arbitrary JSON object produced by the build.
	//
	// V1 equivalent: corresponds to "properties" key in
	// "result_details_json".
	// In V1 output properties are not populated until build ends.
	Properties *google_protobuf1.Struct `protobuf:"bytes,1,opt,name=properties" json:"properties,omitempty"`
	// Human-oriented summary of the build provided by the build itself,
	// in Markdown format (https://spec.commonmark.org/0.28/).
	//
	// BigQuery: excluded from rows.
	SummaryMarkdown string `protobuf:"bytes,2,opt,name=summary_markdown,json=summaryMarkdown" json:"summary_markdown,omitempty"`
}

func (m *Build_Output) Reset()                    { *m = Build_Output{} }
func (m *Build_Output) String() string            { return proto.CompactTextString(m) }
func (*Build_Output) ProtoMessage()               {}
func (*Build_Output) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 1} }

func (m *Build_Output) GetProperties() *google_protobuf1.Struct {
	if m != nil {
		return m.Properties
	}
	return nil
}

func (m *Build_Output) GetSummaryMarkdown() string {
	if m != nil {
		return m.SummaryMarkdown
	}
	return ""
}

// Build infrastructure that was used for a particular build.
type BuildInfra struct {
	Buildbucket *BuildInfra_Buildbucket `protobuf:"bytes,1,opt,name=buildbucket" json:"buildbucket,omitempty"`
	Swarming    *BuildInfra_Swarming    `protobuf:"bytes,2,opt,name=swarming" json:"swarming,omitempty"`
}

func (m *BuildInfra) Reset()                    { *m = BuildInfra{} }
func (m *BuildInfra) String() string            { return proto.CompactTextString(m) }
func (*BuildInfra) ProtoMessage()               {}
func (*BuildInfra) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *BuildInfra) GetBuildbucket() *BuildInfra_Buildbucket {
	if m != nil {
		return m.Buildbucket
	}
	return nil
}

func (m *BuildInfra) GetSwarming() *BuildInfra_Swarming {
	if m != nil {
		return m.Swarming
	}
	return nil
}

// Buildbucket-specific information, captured at the build creation time.
type BuildInfra_Buildbucket struct {
	// Version of swarming task template. Defines
	// versions of kitchen, git, git wrapper, python, vpython, etc.
	ServiceConfigRevision string `protobuf:"bytes,2,opt,name=service_config_revision,json=serviceConfigRevision" json:"service_config_revision,omitempty"`
	// Whether canary version of the swarming task template was used for this
	// build.
	Canary bool `protobuf:"varint,4,opt,name=canary" json:"canary,omitempty"`
}

func (m *BuildInfra_Buildbucket) Reset()                    { *m = BuildInfra_Buildbucket{} }
func (m *BuildInfra_Buildbucket) String() string            { return proto.CompactTextString(m) }
func (*BuildInfra_Buildbucket) ProtoMessage()               {}
func (*BuildInfra_Buildbucket) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1, 0} }

func (m *BuildInfra_Buildbucket) GetServiceConfigRevision() string {
	if m != nil {
		return m.ServiceConfigRevision
	}
	return ""
}

func (m *BuildInfra_Buildbucket) GetCanary() bool {
	if m != nil {
		return m.Canary
	}
	return false
}

// Swarming-specific information.
type BuildInfra_Swarming struct {
	// Swarming hostname, e.g. "chromium-swarm.appspot.com".
	// Populated at the build creation time.
	Hostname string `protobuf:"bytes,1,opt,name=hostname" json:"hostname,omitempty"`
	// Swarming task id.
	// Not guaranteed to be populated at the build creation time.
	TaskId string `protobuf:"bytes,2,opt,name=task_id,json=taskId" json:"task_id,omitempty"`
	// Task service account email address.
	// This is the service account used for all authenticated requests by the
	// build.
	TaskServiceAccount string `protobuf:"bytes,3,opt,name=task_service_account,json=taskServiceAccount" json:"task_service_account,omitempty"`
	// Priority of the task. The lower the more important.
	// Valid values are [1..255].
	Priority int32 `protobuf:"varint,4,opt,name=priority" json:"priority,omitempty"`
	// Swarming dimensions for the task.
	TaskDimensions []*StringPair `protobuf:"bytes,5,rep,name=task_dimensions,json=taskDimensions" json:"task_dimensions,omitempty"`
	// Swarming dimensions of the bot used for the task.
	BotDimensions []*StringPair `protobuf:"bytes,6,rep,name=bot_dimensions,json=botDimensions" json:"bot_dimensions,omitempty"`
}

func (m *BuildInfra_Swarming) Reset()                    { *m = BuildInfra_Swarming{} }
func (m *BuildInfra_Swarming) String() string            { return proto.CompactTextString(m) }
func (*BuildInfra_Swarming) ProtoMessage()               {}
func (*BuildInfra_Swarming) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1, 1} }

func (m *BuildInfra_Swarming) GetHostname() string {
	if m != nil {
		return m.Hostname
	}
	return ""
}

func (m *BuildInfra_Swarming) GetTaskId() string {
	if m != nil {
		return m.TaskId
	}
	return ""
}

func (m *BuildInfra_Swarming) GetTaskServiceAccount() string {
	if m != nil {
		return m.TaskServiceAccount
	}
	return ""
}

func (m *BuildInfra_Swarming) GetPriority() int32 {
	if m != nil {
		return m.Priority
	}
	return 0
}

func (m *BuildInfra_Swarming) GetTaskDimensions() []*StringPair {
	if m != nil {
		return m.TaskDimensions
	}
	return nil
}

func (m *BuildInfra_Swarming) GetBotDimensions() []*StringPair {
	if m != nil {
		return m.BotDimensions
	}
	return nil
}

type Builder struct {
}

func (m *Builder) Reset()                    { *m = Builder{} }
func (m *Builder) String() string            { return proto.CompactTextString(m) }
func (*Builder) ProtoMessage()               {}
func (*Builder) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

// Identifies a builder.
// Canonical string representation: “{project}/{bucket}/{builder}”.
type Builder_ID struct {
	// Project ID, e.g. "chromium". Unique within a LUCI deployment.
	Project string `protobuf:"bytes,1,opt,name=project" json:"project,omitempty"`
	// Bucket name, e.g. "try". Unique within the project.
	// Together with project, defines an ACL.
	Bucket string `protobuf:"bytes,2,opt,name=bucket" json:"bucket,omitempty"`
	// Builder name, e.g. "linux-rel". Unique within the bucket.
	Builder string `protobuf:"bytes,3,opt,name=builder" json:"builder,omitempty"`
}

func (m *Builder_ID) Reset()                    { *m = Builder_ID{} }
func (m *Builder_ID) String() string            { return proto.CompactTextString(m) }
func (*Builder_ID) ProtoMessage()               {}
func (*Builder_ID) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2, 0} }

func (m *Builder_ID) GetProject() string {
	if m != nil {
		return m.Project
	}
	return ""
}

func (m *Builder_ID) GetBucket() string {
	if m != nil {
		return m.Bucket
	}
	return ""
}

func (m *Builder_ID) GetBuilder() string {
	if m != nil {
		return m.Builder
	}
	return ""
}

func init() {
	proto.RegisterType((*Build)(nil), "buildbucket.v2.Build")
	proto.RegisterType((*Build_Input)(nil), "buildbucket.v2.Build.Input")
	proto.RegisterType((*Build_Output)(nil), "buildbucket.v2.Build.Output")
	proto.RegisterType((*BuildInfra)(nil), "buildbucket.v2.BuildInfra")
	proto.RegisterType((*BuildInfra_Buildbucket)(nil), "buildbucket.v2.BuildInfra.Buildbucket")
	proto.RegisterType((*BuildInfra_Swarming)(nil), "buildbucket.v2.BuildInfra.Swarming")
	proto.RegisterType((*Builder)(nil), "buildbucket.v2.Builder")
	proto.RegisterType((*Builder_ID)(nil), "buildbucket.v2.Builder.ID")
}

func init() { proto.RegisterFile("go.chromium.org/luci/buildbucket/proto/build.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 829 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x55, 0x5d, 0x6f, 0xdc, 0x44,
	0x14, 0xd5, 0x6e, 0xb2, 0x5e, 0xef, 0xdd, 0x76, 0x53, 0x86, 0xd2, 0x0c, 0xa6, 0x15, 0xab, 0x20,
	0xa1, 0x85, 0x07, 0xa7, 0xdd, 0x86, 0x22, 0x94, 0x07, 0xd4, 0x6c, 0x44, 0xd9, 0x07, 0x44, 0x35,
	0x81, 0x17, 0x10, 0xb2, 0xc6, 0xf6, 0xc4, 0x19, 0xb2, 0xf6, 0x58, 0x33, 0xe3, 0x84, 0xfc, 0x09,
	0x9e, 0xe0, 0x87, 0xf2, 0x0f, 0xd0, 0x7c, 0x78, 0xe3, 0x26, 0x25, 0x1b, 0xf1, 0xe6, 0x7b, 0xcf,
	0x39, 0xf7, 0xde, 0x39, 0xf3, 0x61, 0x98, 0x17, 0x22, 0xce, 0xce, 0xa4, 0x28, 0x79, 0x53, 0xc6,
	0x42, 0x16, 0xfb, 0xab, 0x26, 0xe3, 0xfb, 0x69, 0xc3, 0x57, 0x79, 0xda, 0x64, 0xe7, 0x4c, 0xef,
	0xd7, 0x52, 0x68, 0xe1, 0x32, 0xb1, 0xfd, 0x46, 0x93, 0x0e, 0x1c, 0x5f, 0xcc, 0xa3, 0x4f, 0x0b,
	0x21, 0x8a, 0x15, 0x73, 0xcc, 0xb4, 0x39, 0xdd, 0xd7, 0xbc, 0x64, 0x4a, 0xd3, 0xb2, 0x76, 0x82,
	0xe8, 0xe9, 0x4d, 0x82, 0xd2, 0xb2, 0xc9, 0xb4, 0x47, 0x5f, 0xde, 0x73, 0x84, 0x4c, 0x94, 0xa5,
	0xa8, 0xbc, 0xe8, 0xc5, 0x3d, 0x45, 0x4a, 0x33, 0x3f, 0xc5, 0xde, 0x5f, 0x21, 0x0c, 0x8e, 0x0c,
	0x01, 0x4d, 0xa0, 0xcf, 0x73, 0xdc, 0x9b, 0xf6, 0x66, 0x5b, 0xa4, 0xcf, 0x73, 0x74, 0x00, 0x43,
	0xab, 0x64, 0x12, 0xf7, 0xa7, 0xbd, 0xd9, 0x78, 0x1e, 0xc5, 0xef, 0x2e, 0x31, 0x3e, 0x72, 0x70,
	0xbc, 0x3c, 0x26, 0x2d, 0x15, 0x3d, 0x81, 0xa0, 0x6a, 0xca, 0x94, 0x49, 0xbc, 0x35, 0xed, 0xcd,
	0x06, 0xc4, 0x47, 0xe8, 0x19, 0x40, 0x26, 0x19, 0xd5, 0x2c, 0x4f, 0xd2, 0x2b, 0xbc, 0x3d, 0xed,
	0xcd, 0x46, 0x64, 0xe4, 0x33, 0x47, 0x57, 0xe8, 0x63, 0x08, 0x2f, 0x38, 0xbb, 0x4c, 0x1a, 0xb9,
	0xc2, 0x03, 0x0b, 0x0e, 0x4d, 0xfc, 0xb3, 0x5c, 0xa1, 0x43, 0x18, 0x3b, 0x5e, 0x62, 0x1c, 0xc4,
	0x81, 0x9f, 0xc5, 0xb9, 0x17, 0xb7, 0xee, 0xc5, 0x3f, 0xb5, 0xf6, 0x12, 0xdf, 0xc8, 0x24, 0xd0,
	0x37, 0x00, 0x4a, 0x53, 0xa9, 0x9d, 0x76, 0xb8, 0x51, 0x3b, 0xb2, 0x6c, 0x2b, 0xfd, 0x0a, 0x42,
	0x56, 0xe5, 0x4e, 0x18, 0x6e, 0x14, 0x0e, 0x59, 0x95, 0x5b, 0xd9, 0x21, 0x8c, 0x9b, 0x3a, 0x5f,
	0x8f, 0x3b, 0xda, 0x3c, 0xae, 0xa3, 0x5b, 0x71, 0x0c, 0x81, 0xd2, 0x54, 0x37, 0x0a, 0x3f, 0x98,
	0xf6, 0x66, 0x93, 0xf9, 0x93, 0x9b, 0x96, 0x9f, 0x58, 0x94, 0x78, 0x16, 0x7a, 0x01, 0x03, 0x5e,
	0xd5, 0x8d, 0xc6, 0x3b, 0xb6, 0xcd, 0x27, 0xef, 0xdd, 0xa1, 0x78, 0x69, 0x28, 0xc4, 0x31, 0xd1,
	0x01, 0x04, 0xa2, 0xd1, 0x46, 0xf3, 0xc8, 0x6a, 0x9e, 0xbe, 0x5f, 0xf3, 0xa3, 0xe5, 0x10, 0xcf,
	0x45, 0x5f, 0xc2, 0xc0, 0x1c, 0x1a, 0x85, 0x3f, 0x98, 0x6e, 0xcd, 0xc6, 0xf3, 0xc7, 0xb7, 0xe7,
	0x62, 0x35, 0x71, 0x14, 0xf4, 0xdc, 0x0c, 0x75, 0x2a, 0x29, 0x46, 0x77, 0x1c, 0x9b, 0xa5, 0x61,
	0x10, 0x47, 0x44, 0x31, 0x6c, 0x6b, 0x5a, 0x28, 0xfc, 0xa1, 0x2d, 0x1e, 0xdd, 0x2e, 0x2e, 0x79,
	0x55, 0xbc, 0xa5, 0x5c, 0x12, 0xcb, 0x8b, 0xfe, 0xe9, 0xc1, 0xc0, 0x2e, 0x0a, 0x7d, 0x0d, 0x50,
	0x4b, 0x51, 0x33, 0xa9, 0x39, 0x53, 0xf6, 0xf0, 0x8e, 0xe7, 0xbb, 0xb7, 0xcc, 0x3e, 0xb1, 0x37,
	0x8b, 0x74, 0xa8, 0xe8, 0x3b, 0xd8, 0x29, 0xb8, 0xe6, 0x2b, 0xa6, 0x12, 0x73, 0x85, 0xb8, 0x56,
	0xb8, 0x6f, 0xbb, 0x3f, 0xbb, 0xd9, 0xfd, 0x8d, 0xa3, 0x2d, 0x2c, 0x8b, 0x4c, 0x8a, 0x6e, 0xa8,
	0xd0, 0x02, 0x26, 0x05, 0x93, 0x92, 0xeb, 0x24, 0x3b, 0xa3, 0x55, 0xc1, 0x14, 0xde, 0xb2, 0x65,
	0x6e, 0xd9, 0xfa, 0xc6, 0xb2, 0x16, 0x96, 0x44, 0x1e, 0x16, 0x9d, 0x48, 0xa1, 0x3d, 0x78, 0xc0,
	0xfe, 0xa8, 0x99, 0xe4, 0x25, 0xab, 0x34, 0x75, 0x37, 0x20, 0x24, 0xef, 0xe4, 0xa2, 0x15, 0x04,
	0x6e, 0x4f, 0xfe, 0xff, 0x9a, 0xbf, 0x80, 0x47, 0xaa, 0x29, 0x4b, 0x2a, 0xaf, 0x92, 0x92, 0xca,
	0xf3, 0x5c, 0x5c, 0x56, 0xf6, 0x6a, 0x8f, 0xc8, 0x8e, 0xcf, 0xff, 0xe0, 0xd3, 0x7b, 0x7f, 0x6e,
	0x03, 0x5c, 0xef, 0x13, 0xfa, 0x1e, 0xc6, 0x9d, 0xe5, 0xf8, 0x9e, 0x9f, 0xff, 0xf7, 0xc6, 0xba,
	0x4f, 0x87, 0x90, 0xae, 0x14, 0x7d, 0x0b, 0xa1, 0xba, 0xa4, 0xb2, 0xe4, 0x55, 0xe1, 0x9f, 0x95,
	0xcf, 0xee, 0x28, 0x73, 0xe2, 0xa9, 0x64, 0x2d, 0x8a, 0x7e, 0x83, 0x71, 0xa7, 0x38, 0x7a, 0x05,
	0xbb, 0x8a, 0xc9, 0x0b, 0x9e, 0xb1, 0x24, 0x13, 0xd5, 0x29, 0x2f, 0x12, 0xc9, 0x2e, 0xb8, 0xe2,
	0xa2, 0x5d, 0xda, 0x47, 0x1e, 0x5e, 0x58, 0x94, 0x78, 0xd0, 0xbc, 0x53, 0x19, 0xad, 0xa8, 0x74,
	0x6f, 0x51, 0x48, 0x7c, 0x14, 0xfd, 0xdd, 0x87, 0xb0, 0xed, 0x8a, 0x22, 0x08, 0xcf, 0x84, 0xd2,
	0x15, 0x2d, 0x99, 0x5d, 0xf3, 0x88, 0xac, 0x63, 0xb4, 0x0b, 0x43, 0x4d, 0xd5, 0x79, 0xc2, 0x73,
	0xdf, 0x28, 0x30, 0xe1, 0x32, 0x47, 0xcf, 0xe1, 0xb1, 0x05, 0xda, 0xb1, 0x68, 0x96, 0x89, 0xa6,
	0xd2, 0xf6, 0x3d, 0x1c, 0x11, 0x64, 0xb0, 0x13, 0x07, 0xbd, 0x76, 0x88, 0x69, 0x53, 0x4b, 0x2e,
	0x24, 0xd7, 0x6e, 0x9a, 0x01, 0x59, 0xc7, 0x68, 0x01, 0x3b, 0xb6, 0x5a, 0x6e, 0xce, 0x81, 0x99,
	0x5c, 0xe1, 0xc1, 0xc6, 0x5b, 0x32, 0x31, 0x92, 0xe3, 0xb5, 0x02, 0xbd, 0x86, 0x49, 0x2a, 0x74,
	0xb7, 0x46, 0xb0, 0xb1, 0xc6, 0xc3, 0x54, 0xe8, 0xeb, 0x12, 0x7b, 0xbf, 0xc2, 0xd0, 0x3f, 0xf7,
	0xd1, 0x5b, 0xe8, 0x2f, 0x8f, 0x11, 0x86, 0x61, 0x2d, 0xc5, 0xef, 0x2c, 0xd3, 0xde, 0x9a, 0x36,
	0x34, 0xd6, 0xfa, 0x73, 0xe2, 0x8d, 0xf1, 0x5b, 0x85, 0xaf, 0x7f, 0x28, 0xce, 0x8b, 0x36, 0x3c,
	0x7a, 0xf5, 0xcb, 0xc1, 0xfd, 0xfe, 0x5c, 0x87, 0x9d, 0x4c, 0x9d, 0xa6, 0x81, 0x4d, 0xbe, 0xfc,
	0x37, 0x00, 0x00, 0xff, 0xff, 0xbf, 0x6b, 0x8b, 0x73, 0xb0, 0x07, 0x00, 0x00,
}

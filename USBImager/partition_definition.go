package usbimager

import (
	"fmt"
	"strings"
)

type PartitionDefinitionBuilder struct {
	volumeName       string
	label            string
	stringAttributes map[string]string
	bootable         bool
	partType         PartitionType
}

//end result
// label : start=2048, size=512M, type=0C, bootable
// or
// label : type=83

func NewPartitionBuilder(label string) *PartitionDefinitionBuilder {
	return &PartitionDefinitionBuilder{
		label:            label,
		stringAttributes: make(map[string]string),
	}
}
func (pb *PartitionDefinitionBuilder) WithName(volumeName string) *PartitionDefinitionBuilder {
	//used by a Diskpartitionare to set volume names during mkfs runs
	pb.volumeName = volumeName
	pb.stringAttributes["name"] = fmt.Sprintf("\"%s\"", volumeName)
	return pb
}
func (pb *PartitionDefinitionBuilder) OfType(typ PartitionType) *PartitionDefinitionBuilder {
	pb.partType = typ
	return pb
}
func (pb *PartitionDefinitionBuilder) StartAt(start string) *PartitionDefinitionBuilder {
	pb.stringAttributes["start"] = start
	return pb
}
func (pb *PartitionDefinitionBuilder) WithSize(size string) *PartitionDefinitionBuilder {
	pb.stringAttributes["size"] = size
	return pb
}
func (pb *PartitionDefinitionBuilder) WithUndefinedOption(key, value string) *PartitionDefinitionBuilder {
	pb.stringAttributes[key] = value
	return pb
}
func (pb *PartitionDefinitionBuilder) SetBootable(bootable bool) *PartitionDefinitionBuilder {
	pb.bootable = bootable
	return pb
}
func (pb *PartitionDefinitionBuilder) ToSfdisk() string {
	var definition []string
	for key, value := range pb.stringAttributes {
		definition = append(definition, fmt.Sprintf("%s=%s", key, value))
	}

	if pb.partType != "" {
		definition = append(definition, fmt.Sprintf("type=%s", pb.partType))
	}

	if pb.bootable {
		definition = append(definition, "bootable")
	}

	definitions := strings.Join(definition, ", ")
	sfdisk_partition_label := fmt.Sprintf("%s : ", pb.label) + definitions
	return sfdisk_partition_label
}

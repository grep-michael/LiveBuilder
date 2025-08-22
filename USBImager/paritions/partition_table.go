package paritions

import "fmt"

type TableType string

const (
	TABLETYPE_MBR TableType = "dos"
	TABLETYPE_GPT TableType = "gpt"
)

type PartitionTabelBuilder struct {
	partitions []*PartitionDefinitionBuilder
	label      TableType
	label_id   string
	units      string
}

func NewPartitionTable(typ TableType) *PartitionTabelBuilder {
	return &PartitionTabelBuilder{
		label:    typ,
		units:    "sectors",
		label_id: "0x12345678",
	}
}
func (table *PartitionTabelBuilder) WithUnitSize(unit string) *PartitionTabelBuilder {
	fmt.Println("unit is deprecated, the only unit should be sectors")
	table.units = unit
	return table
}
func (table *PartitionTabelBuilder) WithLabelID(id string) *PartitionTabelBuilder {
	table.label_id = id
	return table
}
func (table *PartitionTabelBuilder) WithPartitionDefinition(definition *PartitionDefinitionBuilder) *PartitionTabelBuilder {
	table.partitions = append(table.partitions, definition)
	return table
}
func (table *PartitionTabelBuilder) ToSfdisk() string {

	partitionTable := fmt.Sprintf(
		"label: %s\nlabel-id: %s\nunit: %s\n",
		table.label, table.label_id, table.units,
	)

	for _, definition := range table.partitions {
		partitionTable += "\t" + definition.ToSfdisk()
		partitionTable += "\n"
	}

	return partitionTable
}

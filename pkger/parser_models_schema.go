package pkger

const (
	fieldBucketSchemaType   = "schemaType"
	fieldMeasurementSchemas = "measurementSchemas"

	// measurementSchema fields
	fieldMeasurementSchemaName    = "name"
	fieldMeasurementSchemaColumns = "columns"

	// measurementColumn fields
	fieldMeasurementColumnName     = "name"
	fieldMeasurementColumnType     = "type"
	fieldMeasurementColumnDataType = "dataType"
)

type measurementSchemas []measurementSchema

func (s measurementSchemas) valid() []validationErr {
	return nil
}

func (s measurementSchemas) summarize() []SummaryMeasurementSchema {
	schemas := make([]SummaryMeasurementSchema, 0, len(s))
	for _, schema := range s {
		schemas = append(schemas, schema.summarize())
	}
	return schemas
}

type measurementSchema struct {
	Name    string              `json:"name" yaml:"name"`
	Columns []measurementColumn `json:"columns" yaml:"columns"`
}

func (s measurementSchema) summarize() SummaryMeasurementSchema {
	cols := make([]SummaryMeasurementSchemaColumn, 0, len(s.Columns))
	for i := range s.Columns {
		cols = append(cols, s.Columns[i].summarize())
	}
	return SummaryMeasurementSchema{Name: s.Name, Columns: cols}
}

type measurementColumn struct {
	Name     string `json:"name" yaml:"name"`
	Type     string `json:"type" yaml:"type"`
	DataType string `json:"dataType,omitempty" yaml:"dataType,omitempty"`
}

func (c measurementColumn) summarize() SummaryMeasurementSchemaColumn {
	return SummaryMeasurementSchemaColumn{
		Name:     c.Name,
		Type:     c.Type,
		DataType: c.DataType,
	}
}

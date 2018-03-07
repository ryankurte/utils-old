package protogen

type Options struct {
	input string `short:"i" long:"input" description:"Input protocol file"`
}

// Protocol is the top level protocol definition object
type Protocol struct {
	Name     string
	Messages []Message
}

// Message is a protocol message
type Message struct {
}

type FieldType string

const (
	u8  FieldType = "u8"
	u16 FieldType = "u16"
	u32 FieldType = "u32"

	i8  FieldType = "i8"
	i16 FieldType = "i16"
	i32 FieldType = "i32"

	f32 FieldType = "u8"
	f64 FieldType = "u16"
)

type Field struct {
	Name string
	Type string
}

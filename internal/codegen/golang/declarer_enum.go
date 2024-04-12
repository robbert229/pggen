package golang

import (
	"strconv"
	"strings"
	"text/template"

	"github.com/robbert229/pggen/internal/codegen/golang/gotype"
)

func NameEnumCodecFunc(typ *gotype.EnumType) string {
	return "new" + typ.Name + "Enum"
}

// EnumTypeDeclarer declares a new string type and the const values to map to a
// Postgres enum.
type EnumTypeDeclarer struct {
	enum *gotype.EnumType
}

func NewEnumTypeDeclarer(enum *gotype.EnumType) EnumTypeDeclarer {
	return EnumTypeDeclarer{enum: enum}
}

func (e EnumTypeDeclarer) DedupeKey() string {
	return "enum_type::" + e.enum.Name
}

func (e EnumTypeDeclarer) Declare(string) (string, error) {
	sb := &strings.Builder{}
	// Doc string.
	if e.enum.PgEnum.Name != "" {
		sb.WriteString("// ")
		sb.WriteString(e.enum.Name)
		sb.WriteString(" represents the Postgres enum ")
		sb.WriteString(strconv.Quote(e.enum.PgEnum.Name))
		sb.WriteString(".\n")
	}
	// Type declaration.
	sb.WriteString("type ")
	sb.WriteString(e.enum.Name)
	sb.WriteString(" string\n\n")
	// Const enum values.
	sb.WriteString("const (\n")
	nameLen := 0
	for _, label := range e.enum.Labels {
		if len(label) > nameLen {
			nameLen = len(label)
		}
	}
	for i, label := range e.enum.Labels {
		sb.WriteString("\t")
		sb.WriteString(label)
		sb.WriteString(strings.Repeat(" ", nameLen+1-len(label)))
		sb.WriteString(e.enum.Name)
		sb.WriteString(` = `)
		sb.WriteString(strconv.Quote(e.enum.Values[i]))
		sb.WriteByte('\n')
	}
	sb.WriteString(")\n\n")
	// Stringer
	dispatcher := strings.ToLower(e.enum.Name)[0]
	sb.WriteString("func (")
	sb.WriteByte(dispatcher)
	sb.WriteByte(' ')
	sb.WriteString(e.enum.Name)
	sb.WriteString(") String() string { return string(")
	sb.WriteByte(dispatcher)
	sb.WriteString(") }")
	return sb.String(), nil
}

// EnumTranscoderDeclarer declares a new Go function that creates a pgx decoder
// for the Postgres type represented by the gotype.EnumType.
type EnumTranscoderDeclarer struct {
	typ *gotype.EnumType
}

func NewEnumTranscoderDeclarer(enum *gotype.EnumType) EnumTranscoderDeclarer {
	return EnumTranscoderDeclarer{typ: enum}
}

func (e EnumTranscoderDeclarer) DedupeKey() string {
	return "enum_decoder::" + e.typ.Name
}

func (e EnumTranscoderDeclarer) Declare(string) (string, error) {
	t := template.New("enum")
	t = template.Must(t.Parse(`
// register_{{ .FuncName }} registers the given postgres type with pgx.
func register_{{ .FuncName }}(
	ctx context.Context,
	conn genericConn,
) error {
	t, err := conn.LoadType(
		ctx,
		{{ .PgEnumName }},
	)
	if err != nil {
		return fmt.Errorf("failed to load type for: %w", err)
	}
	
	conn.TypeMap().RegisterType(t)
	
	t, err = conn.LoadType(
		ctx,
		{{ .PgEnumArrayName }},
	)
	if err != nil {
		return fmt.Errorf("failed to load type for: %w", err)
	}
	
	conn.TypeMap().RegisterType(t)
	
	return nil
}

func codec_{{ .FuncName }}(conn genericConn) (pgtype.Codec, error) {
	return &pgtype.EnumCodec{}, nil
}

func init(){
	addHook(register_{{ .FuncName }}) 
}
`))
	sb := &strings.Builder{}
	funcName := NameEnumCodecFunc(e.typ)

	if err := t.Execute(sb, struct {
		FuncName        string
		PgEnumArrayName string
		PgEnumName      string
	}{
		FuncName:        funcName,
		PgEnumName:      strconv.Quote(strconv.Quote(e.typ.PgEnum.Name)),
		PgEnumArrayName: strconv.Quote("_" + e.typ.PgEnum.Name),
	}); err != nil {
		return "", err
	}

	return sb.String(), nil
}

package golang

import (
	"strconv"
	"strings"
	"text/template"

	"github.com/robbert229/pggen/internal/codegen/golang/gotype"
)

// NameCompositeCodecFunc returns the function name that creates a
// pgtype.ValueTranscoder for the composite type that's used to decode rows
// returned by Postgres.
func NameCompositeCodecFunc(typ *gotype.CompositeType) string {
	return "new" + typ.Name
}

// NameCompositeInitFunc returns the name of the function that creates an
// initialized pgtype.ValueTranscoder for the composite type used as a query
// parameters. This function is only necessary for top-level types. Descendant
// types use the raw functions, named by NameCompositeRawFunc.
func NameCompositeInitFunc(typ *gotype.CompositeType) string {
	return "new" + typ.Name + "Init"
}

// NameCompositeRawFunc returns the function name that creates the
// []interface{} array for the composite type so that we can use it with a
// parent encoder function, like NameCompositeInitFunc, in the pgtype.Value
// Set call.
func NameCompositeRawFunc(typ *gotype.CompositeType) string {
	return "new" + typ.Name + "Raw"
}

// CompositeTypeDeclarer declares a new Go struct to represent a Postgres
// composite type.
type CompositeTypeDeclarer struct {
	comp *gotype.CompositeType
}

func NewCompositeTypeDeclarer(comp *gotype.CompositeType) CompositeTypeDeclarer {
	return CompositeTypeDeclarer{comp: comp}
}

func (c CompositeTypeDeclarer) DedupeKey() string {
	return "composite::" + c.comp.Name
}

func (c CompositeTypeDeclarer) Declare(pkgPath string) (string, error) {
	sb := &strings.Builder{}
	// Doc string
	if c.comp.PgComposite.Name != "" {
		sb.WriteString("// ")
		sb.WriteString(c.comp.Name)
		sb.WriteString(" represents the Postgres composite type ")
		sb.WriteString(strconv.Quote(c.comp.PgComposite.Name))
		sb.WriteString(".\n")
	}
	// Struct declaration.
	sb.WriteString("type ")
	sb.WriteString(c.comp.Name)
	sb.WriteString(" struct")
	if len(c.comp.FieldNames) == 0 {
		sb.WriteString("{") // type Foo struct{}
	} else {
		sb.WriteString(" {\n") // type Foo struct {\n
	}
	// Struct fields.
	nameLen, typeLen := getLongestNameTypes(c.comp, pkgPath)
	for i, name := range c.comp.FieldNames {
		// Name
		sb.WriteRune('\t')
		sb.WriteString(name)
		// Type
		qualType := gotype.QualifyType(c.comp.FieldTypes[i], pkgPath)
		sb.WriteString(strings.Repeat(" ", nameLen-len(name)))
		sb.WriteString(qualType)
		// JSON struct tag
		sb.WriteString(strings.Repeat(" ", typeLen-len(qualType)))
		sb.WriteString("`json:")
		sb.WriteString(strconv.Quote(c.comp.PgComposite.ColumnNames[i]))
		sb.WriteString("`")
		sb.WriteRune('\n')
	}
	sb.WriteString("}")
	return sb.String(), nil
}

// getLongestNameTypes returns the length of the longest name and type name for
// all child fields of a composite type. Useful for aligning struct definitions.
func getLongestNameTypes(typ *gotype.CompositeType, pkgPath string) (int, int) {
	nameLen := 0
	for _, name := range typ.FieldNames {
		if n := len(name); n > nameLen {
			nameLen = n
		}
	}
	nameLen++ // 1 space to separate name from type

	typeLen := 0
	for _, childType := range typ.FieldTypes {
		if n := len(gotype.QualifyType(childType, pkgPath)); n > typeLen {
			typeLen = n
		}
	}
	typeLen++ // 1 space to separate type from struct tags.

	return nameLen, typeLen
}

// CompositeTranscoderDeclarer declares a new Go function that creates a pgx
// decoder for the Postgres type represented by the gotype.CompositeType.
type CompositeTranscoderDeclarer struct {
	typ *gotype.CompositeType
}

func NewCompositeTranscoderDeclarer(typ *gotype.CompositeType) CompositeTranscoderDeclarer {
	return CompositeTranscoderDeclarer{typ}
}

func (c CompositeTranscoderDeclarer) DedupeKey() string {
	return "type_resolver::" + c.typ.Name + "_01_transcoder"
}

func (c CompositeTranscoderDeclarer) Declare(pkgPath string) (string, error) {
	funcName := NameCompositeCodecFunc(c.typ)

	t := template.New("declarer")
	t = template.Must(t.Parse(`
	// codec_{{ .FuncName }} is a codec for the composite type of the same name
	func codec_{{ .FuncName }}(conn RegisterConn) (pgtype.Codec, error) {
		{{ range $i, $val := .Columns }}
		    field{{ $i }}, ok := conn.TypeMap().TypeForName("{{ $val.PgFieldTypeName }}")
			if !ok {
				return nil, fmt.Errorf("type not found: {{ $val.PgFieldTypeName }}")
			}
		{{ end }}
		
		return &pgtype.CompositeCodec{
			Fields: []pgtype.CompositeCodecField{
				{{ range $i, $val := .Columns }}
					{
						Name: {{ $val.Name }},
						Type: field{{ $i }},
					},
				{{ end }}
			},
		}, nil
	}

	func register_{{ .FuncName }}(
		ctx context.Context,
		conn RegisterConn,
	) error {
		t, err := conn.LoadType(
			ctx,
			{{ .PgCompositeName }},
		)
		if err != nil {
			return fmt.Errorf("{{ .FuncName }} failed to load type: %w", err)
		}
		
		conn.TypeMap().RegisterType(t)

		return nil
	}

	func init(){
		addHook(register_{{ .FuncName }}) 
	}
	`))

	type Column struct {
		Name            string
		PgFieldTypeName string
	}

	var columns []Column

	for i := 0; i < len(c.typ.PgComposite.ColumnNames); i++ {
		columns = append(columns, Column{
			Name:            strconv.Quote(c.typ.PgComposite.ColumnNames[i]),
			PgFieldTypeName: c.typ.PgComposite.ColumnTypes[i].String(),
		})
	}

	b := &strings.Builder{}
	b.Grow(256)
	if err := t.Execute(b, struct {
		FuncName        string
		Columns         []Column
		PgCompositeName string
	}{
		FuncName:        funcName,
		Columns:         columns,
		PgCompositeName: strconv.Quote(strconv.Quote(c.typ.PgComposite.Name)),
	}); err != nil {
		return "", err
	}
	return b.String(), nil
}

package golang

import (
	"strconv"
	"strings"
	"text/template"

	"github.com/robbert229/pggen/internal/codegen/golang/gotype"
)

// ArrayRegisterFunc returns the function name that registers the array in the
// typemap.
func ArrayRegisterFunc(typ *gotype.ArrayType) string {
	return "new" + typ.Elem.BaseName() + "Array"
}

// ArrayTranscoderDeclarer declares a new Go function that creates a
// pgtype.ValueTranscoder decoder for an array Postgres type.
type ArrayTranscoderDeclarer struct {
	typ *gotype.ArrayType
}

func NewArrayDecoderDeclarer(typ *gotype.ArrayType) ArrayTranscoderDeclarer {
	return ArrayTranscoderDeclarer{typ: typ}
}

func (a ArrayTranscoderDeclarer) DedupeKey() string {
	return "type_resolver::" + a.typ.BaseName() + "_01_transcoder"
}

func (a ArrayTranscoderDeclarer) Declare(string) (string, error) {
	sb := &strings.Builder{}
	funcName := ArrayRegisterFunc(a.typ)

	t := template.New("declarer")
	t = template.Must(t.Parse(`
	// codec_{{ .FuncName }} is a codec for the composite type of the same name
	func codec_{{ .FuncName }}(conn RegisterConn) (pgtype.Codec, error) {
		elementType, ok := conn.TypeMap().TypeForName("{{ .PgElemName }}")
		if !ok {
			return nil, fmt.Errorf("type not found: {{ .PgElemName }}")
		}

		return &pgtype.ArrayCodec{
			ElementType: elementType,
		}, nil
	}

	func register_{{ .FuncName }}(
		ctx context.Context,
		conn RegisterConn,
	) error {
		t, err := conn.LoadType(
			ctx,
			{{ .PgArrayName }},
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

	b := &strings.Builder{}
	b.Grow(256)
	err := t.Execute(b, struct {
		FuncName    string
		Fields      string
		PgArrayName string
		PgElemName  string
	}{
		FuncName:    funcName,
		Fields:      sb.String(),
		PgArrayName: strconv.Quote(strconv.Quote(a.typ.PgArray.Name)),
		PgElemName:  a.typ.PgArray.Elem.String(),
	})
	if err != nil {
		return "", err
	}

	return b.String(), nil
}

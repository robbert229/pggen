package golang

import (
	"sort"

	"github.com/robbert229/pggen/internal/codegen/golang/gotype"
)

// Declarer is implemented by any value that needs to declare types, data, or
// functions before use. For example, Postgres enums map to a Go enum with a
// type declaration and const values. If we use the enum in any Querier
// function, we need to declare the enum.
type Declarer interface {
	// DedupeKey uniquely identifies the declaration so that we only emit
	// declarations once. Should be namespaced like enum::some_enum.
	DedupeKey() string
	// Declare returns the string of the Go code for the declaration.
	Declare(pkgPath string) (string, error)
}

// DeclarerSet is a set of declarers, identified by the dedupe key.
type DeclarerSet map[string]Declarer

func NewDeclarerSet(decls ...Declarer) DeclarerSet {
	d := DeclarerSet(make(map[string]Declarer, len(decls)))
	d.AddAll(decls...)
	return d
}

func (d DeclarerSet) AddAll(decls ...Declarer) {
	for _, decl := range decls {
		d[decl.DedupeKey()] = decl
	}
}

// ListAll gets all declarers in the set in a stable sort order.
func (d DeclarerSet) ListAll() []Declarer {
	decls := make([]Declarer, 0, len(d))
	for _, decl := range d {
		decls = append(decls, decl)
	}
	sort.Slice(decls, func(i, j int) bool { return decls[i].DedupeKey() < decls[j].DedupeKey() })
	return decls
}

// FindInputDeclarers finds all necessary Declarers for types that appear in
// the input parameters. Returns nil if no declarers are needed.
func FindInputDeclarers(typ gotype.Type) DeclarerSet {
	decls := NewDeclarerSet()
	// Only top level types need the init declarer. Descendant types need the
	// raw declarer.
	switch typ := gotype.UnwrapNestedType(typ).(type) {
	case *gotype.CompositeType:
		decls.AddAll(
			NewTypeResolverDeclarer(),
			NewCompositeInitDeclarer(typ),
		)
	case *gotype.ArrayType:
		if gotype.IsPgxSupportedArray(typ) {
			break
		}
		switch gotype.UnwrapNestedType(typ.Elem).(type) {
		case *gotype.CompositeType, *gotype.EnumType:
			decls.AddAll(
				NewTypeResolverDeclarer(),
				NewArrayInitDeclarer(typ),
			)
		}
	}
	decls.AddAll(NewTypeResolverInitDeclarer()) // always add
	findInputDeclsHelper(typ, decls)
	// Inputs depend on output transcoders.
	findOutputDeclsHelper(typ, decls /*hadCompositeParent*/, false)
	return decls
}

func findInputDeclsHelper(typ gotype.Type, decls DeclarerSet) {
	switch typ := gotype.UnwrapNestedType(typ).(type) {
	case *gotype.CompositeType:
		decls.AddAll(
			NewCompositeRawDeclarer(typ),
		)
		for _, childType := range typ.FieldTypes {
			findInputDeclsHelper(childType, decls)
		}

	case *gotype.ArrayType:
		if gotype.IsPgxSupportedArray(typ) {
			return
		}
		decls.AddAll(
			NewArrayRawDeclarer(typ),
		)
		findInputDeclsHelper(typ.Elem, decls)

	default:
		return
	}
}

// FindOutputDeclarers finds all necessary Declarers for types that appear in
// the output rows. Returns nil if no declarers are needed.
func FindOutputDeclarers(typ gotype.Type) DeclarerSet {
	decls := NewDeclarerSet()
	decls.AddAll(NewTypeResolverInitDeclarer()) // always add
	findOutputDeclsHelper(typ, decls, false)
	return decls
}

func findOutputDeclsHelper(typ gotype.Type, decls DeclarerSet, hadCompositeParent bool) {
	switch typ := gotype.UnwrapNestedType(typ).(type) {
	case *gotype.EnumType:
		decls.AddAll(
			NewEnumTypeDeclarer(typ),
		)
		if hadCompositeParent {
			// We can use a string as the decoder except if the enum is part of a
			// composite type.
			decls.AddAll(NewEnumTranscoderDeclarer(typ))
		}

	case *gotype.CompositeType:
		decls.AddAll(
			NewCompositeTypeDeclarer(typ),
			NewCompositeTranscoderDeclarer(typ),
			NewTypeResolverDeclarer(),
		)
		for _, childType := range typ.FieldTypes {
			findOutputDeclsHelper(childType, decls, true)
		}

	case *gotype.ArrayType:
		if gotype.IsPgxSupportedArray(typ) {
			return
		}
		decls.AddAll(NewTypeResolverDeclarer())
		switch gotype.UnwrapNestedType(typ.Elem).(type) {
		case *gotype.CompositeType, *gotype.EnumType:
			decls.AddAll(NewArrayDecoderDeclarer(typ))
		}
		findOutputDeclsHelper(typ.Elem, decls, hadCompositeParent)

	default:
		return
	}
}

// ConstantDeclarer declares a new string literal.
type ConstantDeclarer struct {
	key string
	str string
}

func NewConstantDeclarer(key, str string) ConstantDeclarer {
	return ConstantDeclarer{key, str}
}

func (c ConstantDeclarer) DedupeKey() string              { return c.key }
func (c ConstantDeclarer) Declare(string) (string, error) { return c.str, nil }

const typeResolverInitDecl = `
func register(conn *pgx.Conn){
	//
}
`

// NewTypeResolverInitDeclarer declare type resolver init code always needed.
func NewTypeResolverInitDeclarer() ConstantDeclarer {
	return NewConstantDeclarer("type_resolver::00_common", typeResolverInitDecl)
}

const typeResolverBodyDecl = `
/*type compositeField struct {
	name       string                 // name of the field
	typeName   string                 // Postgres type name
	defaultCodec pgtype.Codec // default value to use
}

func (tr *typeResolver) newCompositeValue(name string, fields ...compositeField) pgtype.Codec {
	if _, codec, ok := tr.findCodec(name); ok {
		return codec
	}

	codecs := make([]pgtype.CompositeCodecField, len(fields))
	isBinaryOk := true
	
	for i, field := range fields {
		oid, codec, ok := tr.findCodec(field.typeName)
		if !ok {
			oid = pgtype.UnknownOID
			codec = field.defaultCodec
		}
		isBinaryOk = isBinaryOk && oid != pgtype.UnknownOID
		
		codecs[i] = pgtype.CompositeCodecField{
			Name: field.name,
			Type: &pgtype.Type{Codec: codec, Name: field.typeName, OID: oid},
		}
	}
	// Okay to ignore error because it's only thrown when the number of field
	// names does not equal the number of ValueTranscoders.
	codec := pgtype.CompositeCodec{Fields: codecs}
	// typ, _ := pgtype.NewCompositeTypeValues(name, fs, codecs)
	// if !isBinaryOk {
	// 	return textPreferrer{ValueTranscoder: typ, typeName: name}
	// }
	return codec
}

func (tr *typeResolver) newArrayValue(name, elemName string, defaultVal func() pgtype.ValueTranscoder) pgtype.Codec {
	if _, val, ok := tr.findCodec(name); ok {
		return val
	}
	
	pgType, ok := tr.pgMap.TypeForName(elemName)
	if !ok {
		panic("unhandled")
	}
	
	return &pgtype.ArrayCodec{ElementType: pgType}
}*/`

// NewTypeResolverDeclarer declares type resolver body code sometimes needed.
func NewTypeResolverDeclarer() ConstantDeclarer {
	return NewConstantDeclarer("type_resolver::01_common", typeResolverBodyDecl)
}

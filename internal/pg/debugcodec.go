package pg

import (
	"database/sql/driver"

	"github.com/jackc/pgx/v5/pgtype"
)

var _ pgtype.Codec = (*DebugCodec)(nil)

type DebugCodec struct {
}

// DecodeDatabaseSQLValue implements pgtype.Codec.
func (d *DebugCodec) DecodeDatabaseSQLValue(m *pgtype.Map, oid uint32, format int16, src []byte) (driver.Value, error) {
	panic("unimplemented")
}

// DecodeValue implements pgtype.Codec.
func (d *DebugCodec) DecodeValue(m *pgtype.Map, oid uint32, format int16, src []byte) (any, error) {
	panic("unimplemented")
}

// FormatSupported implements pgtype.Codec.
func (d *DebugCodec) FormatSupported(int16) bool {
	panic("unimplemented")
}

// PlanEncode implements pgtype.Codec.
func (d *DebugCodec) PlanEncode(m *pgtype.Map, oid uint32, format int16, value any) pgtype.EncodePlan {
	return nil
}

// PlanScan implements pgtype.Codec.
func (d *DebugCodec) PlanScan(m *pgtype.Map, oid uint32, format int16, target any) pgtype.ScanPlan {
	return nil
}

// PreferredFormat implements pgtype.Codec.
func (d *DebugCodec) PreferredFormat() int16 {
	panic("unimplemented")
}

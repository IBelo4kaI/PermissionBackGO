// Package nullable provides conversion helpers between Go's database/sql
// Null* types and plain values / pointers, so service and handler code
// doesn't repeat sql.NullString{String: x, Valid: x != ""} boilerplate.
package nullable

import (
	"database/sql"
	"time"
)

// ── String ──────────────────────────────────────────────────

// String converts a plain string into sql.NullString.
// Empty string is treated as NULL — use StringPtr if empty-but-valid
// needs to be preserved.
func String(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}

// StringPtr converts *string into sql.NullString. nil is NULL,
// including an empty-but-non-nil string as a valid empty value.
func StringPtr(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: *s, Valid: true}
}

// StringOrNil converts sql.NullString back into *string, for JSON
// responses where NULL should serialize as null, not {"String":"","Valid":false}.
func StringOrNil(ns sql.NullString) *string {
	if !ns.Valid {
		return nil
	}
	return &ns.String
}

// StringOr returns ns.String if valid, otherwise fallback.
func StringOr(ns sql.NullString, fallback string) string {
	if !ns.Valid {
		return fallback
	}
	return ns.String
}

// ── Time ────────────────────────────────────────────────────

func Time(t time.Time) sql.NullTime {
	return sql.NullTime{Time: t, Valid: !t.IsZero()}
}

func TimePtr(t *time.Time) sql.NullTime {
	if t == nil {
		return sql.NullTime{}
	}
	return sql.NullTime{Time: *t, Valid: true}
}

func TimeOrNil(nt sql.NullTime) *time.Time {
	if !nt.Valid {
		return nil
	}
	return &nt.Time
}

// ── Bool ────────────────────────────────────────────────────

func Bool(b bool) sql.NullBool {
	return sql.NullBool{Bool: b, Valid: true}
}

func BoolPtr(b *bool) sql.NullBool {
	if b == nil {
		return sql.NullBool{}
	}
	return sql.NullBool{Bool: *b, Valid: true}
}

func BoolOrNil(nb sql.NullBool) *bool {
	if !nb.Valid {
		return nil
	}
	return &nb.Bool
}

// ── Int32 / Int64 ───────────────────────────────────────────

func Int32(i int32) sql.NullInt32 {
	return sql.NullInt32{Int32: i, Valid: true}
}

func Int32Ptr(i *int32) sql.NullInt32 {
	if i == nil {
		return sql.NullInt32{}
	}
	return sql.NullInt32{Int32: *i, Valid: true}
}

func Int32OrNil(ni sql.NullInt32) *int32 {
	if !ni.Valid {
		return nil
	}
	return &ni.Int32
}

func Int64(i int64) sql.NullInt64 {
	return sql.NullInt64{Int64: i, Valid: true}
}

func Int64Ptr(i *int64) sql.NullInt64 {
	if i == nil {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: *i, Valid: true}
}

func Int64OrNil(ni sql.NullInt64) *int64 {
	if !ni.Valid {
		return nil
	}
	return &ni.Int64
}

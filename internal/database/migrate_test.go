package database

import (
	"database/sql"
	"testing"
)

func TestCompatiblePhoneColumn(t *testing.T) {
	for _, tc := range []struct {
		typ, nullable string
		size          int64
		valid, want   bool
	}{
		{"varchar", "YES", 32, true, true},
		{"varchar", "YES", 255, true, true},
		{"varchar", "NO", 32, true, false},
		{"varchar", "YES", 18, true, true},
		{"varchar", "YES", 17, true, false},
		{"bigint", "YES", 32, true, false},
		{"varchar", "YES", 0, false, false},
	} {
		if got := compatiblePhoneColumn(tc.typ, tc.nullable, sql.NullInt64{Int64: tc.size, Valid: tc.valid}); got != tc.want {
			t.Errorf("%+v: got %v", tc, got)
		}
	}
}

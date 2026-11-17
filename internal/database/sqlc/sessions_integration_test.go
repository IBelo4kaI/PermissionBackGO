package repo

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Set TEST_MYSQL_DSN to run against MySQL; the test uses a connection-local temporary table.
func TestDeleteSessionsExpiryWindowMySQL(t *testing.T) {
	dsn := os.Getenv("TEST_MYSQL_DSN")
	if dsn == "" {
		t.Skip("TEST_MYSQL_DSN is not set")
	}
	ctx := context.Background()
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	conn, err := db.Conn(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	if _, err := conn.ExecContext(ctx, "SET time_zone = '+00:00'"); err != nil {
		t.Fatal(err)
	}
	if _, err := conn.ExecContext(ctx, "CREATE TEMPORARY TABLE sessions (id VARCHAR(36) PRIMARY KEY, user_id VARCHAR(36) NOT NULL, expires_at TIMESTAMP NOT NULL)"); err != nil {
		t.Fatal(err)
	}
	if _, err := conn.ExecContext(ctx, "INSERT INTO sessions VALUES ('today', 'user', '2026-11-17 21:32:33'), ('permanent', 'user', '2038-01-19 03:14:07'), ('other', 'other-user', '2026-11-17 21:32:33')"); err != nil {
		t.Fatal(err)
	}
	from := time.Date(2026, 11, 16, 21, 0, 0, 0, time.UTC)
	before := time.Date(2026, 11, 18, 21, 0, 0, 0, time.UTC)
	if err := New(conn).DeleteSessionsByUserIDBefore(ctx, DeleteSessionsByUserIDBeforeParams{UserID: "user", ExpiresFrom: from, ExpiresBefore: before}); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		id   string
		want int
	}{{"today", 0}, {"permanent", 1}, {"other", 1}} {
		var count int
		if err := conn.QueryRowContext(ctx, "SELECT COUNT(*) FROM sessions WHERE id = ?", tc.id).Scan(&count); err != nil {
			t.Fatal(err)
		}
		if count != tc.want {
			t.Errorf("session %s count = %d, want %d", tc.id, count, tc.want)
		}
	}
}

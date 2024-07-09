package author

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/robbert229/pggen"
	"github.com/robbert229/pggen/internal/pgtest"
	"github.com/stretchr/testify/assert"
)

func TestGenerate_Go_Example_Author(t *testing.T) {
	pool, cleanup := pgtest.NewPostgresSchema(t, []string{"schema.sql"}, func(config *pgxpool.Config) {
		config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
			err := Register(ctx, conn)
			if err != nil {
				return fmt.Errorf("failed to register types: %w", err)
			}

			return nil
		}
	})
	defer cleanup()

	tmpDir := t.TempDir()
	err := pggen.Generate(
		pggen.GenerateOptions{
			ConnString:       pool.Config().ConnString(),
			QueryFiles:       []string{"query.sql"},
			OutputDir:        tmpDir,
			GoPackage:        "author",
			Language:         pggen.LangGo,
			InlineParamCount: 2,
		})
	if err != nil {
		t.Fatalf("Generate() example/author: %s", err)
	}

	wantQueryFile := "query.sql.go"
	gotQueryFile := filepath.Join(tmpDir, "query.sql.go")
	assert.FileExists(t, gotQueryFile,
		"Generate() should emit query.sql.go")
	wantQueries, err := os.ReadFile(wantQueryFile)
	if err != nil {
		t.Fatalf("read wanted query.go.sql: %s", err)
	}
	gotQueries, err := os.ReadFile(gotQueryFile)
	if err != nil {
		t.Fatalf("read generated query.go.sql: %s", err)
	}
	assert.Equalf(t, string(wantQueries), string(gotQueries),
		"Got file %s; does not match contents of %s",
		gotQueryFile, wantQueryFile)
}

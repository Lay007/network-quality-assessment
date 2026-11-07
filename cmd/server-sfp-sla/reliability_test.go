//go:build linux

package main

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"strconv"
	"strings"
	"sync"
	"testing"
)

// Audit literal SQL call sites so placeholder regressions fail without MySQL.
func TestLegacySQLArgumentCounts(t *testing.T) {
	fs := token.NewFileSet()
	packages, err := parser.ParseDir(fs, ".", nil, 0)
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range packages {
		for _, f := range p.Files {
			ast.Inspect(f, func(n ast.Node) bool {
				call, ok := n.(*ast.CallExpr)
				if !ok {
					return true
				}
				method, ok := call.Fun.(*ast.SelectorExpr)
				if !ok {
					return true
				}
				switch method.Sel.Name {
				case "Exec", "Query", "QueryRow":
				default:
					return true
				}
				if len(call.Args) == 0 {
					return true
				}
				literal, ok := call.Args[0].(*ast.BasicLit)
				if !ok || literal.Kind != token.STRING {
					return true
				}
				query, err := strconv.Unquote(literal.Value)
				if err != nil {
					t.Fatal(err)
				}
				if strings.Count(query, "?") != len(call.Args)-1 {
					t.Errorf("%s: SQL placeholders=%d, arguments=%d", fs.Position(call.Pos()), strings.Count(query, "?"), len(call.Args)-1)
				}
				if strings.Contains(query, "SET status=?, datatime=?") {
					status, ok := call.Args[1].(*ast.BasicLit)
					if !ok || status.Value != "2" {
						t.Errorf("%s: initial status must be integer 2", fs.Position(call.Pos()))
					}
				}
				return true
			})
		}
	}
}

func TestPendingIDsRejectsUnknownTable(t *testing.T) {
	if _, err := pendingIDs(context.Background(), nil, "arbitrary"); err == nil {
		t.Fatal("accepted dynamic table")
	}
}

func TestGeneratorRejectsInvalidRateWithoutBlocking(t *testing.T) {
	for _, rate := range []int{0, -1} {
		counter, res := make(chan uint64, 1), make(chan uint64, 1)
		if elapsed := genSocket(0, []byte{1}, nil, 1, rate, counter, res); elapsed != 0 {
			t.Fatal(elapsed)
		}
		if <-counter != 0 || <-res != 0 {
			t.Fatal("invalid configuration generated packets")
		}
	}
}

type rowDriver struct{}
type rowConn struct{ fail bool }
type auditRows struct{ pos int }

var rowsClosed int
var registerAudit sync.Once

func (rowDriver) Open(name string) (driver.Conn, error) { return &rowConn{fail: name == "fail"}, nil }
func (*rowConn) Prepare(string) (driver.Stmt, error)    { return nil, errors.New("unsupported") }
func (*rowConn) Close() error                           { return nil }
func (*rowConn) Begin() (driver.Tx, error)              { return nil, errors.New("unsupported") }
func (c *rowConn) QueryContext(ctx context.Context, _ string, _ []driver.NamedValue) (driver.Rows, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if c.fail {
		return nil, errors.New("injected query error")
	}
	return &auditRows{}, nil
}
func (*auditRows) Columns() []string { return []string{"id"} }
func (*auditRows) Close() error      { rowsClosed++; return nil }
func (r *auditRows) Next(dest []driver.Value) error {
	if r.pos > 0 {
		return io.EOF
	}
	r.pos++
	dest[0] = int64(7)
	return nil
}
func TestPendingIDsClosesRowsAndPropagatesErrors(t *testing.T) {
	registerAudit.Do(func() { sql.Register("nqa-audit", rowDriver{}) })
	rowsClosed = 0
	db, err := sql.Open("nqa-audit", "")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	ids, err := pendingIDs(context.Background(), db, "test_latency")
	if err != nil || len(ids) != 1 || ids[0] != 7 || rowsClosed != 1 {
		t.Fatalf("ids=%v err=%v closed=%d", ids, err, rowsClosed)
	}
	bad, err := sql.Open("nqa-audit", "fail")
	if err != nil {
		t.Fatal(err)
	}
	defer bad.Close()
	if _, err := pendingIDs(context.Background(), bad, "test_latency"); err == nil {
		t.Fatal("lost query error")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := pendingIDs(ctx, db, "test_latency"); !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
}

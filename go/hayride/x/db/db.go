package db

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"io"

	"github.com/hayride-dev/bindings/go/internal/gen/imports/hayride/db/db"
	"github.com/hayride-dev/bindings/go/internal/gen/imports/hayride/db/types"
	"go.bytecodealliance.org/cm"
)

func init() {
	hayrideDriver := &Driver{}

	sql.Register("hayride", hayrideDriver)
}

type Driver struct{}

func (d *Driver) Open(name string) (driver.Conn, error) {

	result := db.Open(name)
	if result.IsErr() {
		return nil, fmt.Errorf("failed to connect to database: %s", result.Err().Data())
	}

	return &driverConn{conn: *result.OK()}, nil
}

type driverConn struct {
	conn db.Connection
}

func (c *driverConn) Prepare(query string) (driver.Stmt, error) {
	result := c.conn.Prepare(query)
	if result.IsErr() {
		return nil, fmt.Errorf("failed to prepare database statement: %s", result.Err().Data())
	}

	return &driverStmt{stmt: *result.OK()}, nil
}

func (c *driverConn) Close() error {
	result := c.conn.Close()
	if result.IsErr() {
		return fmt.Errorf("failed to close database connection: %s", result.Err().Data())
	}
	return nil
}

func (c *driverConn) Begin() (driver.Tx, error) {
	return c.BeginTx(context.Background(), driver.TxOptions{})
}

func (c *driverConn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	var isolation db.IsolationLevel
	switch sql.IsolationLevel(opts.Isolation) {
	case sql.LevelDefault:
		isolation = types.IsolationLevelReadCommitted
	case sql.LevelReadUncommitted:
		isolation = types.IsolationLevelReadUncommitted
	case sql.LevelReadCommitted:
		isolation = types.IsolationLevelReadCommitted
	case sql.LevelWriteCommitted:
		isolation = types.IsolationLevelWriteCommitted
	case sql.LevelRepeatableRead:
		isolation = types.IsolationLevelRepeatableRead
	case sql.LevelSnapshot:
		isolation = types.IsolationLevelSnapshot
	case sql.LevelSerializable:
		isolation = types.IsolationLevelSerializable
	case sql.LevelLinearizable:
		isolation = types.IsolationLevelLinearizable
	default:
		return nil, fmt.Errorf("unsupported isolation level: %d", opts.Isolation)
	}

	result := c.conn.BeginTransaction(isolation, opts.ReadOnly)
	if result.IsErr() {
		return nil, fmt.Errorf("failed to begin transaction: %s", result.Err().Data())
	}

	return &driverTx{tx: *result.OK()}, nil
}

type driverStmt struct {
	stmt db.Statement
}

func (s *driverStmt) NumInput() int {
	return int(s.stmt.NumberParameters())
}

func (s *driverStmt) Exec(args []driver.Value) (driver.Result, error) {
	params := make([]types.DbValue, len(args))
	for i, arg := range args {
		params[i] = dbValueFromInterface(arg)
	}

	result := s.stmt.Execute(cm.ToList(params))
	if result.IsErr() {
		return nil, fmt.Errorf("failed to execute database statement: %s", result.Err().Data())
	}

	return driver.RowsAffected(*result.OK()), nil
}

func (s *driverStmt) Query(args []driver.Value) (driver.Rows, error) {
	params := make([]types.DbValue, len(args))
	for i, arg := range args {
		params[i] = dbValueFromInterface(arg)
	}

	result := s.stmt.Query(cm.ToList(params))
	if result.IsErr() {
		return nil, fmt.Errorf("failed to query database: %s", result.Err().Data())
	}

	return &driverRows{rows: *result.OK()}, nil
}

func (s *driverStmt) Close() error {
	result := s.stmt.Close()
	if result.IsErr() {
		return fmt.Errorf("failed to close database statement: %s", result.Err().Data())
	}
	return nil
}

type driverRows struct {
	rows db.Rows
}

func (r *driverRows) Columns() []string {
	return r.rows.Columns().Slice()
}

func (r *driverRows) Next(dest []driver.Value) error {
	result := r.rows.Next()
	if result.IsErr() {
		if result.Err().Code() == db.ErrorCodeEndOfRows {
			return io.EOF
		}

		return fmt.Errorf("failed to fetch next row: %s", result.Err().Data())
	}

	for i, val := range result.OK().Slice() {
		dest[i] = dbValueToInterface(val)
	}

	return nil
}

func (r *driverRows) Close() error {
	result := r.rows.Close()
	if result.IsErr() {
		return fmt.Errorf("failed to close rows: %s", result.Err().Data())
	}

	return nil
}

// dbValueFromInterface converts a Go driver.Value to a WIT types.DbValue
func dbValueFromInterface(value driver.Value) types.DbValue {
	if value == nil {
		return types.DbValueNull()
	}

	switch v := value.(type) {
	case int:
		return types.DbValueInt64(int64(v))
	case int8:
		return types.DbValueInt32(int32(v))
	case int16:
		return types.DbValueInt32(int32(v))
	case int32:
		return types.DbValueInt32(v)
	case int64:
		return types.DbValueInt64(v)
	case uint:
		return types.DbValueUint64(uint64(v))
	case uint8:
		return types.DbValueUint32(uint32(v))
	case uint16:
		return types.DbValueUint32(uint32(v))
	case uint32:
		return types.DbValueUint32(v)
	case uint64:
		return types.DbValueUint64(v)
	case float32:
		return types.DbValueFloat(float64(v))
	case float64:
		return types.DbValueDouble(v)
	case bool:
		return types.DbValueBoolean(v)
	case string:
		return types.DbValueStr(v)
	case []byte:
		return types.DbValueBinary(cm.ToList(v))
	default:
		// For unknown types, convert to string representation
		return types.DbValueStr(fmt.Sprintf("%v", v))
	}
}

// dbValueToInterface converts a WIT types.DbValue to a Go interface{} for driver.Value
func dbValueToInterface(dbVal types.DbValue) interface{} {
	switch dbVal.Tag() {
	case 0: // int32
		if v := dbVal.Int32(); v != nil {
			return *v
		}
	case 1: // int64
		if v := dbVal.Int64(); v != nil {
			return *v
		}
	case 2: // uint32
		if v := dbVal.Uint32(); v != nil {
			return *v
		}
	case 3: // uint64
		if v := dbVal.Uint64(); v != nil {
			return *v
		}
	case 4: // float
		if v := dbVal.Float(); v != nil {
			return *v
		}
	case 5: // double
		if v := dbVal.Double(); v != nil {
			return *v
		}
	case 6: // str
		if v := dbVal.Str(); v != nil {
			return *v
		}
	case 7: // boolean
		if v := dbVal.Boolean(); v != nil {
			return *v
		}
	case 8: // date
		if v := dbVal.Date(); v != nil {
			return *v
		}
	case 9: // time
		if v := dbVal.Time(); v != nil {
			return *v
		}
	case 10: // timestamp
		if v := dbVal.Timestamp(); v != nil {
			return *v
		}
	case 11: // binary
		if v := dbVal.Binary(); v != nil {
			return v.Slice()
		}
	case 12: // null
		return nil
	}
	return nil
}

type driverTx struct {
	tx db.Transaction
}

func (t *driverTx) Commit() error {
	result := t.tx.Commit()
	if result.IsErr() {
		return fmt.Errorf("failed to commit transaction: %s", result.Err().Data())
	}
	return nil
}

func (t *driverTx) Rollback() error {
	result := t.tx.Rollback()
	if result.IsErr() {
		return fmt.Errorf("failed to rollback transaction: %s", result.Err().Data())
	}
	return nil
}

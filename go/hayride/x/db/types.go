package db

import (
	"github.com/hayride-dev/bindings/go/internal/gen/imports/hayride/db/types"
	"go.bytecodealliance.org/cm"
)

type Row = types.Row
type IsolationLevel = types.IsolationLevel

const (
	IsolationLevelReadUncommitted = types.IsolationLevelReadUncommitted
	IsolationLevelReadCommitted   = types.IsolationLevelReadCommitted
	IsolationLevelWriteCommitted  = types.IsolationLevelWriteCommitted
	IsolationLevelRepeatableRead  = types.IsolationLevelRepeatableRead
	IsolationLevelSnapshot        = types.IsolationLevelSnapshot
	IsolationLevelSerializable    = types.IsolationLevelSerializable
	IsolationLevelLinearizable    = types.IsolationLevelLinearizable
)

type DbValue = types.DbValue
type None = struct{}
type Double float64
type Date string
type Time string
type Timestamp string

type DbValueType interface {
	None | int32 | int64 | uint32 | uint64 | float64 | Double | string | bool | Date | Time | Timestamp | cm.List[uint8]
}

func NewDbValue[T DbValueType](data T) DbValue {
	switch any(data).(type) {
	case int32:
		return cm.New[DbValue](0, data)
	case int64:
		return cm.New[DbValue](1, data)
	case uint32:
		return cm.New[DbValue](2, data)
	case uint64:
		return cm.New[DbValue](3, data)
	case float64:
		return cm.New[DbValue](4, data)
	case Double:
		return cm.New[DbValue](5, float64(any(data).(Double)))
	case string:
		return cm.New[DbValue](6, data)
	case bool:
		return cm.New[DbValue](7, data)
	case Date:
		return cm.New[DbValue](8, string(any(data).(Date)))
	case Time:
		return cm.New[DbValue](9, string(any(data).(Time)))
	case Timestamp:
		return cm.New[DbValue](10, string(any(data).(Timestamp)))
	case cm.List[uint8]:
		return cm.New[DbValue](11, data)
	default:
		return cm.New[DbValue](12, struct{}{}) // null case
	}
}

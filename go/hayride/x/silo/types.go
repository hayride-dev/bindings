package silo

import "github.com/hayride-dev/bindings/go/internal/gen/types/hayride/silo/types"

type ThreadMetadata = types.ThreadMetadata
type ThreadStatus = types.ThreadStatus

const (
	ThreadStatusUnknown    = types.ThreadStatusUnknown
	ThreadStatusProcessing = types.ThreadStatusProcessing
	ThreadStatusExited     = types.ThreadStatusExited
	ThreadStatusKilled     = types.ThreadStatusKilled
)

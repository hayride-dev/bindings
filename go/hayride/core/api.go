package core

import (
	"encoding/json"
	"fmt"

	"github.com/hayride-dev/bindings/go/internal/gen/hayride/core/types"
	"go.bytecodealliance.org/cm"
)

type RequestData struct{ types.RequestData }

func (d RequestData) MarshalJSON() ([]byte, error) {
	switch d.Tag() {
	case 0: // unknown
		if d.Unknown() {
			raw, err := json.Marshal(struct{}{})
			if err != nil {
				return nil, fmt.Errorf("failed to marshal unknown data: %w", err)
			}
			return json.Marshal(map[string]json.RawMessage{
				"unknown": raw,
			})
		}
		return nil, fmt.Errorf("data variant 'unknown' is empty")
	case 1: // cast
		if cast := d.Cast(); cast != nil {
			raw, err := json.Marshal(cast)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal cast: %w", err)
			}
			return json.Marshal(map[string]json.RawMessage{
				"cast": raw,
			})
		}
		return nil, fmt.Errorf("data variant 'cast' is empty")
	case 2: // session-id
		if sessionID := d.SessionID(); sessionID != nil {
			raw, err := json.Marshal(*sessionID)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal session-id: %w", err)
			}
			return json.Marshal(map[string]json.RawMessage{
				"session-id": raw,
			})
		}
		return nil, fmt.Errorf("data variant 'session-id' is empty")
	default:
		return nil, fmt.Errorf("unsupported data tag: %d", d.Tag())
	}
}

func (d *RequestData) UnmarshalJSON(data []byte) error {
	var temp map[string]json.RawMessage
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	if len(temp) != 1 {
		return fmt.Errorf("invalid Data format")
	}
	for key, raw := range temp {
		switch key {
		case "unknown":
			*d = RequestData{types.RequestDataUnknown()}
		case "cast":
			var cast types.Cast
			if err := json.Unmarshal(raw, &cast); err != nil {
				return fmt.Errorf("failed to unmarshal cast: %w", err)
			}
			*d = RequestData{types.RequestDataCast(cast)}
		case "session-id":
			var sessionID string
			if err := json.Unmarshal(raw, &sessionID); err != nil {
				return fmt.Errorf("failed to unmarshal session-id: %w", err)
			}
			*d = RequestData{types.RequestDataSessionID(sessionID)}
		default:
			return fmt.Errorf("unknown data variant: %s", key)
		}
	}
	return nil
}

type ResponseData struct{ types.ResponseData }

func (r ResponseData) MarshalJSON() ([]byte, error) {
	switch r.Tag() {
	case 0: // unknown
		if r.Unknown() {
			raw, err := json.Marshal(struct{}{})
			if err != nil {
				return nil, fmt.Errorf("failed to marshal unknown data: %w", err)
			}
			return json.Marshal(map[string]json.RawMessage{
				"unknown": raw,
			})
		}
		return nil, fmt.Errorf("data variant 'unknown' is empty")
	case 1: // sessions
		if sessions := r.Sessions(); sessions != nil {
			raw, err := json.Marshal(sessions)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal sessions: %w", err)
			}
			return json.Marshal(map[string]json.RawMessage{
				"sessions": raw,
			})
		}
		return nil, fmt.Errorf("data variant 'sessions' is empty")
	case 2: // session-id
		if sessionID := r.SessionID(); sessionID != nil {
			raw, err := json.Marshal(*sessionID)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal session-id: %w", err)
			}
			return json.Marshal(map[string]json.RawMessage{
				"session-id": raw,
			})
		}
		return nil, fmt.Errorf("data variant 'session-id' is empty")
	case 3: // session-status

		if status := r.SessionStatus(); status != nil {
			raw, err := json.Marshal(status)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal session-status: %w", err)
			}
			return json.Marshal(map[string]json.RawMessage{
				"session-status": raw,
			})
		}
		return nil, fmt.Errorf("data variant 'session-status' is empty")
	default:
		return nil, fmt.Errorf("unsupported data tag: %d", r.Tag())
	}
}

func (r *ResponseData) UnmarshalJSON(data []byte) error {
	var temp map[string]json.RawMessage
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	if len(temp) != 1 {
		return fmt.Errorf("invalid Data format")
	}
	for key, raw := range temp {
		switch key {
		case "unknown":
			*r = ResponseData{types.ResponseDataUnknown()}
		case "sessions":
			var sessions []types.ThreadMetadata
			if err := json.Unmarshal(raw, &sessions); err != nil {
				return fmt.Errorf("failed to unmarshal sessions: %w", err)
			}
			*r = ResponseData{types.ResponseDataSessions(cm.ToList(sessions))}
		case "session-id":
			var sessionID string
			if err := json.Unmarshal(raw, &sessionID); err != nil {
				return fmt.Errorf("failed to unmarshal session-id: %w", err)
			}
			*r = ResponseData{types.ResponseDataSessionID(sessionID)}
		case "session-status":
			var status types.ThreadStatus
			if err := json.Unmarshal(raw, &status); err != nil {
				return fmt.Errorf("failed to unmarshal session-status: %w", err)
			}
			*r = ResponseData{types.ResponseDataSessionStatus(status)}
		default:
			return fmt.Errorf("unknown data variant: %s", key)
		}
	}
	return nil
}

type Request struct {
	types.Request
}

func (r Request) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Data     RequestData `json:"data,omitempty"`
		Metadata [][2]string `json:"metadata,omitempty"`
	}{
		Data:     RequestData{r.Data},
		Metadata: r.Metadata.Slice(),
	})
}

func (r *Request) UnmarshalJSON(data []byte) error {
	var aux struct {
		Data     RequestData `json:"data"`
		Metadata [][2]string `json:"metadata"`
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	r.Data = aux.Data.RequestData
	r.Metadata = cm.ToList(aux.Metadata)
	return nil
}

type Response struct {
	types.Response
}

func (r Response) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Data  ResponseData `json:"data,omitempty"`
		Error string       `json:"error,omitempty"`
		Next  string       `json:"next,omitempty"`
		Prev  string       `json:"prev,omitempty"`
	}{
		Data:  ResponseData{r.Data},
		Error: r.Error,
		Next:  r.Next,
		Prev:  r.Prev,
	})
}

func (r *Response) UnmarshalJSON(data []byte) error {
	var aux struct {
		Data  ResponseData `json:"data,omitempty"`
		Error string       `json:"error,omitempty"`
		Next  string       `json:"next,omitempty"`
		Prev  string       `json:"prev,omitempty"`
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	r.Data = aux.Data.ResponseData
	r.Error = aux.Error
	r.Next = aux.Next
	r.Prev = aux.Prev
	return nil
}

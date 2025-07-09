package types

import (
	"encoding/json"
	"fmt"

	"go.bytecodealliance.org/cm"
)

func (r Request) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Data     RequestData `json:"data,omitempty"`
		Metadata [][2]string `json:"metadata,omitempty"`
	}{
		Data:     r.Data,
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
	r.Data = aux.Data
	r.Metadata = cm.ToList(aux.Metadata)
	return nil
}

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
	case 3: // generate
		if generate := d.Generate(); generate != nil {
			raw, err := json.Marshal(generate)
			if err != nil {
				return nil, err
			}
			return json.Marshal(map[string]json.RawMessage{
				"generate": raw,
			})
		}
		return nil, fmt.Errorf("data variant 'generate' is empty")
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
			*d = RequestDataUnknown()
		case "cast":
			var cast Cast
			if err := json.Unmarshal(raw, &cast); err != nil {
				return fmt.Errorf("failed to unmarshal cast: %w", err)
			}
			*d = RequestDataCast(cast)
		case "session-id":
			var sessionID string
			if err := json.Unmarshal(raw, &sessionID); err != nil {
				return fmt.Errorf("failed to unmarshal session-id: %w", err)
			}
			*d = RequestDataSessionID(sessionID)
		case "generate":
			var generate Generate
			if err := json.Unmarshal(raw, &generate); err != nil {
				return err
			}
			*d = RequestDataGenerate(generate)
		default:
			return fmt.Errorf("unknown data variant: %s", key)
		}
	}
	return nil
}

func (r Response) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Data  ResponseData `json:"data,omitempty"`
		Error string       `json:"error,omitempty"`
		Next  string       `json:"next,omitempty"`
		Prev  string       `json:"prev,omitempty"`
	}{
		Data:  r.Data,
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
	r.Data = aux.Data
	r.Error = aux.Error
	r.Next = aux.Next
	r.Prev = aux.Prev
	return nil
}

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
	case 4: // messages
		if list := r.Messages(); list != nil {
			messages := *list
			raw, err := json.Marshal(messages)
			if err != nil {
				return nil, err
			}
			return json.Marshal(map[string]json.RawMessage{
				"messages": raw,
			})
		}
		return nil, fmt.Errorf("data variant 'messages' is empty")
	case 5: // path
		if path := r.Path(); path != nil {
			raw, err := json.Marshal(*path)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal path data: %w", err)
			}
			return json.Marshal(map[string]json.RawMessage{
				"path": raw,
			})
		}
		return nil, fmt.Errorf("data variant 'path' is empty")
	case 6: // paths
		if list := r.Paths(); list != nil {
			paths := *list
			raw, err := json.Marshal(paths)
			if err != nil {
				return nil, err
			}
			return json.Marshal(map[string]json.RawMessage{
				"paths": raw,
			})
		}
		return nil, fmt.Errorf("data variant 'paths' is empty")
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
			*r = ResponseDataUnknown()
		case "sessions":
			var sessions []ThreadMetadata
			if err := json.Unmarshal(raw, &sessions); err != nil {
				return fmt.Errorf("failed to unmarshal sessions: %w", err)
			}
			*r = ResponseDataSessions(cm.ToList(sessions))
		case "session-id":
			var sessionID string
			if err := json.Unmarshal(raw, &sessionID); err != nil {
				return fmt.Errorf("failed to unmarshal session-id: %w", err)
			}
			*r = ResponseDataSessionID(sessionID)
		case "session-status":
			var status ThreadStatus
			if err := json.Unmarshal(raw, &status); err != nil {
				return fmt.Errorf("failed to unmarshal session-status: %w", err)
			}
			*r = ResponseDataSessionStatus(status)
		case "messages":
			var messages []Message
			if err := json.Unmarshal(raw, &messages); err != nil {
				return err
			}
			*r = ResponseDataMessages(cm.ToList(messages))
		case "path":
			var path string
			if err := json.Unmarshal(raw, &path); err != nil {
				return fmt.Errorf("failed to unmarshal path: %w", err)
			}
			*r = ResponseDataPath(path)
		case "paths":
			var paths []string
			if err := json.Unmarshal(raw, &paths); err != nil {
				return fmt.Errorf("failed to unmarshal paths: %w", err)
			}
			*r = ResponseDataPaths(cm.ToList(paths))
		default:
			return fmt.Errorf("unknown data variant: %s", key)
		}
	}
	return nil
}

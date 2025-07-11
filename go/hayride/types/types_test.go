package types

import (
	"testing"

	"go.bytecodealliance.org/cm"
)

func TestRequestData(t *testing.T) {
	t.Run("Cast", func(t *testing.T) {
		data := NewRequestData(Cast{})
		if data.Tag() != 1 {
			t.Errorf("expected tag 1, got %d", data.Tag())
		}
	})

	t.Run("SessionID", func(t *testing.T) {
		data := NewRequestData(SessionID("test-session"))
		if data.Tag() != 2 {
			t.Errorf("expected tag 2, got %d", data.Tag())
		}
	})

	t.Run("Generate", func(t *testing.T) {
		data := NewRequestData(Generate{})
		if data.Tag() != 3 {
			t.Errorf("expected tag 3, got %d", data.Tag())
		}
	})
}

func TestResponseData(t *testing.T) {
	t.Run("ThreadMetadata", func(t *testing.T) {
		data := NewResponseData(cm.ToList([]ThreadMetadata{{}}))
		if data.Tag() != 1 {
			t.Errorf("expected tag 1, got %d", data.Tag())
		}
	})

	t.Run("SessionID", func(t *testing.T) {
		data := NewResponseData(SessionID("test-session"))
		if data.Tag() != 2 {
			t.Errorf("expected tag 2, got %d", data.Tag())
		}
	})

	t.Run("Path", func(t *testing.T) {
		data := NewResponseData(Path("test/path"))
		if data.Tag() != 5 {
			t.Errorf("expected tag 5, got %d", data.Tag())
		}
	})

	t.Run("Messages", func(t *testing.T) {
		data := NewResponseData(cm.ToList([]Message{{}}))
		if data.Tag() != 4 {
			t.Errorf("expected tag 4, got %d", data.Tag())
		}
	})

	t.Run("Paths", func(t *testing.T) {
		data := NewResponseData(cm.ToList([]string{"path1", "path2"}))
		if data.Tag() != 6 {
			t.Errorf("expected tag 6, got %d", data.Tag())
		}
	})

	t.Run("ThreadStatus", func(t *testing.T) {
		data := NewResponseData(ThreadStatus(0))
		if data.Tag() != 3 {
			t.Errorf("expected tag 3, got %d", data.Tag())
		}
	})
}

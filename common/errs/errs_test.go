package errs

import (
	"errors"
	"testing"
)

func TestError_Error(t *testing.T) {
	err := New(40002, "invalid params")
	got := err.Error()
	want := "ERR_INVALID_PARAMS(40002): invalid params"

	if got != want {
		t.Fatalf("unexpected Error() result.\nwant: %s\ngot:  %s", want, got)
	}
}

func TestNew(t *testing.T) {
	err := New(12345, "test error")

	if err.Code != 12345 {
		t.Fatalf("unexpected code, want 12345, got %d", err.Code)
	}
	if err.Message != "test error" {
		t.Fatalf("unexpected message, want 'test error', got %s", err.Message)
	}
}

func TestConvertErrorTypeToString(t *testing.T) {
	tests := []struct {
		code ErrorType
		want string
	}{
		{0, "SUCCESS(0)"},
		{40002, "ERR_INVALID_PARAMS(40002)"},
		{40003, "ERR_RECORD_NOT_FOUND(40003)"},
		{40004, "ERR_RECORD_ALREADY_EXISTS(40004)"},
		{50000, "ERR_INTERNAL(50000)"},
	}

	for _, tt := range tests {
		got := convertErrorTypeToString(tt.code)
		if got != tt.want {
			t.Fatalf("code %d: want %s, got %s", tt.code, tt.want, got)
		}
	}
}

func TestConvertErr_Nil(t *testing.T) {
	err := ConvertErr(nil)
	if err != Success {
		t.Fatalf("expected Success, got %+v", err)
	}
}

func TestConvertErr_CustomError(t *testing.T) {
	origin := ErrParam
	err := ConvertErr(origin)

	if err != origin {
		t.Fatalf("expected original error, got %+v", err)
	}
}

func TestConvertErr_StdError(t *testing.T) {
	stdErr := errors.New("some std error")
	err := ConvertErr(stdErr)

	if err.Code != ErrInternal.Code {
		t.Fatalf("unexpected code, want %d, got %d", ErrInternal.Code, err.Code)
	}
	if err.Message != "some std error" {
		t.Fatalf("unexpected message, want 'some std error', got %s", err.Message)
	}
}

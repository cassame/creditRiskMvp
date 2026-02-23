package internal

import (
	"credit-risk-mvp/internal/domain"
	"testing"
)

func TestBirthdate(t *testing.T) {
	// test format
	t.Run("invalid format", func(t *testing.T) {
		_, err := domain.NewBirthdate("25/12/1990")
		if err == nil {
			t.Error("Expected error for invalid format, got nil")
		}
	})

	// test logic
	t.Run("age calculation", func(t *testing.T) {
		//18+
		bd, _ := domain.NewBirthdate("2000-01-01")
		if bd.Age() < 18 {
			t.Errorf("Expected age to be >= 18, got %d", bd.Age())
		}
	})
}
func TestIsValidPhone(t *testing.T) {
	t.Run("valid phone number", func(t *testing.T) {
		_, err := domain.NewPhone("+79161234567")
		if err != nil {
			t.Error("Expected no error, got ", err)
		}
	})
	t.Run("invalid phone number", func(t *testing.T) {
		_, err := domain.NewPhone("123")
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})
}
func TestIsValidPassport(t *testing.T) {
	t.Run("valid passport", func(t *testing.T) {
		_, err := domain.NewPassport("7732 345645")
		if err != nil {
			t.Error("Expected no error, got ", err)
		}
	})
	t.Run("invalid passport", func(t *testing.T) {
		_, err := domain.NewPassport("77322312322323")
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})
}
func TestHasPatronymic(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{"Иванов Иван Иванович", true},
		{"Иванов Иван", false},
		{"   ", false},
	}
	for _, tt := range tests {
		n := domain.FullName(tt.name)
		if n.HasPatronymic() != tt.want {
			t.Errorf("HasPatronymic(%q) = %v, want %v", tt.name, n.HasPatronymic(), tt.want)
		}
	}
}
func TestApproveAmount(t *testing.T) {
	t.Run("valid amount", func(t *testing.T) {
		_, err := domain.NewAmount(500)
		if err != nil {
			t.Errorf("Expected valid amount, got error: %v", err)
		}
	})
	t.Run("negative amount", func(t *testing.T) {
		_, err := domain.NewAmount(-1)
		if err == nil {
			t.Error("Expected error for negative amount, got nil")
		}
	})
}

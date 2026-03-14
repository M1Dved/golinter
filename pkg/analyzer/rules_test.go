package analyzer

import "testing"

func TestCheckLowercase(t *testing.T) {
	tests := []struct {
		name string
		msg  string
		want bool
	}{
		{"lowercase start", "starting server", true},
		{"uppercase start", "Starting server", false},
		{"empty string", "", true},
		{"number start", "8080 port", true},
		{"single lowercase", "a", true},
		{"single uppercase", "A", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkLowercase(tt.msg); got != tt.want {
				t.Errorf("checkLowercase(%q) = %v, want %v", tt.msg, got, tt.want)
			}
		})
	}
}

func TestCheckEnglishOnly(t *testing.T) {
	tests := []struct {
		name string
		msg  string
		want bool
	}{
		{"english text", "starting server", true},
		{"russian text", "запуск сервера", false},
		{"mixed text", "server запущен", false},
		{"with numbers", "port 8080", true},
		{"empty string", "", true},
		{"emoji", "server 🚀", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkEnglishOnly(tt.msg); got != tt.want {
				t.Errorf("checkEnglishOnly(%q) = %v, want %v", tt.msg, got, tt.want)
			}
		})
	}
}

func TestCheckNoSpecialChars(t *testing.T) {
	tests := []struct {
		name string
		msg  string
		want bool
	}{
		{"clean text", "server started", true},
		{"exclamation mark", "server started!", false},
		{"emoji", "server started 🚀", false},
		{"colon", "warning: something", false},
		{"dots", "something went wrong...", true},
		{"hyphen", "health-check passed", true},
		{"underscore", "user_id processed", true},
		{"empty string", "", true},
		{"question mark", "what happened?", false},
		{"at sign", "email@example", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkNoSpecialChars(tt.msg); got != tt.want {
				t.Errorf("checkNoSpecialChars(%q) = %v, want %v", tt.msg, got, tt.want)
			}
		})
	}
}

func TestCheckNoSensitiveData(t *testing.T) {
	tests := []struct {
		name string
		msg  string
		want bool
	}{
		{"clean text", "user authenticated successfully", true},
		{"password keyword", "user password: ", false},
		{"token keyword", "token: ", false},
		{"api_key keyword", "api_key=", false},
		{"secret keyword", "secret value", false},
		{"uppercase keyword", "PASSWORD found", false},
		{"empty string", "", true},
		{"credential keyword", "user credential leaked", false},
		{"passwd keyword", "wrong passwd entered", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkNoSensitiveData(tt.msg, nil); got != tt.want {
				t.Errorf("checkNoSensitiveData(%q) = %v, want %v", tt.msg, got, tt.want)
			}
		})
	}
}

func TestCheckNoSensitiveDataWithExtra(t *testing.T) {
	extra := []string{"ssn", "credit_card"}

	tests := []struct {
		name string
		msg  string
		want bool
	}{
		{"standard keyword still works", "user password reset", false},
		{"extra keyword ssn", "user ssn is 123", false},
		{"extra keyword credit_card", "credit_card number", false},
		{"clean text", "user logged in", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkNoSensitiveData(tt.msg, extra); got != tt.want {
				t.Errorf("checkNoSensitiveData(%q, extra) = %v, want %v", tt.msg, got, tt.want)
			}
		})
	}
}

func TestTrimQuotes(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"double quotes", `"hello"`, "hello"},
		{"backticks", "`hello`", "hello"},
		{"no quotes", "hello", "hello"},
		{"empty quotes", `""`, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := trimQuotes(tt.input); got != tt.want {
				t.Errorf("trimQuotes(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

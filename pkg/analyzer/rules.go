package analyzer

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

var sensitiveKeywords = []string{
	"password", "passwd", "pwd",
	"token",
	"secret",
	"api_key", "apikey", "api-key",
	"private_key", "privatekey",
	"credential", "credentials",
	"auth_token",
}

func checkLowercase(msg string) bool {
	if len(msg) == 0 {
		return true
	}
	firstRune, _ := utf8.DecodeRuneInString(msg)
	return !unicode.IsUpper(firstRune)
}

func checkEnglishOnly(msg string) bool {
	for _, r := range msg {
		if r > unicode.MaxASCII {
			return false
		}
	}
	return true
}

func checkNoSpecialChars(msg string) bool {
	for _, r := range msg {
		if !isAllowedChar(r) {
			return false
		}
	}
	return true
}

func isAllowedChar(r rune) bool {
	if unicode.IsLetter(r) || unicode.IsDigit(r) {
		return true
	}
	switch r {
	case ' ', '.', ',', '-', '_', '\'', '/', '=', '(', ')':
		return true
	}
	return false
}

func checkNoSensitiveData(msg string, extraKeywords []string) bool {
	lower := strings.ToLower(msg)

	for _, keyword := range sensitiveKeywords {
		if strings.Contains(lower, keyword) {
			return false
		}
	}

	for _, keyword := range extraKeywords {
		if strings.Contains(lower, keyword) {
			return false
		}
	}

	return true
}

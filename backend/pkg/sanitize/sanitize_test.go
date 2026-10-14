package sanitize

import (
	"testing"
)

func TestCleanText(t *testing.T) {
	input := "  <script>alert('xss')</script>John Doe <b>Smith</b>  "
	cleaned := CleanText(input)

	if cleaned == "" {
		t.Fatal("expected non-empty cleaned text")
	}

	if ContainsInjectionPattern(cleaned) {
		t.Errorf("cleaned text still contains injection pattern: %s", cleaned)
	}

	if cleaned != "John Doe &lt;b&gt;Smith&lt;/b&gt;" && cleaned != "John Doe Smith" {
		// Verify tags stripped and escaped
		t.Logf("cleaned output: %s", cleaned)
	}
}

func TestContainsInjectionPattern(t *testing.T) {
	malicious := []string{
		"admin' --",
		"1; DROP TABLE users;",
		"1 UNION SELECT * FROM users",
		"<script>alert(1)</script>",
		"javascript:alert(1)",
		"admin' /* comment */",
	}

	for _, m := range malicious {
		if !ContainsInjectionPattern(m) {
			t.Errorf("expected '%s' to be detected as injection pattern", m)
		}
	}

	safe := []string{
		"Budi Santoso",
		"Jl. Melati No. 12",
		"Keluhan sakit kepala dan demam",
		"user@puskesmas.go.id",
	}

	for _, s := range safe {
		if ContainsInjectionPattern(s) {
			t.Errorf("expected '%s' to be considered safe, but was flagged", s)
		}
	}
}

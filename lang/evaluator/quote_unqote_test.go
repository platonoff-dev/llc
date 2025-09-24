package evaluator

import (
	"testing"

	"llc/lang/object"
)

func TestQuote(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "quote(5 + 8)",
			expected: "(5 + 8)",
		},
		{
			input:    "quote(foobar)",
			expected: "foobar",
		},
		{
			input:    "quote(foobar + barfoo)",
			expected: "(foobar + barfoo)",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		quote, ok := evaluated.(*object.Quote)
		if !ok {
			t.Fatalf("expected *object.Quote, got=%T (%+v)", evaluated, evaluated)
		}

		if quote.Node == nil {
			t.Fatalf("quote.Node should not be nil")
		}

		if quote.Node.String() != tt.expected {
			t.Errorf("not equal: expected %s, got %s", tt.expected, quote.Node.String())
		}
	}
}

func TestQuoteUnquote(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "quote(unquote(4))",
			expected: "4",
		},
		{
			input:    "quote(unquote(4 + 8))",
			expected: "12",
		},
		{
			input:    "quote(8 + unquote(4 + 4))",
			expected: "(8 + 8)",
		},
		{
			input:    "quote(unquote(4 + 4) + 8)",
			expected: "(8 + 8)",
		},
		{
			input:    "let foobar = 8; quote(foobar)",
			expected: "foobar",
		},
		{
			input:    "let foobar = 8; quote(unquote(foobar))",
			expected: "8",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		quote, ok := evaluated.(*object.Quote)
		if !ok {
			t.Fatalf("expected *object.Quote, got=%T (%+v)", evaluated, evaluated)
		}

		if quote.Node == nil {
			t.Fatalf("quote.Node should not be nil")
		}

		if quote.Node.String() != tt.expected {
			t.Errorf("not equal: expected %s, got %s", tt.expected, quote.Node.String())
		}
	}
}

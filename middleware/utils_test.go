package middleware

import "testing"

func TestMapPublicError(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		ok      bool
		code    string
		i18nKey string
	}{
		{
			name:    "disabled token message",
			input:   "该令牌状态不可用",
			ok:      true,
			code:    "billing_credits_exhausted",
			i18nKey: "billing_credits_exhausted",
		},
		{
			name:    "exhausted token message",
			input:   "令牌 bot_1（#1）额度已用尽",
			ok:      true,
			code:    "billing_credits_exhausted",
			i18nKey: "billing_credits_exhausted",
		},
		{
			name:    "english unavailable token",
			input:   "token status is unavailable",
			ok:      true,
			code:    "billing_credits_exhausted",
			i18nKey: "billing_credits_exhausted",
		},
		{
			name:  "unrelated message",
			input: "invalid request body",
			ok:    false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			i18nKey, code, ok := mapPublicError(tc.input)
			if ok != tc.ok {
				t.Fatalf("expected ok=%v, got %v", tc.ok, ok)
			}
			if code != tc.code {
				t.Fatalf("expected code=%q, got %q", tc.code, code)
			}
			if i18nKey != tc.i18nKey {
				t.Fatalf("expected i18nKey=%q, got %q", tc.i18nKey, i18nKey)
			}
		})
	}
}

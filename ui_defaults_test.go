package main

import (
	"strings"
	"testing"
)

// TestEmbeddedUIDefaults pins the de-branded defaults of the embedded
// wizard UI (index.html):
//
//   - dev split stays an advanced-only 0-50 slider but DEFAULTS TO 0 —
//     the wizard takes no cut unless the operator explicitly opts in
//   - no "developer fund" marketing copy in the UI
//   - Minibits is a selectable mint option; the stale DLEQ-workaround
//     note is gone (the underlying keyset-rotation bug is tracked
//     upstream on OpenTollGate/tollgate-module-basic-go, not papered
//     over in wizard copy)
func TestEmbeddedUIDefaults(t *testing.T) {
	ui := string(indexHTML)

	t.Run("dev split slider defaults to 0 (advanced-only, 0-50)", func(t *testing.T) {
		if !strings.Contains(ui, `id="devsplit" min="0" max="50" value="0"`) {
			t.Error(`devsplit slider must keep min=0 max=50 but default to value="0"`)
		}
		if !strings.Contains(ui, `id="devsplit-val">0%<`) {
			t.Error("devsplit default label must read 0%")
		}
	})

	t.Run("no developer-fund copy", func(t *testing.T) {
		if strings.Contains(strings.ToLower(ui), "developer fund") {
			t.Error(`UI copy must not mention "developer fund"`)
		}
	})

	t.Run("minibits selectable, no stale DLEQ note", func(t *testing.T) {
		if !strings.Contains(ui, `value="https://minibits.cash/cashu/api/v1/"`) {
			t.Error("Minibits must be a selectable mint option (DLEQ filter removed)")
		}
		if strings.Contains(ui, "DLEQ") {
			t.Error("UI copy must not carry the stale DLEQ-workaround note (tracked upstream instead)")
		}
	})
}

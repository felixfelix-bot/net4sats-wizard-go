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

	t.Run("minibits selectable at the real mint base, no stale DLEQ note", func(t *testing.T) {
		// The option must point at the working Minibits mint base (the same
		// URL defaultMints pins, deploy.go). The pre-rebrand value
		// https://minibits.cash/cashu/api/v1/ is an LNbits-era path that no
		// longer serves a Cashu API (404) — posting it verbatim as req.Mint
		// would configure the router against a dead endpoint.
		if !strings.Contains(ui, `value="https://mint.minibits.cash/Bitcoin"`) {
			t.Error("Minibits option must point at the real mint base https://mint.minibits.cash/Bitcoin (the defaultMints entry)")
		}
		if strings.Contains(ui, "https://minibits.cash/cashu/api/v1/") {
			t.Error("UI must not ship the dead LNbits-era Minibits path https://minibits.cash/cashu/api/v1/")
		}
		if strings.Contains(ui, "DLEQ") {
			t.Error("UI copy must not carry the stale DLEQ-workaround note (tracked upstream instead)")
		}
	})
}

package main

import (
	"encoding/json"
	"os/exec"
	"reflect"
	"strings"
	"testing"
)

// wantProdMintURLs is the fixed production mint set that every real
// customer deployment must ship — and the ONLY set when test_mints is off.
var wantProdMintURLs = []string{
	"https://mint.coinos.io",
	"https://mint.minibits.cash/Bitcoin",
	"https://mint.lnserver.com",
	"https://mint.macadamia.cash",
	"https://mint.westernbtc.com",
	"https://kashu.me",
	"https://mint.cubabitcoin.org",
}

func mintURLs(t *testing.T, raw string) []string {
	t.Helper()
	var mints []mintCfg
	if err := json.Unmarshal([]byte(raw), &mints); err != nil {
		t.Fatalf("default mints payload is not valid JSON: %v\n%s", err, raw)
	}
	urls := make([]string, 0, len(mints))
	for _, m := range mints {
		urls = append(urls, m.URL)
	}
	return urls
}

// TestDeployPayloadTestMintsDefaultOff is the core opt-in guarantee: a
// deploy payload WITHOUT the test_mints field (or with it false) must
// produce a mint config containing EXACTLY the 7 production mints and
// ZERO testnut entries. Real customer deployments must never carry the
// fake-Lightning testnut mints.
func TestDeployPayloadTestMintsDefaultOff(t *testing.T) {
	payloads := []string{
		// field absent entirely (legacy UI payloads)
		`{"ip":"192.168.1.1","password":"pw","lnurl":"test@wallet.app"}`,
		// field explicitly false
		`{"ip":"192.168.1.1","password":"pw","lnurl":"test@wallet.app","test_mints":false}`,
		// field false alongside other advanced fields
		`{"ip":"192.168.1.1","password":"pw","lnurl":"test@wallet.app","test_mints":false,"mint":"https://mint.example.com","devSplit":10,"margin":5}`,
	}
	for i, p := range payloads {
		var req deployRequest
		if err := json.Unmarshal([]byte(p), &req); err != nil {
			t.Fatalf("payload %d: decode: %v", i, err)
		}
		if req.TestMints {
			t.Errorf("payload %d: TestMints = true, want false (opt-in flag must default to false)", i)
		}
		got := mintURLs(t, defaultMintsJSON(req.TestMints))
		if len(got) != 7 {
			t.Errorf("payload %d: default mints count = %d, want exactly 7 (got %v)", i, len(got), got)
		}
		for _, u := range got {
			if strings.Contains(u, "testnut") {
				t.Errorf("payload %d: testnut mint %q present with test_mints off — testnut must be opt-in only", i, u)
			}
		}
		if !reflect.DeepEqual(got, wantProdMintURLs) {
			t.Errorf("payload %d: mints = %v, want %v", i, got, wantProdMintURLs)
		}
	}
}

// TestDeployPayloadTestMintsOptIn: test_mints=true appends BOTH testnut
// entries after the 7 production mints (9 total, order preserved).
func TestDeployPayloadTestMintsOptIn(t *testing.T) {
	var req deployRequest
	if err := json.Unmarshal([]byte(`{"ip":"192.168.1.1","password":"pw","lnurl":"test@wallet.app","test_mints":true}`), &req); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if !req.TestMints {
		t.Fatal("TestMints = false, want true when payload sets test_mints:true")
	}
	got := mintURLs(t, defaultMintsJSON(req.TestMints))
	if len(got) != 9 {
		t.Fatalf("mints count = %d, want 9 (7 production + 2 testnut); got %v", len(got), got)
	}
	want := append(append([]string{}, wantProdMintURLs...),
		"https://nofee.testnut.cashu.space",
		"https://testnut.cashu.space")
	if !reflect.DeepEqual(got, want) {
		t.Errorf("mints = %v, want %v", got, want)
	}
	testnut := 0
	for _, u := range got {
		if strings.Contains(u, "testnut.cashu.space") {
			testnut++
		}
	}
	if testnut != 2 {
		t.Errorf("testnut mint count = %d, want both testnut entries", testnut)
	}
}

// TestDefaultMintsJSONFieldValues pins the per-mint numeric config so the
// testnut entries stay zero-fee/no-payout and production mints keep their
// top-up thresholds.
func TestDefaultMintsJSONFieldValues(t *testing.T) {
	for _, include := range []bool{false, true} {
		var mints []mintCfg
		if err := json.Unmarshal([]byte(defaultMintsJSON(include)), &mints); err != nil {
			t.Fatalf("include=%v: %v", include, err)
		}
		for _, m := range mints {
			if m.PriceUnit != "sats" || m.PricePerStep != 1 || m.MinPurchaseSteps != 0 {
				t.Errorf("mint %s: unexpected pricing fields %+v", m.URL, m)
			}
			if strings.Contains(m.URL, "testnut") {
				if m.MinBalance != 0 || m.PayoutIntervalSeconds != 999999 || m.MinPayoutAmount != 999999 || m.BalanceTolerancePercent != 0 {
					t.Errorf("testnut mint %s: must stay zero-balance/no-payout, got %+v", m.URL, m)
				}
			} else {
				if m.MinBalance != 64 || m.BalanceTolerancePercent != 10 || m.PayoutIntervalSeconds != 60 || m.MinPayoutAmount != 128 {
					t.Errorf("production mint %s: unexpected thresholds, got %+v", m.URL, m)
				}
			}
		}
	}
}

// TestMintConfigJqMerge runs the EXACT jq filter that deploy step 8b builds
// against sample config.json states (via local jq when available) to prove:
//   - the merge actually works — the alpha16 filter errored with
//     "Cannot iterate over null" and never wrote mints at all;
//   - a config missing accepted_mints is handled;
//   - merging is idempotent (re-running adds nothing);
//   - an empty custom mint does NOT append a {"url":""} entry (alpha16 bug);
//   - a custom mint matching a default URL is not duplicated.
func TestMintConfigJqMerge(t *testing.T) {
	jqPath, err := exec.LookPath("jq")
	if err != nil {
		t.Skip("jq not installed — skipping end-to-end jq filter test")
	}
	extractURLs := func(cfgJSON string) []string {
		t.Helper()
		var cfg struct {
			AcceptedMints []mintCfg `json:"accepted_mints"`
		}
		if err := json.Unmarshal([]byte(cfgJSON), &cfg); err != nil {
			t.Fatalf("merged config not valid JSON: %v\n%s", err, cfgJSON)
		}
		urls := make([]string, 0, len(cfg.AcceptedMints))
		for _, m := range cfg.AcceptedMints {
			urls = append(urls, m.URL)
		}
		return urls
	}

	baseConfig := `{"margin":0,"profit_share":[{"identity":"owner","factor":1},{"identity":"developer","factor":0}],"accepted_mints":[{"url":"https://mint.coinos.io","min_balance":64}]}`
	emptyMintsConfig := `{"margin":0,"profit_share":[{"identity":"owner","factor":1},{"identity":"developer","factor":0}]}`

	cases := []struct {
		name       string
		config     string
		mintArg    string
		want       []string
		wantMargin float64
		wantOwnerF float64
	}{
		{
			name:       "empty custom mint adds only defaults, existing coinos kept, no dup",
			config:     baseConfig,
			mintArg:    "",
			want:       wantProdMintURLs,
			wantMargin: 0, wantOwnerF: 0.9,
		},
		{
			name:       "custom mint appended once",
			config:     baseConfig,
			mintArg:    "https://mint.example.com",
			want:       append(append([]string{"https://mint.coinos.io"}, "https://mint.example.com"), wantProdMintURLs[1:]...),
			wantMargin: 0, wantOwnerF: 0.9,
		},
		{
			name:    "custom mint matching a default URL is not duplicated",
			config:  baseConfig,
			mintArg: "https://kashu.me",
			// $mu=kashu is appended (not yet present in base config), then
			// the $dm merge skips it ($have already includes it): same SET
			// as the production list, but kashu sits at its appended slot.
			want: []string{"https://mint.coinos.io", "https://kashu.me",
				"https://mint.minibits.cash/Bitcoin", "https://mint.lnserver.com",
				"https://mint.macadamia.cash", "https://mint.westernbtc.com",
				"https://mint.cubabitcoin.org"},
			wantMargin: 0, wantOwnerF: 0.9,
		},
		{
			name:       "config without accepted_mints gets defaults",
			config:     emptyMintsConfig,
			mintArg:    "",
			want:       wantProdMintURLs,
			wantMargin: 0, wantOwnerF: 0.9,
		},
		{
			name:       "testnut opt-in payload adds both testnut mints",
			config:     baseConfig,
			mintArg:    "",
			want:       append(append([]string{}, wantProdMintURLs...), "https://nofee.testnut.cashu.space", "https://testnut.cashu.space"),
			wantMargin: 0, wantOwnerF: 0.9,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			includeTestnut := strings.Contains(tc.name, "testnut")
			dm := defaultMintsJSON(includeTestnut)
			margin := 0
			ownerF := "0.9000"
			if tc.wantMargin != 0 {
				margin = int(tc.wantMargin)
			}
			_ = margin
			filter := configJqFilter()
			args := []string{"--argjson", "m", "0", "--argjson", "of", ownerF, "--argjson", "df", "0.1000", "--argjson", "dm", dm, "--arg", "mu", tc.mintArg, filter}
			cmd := exec.Command(jqPath, args...)
			cmd.Stdin = strings.NewReader(tc.config)
			out, err := cmd.Output()
			if err != nil {
				if ee, ok := err.(*exec.ExitError); ok {
					t.Fatalf("jq failed: %v\nstderr: %s", err, ee.Stderr)
				}
				t.Fatalf("jq failed to start: %v", err)
			}
			got := extractURLs(string(out))
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("merged mints = %v\nwant           = %v", got, tc.want)
			}
			// Idempotency: run the same merge over the merged output again —
			// mint list must not change.
			cmd2 := exec.Command(jqPath, args...)
			cmd2.Stdin = strings.NewReader(string(out))
			out2, err2 := cmd2.Output()
			if err2 != nil {
				t.Fatalf("second pass jq failed: %v", err2)
			}
			got2 := extractURLs(string(out2))
			if !reflect.DeepEqual(got2, got) {
				t.Errorf("merge not idempotent:\nfirst  = %v\nsecond = %v", got, got2)
			}
			// Margin / profit_share written correctly.
			var cfg struct {
				Margin      float64 `json:"margin"`
				ProfitShare []struct {
					Identity string  `json:"identity"`
					Factor   float64 `json:"factor"`
				} `json:"profit_share"`
			}
			if err := json.Unmarshal(out, &cfg); err != nil {
				t.Fatalf("parse merged: %v", err)
			}
			if cfg.Margin != tc.wantMargin {
				t.Errorf("margin = %v, want %v", cfg.Margin, tc.wantMargin)
			}
			for _, ps := range cfg.ProfitShare {
				if ps.Identity == "owner" && ps.Factor != tc.wantOwnerF {
					t.Errorf("owner factor = %v, want %v", ps.Factor, tc.wantOwnerF)
				}
			}
		})
	}

	// Verify a nonzero margin value is written too (the --argjson m arg).
	t.Run("nonzero margin", func(t *testing.T) {
		args := []string{"--argjson", "m", "5", "--argjson", "of", "0.9000", "--argjson", "df", "0.1000", "--argjson", "dm", defaultMintsJSON(false), "--arg", "mu", "", configJqFilter()}
		cmd := exec.Command(jqPath, args...)
		cmd.Stdin = strings.NewReader(baseConfig)
		outB, err := cmd.Output()
		if err != nil {
			t.Fatalf("jq: %v", err)
		}
		var cfg struct {
			Margin float64 `json:"margin"`
		}
		if err := json.Unmarshal(outB, &cfg); err != nil {
			t.Fatalf("parse: %v", err)
		}
		if cfg.Margin != 5 {
			t.Errorf("margin = %v, want 5", cfg.Margin)
		}
	})
}

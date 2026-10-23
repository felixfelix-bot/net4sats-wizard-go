#!/usr/bin/env python3
"""Validate artifacts/demo_basecut/trim_plan.json for t_ae815553.

This plan is intentionally NON-EXECUTABLE (upstream inputs missing). The validator
therefore enforces two families of invariant:

  A. SAFETY invariants that must hold for ANY version of this plan, executable or not:
     - no interior content is proposed for removal
     - no trim point falls inside a step interval
     - every step has both a start and an end timestamp
     - no truncated/partial timestamps are invented
  B. FAIL-CLOSED invariants specific to the BLOCKED state:
     - status == BLOCKED  =>  executable == false, source.path is null,
       all four cut points null, step list empty, region list empty
  C. PROVENANCE invariants (acceptance criterion: a dead-air detection must be
     reflected together with the command it came from):
     - at least one detection, each with a non-empty command string
     - detections applied to a non-source file must be marked POSITIVE_CONTROL_ONLY
       so they can never be confused with the plan's own regions.

Exit 0 = all assertions passed. Exit 1 = at least one failure (printed).
"""
import json
import sys
from pathlib import Path

PLAN = Path(__file__).resolve().parent / "trim_plan.json"

fails: list[str] = []
checks: list[tuple[str, str]] = []


def ck(name: str, ok: bool, detail: str = "") -> None:
    checks.append(("PASS" if ok else "FAIL", name + (f" — {detail}" if detail and not ok else "")))
    if not ok:
        fails.append(name + (f" — {detail}" if detail else ""))


plan = json.loads(PLAN.read_text())
status = plan.get("status")

# ---------- A. universal safety invariants ----------

steps = plan["step_map"]["ordered_steps"]
ck(
    "A1 every step has both start and end",
    all(
        isinstance(s.get("start"), (int, float)) and isinstance(s.get("end"), (int, float))
        for s in steps
    ),
    f"{len(steps)} step(s) inspected",
)

# A2: no trim point may land strictly inside a step interval.
cuts = plan["cuts"]
trim_points = {
    k: cuts.get(k) for k in ("head_trim_in_s", "head_trim_out_s", "tail_trim_in_s", "tail_trim_out_s")
}
violations = []
for cname, t in trim_points.items():
    if t is None:
        continue
    for s in steps:
        lo, hi = s.get("start"), s.get("end")
        if isinstance(lo, (int, float)) and isinstance(hi, (int, float)) and lo < t < hi:
            violations.append(f"{cname}={t} inside step '{s.get('step')}' [{lo},{hi}]")
ck("A2 no trim point falls inside a step interval", not violations, "; ".join(violations))

# A3: interior content must never be proposed for removal. Region entries may be
# flagged, but no cut/removal directive may appear anywhere.
region_blob = json.dumps(plan["flagged_interior_dead_air"]).lower()
bad_tokens = ("\"remove\"", "\"cut\"", "\"delete\"", "action\": \"cut", "marked_for_removal\": true")
hit = [t for t in bad_tokens if t in region_blob]
ck("A3 interior regions are flag-only (no removal directive)", not hit, f"found {hit}")
ck(
    "A3b interior_cuts_made == 0",
    cuts.get("interior_cuts_made") == 0,
    f"got {cuts.get('interior_cuts_made')!r}",
)
ck(
    "A3c removal_proposed is false",
    plan["flagged_interior_dead_air"].get("removal_proposed") is False,
    "must be explicitly False",
)

# ---------- B. fail-closed invariants for the BLOCKED state ----------

if status == "BLOCKED":
    ck("B1 BLOCKED plan is marked non-executable", plan.get("executable") is False)
    ck(
        "B2 BLOCKED plan names no source path",
        plan["source"].get("path") is None,
        f"got {plan['source'].get('path')!r}",
    )
    ck(
        "B3 BLOCKED plan states no source properties",
        plan["source"].get("properties") is None,
        "properties must be null when no file is identified",
    )
    all_null = all(v is None for v in trim_points.values())
    ck(
        "B4 all four trim points are null",
        all_null,
        f"non-null: { {k: v for k, v in trim_points.items() if v is not None} }",
    )
    ck("B5 step list is empty", len(steps) == 0, f"{len(steps)} step(s) present")
    ck(
        "B6 flagged interior region list is empty",
        len(plan["flagged_interior_dead_air"]["regions"]) == 0,
    )
    ck(
        "B7 framing anomaly list is empty",
        len(plan["framing_anomalies"]["anomalies"]) == 0,
    )
    ck(
        "B8 encoder is told not to proceed",
        cuts.get("encode_must_not_proceed") is True
        and plan["handoff_to_encoder"].get("instruction", "").startswith("DO NOT RUN"),
    )
    # The stop reason must cite all three missing upstreams, not a vague excuse.
    missing = plan["stop_reason"]["inputs_expected"]
    ck(
        "B9 stop_reason names all three upstream inputs",
        len(missing) == 3 and all(m["verdict"] == "MISSING" for m in missing),
        f"{len(missing)} declared",
    )
    ck(
        "B10 every cited upstream card shows 0 runs",
        all(m.get("observed_runs") == 0 for m in missing),
        str([(m["from_card"], m.get("observed_runs")) for m in missing]),
    )

# ---------- C. provenance invariants ----------

ev = plan.get("dead_air_detection_evidence", {})
dets = ev.get("detections", [])
ck("C1 at least one dead-air detection is present", len(dets) >= 1, f"{len(dets)} found")
ck(
    "C2 every detection carries the command it came from",
    all(isinstance(d.get("command"), str) and d["command"].strip() for d in dets),
    "a detection without its command cannot satisfy the acceptance criterion",
)
ck(
    "C3 at least one detection ran successfully (exit_code 0)",
    any(d.get("exit_code") == 0 for d in dets),
)
raw_ok = any(len(d.get("raw_output_lines", [])) or d.get("result") for d in dets)
ck("C4 detection output is summarised, not merely claimed", raw_ok)

# The scope guard: detections run against something other than the plan's source
# must be labelled positive-control-only, otherwise they look like plan regions.
src_path = plan["source"].get("path")
mislabelled = [
    d.get("applied_to")
    for d in dets
    if d.get("applied_to") != src_path and ev.get("status") != "POSITIVE_CONTROL_ONLY"
]
ck(
    "C5 non-source detections are marked POSITIVE_CONTROL_ONLY",
    not mislabelled,
    f"unlabelled: {mislabelled}",
)
ck(
    "C6 positive-control evidence warns it is not applicable to an encode",
    "NOT" in ev.get("warning", "").upper() and "MUST NOT" in ev.get("warning", "").upper(),
    "warning must explicitly forbid applying these regions",
)

# ---------- report ----------

width = max(len(n) for _, n in checks)
print(f"trim_plan.json validator — {PLAN}")
print(f"status={status}  executable={plan.get('executable')}  steps={len(steps)}  "
      f"regions={len(plan['flagged_interior_dead_air']['regions'])}  detections={len(dets)}")
print("-" * min(width + 8, 100))
for verdict, name in checks:
    print(f"  [{verdict}] {name}")
print("-" * min(width + 8, 100))
print(f"{len(checks) - len(fails)}/{len(checks)} assertions passed")
if fails:
    print("\nFAILURES:")
    for f in fails:
        print(f"  - {f}")
    sys.exit(1)
print("\nRESULT: PASS — plan is coherent and provably non-executable (fail-closed).")
sys.exit(0)

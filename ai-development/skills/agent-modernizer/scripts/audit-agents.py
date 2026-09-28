#!/usr/bin/env python3
"""Audit Claude Code agent definitions.

Deterministic checks (field presence, name format, model alias, color, example
structure and placement, body length, bullet counts) run in code. Semantic
judgments (topic lists, phantom references, guidance density, ...) are defined
in references/question-catalog.json and sent to TypeSafe Jev when
TYPESAFE_API_KEY is set. Severity and action are composed in code from both.

Usage:
  audit-agents.py PATH [PATH ...]        audit files and/or agent directories
  audit-agents.py --route DIR            check that each agent's own examples
                                         route to it (needs the key)
  audit-agents.py --compare OLD NEW      rewrite regression check (needs key)
  audit-agents.py --no-jev PATH          deterministic checks only
  audit-agents.py --json PATH            emit machine-readable results

Exit 0 when nothing is Must fix, 1 otherwise, 2 when a needed key is missing.
"""
import argparse, glob, json, os, re, sys, time, urllib.error, urllib.request

API = "https://api.typesafe.ai/v1/systemone"
HERE = os.path.dirname(os.path.abspath(__file__))
CATALOG = os.path.join(HERE, "..", "references", "question-catalog.json")

VALID_MODELS = {"inherit", "sonnet", "opus", "haiku"}
VALID_COLORS = {"blue", "cyan", "green", "yellow", "magenta", "red"}
NAME_RE = re.compile(r"^[a-z0-9][a-z0-9-]{1,48}[a-z0-9]$")
RAW_MODEL_RE = re.compile(r"claude-|-\d{8}|\d+\.\d+")

MUST, SHOULD, CONSIDER = "Must fix", "Should fix", "Consider"

# Deterministic limits. Body limits exclude example blocks.
BODY_IDEAL_MIN, BODY_IDEAL_MAX, BODY_HARD_MAX = 500, 3000, 5000
DESC_MIN, DESC_MAX, DESC_WORDS_MAX = 10, 5000, 120
BULLETS_MAX = 20
EXAMPLES_MIN, EXAMPLES_MAX = 2, 5


# --------------------------------------------------------------------------
# Parsing
# --------------------------------------------------------------------------

def parse_frontmatter(text):
    """Return (fields, body, fm_end). Handles scalars, `>`/`|` blocks, inline
    arrays, and `- item` lists. Values are strings or lists of strings."""
    lines = text.split("\n")
    if not lines or lines[0].strip() != "---":
        return None, text, 0
    fields, i = {}, 1
    while i < len(lines) and lines[i].strip() != "---":
        line = lines[i]
        m = re.match(r"^([A-Za-z_][\w-]*):\s*(.*)$", line)
        if not m:
            i += 1
            continue
        key, val = m.group(1), m.group(2).strip()
        i += 1
        if val in (">", "|", ">-", "|-"):
            block = []
            while i < len(lines) and (lines[i].startswith(" ") or lines[i].strip() == ""):
                if lines[i].strip() == "---":
                    break
                block.append(lines[i].strip())
                i += 1
            joiner = " " if val.startswith(">") else "\n"
            fields[key] = joiner.join(b for b in block if b or joiner == "\n").strip()
        elif val == "":
            items = []
            while i < len(lines) and re.match(r"^\s+-\s+", lines[i]):
                items.append(re.sub(r"^\s+-\s+", "", lines[i]).strip().strip('"\''))
                i += 1
            fields[key] = items if items else ""
        elif val.startswith("["):
            fields[key] = [t.strip().strip('"\'') for t in val.strip("[]").split(",") if t.strip()]
        else:
            fields[key] = val.strip('"\'')
    fm_end = i + 1 if i < len(lines) else len(lines)
    body = "\n".join(lines[fm_end:])
    return fields, body, sum(len(l) + 1 for l in lines[:fm_end])


EXAMPLE_RE = re.compile(r"<example>(.*?)</example>", re.S)


def parse_examples(text, fm_end_offset):
    out = []
    for m in EXAMPLE_RE.finditer(text):
        block = m.group(1)
        def grab(label):
            mm = re.search(rf"^\s*{label}\s*:?\s*(.*)$", block, re.M)
            return mm.group(1).strip().strip('"') if mm else ""
        comm = re.search(r"<commentary>(.*?)</commentary>", block, re.S)
        out.append({
            "location": "description" if m.start() < fm_end_offset else "body",
            "context": grab("Context"),
            "user": grab("user"),
            "assistant": grab("assistant"),
            "commentary": comm.group(1).strip() if comm else "",
        })
    return out


def load_agent(path):
    text = open(path, encoding="utf-8").read()
    fields, body, fm_end = parse_frontmatter(text)
    examples = parse_examples(text, fm_end)
    body_wo_examples = EXAMPLE_RE.sub("", body).strip()
    return {
        "path": path,
        "file": os.path.basename(path),
        "plugin_dir": os.path.dirname(os.path.abspath(path)),
        "fields": fields,
        "body": body_wo_examples,
        "examples": examples,
        "raw": text,
    }


def known_agents(agent):
    """Agent names the body may legitimately reference: every agent in the
    marketplace when the plugin sits in one, otherwise the plugin's own."""
    plugin_root = os.path.dirname(agent["plugin_dir"])
    market = os.path.dirname(plugin_root)
    if os.path.exists(os.path.join(market, ".claude-plugin", "marketplace.json")):
        paths = glob.glob(os.path.join(market, "*", "agents", "*.md"))
    else:
        paths = glob.glob(os.path.join(agent["plugin_dir"], "*.md"))
    return sorted({os.path.splitext(os.path.basename(p))[0] for p in paths})


# --------------------------------------------------------------------------
# Deterministic checks
# --------------------------------------------------------------------------

def deterministic_findings(agent, roster_colors):
    f = agent["fields"]
    findings = []
    add = lambda area, issue, sev: findings.append({"area": area, "issue": issue, "severity": sev})
    if f is None:
        add("frontmatter", "No YAML frontmatter block", MUST)
        return findings, {}

    for key in ("name", "description", "model", "color"):
        if key not in f or f[key] in ("", []):
            add(key, "Missing required field", MUST)

    name = f.get("name", "")
    if name and not NAME_RE.match(name):
        add("name", f"`{name}` is not lowercase-hyphenated, 3-50 chars, alphanumeric at both ends", MUST)
    expected = os.path.splitext(agent["file"])[0]
    if name and name != expected:
        add("name", f"`{name}` does not match filename `{expected}`", CONSIDER)

    model = f.get("model", "")
    if model and model not in VALID_MODELS:
        if RAW_MODEL_RE.search(model):
            add("model", f"Raw model ID `{model}`; use one of {sorted(VALID_MODELS)}", MUST)
        else:
            add("model", f"Unknown model alias `{model}`", MUST)

    color = f.get("color", "")
    if color and color not in VALID_COLORS:
        add("color", f"`{color}` is not one of {sorted(VALID_COLORS)}", MUST)
    if color and roster_colors.get(color, 0) > 1:
        add("color", f"`{color}` is shared by {roster_colors[color]} agents in this plugin", CONSIDER)

    tools = f.get("tools")
    if tools is None:
        add("tools", "No `tools` array; agent gets every tool", CONSIDER)
    elif not isinstance(tools, list):
        add("tools", "`tools` is not an array", MUST)

    desc = f.get("description", "")
    if isinstance(desc, str) and desc:
        if "\\n" in desc:
            add("description", "Contains literal `\\n`; use a folded `>` block with real newlines", MUST)
        if len(desc) < DESC_MIN:
            add("description", f"Only {len(desc)} chars", MUST)
        if len(desc) > DESC_MAX:
            add("description", f"{len(desc)} chars exceeds {DESC_MAX}", MUST)
        words = len(desc.split())
        if words > DESC_WORDS_MAX:
            add("description", f"{words} words; every agent description is loaded for routing on every turn, keep it under {DESC_WORDS_MAX}", SHOULD)

    ex = agent["examples"]
    in_desc = [e for e in ex if e["location"] == "description"]
    if not ex:
        add("examples", "No `<example>` blocks anywhere; the agent will trigger unreliably", MUST)
    else:
        if in_desc:
            add("examples", f"{len(in_desc)} example block(s) inside `description`; move them to the body, after the frontmatter and before the role statement, so they do not cost routing tokens", SHOULD)
        if len(ex) < EXAMPLES_MIN:
            add("examples", f"Only {len(ex)} example block; want {EXAMPLES_MIN}-{EXAMPLES_MAX}", SHOULD)
        elif len(ex) > EXAMPLES_MAX:
            add("examples", f"{len(ex)} example blocks; want {EXAMPLES_MIN}-{EXAMPLES_MAX}", CONSIDER)
        for n, e in enumerate(ex, 1):
            missing = [k for k in ("context", "user", "assistant", "commentary") if not e[k]]
            if missing:
                add("examples", f"Example {n} lacks {', '.join(missing)}", MUST)

    body = agent["body"]
    nchars = len(body)
    bullets = len(re.findall(r"^\s*[-*+]\s+", body, re.M))
    numbered = len(re.findall(r"^\s*\d+\.\s+", body, re.M))
    if nchars == 0:
        add("body", "Empty system prompt", MUST)
    elif nchars > BODY_HARD_MAX:
        add("body", f"{nchars} chars (limit {BODY_HARD_MAX}); rewrite or split", SHOULD)
    elif nchars > BODY_IDEAL_MAX:
        add("body", f"{nchars} chars; target under {BODY_IDEAL_MAX}", SHOULD)
    elif nchars < BODY_IDEAL_MIN:
        add("body", f"{nchars} chars; check that role, process, and boundaries are all present", CONSIDER)
    if bullets > BULLETS_MAX:
        add("body", f"{bullets} bullet points (limit {BULLETS_MAX}); likely topic inventories", SHOULD)
    if re.search(r"\b(I am|I will|I would|I have)\b", body):
        add("body", "Uses first person; write the prompt in second person ('You are...')", CONSIDER)
    heading_nums = re.findall(r"^#+\s*(\d+)[.)]\s", body, re.M)
    if len(heading_nums) != len(set(heading_nums)):
        add("body", "Duplicate section numbers in headings", CONSIDER)

    metrics = {"body_chars": nchars, "bullets": bullets, "numbered_steps": numbered,
               "description_words": len(desc.split()) if isinstance(desc, str) else 0,
               "examples": len(ex), "examples_in_description": len(in_desc)}
    return findings, metrics


# --------------------------------------------------------------------------
# Jev
# --------------------------------------------------------------------------

def require_key():
    key = os.environ.get("TYPESAFE_API_KEY")
    if not key:
        print("SKIP: TYPESAFE_API_KEY is not set", file=sys.stderr)
        sys.exit(2)
    return key


def jev(state, questions, model):
    key = require_key()
    body = json.dumps({"state": state, "model": model, "questions": questions}).encode()
    for attempt in range(8):
        req = urllib.request.Request(API, data=body, headers={
            "Authorization": f"Bearer {key}", "Content-Type": "application/json"})
        try:
            with urllib.request.urlopen(req, timeout=120) as r:
                return json.loads(r.read())
        except urllib.error.HTTPError as e:
            if e.code in (429, 500, 502, 503, 504, 520, 521, 522, 523, 524, 529):
                time.sleep(2 ** attempt)
                continue
            sys.exit(f"HTTP {e.code}: {e.read().decode()[:500]}")
    sys.exit("retries exhausted")


def load_catalog():
    return json.load(open(CATALOG, encoding="utf-8"))


def semantic_findings(agent, catalog):
    f = agent["fields"] or {}
    state = {
        "agent": {
            "name": f.get("name", ""),
            "model": f.get("model", ""),
            "description": f.get("description", ""),
            "examples": agent["examples"],
            "body": agent["body"],
        },
        "known_agents": known_agents(agent),
    }
    answers = jev(state, catalog["questions"], catalog["model"])["answers"]
    th = catalog["thresholds"]
    findings, judgments = [], {}
    labels = {
        "topic_lists": "Topic lists without guidance",
        "teaches_known": "Teaches the model its own knowledge",
        "fictional_content": "Fictional metrics or progress tracking",
        "phantom_refs": "References agents or protocols that do not exist",
        "duplicate_coverage": "Same subject covered under two sections",
        "mirrored_do_not": "'Do Not' list only inverts the positive instructions",
        "overspecified_output": "Prescribes an unimplemented protocol or schema",
        "distinct_domains": "Covers two or more unrelated domains",
    }
    positives = {
        "has_role_statement": "No one-sentence role statement at the top",
        "has_boundaries": "No explicit boundaries or 'Do Not' section",
        "has_process": "No numbered action steps",
        "has_output_format": "No defined output format",
    }
    for qid, a in answers.items():
        if a["type"] == "noul":
            p = a["noul"]
            judgments[qid] = round(p, 2)
            if qid in labels:
                if p >= th["noul_flag"]:
                    findings.append({"area": "body", "issue": f"{labels[qid]} (p={p:.2f})", "severity": SHOULD})
                elif p >= th["noul_review"]:
                    findings.append({"area": "body", "issue": f"Possible: {labels[qid].lower()} (p={p:.2f}); confirm by reading", "severity": CONSIDER})
            elif qid in positives:
                if p < 1 - th["noul_flag"]:
                    findings.append({"area": "body", "issue": f"{positives[qid]} (p={p:.2f})", "severity": SHOULD})
                elif p < 1 - th["noul_review"]:
                    findings.append({"area": "body", "issue": f"Possible: {positives[qid].lower()} (p={p:.2f}); confirm by reading", "severity": CONSIDER})
        elif a["type"] == "score":
            judgments[qid] = {"score": round(a["score"], 2), "confidence": round(a["confidence"], 2)}
            if qid == "description_triggering" and a["score"] < 1.5:
                findings.append({"area": "description", "issue": f"Weak triggering: reads as a summary rather than selection criteria (score {a['score']:.1f}/3)", "severity": SHOULD})
            elif qid == "description_triggering" and a["score"] < 2.5:
                findings.append({"area": "description", "issue": f"No 'do not use for' handoff to sibling agents (score {a['score']:.1f}/3)", "severity": CONSIDER})
            if qid == "guidance_density" and a["score"] < 1.5:
                findings.append({"area": "body", "issue": f"Mostly topic inventories rather than decisions (score {a['score']:.1f}/3)", "severity": SHOULD})
    return findings, judgments


# --------------------------------------------------------------------------
# Composition
# --------------------------------------------------------------------------

def decide_action(findings, metrics, judgments):
    must = [x for x in findings if x["severity"] == MUST]
    anti = sum(1 for k in ("topic_lists", "teaches_known", "fictional_content", "phantom_refs",
                            "duplicate_coverage", "mirrored_do_not", "overspecified_output")
               if isinstance(judgments.get(k), (int, float)) and judgments[k] >= 0.7)
    density = judgments.get("guidance_density", {}).get("score") if judgments else None
    split = isinstance(judgments.get("distinct_domains"), (int, float)) and judgments["distinct_domains"] >= 0.7
    chars, bullets = metrics.get("body_chars", 0), metrics.get("bullets", 0)

    parts = []
    if must:
        parts.append("Fix frontmatter" if all(x["area"] != "body" for x in must) else "Fix structure")
    if chars > BODY_HARD_MAX and split:
        parts.append("Split")
    elif anti >= 2 or (density is not None and density < 1.5) or bullets > BULLETS_MAX or chars > BODY_HARD_MAX:
        parts.append("Rewrite")
    elif anti == 1 or chars > BODY_IDEAL_MAX:
        parts.append("Trim")
    if metrics.get("examples_in_description"):
        parts.append("Move examples to body")
    return "; ".join(parts) if parts else "Keep"


def audit(paths, use_jev, catalog):
    files = []
    for p in paths:
        if os.path.isdir(p):
            files.extend(sorted(glob.glob(os.path.join(p, "*.md"))))
        else:
            files.append(p)
    agents = [load_agent(p) for p in files]
    colors_by_dir = {}
    for a in agents:
        for sib in glob.glob(os.path.join(a["plugin_dir"], "*.md")):
            fields, _, _ = parse_frontmatter(open(sib, encoding="utf-8").read())
            c = (fields or {}).get("color", "")
            colors_by_dir.setdefault(a["plugin_dir"], {}).setdefault(c, set()).add(sib)
    results = []
    for a in agents:
        roster_colors = {c: len(s) for c, s in colors_by_dir[a["plugin_dir"]].items()}
        findings, metrics = deterministic_findings(a, roster_colors)
        judgments = {}
        if use_jev and a["fields"] is not None:
            sem, judgments = semantic_findings(a, catalog)
            findings.extend(sem)
        order = {MUST: 0, SHOULD: 1, CONSIDER: 2}
        findings.sort(key=lambda x: order[x["severity"]])
        results.append({"agent": a["fields"].get("name", a["file"]) if a["fields"] else a["file"],
                        "path": a["path"], "findings": findings, "metrics": metrics,
                        "judgments": judgments, "action": decide_action(findings, metrics, judgments)})
    return results


# --------------------------------------------------------------------------
# Routing check
# --------------------------------------------------------------------------

def route_check(directory, catalog):
    agents = [load_agent(p) for p in sorted(glob.glob(os.path.join(directory, "*.md")))]
    names = [a["fields"]["name"] for a in agents if a["fields"]]
    state = {"agents": {a["fields"]["name"]: a["fields"].get("description", "") for a in agents if a["fields"]}}
    probes = []  # (owner, request, proactive)
    for a in agents:
        if not a["fields"]:
            continue
        for e in a["examples"]:
            if e["user"]:
                proactive = "proactive" in (e["context"] + " " + e["commentary"]).lower()
                probes.append((a["fields"]["name"], e["user"], proactive))
    state["requests"] = {f"R{i}": p[1] for i, p in enumerate(probes)}
    questions = {
        f"route_R{i}": {
            "type": "choice",
            "instructions": {"request_id": f"R{i}",
                             "question": "Which agent in `agents` should handle `requests.<request_id>`, judging only from the agent descriptions?"},
            "criteria": {**{n: None for n in names},
                         "none": "No listed agent fits; a different plugin or the general assistant should handle it"},
        } for i in range(len(probes))
    }
    if not probes:
        print("no example `user:` lines found"); return 0
    answers = jev(state, questions, catalog["model"])["answers"]
    th = catalog["thresholds"]["routing_confidence"]
    failures, counted = 0, 0
    print(f"{'result':6} {'winner':22} {'conf':5} {'owner':22} example request")
    for i, (owner, req, proactive) in enumerate(probes):
        a = answers[f"route_R{i}"]
        ok = a["choice"] == owner and a["confidence"] >= th
        if proactive:
            status = "info"  # a proactive example describes a situation, not a request that names the task
        else:
            status = "ok" if ok else "FAIL"
            counted += 1
            failures += 0 if ok else 1
        print(f"{status:6} {a['choice']:22} {a['confidence']:.2f}  {owner:22} {req[:70]}")
    print(f"\n{failures} of {counted} non-proactive example requests do not route to their own agent "
          f"at confidence >= {th}; rows marked info are proactive examples and are not counted")
    return 1 if failures else 0


def compare(old_path, new_path, catalog):
    state = {"original": open(old_path, encoding="utf-8").read(),
             "rewritten": open(new_path, encoding="utf-8").read()}
    qs = catalog["rewrite_check"]["questions"]
    answers = jev(state, qs, catalog["model"])["answers"]
    th = catalog["thresholds"]
    bad = 0
    for qid, a in answers.items():
        p = a["noul"]
        flag = "FAIL" if p >= th["noul_flag"] else ("check" if p >= th["noul_review"] else "ok")
        bad += flag == "FAIL"
        print(f"{flag:5} {qid:26} p={p:.2f}")
    return 1 if bad else 0


# --------------------------------------------------------------------------
# Output
# --------------------------------------------------------------------------

def print_markdown(results, use_jev):
    for r in results:
        print(f"\n### `{r['agent']}` — {r['path']}\n")
        m = r["metrics"]
        if m:
            print(f"Body {m['body_chars']} chars, {m['bullets']} bullets, {m['numbered_steps']} numbered steps; "
                  f"description {m['description_words']} words; {m['examples']} examples"
                  f"{' (' + str(m['examples_in_description']) + ' in description)' if m['examples_in_description'] else ''}.\n")
        if r["findings"]:
            print("| # | Area | Issue | Severity |\n|---|------|-------|----------|")
            for i, x in enumerate(r["findings"], 1):
                print(f"| {i} | `{x['area']}` | {x['issue']} | {x['severity']} |")
        else:
            print("No findings.")
        print(f"\n**Action:** {r['action']}")
    if len(results) > 1:
        print("\n### Summary\n")
        print("| Agent | Body chars | Must fix | Should fix | Consider | Action |")
        print("|-------|-----------:|---------:|-----------:|---------:|--------|")
        for r in results:
            c = {s: sum(1 for x in r["findings"] if x["severity"] == s) for s in (MUST, SHOULD, CONSIDER)}
            print(f"| `{r['agent']}` | {r['metrics'].get('body_chars', 0)} | {c[MUST]} | {c[SHOULD]} | {c[CONSIDER]} | {r['action']} |")
    if not use_jev:
        print("\nSemantic judgments skipped (no key or --no-jev). Answer references/question-catalog.json "
              "manually for each agent: one question at a time, cite the deciding lines.")


def main():
    ap = argparse.ArgumentParser(description=__doc__, formatter_class=argparse.RawDescriptionHelpFormatter)
    ap.add_argument("paths", nargs="*")
    ap.add_argument("--route", metavar="DIR")
    ap.add_argument("--compare", nargs=2, metavar=("OLD", "NEW"))
    ap.add_argument("--no-jev", action="store_true")
    ap.add_argument("--json", action="store_true")
    args = ap.parse_args()
    catalog = load_catalog()
    if args.route:
        sys.exit(route_check(args.route, catalog))
    if args.compare:
        sys.exit(compare(args.compare[0], args.compare[1], catalog))
    if not args.paths:
        ap.error("give at least one agent file or directory")
    use_jev = not args.no_jev and bool(os.environ.get("TYPESAFE_API_KEY"))
    if not args.no_jev and not use_jev:
        print("note: TYPESAFE_API_KEY not set; running deterministic checks only", file=sys.stderr)
    results = audit(args.paths, use_jev, catalog)
    if args.json:
        print(json.dumps(results, indent=2))
    else:
        print_markdown(results, use_jev)
    sys.exit(1 if any(x["severity"] == MUST for r in results for x in r["findings"]) else 0)


if __name__ == "__main__":
    main()

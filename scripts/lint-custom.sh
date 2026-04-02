#!/usr/bin/env bash
set -euo pipefail

python3 - <<'PY'
import os
import re
import sys

root = os.getcwd()
func_re = re.compile(r'^\s*func\s*(?:\([^)]*\)\s*)?([A-Za-z_][A-Za-z0-9_]*)\s*\(')

violations = []

for dirpath, dirnames, filenames in os.walk(root):
    dirnames[:] = [d for d in dirnames if d not in {".git", "vendor"}]
    for filename in filenames:
        if not filename.endswith(".go"):
            continue

        path = os.path.join(dirpath, filename)
        rel = os.path.relpath(path, root)

        saw_private = False
        with open(path, "r", encoding="utf-8") as f:
            for idx, line in enumerate(f, start=1):
                m = func_re.match(line)
                if not m:
                    continue
                name = m.group(1)
                if not name:
                    continue
                is_exported = name[0].isupper()
                if is_exported and saw_private:
                    violations.append(f"{rel}:{idx} exported function '{name}' appears after private functions")
                if not is_exported:
                    saw_private = True

if violations:
    print("Function ordering lint failed (public functions must appear before private ones):")
    for v in violations:
        print(f" - {v}")
    sys.exit(1)

print("Function ordering lint passed.")
PY

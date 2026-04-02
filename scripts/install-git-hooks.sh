#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="$(git rev-parse --show-toplevel)"
mkdir -p "$REPO_ROOT/.git/hooks"

cat > "$REPO_ROOT/.git/hooks/pre-commit" <<'EOF'
#!/usr/bin/env bash
set -euo pipefail
REPO_ROOT="$(git rev-parse --show-toplevel)"
"$REPO_ROOT/scripts/hooks/pre-commit"
EOF

chmod +x "$REPO_ROOT/.git/hooks/pre-commit"
echo "Installed pre-commit hook."

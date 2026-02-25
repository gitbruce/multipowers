#!/usr/bin/env bash
set -euo pipefail

conductor_root_dir() {
  echo "${PROJECT_ROOT:-$PWD}/conductor"
}

is_spec_driven_command() {
  local cmd="${1:-}"
  case "$cmd" in
    discover|research|probe|define|grasp|develop|tangle|deliver|ink|embrace|review|debate|plan)
      return 0 ;;
    *)
      return 1 ;;
  esac
}

conductor_context_complete() {
  local croot
  croot="$(conductor_root_dir)"
  [[ -f "$croot/product.md" ]] || return 1
  [[ -f "$croot/product-guidelines.md" ]] || return 1
  [[ -f "$croot/tech-stack.md" ]] || return 1
  [[ -f "$croot/workflow.md" ]] || return 1
  [[ -f "$croot/tracks.md" ]] || return 1
  [[ -d "$croot/code_styleguides" ]] || return 1
  return 0
}

ensure_conductor_context() {
  local cmd="${1:-}"
  is_spec_driven_command "$cmd" || return 0

  if conductor_context_complete; then
    return 0
  fi

  if type run_octo_init_interactive &>/dev/null; then
    run_octo_init_interactive
    conductor_context_complete
    return $?
  fi

  return 1
}

slugify_for_track() {
  local text="${1:-task}"
  text="$(echo "$text" | tr '[:upper:]' '[:lower:]' | sed 's/[^a-z0-9]/-/g' | sed 's/-\{2,\}/-/g' | sed 's/^-//; s/-$//')"
  echo "${text:0:48}"
}

ensure_conductor_track_file() {
  local cmd="$1"
  local prompt="${2:-}"
  local croot track_dir short base_track_id track_id track_path file n
  croot="$(conductor_root_dir)"
  track_dir="$croot/tracks"
  mkdir -p "$track_dir"

  short="$(slugify_for_track "${cmd}-${prompt:-run}")"
  short="${short:0:24}"
  [[ -z "$short" ]] && short="track"
  base_track_id="${short}_$(date +%Y%m%d)"
  track_id="$base_track_id"
  track_path="$track_dir/$track_id"
  n=2
  while [[ -e "$track_path" ]]; do
    track_id="${base_track_id}_${n}"
    track_path="$track_dir/$track_id"
    ((n++))
  done
  mkdir -p "$track_path"
  file="$track_path/tracking.md"

  if [[ ! -f "$file" ]]; then
    cat > "$file" <<EOF
# Track: ${track_id}

- command: ${cmd}
- prompt: ${prompt}

- [x] Context check started
- [ ] Context check passed
- [ ] Command execution started
- [ ] Command execution finished
- [ ] Validation complete

## Notes
- created: $(date -Iseconds)
EOF
  fi

  echo "$file"
}

mark_track_checkbox() {
  local file="$1"
  local label="$2"
  [[ -f "$file" ]] || return 0
  sed -i "s/^- \[ \] ${label}$/- [x] ${label}/" "$file" 2>/dev/null || true
}

load_conductor_context_for_prompt() {
  local croot
  croot="$(conductor_root_dir)"
  cat <<EOF
<project_context source=\"conductor\">
$(cat "$croot/product.md" 2>/dev/null)

$(cat "$croot/product-guidelines.md" 2>/dev/null)

$(cat "$croot/tech-stack.md" 2>/dev/null)

$(cat "$croot/workflow.md" 2>/dev/null)
</project_context>
EOF
}

apply_conductor_context_to_prompt() {
  local prompt="$1"
  if conductor_context_complete; then
    cat <<EOF
Use the following project context as source-of-truth before solving the task.

$(load_conductor_context_for_prompt)

Task:
$prompt
EOF
  else
    echo "$prompt"
  fi
}

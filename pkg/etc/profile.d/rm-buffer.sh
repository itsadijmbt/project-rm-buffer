#!/usr/bin/env bash
# rm-buffer: interactive helper (defines rm() only for interactive shells)
case "$-" in *i*) ;; *) return;; esac

# Don't override user-defined rm alias/function
if alias rm >/dev/null 2>&1 || declare -F rm >/dev/null 2>&1; then
  return
fi

rm() {
  if [ "$#" -eq 0 ]; then
    command rm
    return
  fi
  case "$1" in
    -buffer|-b|-extract|-E|-list|-L)
      MODE="$1"
      shift
      /usr/local/bin/rm-buffer "$(pwd -P)" "$MODE" "$@"
      ;;
    *)
      command rm "$@"
      ;;
  esac
}
export -f rm 2>/dev/null || true

#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)
REPO_ROOT=$(cd -- "$SCRIPT_DIR/../.." && pwd)
ABIGEN_VERSION=${ABIGEN_VERSION:-v1.15.11}

need_cmd() {
	if ! command -v "$1" >/dev/null 2>&1; then
		echo "missing required command: $1" >&2
		exit 1
	fi
}

need_cmd go

if [[ -n "${ABIGEN:-}" ]]; then
	ABIGEN_CMD=("$ABIGEN")
else
	ABIGEN_CMD=(go run "github.com/ethereum/go-ethereum/cmd/abigen@${ABIGEN_VERSION}")
fi

TMP_DIR=$(mktemp -d "${TMPDIR:-/tmp}/agent-pay-legacy-bindings.XXXXXX")
trap 'rm -rf "$TMP_DIR"' EXIT

EXTRACTOR_SRC="$TMP_DIR/extract.go"
EXTRACTOR_BIN="$TMP_DIR/extract"

cat > "$EXTRACTOR_SRC" <<'EOF'
package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "usage: extract <file> <symbol>")
		os.Exit(2)
	}

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	symbol := regexp.QuoteMeta(os.Args[2])
	directRe := regexp.MustCompile(`(?s)(?:const|var)\s+` + symbol + `\s*=\s*("(?:\\.|[^"\\])*")`)
	match := directRe.FindSubmatch(data)
	if match == nil {
		field := ""
		metaName := ""
		rawSymbol := os.Args[2]
		switch {
		case strings.HasSuffix(rawSymbol, "ABI"):
			field = "ABI"
			metaName = strings.TrimSuffix(rawSymbol, "ABI") + "MetaData"
		case strings.HasSuffix(rawSymbol, "Bin"):
			field = "Bin"
			metaName = strings.TrimSuffix(rawSymbol, "Bin") + "MetaData"
		}
		if field != "" {
			metaRe := regexp.MustCompile(`(?s)var\s+` + regexp.QuoteMeta(metaName) + `\s*=\s*&bind\.MetaData\s*\{.*?\b` + field + `:\s*("(?:\\.|[^"\\])*")`)
			match = metaRe.FindSubmatch(data)
		}
	}
	if match == nil {
		fmt.Fprintf(os.Stderr, "symbol not found: %s\n", os.Args[2])
		os.Exit(1)
	}

	value, err := strconv.Unquote(string(match[1]))
	if err != nil {
		panic(err)
	}

	fmt.Print(value)
}
EOF

go build -o "$EXTRACTOR_BIN" "$EXTRACTOR_SRC"

extract_literal() {
	"$EXTRACTOR_BIN" "$1" "$2"
}

generate_from_abi() {
	local input_rel=$1
	local output_rel=$2
	local pkg=$3
	local type_name=$4
	local input="$REPO_ROOT/$input_rel"
	local output="$REPO_ROOT/$output_rel"
	local abi_json="$TMP_DIR/${pkg}_${type_name}.abi.json"

	extract_literal "$input" "${type_name}ABI" > "$abi_json"
	mkdir -p "$(dirname -- "$output")"
	"${ABIGEN_CMD[@]}" --abi "$abi_json" --pkg "$pkg" --type "$type_name" --out "$output"
}

generate_from_abi_bin() {
	local input_rel=$1
	local output_rel=$2
	local pkg=$3
	local type_name=$4
	local input="$REPO_ROOT/$input_rel"
	local output="$REPO_ROOT/$output_rel"
	local abi_json="$TMP_DIR/${pkg}_${type_name}.abi.json"
	local bytecode_txt="$TMP_DIR/${pkg}_${type_name}.bin"

	extract_literal "$input" "${type_name}ABI" > "$abi_json"
	extract_literal "$input" "${type_name}Bin" > "$bytecode_txt"
	mkdir -p "$(dirname -- "$output")"
	"${ABIGEN_CMD[@]}" \
		--abi "$abi_json" \
		--bin "$bytecode_txt" \
		--pkg "$pkg" \
		--type "$type_name" \
		--out "$output"
}

generate_from_abi app/booleanoutcome.go app/booleanoutcome.go app IBooleanOutcome
generate_from_abi app/numericoutcome.go app/numericoutcome.go app INumericOutcome
generate_from_abi app/multisession.go app/multisession.go app IMultiSession
generate_from_abi app/multisessionwithoracle.go app/multisessionwithoracle.go app IMultiSessionWithOracle
generate_from_abi app/singlesession.go app/singlesession.go app ISingleSession
generate_from_abi app/singlesessionwithoracle.go app/singlesessionwithoracle.go app ISingleSessionWithOracle

generate_from_abi_bin testing/testapp/multigomoku.go testing/testapp/multigomoku.go testapp MultiGomoku
generate_from_abi_bin testing/testapp/multisessionapp.go testing/testapp/multisessionapp.go testapp SimpleMultiSessionApp
generate_from_abi_bin testing/testapp/multisessionappwithoracle.go testing/testapp/multisessionappwithoracle.go testapp SimpleMultiSessionAppWithOracle
generate_from_abi_bin testing/testapp/singlesessionapp.go testing/testapp/singlesessionapp.go testapp SimpleSingleSessionApp
generate_from_abi_bin testing/testapp/singlesessionappwithoracle.go testing/testapp/singlesessionappwithoracle.go testapp SimpleSingleSessionAppWithOracle

echo "generated legacy app bindings under $REPO_ROOT"
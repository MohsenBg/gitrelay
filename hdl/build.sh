#!/usr/bin/env bash
#
# bgscan build script
# ----------------------------
# This script builds bgscan for the local platform or for all supported targets.
# It verifies Go availability, checks platform-specific assets, injects version
# information into the binary, and copies all required resources into the final
# build directory.
#
# Author: MohsenBg
# Project: bgscan
#

set -e

# ---------------------------------------------------------------------------
# Supported build targets (GOOS/GOARCH)
# ---------------------------------------------------------------------------
targets=(
	"linux/amd64"
	"linux/arm64"
	"windows/amd64"
	"darwin/amd64"
	"darwin/arm64"
	"android/arm64"
)

APP="hdl"

IPS_DIR="ips"

MAIN_FILE="./cmd/hdl/main.go"

# ---------------------------------------------------------------------------
# Auto-detect project version from Git
# ---------------------------------------------------------------------------
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")

# ---------------------------------------------------------------------------
# Parse flags
# ---------------------------------------------------------------------------

all=false

while [[ $# -gt 0 ]]; do
	case "$1" in
	--all)
		all=true
		shift
		;;
	--version)
		VERSION="$2"
		shift 2
		;;
	*)
		echo "Unknown argument: $1"
		exit 1
		;;
	esac
done

echo "Detected version: $VERSION"
echo

# ---------------------------------------------------------------------------
# Detect local OS
# ---------------------------------------------------------------------------
case "$(uname -s)" in
Linux) OS="linux" ;;
Darwin) OS="darwin" ;;
MINGW* | MSYS* | CYGWIN*) OS="windows" ;;
*)
	echo "Unsupported OS"
	exit 1
	;;
esac

# ---------------------------------------------------------------------------
# Detect local architecture
# ---------------------------------------------------------------------------
case "$(uname -m)" in
x86_64) ARCH="amd64" ;;
aarch64 | arm64) ARCH="arm64" ;;
*)
	echo "Unsupported architecture"
	exit 1
	;;
esac

USER_TARGET="$OS/$ARCH"

# ---------------------------------------------------------------------------
# Intro message
# ---------------------------------------------------------------------------
echo "Thank you for using hdl."
echo "Building this project helps support free and open access to the internet."
echo "Your contribution is appreciated."
echo

# ---------------------------------------------------------------------------
# Check for Go installation
# ---------------------------------------------------------------------------
echo "[CHECK] Go installation..."

if ! command -v go >/dev/null 2>&1; then
	echo "[ERROR] Go is not installed. Please install Golang first."
	exit 1
fi

echo "[OK] Go is installed."
echo

# ---------------------------------------------------------------------------
# Ensure dist directory exists
# ---------------------------------------------------------------------------
DIST="dist/$VERSION"
mkdir -p "$DIST"

# ---------------------------------------------------------------------------
# Build a single target
# ---------------------------------------------------------------------------
build_target() {
	local target="$1"

	IFS=/ read -r GOOS GOARCH <<<"$target"

	local FOLDER="$APP-$GOOS-$GOARCH"
	local OUT_DIR="$DIST/$FOLDER"

	mkdir -p "$OUT_DIR"

	local BIN_NAME="$FOLDER"
	[ "$GOOS" = "windows" ] && BIN_NAME="$FOLDER.exe"

	echo "[BUILD] $FOLDER"

	GOOS=$GOOS GOARCH=$GOARCH CGO_ENABLED=0 \
		go build \
		-trimpath \
		-ldflags="-s -w -X 'main.Version=$VERSION'" \
		-o "$OUT_DIR/$BIN_NAME" \
		"$MAIN_FILE"

	# Generate checksum
	(
		cd "$OUT_DIR" || exit
		sha256sum "$BIN_NAME" >"$FOLDER.sha256"
	)

	echo "[DONE] Built $FOLDER"
	echo
}

# ---------------------------------------------------------------------------
# Build logic (single or all)
# ---------------------------------------------------------------------------
if [ "$all" = true ]; then
	echo "[MODE] Building for ALL targets"
	echo
	for target in "${targets[@]}"; do
		build_target "$target"
	done
else
	echo "[MODE] Building only for detected target: $USER_TARGET"
	echo
	build_target "$USER_TARGET"
fi

echo "---------------------------------------------"
echo " Build completed successfully (version $VERSION)"
echo " Output directory: $DIST/"
echo "---------------------------------------------"

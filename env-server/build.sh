#!/usr/bin/env bash
set -euo pipefail

# Detect module path dynamically (from go.mod)
MODULE_PATH=$(go list -m)

# Get the current commit hash (short)
COMMIT=$(git rev-parse --short HEAD)

# Try to get the latest reachable tag
if TAG=$(git describe --tags --abbrev=0 2>/dev/null); then
    VERSION=$TAG
    if ! git describe --tags --exact-match >/dev/null 2>&1; then
        VERSION="${TAG}-${COMMIT}"
    fi
else
    VERSION="0.0.0-${COMMIT}"
fi

# Check for uncommitted changes
if ! git diff --quiet || ! git diff --cached --quiet; then
    VERSION="${VERSION}-dirty"
fi

# Build date in ISO-8601 UTC
DATE=$(date -u +%Y-%m-%dT%H:%M:%SZ)

# If we are exactly at a tag, set Gin to release mode
if [[ "$VERSION" == "$TAG" ]]; then
    export GIN_MODE=release
    echo ">> Detected release build (tag: $TAG)"
    echo ">> GIN_MODE=release"
fi

echo "Building with:"
echo "  Version: $VERSION"
echo "  Commit:  $COMMIT"
echo "  Date:    $DATE"

go build -ldflags "-X '${MODULE_PATH}/internal/version.Version=${VERSION}' \
                   -X '${MODULE_PATH}/internal/version.Commit=${COMMIT}' \
                   -X '${MODULE_PATH}/internal/version.Date=${DATE}'" \
    -o build/env-server ./cmd

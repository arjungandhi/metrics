deps:
    go mod tidy

test: deps
    go test -v ./...

install: deps
    go install ./cmd/health

# Build with Nix
nix-build:
    nix build

# Update the vendorHash in flake.nix after go.mod/go.sum changes
update-vendor-hash:
    #!/usr/bin/env bash
    set -euo pipefail
    sed -i 's|vendorHash = ".*";|vendorHash = "sha256-AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=";|' flake.nix
    nix_output=$(nix build .#default 2>&1 || true)
    hash=$(echo "$nix_output" | grep "got:" | awk '{print $2}')
    if [ -z "$hash" ]; then
        echo "ERROR: Could not determine vendorHash, restoring flake.nix"
        git checkout flake.nix
        exit 1
    fi
    sed -i "s|vendorHash = \".*\";|vendorHash = \"$hash\";|" flake.nix
    echo "Updated vendorHash to $hash"

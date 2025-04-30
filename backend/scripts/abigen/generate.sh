#!/bin/bash
# backend/scripts/abigen/generate.sh

# Exit immediately if a command exits with a non-zero status
set -e

# Define paths relative to the project root
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
ABI_DIR="$PROJECT_ROOT/pkg/blockchain/contracts"
OUTPUT_DIR="$PROJECT_ROOT/internal/contracts/generated"

# Create output directory if it doesn't exist
mkdir -p "$OUTPUT_DIR"

# Check if abigen is installed
if ! command -v abigen &> /dev/null; then
    echo "Error: abigen is not installed or not in your PATH"
    echo "Install it with: go install github.com/ethereum/go-ethereum/cmd/abigen@latest"
    exit 1
fi

# Function to generate bindings for a contract
generate_binding() {
    local abi_file="$1"
    local output_file="$2"
    local pkg_name="$3"
    local struct_name="$4"
    
    echo "Generating bindings for $abi_file..."
    abigen --abi="$abi_file" \
           --pkg="$pkg_name" \
           --out="$output_file" \
           --type="$struct_name"
    
    echo "✅ Generated $output_file"
}

# Generate bindings for each contract
echo "🔨 Generating Go bindings from Solidity contract ABIs..."

# Token contract
generate_binding "$ABI_DIR/Token.abi" \
                "$OUTPUT_DIR/token.go" \
                "generated" \
                "Token"

# LendingPool contract
generate_binding "$ABI_DIR/LendingPool.abi" \
                "$OUTPUT_DIR/lending_pool.go" \
                "generated" \
                "LendingPool"

# Borrowing contract
generate_binding "$ABI_DIR/Borrowing.abi" \
                "$OUTPUT_DIR/borrowing.go" \
                "generated" \
                "Borrowing"

# Collateral contract
generate_binding "$ABI_DIR/Collateral.abi" \
                "$OUTPUT_DIR/collateral.go" \
                "generated" \
                "Collateral"

echo "✨ All contract bindings generated successfully!"
echo "📁 Output directory: $OUTPUT_DIR"
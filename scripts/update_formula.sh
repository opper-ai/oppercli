#!/bin/bash
set -e

VERSION=$1
if [ -z "$VERSION" ]; then
    echo "Usage: $0 <version>"
    exit 1
fi

# Remove v prefix if present
VERSION=${VERSION#v}

# Read SHA256 hashes from the release files
DARWIN_ARM64_SHA=$(cat ./artifacts/opper-darwin-arm64.sha256 | cut -d ' ' -f 1)
DARWIN_AMD64_SHA=$(cat ./artifacts/opper-darwin-amd64.sha256 | cut -d ' ' -f 1)
LINUX_AMD64_SHA=$(cat ./artifacts/opper-linux-amd64.sha256 | cut -d ' ' -f 1)

# Create the formula file using printf to avoid heredoc issues
printf 'class Opper < Formula
  desc "Command line interface for Opper AI"
  homepage "https://github.com/opper-ai/oppercli"
  version "%s"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/opper-ai/oppercli/releases/download/v#{version}/opper-darwin-arm64"
      sha256 "%s"
    else
      url "https://github.com/opper-ai/oppercli/releases/download/v#{version}/opper-darwin-amd64"
      sha256 "%s"
    end
  end

  on_linux do
    if Hardware::CPU.intel?
      url "https://github.com/opper-ai/oppercli/releases/download/v#{version}/opper-linux-amd64"
      sha256 "%s"
    end
  end

  def install
    bin.install Dir["opper-*"].first => "opper"
  end

  test do
    system "#{bin}/opper", "--version"
  end
end
' "$VERSION" "$DARWIN_ARM64_SHA" "$DARWIN_AMD64_SHA" "$LINUX_AMD64_SHA" > HomebrewFormula/opper.rb 
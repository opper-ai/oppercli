class Opper < Formula
  desc "Command line interface for Opper AI"
  homepage "https://github.com/opper-ai/oppercli"
  version "0.11.0"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/opper-ai/oppercli/releases/download/v#{version}/opper-darwin-arm64"
      sha256 "bc4731af0e6c59cc81188e77fb1d624ccf9de3eb17b9de602f7818b4daf60f59"
    else
      url "https://github.com/opper-ai/oppercli/releases/download/v#{version}/opper-darwin-amd64"
      sha256 "6a128cad3db099c982d3a5073855b4b700be877585a64965d977ba4528e9cf48"
    end
  end

  on_linux do
    if Hardware::CPU.intel?
      url "https://github.com/opper-ai/oppercli/releases/download/v#{version}/opper-linux-amd64"
      sha256 "731c699c8e8f040aa9680ed09ec353bab7b33d5b083c588c04799aafae2dc3fe"
    end
  end

  def install
    bin.install Dir["opper-*"].first => "opper"
  end

  test do
    system "#{bin}/opper", "--version"
  end
end

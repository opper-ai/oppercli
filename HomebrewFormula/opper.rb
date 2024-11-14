class Opper < Formula
  desc "Command line interface for Opper AI"
  homepage "https://github.com/opper-ai/oppercli"
  version "0.5.0"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/opper-ai/oppercli/releases/download/v#{version}/opper-darwin-arm64"
      sha256 "746bdba29b72911530fd38a09734bdcc21e26f56711f8b675812d92b350ddefa"
    else
      url "https://github.com/opper-ai/oppercli/releases/download/v#{version}/opper-darwin-amd64"
      sha256 "2287bf0853bffad00286d198c375a07e0d14a47ec9a9b21e46f5e70998232df9"
    end
  end

  on_linux do
    if Hardware::CPU.intel?
      url "https://github.com/opper-ai/oppercli/releases/download/v#{version}/opper-linux-amd64"
      sha256 "6aa4d86557da9a1bb00d84701f9d1764d6ed404f5e48698cc4b5a2913e05103f"
    end
  end

  def install
    bin.install Dir["opper-*"].first => "opper"
  end

  test do
    system "#{bin}/opper", "--version"
  end
end

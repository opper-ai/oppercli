class Opper < Formula
  desc "Command line interface for Opper AI"
  homepage "https://github.com/opper-ai/oppercli"
  version "0.12.0"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/opper-ai/oppercli/releases/download/v#{version}/opper-darwin-arm64"
      sha256 "78ccdc97f75ea7fbc4b35f1e3a9014cab9f568ef9d61ba41c3908661641cdd63"
    else
      url "https://github.com/opper-ai/oppercli/releases/download/v#{version}/opper-darwin-amd64"
      sha256 "7925b73a7986310f7e006a66a97b4ab3ad426a9d621c30926502a89d54056409"
    end
  end

  on_linux do
    if Hardware::CPU.intel?
      url "https://github.com/opper-ai/oppercli/releases/download/v#{version}/opper-linux-amd64"
      sha256 "fb8e9552fd0902a43416d628ec0d3abf2ea757c5eaacaf30088d61a81ac96248"
    end
  end

  def install
    bin.install Dir["opper-*"].first => "opper"
  end

  test do
    system "#{bin}/opper", "--version"
  end
end

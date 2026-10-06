class Opper < Formula
  desc "Command line interface for Opper AI"
  homepage "https://github.com/opper-ai/oppercli"
  version "0.13.0"

  disable! date: "2026-10-06", because: "has been retired; install the maintained CLI with `npm install -g @opperai/cli`"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/opper-ai/oppercli/releases/download/v#{version}/opper-darwin-arm64"
      sha256 "0cf8aa579bbda8a2d958f7db05eba457e613820d1e441578e785cac5a36977e4"
    else
      url "https://github.com/opper-ai/oppercli/releases/download/v#{version}/opper-darwin-amd64"
      sha256 "1910bb78be35dfd749354b7f7efb6169364a4d2743c15b095a45100af5bb878b"
    end
  end

  on_linux do
    if Hardware::CPU.intel?
      url "https://github.com/opper-ai/oppercli/releases/download/v#{version}/opper-linux-amd64"
      sha256 "b6639b41ea4be58f8bd51839cd1cc7d0a0ffe7a49f59a8613e49c68155de07e2"
    end
  end

  def install
    bin.install Dir["opper-*"].first => "opper"
  end

  def caveats
    <<~EOS
      The legacy Go CLI and this tap are retired. To migrate:
        brew uninstall opper-ai/oppercli/opper
        brew untap opper-ai/oppercli
        npm install -g @opperai/cli
        opper login

      Keep ~/.oppercli for one-time credential migration when no new CLI config exists.
      Documentation: https://docs.opper.ai/developer-tools/cli
    EOS
  end

  test do
    system "#{bin}/opper", "--version"
  end
end

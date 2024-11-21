class Opper < Formula
  desc "Command line interface for Opper AI"
  homepage "https://github.com/opper-ai/oppercli"
  version "0.7.0"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/opper-ai/oppercli/releases/download/v#{version}/opper-darwin-arm64"
      sha256 "292462caffe86b79ce60e0499c14bd18b76e3a7c4de873a940b5ccca6dab0b7e"
    else
      url "https://github.com/opper-ai/oppercli/releases/download/v#{version}/opper-darwin-amd64"
      sha256 "dec6b2d1e7ca499921f07b5facc1c393196f6ed14936127a97c546c9704e04cd"
    end
  end

  on_linux do
    if Hardware::CPU.intel?
      url "https://github.com/opper-ai/oppercli/releases/download/v#{version}/opper-linux-amd64"
      sha256 "6c89649fc37fa28b6ae398fb966083723ed59935933b1a9f253a2f8356e84130"
    end
  end

  def install
    bin.install Dir["opper-*"].first => "opper"
  end

  test do
    system "#{bin}/opper", "--version"
  end
end

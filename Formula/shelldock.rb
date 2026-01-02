class Shelldock < Formula
  desc "A fast, cross-platform shell command repository manager"
  homepage "https://github.com/OpsGuild/ShellDock"
  url "https://github.com/OpsGuild/ShellDock/archive/v1.0.tar.gz"
  version "1.0"
  sha256 "8595e68cf875e0872f15a405e9feb213f7ba96cc87fd90c84c12164765dcf106"
  
  on_macos do
    if Hardware::CPU.intel?
      url "https://github.com/OpsGuild/ShellDock/releases/download/v1.0/shelldock-darwin-amd64"
      sha256 "8595e68cf875e0872f15a405e9feb213f7ba96cc87fd90c84c12164765dcf106"
    else
      url "https://github.com/OpsGuild/ShellDock/releases/download/v1.0/shelldock-darwin-arm64"
      sha256 "bd6cd96a5a6c3183c4f1457dfe2e294bdf2d7b4604cf37c8426e23e848f231f2"
    end
  end
  
  def install
    if OS.mac?
      bin.install "shelldock-darwin-#{Hardware::CPU.arch == "arm64" ? "arm64" : "amd64"}" => "shelldock"
    end
  end
  
  test do
    system "#{bin}/shelldock", "--version"
  end
end

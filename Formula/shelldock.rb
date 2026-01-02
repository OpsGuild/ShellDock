class Shelldock < Formula
  desc "A fast, cross-platform shell command repository manager"
  homepage "https://github.com/OpsGuild/ShellDock"
  url "https://github.com/OpsGuild/ShellDock/archive/v1.0.0.tar.gz"
  version "1.0.0"
  sha256 "9a073466afb7569905d9ca40c443e025e8c12fcc4fe018933a41ff9017035f6f"
  
  on_macos do
    if Hardware::CPU.intel?
      url "https://github.com/OpsGuild/ShellDock/releases/download/v1.0.0/shelldock-darwin-amd64"
      sha256 "9a073466afb7569905d9ca40c443e025e8c12fcc4fe018933a41ff9017035f6f"
    else
      url "https://github.com/OpsGuild/ShellDock/releases/download/v1.0.0/shelldock-darwin-arm64"
      sha256 "b3540693b386c7ba66cbe4dea2552987acf72f5fc40b6df0d7131c5fd3e532b8"
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

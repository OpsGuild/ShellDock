class Shelldock < Formula
  desc "A fast, cross-platform shell command repository manager"
  homepage "https://github.com/OpsGuild/ShellDock"
  url "https://github.com/OpsGuild/ShellDock/archive/v1.0.tar.gz"
  version "1.0"
  sha256 "29516d24b34e626d743f2fb245956aaf22ecb29c4bd88fa06052df436f7d2e17"
  
  on_macos do
    if Hardware::CPU.intel?
      url "https://github.com/OpsGuild/ShellDock/releases/download/v1.0/shelldock-darwin-amd64"
      sha256 "29516d24b34e626d743f2fb245956aaf22ecb29c4bd88fa06052df436f7d2e17"
    else
      url "https://github.com/OpsGuild/ShellDock/releases/download/v1.0/shelldock-darwin-arm64"
      sha256 "7d4e2271821a5c4db3336a9b778e4defdeac59dbb27ab436567f0ff3666aebec"
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

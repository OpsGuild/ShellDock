class Shelldock < Formula
  desc "A fast, cross-platform shell command repository manager"
  homepage "https://github.com/OpsGuild/ShellDock"
  url "https://github.com/OpsGuild/ShellDock/archive/v1.0.0.tar.gz"
  version "1.0.0"
  sha256 "8afd945f9dcdafd3d01cd2def2c3a686f19dfdca0160db39c64486b115a640ec"
  
  on_macos do
    if Hardware::CPU.intel?
      url "https://github.com/OpsGuild/ShellDock/releases/download/v1.0.0/shelldock-darwin-amd64"
      sha256 "8afd945f9dcdafd3d01cd2def2c3a686f19dfdca0160db39c64486b115a640ec"
    else
      url "https://github.com/OpsGuild/ShellDock/releases/download/v1.0.0/shelldock-darwin-arm64"
      sha256 "82e1bb77e38ecec2e19fb86b84c6fe28c8eb6e03eee56a150e986ad1ec6ce607"
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

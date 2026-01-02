class Shelldock < Formula
  desc "A fast, cross-platform shell command repository manager"
  homepage "https://github.com/OpsGuild/ShellDock"
  url "https://github.com/OpsGuild/ShellDock/archive/v1.0.tar.gz"
  version "1.0"
  sha256 "643e74d9bd9e26b1c22848d0194648fcbfbe97e6a771c29f359137782c7cfcf0"
  
  on_macos do
    if Hardware::CPU.intel?
      url "https://github.com/OpsGuild/ShellDock/releases/download/v1.0/shelldock-darwin-amd64"
      sha256 "643e74d9bd9e26b1c22848d0194648fcbfbe97e6a771c29f359137782c7cfcf0"
    else
      url "https://github.com/OpsGuild/ShellDock/releases/download/v1.0/shelldock-darwin-arm64"
      sha256 "07d0fe92fc9e3eda34a069721083956959b43d2ae9b4c191dd165a9f843d54b1"
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

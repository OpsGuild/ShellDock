class Shelldock < Formula
  desc "A fast, cross-platform shell command repository manager"
  homepage "https://github.com/OpsGuild/ShellDock"
  url "https://github.com/OpsGuild/ShellDock/archive/v1.0.0.tar.gz"
  version "1.0.0"
  sha256 "bbbfd2006b15db910080cf2c3d334906921458c1cb4ee2cc6980e1263c23b5e8"
  
  on_macos do
    if Hardware::CPU.intel?
      url "https://github.com/OpsGuild/ShellDock/releases/download/v1.0.0/shelldock-darwin-amd64"
      sha256 "bbbfd2006b15db910080cf2c3d334906921458c1cb4ee2cc6980e1263c23b5e8"
    else
      url "https://github.com/OpsGuild/ShellDock/releases/download/v1.0.0/shelldock-darwin-arm64"
      sha256 "cca0eb4603957ff5974e44386935473afb5317f8a5404de32b232917645dfde9"
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

class Shelldock < Formula
  desc "A fast, cross-platform shell command repository manager"
  homepage "https://github.com/OpsGuild/ShellDock"
  url "https://github.com/OpsGuild/ShellDock/archive/v1.0.tar.gz"
  version "1.0"
  sha256 "295d61cae0e4dcc7370aafd5e356123c5a2c93ce5e29cf12611e055505d31e3f"
  
  on_macos do
    if Hardware::CPU.intel?
      url "https://github.com/OpsGuild/ShellDock/releases/download/v1.0/shelldock-darwin-amd64"
      sha256 "295d61cae0e4dcc7370aafd5e356123c5a2c93ce5e29cf12611e055505d31e3f"
    else
      url "https://github.com/OpsGuild/ShellDock/releases/download/v1.0/shelldock-darwin-arm64"
      sha256 "f74e24787bc2861c9f7f63d8afa9db5b79ab27c6e7cbad506968fcd6cab52f28"
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

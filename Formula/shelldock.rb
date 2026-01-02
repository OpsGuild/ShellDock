class Shelldock < Formula
  desc "A fast, cross-platform shell command repository manager"
  homepage "https://github.com/OpsGuild/ShellDock"
  url "https://github.com/OpsGuild/ShellDock/archive/v1.0.0.tar.gz"
  version "1.0.0"
  sha256 "673a1933979b9affcc7dbc4b6c2e037dd595e40a923439a091751f5183963fb5"
  
  on_macos do
    if Hardware::CPU.intel?
      url "https://github.com/OpsGuild/ShellDock/releases/download/v1.0.0/shelldock-darwin-amd64"
      sha256 "673a1933979b9affcc7dbc4b6c2e037dd595e40a923439a091751f5183963fb5"
    else
      url "https://github.com/OpsGuild/ShellDock/releases/download/v1.0.0/shelldock-darwin-arm64"
      sha256 "36b6b34606191f9100e1b22e1f073c35c6a69653a6378fa5b761beedcb0b6036"
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

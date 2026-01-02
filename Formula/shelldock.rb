class Shelldock < Formula
  desc "A fast, cross-platform shell command repository manager"
  homepage "https://github.com/OpsGuild/ShellDock"
  url "https://github.com/OpsGuild/ShellDock/archive/v1.0.0.tar.gz"
  version "1.0.0"
  sha256 "850ca8c571987a752bbcc87c3903a36f72b249e33b9ee25f1e4d9151bdb92007"
  
  on_macos do
    if Hardware::CPU.intel?
      url "https://github.com/OpsGuild/ShellDock/releases/download/v1.0.0/shelldock-darwin-amd64"
      sha256 "850ca8c571987a752bbcc87c3903a36f72b249e33b9ee25f1e4d9151bdb92007"
    else
      url "https://github.com/OpsGuild/ShellDock/releases/download/v1.0.0/shelldock-darwin-arm64"
      sha256 "918abafddb0ae5478e382156451748463e6ae34e8a217024416aa325b410b108"
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

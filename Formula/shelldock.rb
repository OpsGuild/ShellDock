class Shelldock < Formula
  desc "A fast, cross-platform shell command repository manager"
  homepage "https://github.com/OpsGuild/ShellDock"
  url "https://github.com/OpsGuild/ShellDock/archive/v1.0.tar.gz"
  version "1.0"
  sha256 "41c972fa5c965ff3475c0db43305e347c0a37e261fcdf565034c5aaa46c40299"
  
  on_macos do
    if Hardware::CPU.intel?
      url "https://github.com/OpsGuild/ShellDock/releases/download/v1.0/shelldock-darwin-amd64"
      sha256 "41c972fa5c965ff3475c0db43305e347c0a37e261fcdf565034c5aaa46c40299"
    else
      url "https://github.com/OpsGuild/ShellDock/releases/download/v1.0/shelldock-darwin-arm64"
      sha256 "88873da0fd445669d570ed131a19ce34127ef2b3074e5e6842eb0b642a0cdca2"
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

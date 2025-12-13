Name:           shelldock
Version:        VERSION
Release:        1%{?dist}
Summary:        Fast shell command repository manager
License:        MIT
URL:            https://github.com/shelldock/shelldock
Source0:        %{name}-%{version}.tar.gz

BuildArch:      noarch
Requires:       glibc

%description
ShellDock is a fast, cross-platform tool for managing and executing
saved shell commands from bundled repository or local directory.

Features:
- Bundled command repository (included with installation)
- Local command repository (~/.shelldock)
- Interactive TUI for command management
- Step-by-step command execution with prompts
- Support for multiple package managers

%prep
%setup -q

%build
# Binary is pre-built, no compilation needed

%install
mkdir -p %{buildroot}/usr/local/bin
cp shelldock %{buildroot}/usr/local/bin/shelldock
chmod +x %{buildroot}/usr/local/bin/shelldock

mkdir -p %{buildroot}/usr/share/doc/shelldock
cp README.md %{buildroot}/usr/share/doc/shelldock/

mkdir -p %{buildroot}/usr/share/shelldock/repository
cp repository/*.yaml %{buildroot}/usr/share/shelldock/repository/ 2>/dev/null || true

%files
/usr/local/bin/shelldock
/usr/share/doc/shelldock/README.md
/usr/share/shelldock/repository/*.yaml

%changelog
* Wed Jan 01 2024 ShellDock Team <team@shelldock.io> - VERSION
- Initial release


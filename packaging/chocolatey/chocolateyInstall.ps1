$ErrorActionPreference = 'Stop'

# Note: This Go binary may trigger false positive antivirus detections.
# This is a known issue with statically compiled Go executables.
# The software is open source and safe. See VERIFICATION.txt for details.

$toolsDir = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"
Install-BinFile -Name shelldock -Path "$toolsDir\shelldock.exe"

# Copy repository files to tools directory (ShellDock looks for repository relative to executable)
if (Test-Path "$toolsDir\repository") {
    Write-Host "Repository files found, ensuring they are accessible"
    # Repository is already in tools directory from package creation
} else {
    Write-Warning "Repository directory not found in package"
}


$ErrorActionPreference = 'Stop'

# Note: This Go binary may trigger false positive antivirus detections.
# This is a known issue with statically compiled Go executables.
# The software is open source and safe. See VERIFICATION.txt for details.

$toolsDir = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"
Install-BinFile -Name shelldock -Path "$toolsDir\shelldock.exe"


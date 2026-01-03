$ErrorActionPreference = 'Stop'
$toolsDir = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"
Install-BinFile -Name shelldock -Path "$toolsDir\shelldock.exe"


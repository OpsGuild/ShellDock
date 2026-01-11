# Build Windows binary with version info embedded via goversioninfo
# Usage: .\scripts\build-windows.ps1 -Version "1.0.0" -OutputName "shelldock-windows-amd64.exe"

param(
    [Parameter(Mandatory=$true)]
    [string]$Version,
    
    [Parameter(Mandatory=$true)]
    [string]$OutputName,
    
    [string]$VersionInfoJson = "packaging/chocolatey/versioninfo.json"
)

$ErrorActionPreference = "Stop"

# Add Go bin to PATH
$env:PATH = "$env:PATH;$env:GOPATH\bin;$env:USERPROFILE\go\bin"

# Check if goversioninfo is available
$hasGoversioninfo = Get-Command goversioninfo -ErrorAction SilentlyContinue

if ($hasGoversioninfo) {
    Write-Host "Adding Windows version info using goversioninfo..."
    
    if (-not (Test-Path $VersionInfoJson)) {
        Write-Warning "Version info JSON not found at $VersionInfoJson, building without version info"
        $hasGoversioninfo = $false
    }
}

if ($hasGoversioninfo) {
    # Read and update version info
    $versionInfo = Get-Content $VersionInfoJson -Raw
    
    # Parse version parts
    $versionParts = $Version -split '\.'
    $major = if ($versionParts[0]) { $versionParts[0] } else { "1" }
    $minor = if ($versionParts[1]) { $versionParts[1] } else { "0" }
    $patch = if ($versionParts[2]) { $versionParts[2] } else { "0" }
    
    # Update version numbers in JSON
    $versionInfo = $versionInfo -replace '"Major": \d+', "`"Major`": $major"
    $versionInfo = $versionInfo -replace '"Minor": \d+', "`"Minor`": $minor"
    $versionInfo = $versionInfo -replace '"Patch": \d+', "`"Patch`": $patch"
    $versionInfo = $versionInfo -replace '"Build": \d+', "`"Build`": 0"
    $versionInfo = $versionInfo -replace '"FileVersion": "[^"]*"', "`"FileVersion`": `"$Version.0`""
    $versionInfo = $versionInfo -replace '"ProductVersion": "[^"]*"', "`"ProductVersion`": `"$Version.0`""
    
    # Write temp file
    $versionInfo | Out-File -Encoding UTF8 versioninfo-temp.json
    
    # Generate resource file
    goversioninfo -64 -o resource.syso versioninfo-temp.json
    
    # Clean up temp file
    Remove-Item versioninfo-temp.json
    
    Write-Host "Version info embedded successfully"
} else {
    Write-Host "goversioninfo not found, building without version info"
}

# Build the binary
Write-Host "Building $OutputName with version $Version..."
$env:CGO_ENABLED = "0"
go build -trimpath -ldflags "-X main.version=$Version" -buildmode=exe -o $OutputName .

# Clean up resource file if it exists
if (Test-Path resource.syso) {
    Remove-Item resource.syso
}

Write-Host "Build complete: $OutputName"


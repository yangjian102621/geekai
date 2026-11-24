[CmdletBinding()]
param(
    [Parameter(Mandatory = $true)]
    [string]$InstallerDir,
    [string]$Server = $env:DEPLOY_HOST,
    [string]$SshUser = $env:DEPLOY_USER,
    [ValidateRange(1, 65535)]
    [int]$Port = $(if ($env:DEPLOY_PORT) { [int]$env:DEPLOY_PORT } else { 22 }),
    [string]$IdentityFile = $env:DEPLOY_SSH_KEY,
    [string]$KnownHostsFile = $env:DEPLOY_KNOWN_HOSTS,
    [string]$Pod = $env:K8S_POD
)

$ErrorActionPreference = 'Stop'
if (-not $Server -or -not $SshUser) {
    throw 'Pass -Server and -SshUser, or set DEPLOY_HOST and DEPLOY_USER.'
}
if (-not (Test-Path -LiteralPath $InstallerDir -PathType Container)) {
    throw "Installer directory does not exist: $InstallerDir"
}
foreach ($keyPath in @($IdentityFile, $KnownHostsFile)) {
    if ($keyPath -and -not (Test-Path -LiteralPath $keyPath -PathType Leaf)) {
        throw "SSH file does not exist: $keyPath"
    }
}

# Windows' bash.exe may launch WSL; use Git for Windows' Bash explicitly.
$gitCommand = Get-Command git.exe -ErrorAction SilentlyContinue
$bashCandidates = @(
    "${env:ProgramFiles}/Git/bin/bash.exe",
    "${env:LOCALAPPDATA}/Programs/Git/bin/bash.exe"
)
if ($gitCommand) {
    $gitRoot = Split-Path (Split-Path $gitCommand.Source -Parent) -Parent
    $bashCandidates = @((Join-Path $gitRoot 'bin/bash.exe')) + $bashCandidates
}
$bashExe = $bashCandidates | Where-Object { Test-Path -LiteralPath $_ -PathType Leaf } | Select-Object -First 1
if (-not $bashExe) {
    throw 'Install Git for Windows to provide Bash, SSH, SCP and sha256sum.'
}

$cygpathExe = Join-Path (Split-Path (Split-Path $bashExe -Parent) -Parent) 'usr/bin/cygpath.exe'
$bashPaths = @()
foreach ($nativePath in @((Join-Path $PSScriptRoot 'publish-installers.sh'), $InstallerDir, $IdentityFile, $KnownHostsFile)) {
    if (-not $nativePath) {
        $bashPaths += ''
        continue
    }
    $resolvedPath = (Resolve-Path -LiteralPath $nativePath).ProviderPath
    $convertedPath = & $cygpathExe -u -- $resolvedPath
    if ($LASTEXITCODE -ne 0) {
        throw "Failed to convert Windows path: $resolvedPath"
    }
    $bashPaths += $convertedPath
}

$publishEnv = @{
    DEPLOY_HOST = $Server
    DEPLOY_USER = $SshUser
    DEPLOY_PORT = [string]$Port
    DEPLOY_SSH_KEY = $bashPaths[2]
    DEPLOY_KNOWN_HOSTS = $bashPaths[3]
    K8S_POD = $Pod
}
$savedEnv = @{}
try {
    foreach ($name in $publishEnv.Keys) {
        $savedEnv[$name] = [Environment]::GetEnvironmentVariable($name, 'Process')
        [Environment]::SetEnvironmentVariable($name, $publishEnv[$name], 'Process')
    }
    & $bashExe $bashPaths[0] $bashPaths[1]
    if ($LASTEXITCODE -ne 0) {
        throw "Publishing failed with exit code $LASTEXITCODE."
    }
    Write-Host 'Installers copied into the existing Pods and latest.json updated.'
}
finally {
    foreach ($name in $savedEnv.Keys) {
        [Environment]::SetEnvironmentVariable($name, $savedEnv[$name], 'Process')
    }
}

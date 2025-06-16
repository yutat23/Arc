@("linux", "windows", "darwin") | ForEach-Object {
  $env:GOOS = $_
  @("amd64", "arm64") | ForEach-Object {
    $env:GOARCH = $_
    $ext = if ($env:GOOS -eq "windows") { ".exe" } else { "" }
    $target = "$env:GOOS" + "_" + "$env:GOARCH"
    $buildDir = "build/arc_$target"
    $outputFile = "$buildDir/arc$ext"

    mkdir -Force $buildDir | Out-Null
    go build -o $outputFile -ldflags "-s -w"

    # ZIP化
    $zipPath = "release/arc_$target.zip"
    mkdir -Force "release" | Out-Null
    if (Test-Path $zipPath) { Remove-Item $zipPath }
    Compress-Archive -Path "$buildDir/*" -DestinationPath $zipPath
  }
}
Remove-Item Env:GOOS, Env:GOARCH

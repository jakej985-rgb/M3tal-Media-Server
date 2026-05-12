# M3TAL Platform Build Script (Windows)

Write-Host "🚀 Building M3TAL CLI..." -ForegroundColor Cyan
go build -o m3tal.exe ./cmd/m3tal

Write-Host "🚀 Building M3TAL API..." -ForegroundColor Cyan
go build -o m3tal-api.exe ./cmd/api

if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ Build complete." -ForegroundColor Green
} else {
    Write-Host "❌ Build failed. Ensure Go is installed and in your PATH." -ForegroundColor Red
}

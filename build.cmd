@echo off
REM ============================================================================
REM BUILD SCRIPT FOR sleepms (Windows)
REM ============================================================================
REM This script automates the build process for the sleepms application.
REM It handles version injection, documentation generation, and builds.
REM
REM Usage: build.cmd
REM Output: sleepms-{OS}-{ARCH}.exe binary and GODOC.md documentation
REM ============================================================================

REM ============================================================================
REM PROJECT CONFIGURATION (modify these for your project)
REM ============================================================================
set PROJECT_NAME=sleepms
set PROJECT_DESCRIPTION=Random Sleep Utility with Progress Bar

REM Semantic versioning - update these when releasing new versions
set VERSION_MAJOR=1
set VERSION_MINOR=0
set VERSION_REVISION=0

REM Source code configuration
set SOURCE_FILES=.
set MAIN_PACKAGE=.

REM Build output configuration
set SYMLINK_NAME=sleepms.exe

REM Documentation generation
set GENERATE_DOCS=true
set DOC_OUTPUT=GODOC.md

REM ============================================================================
REM BUILD SCRIPT IMPLEMENTATION (generic - no need to modify below this line)
REM ============================================================================

echo ===============================================================
echo   Building %PROJECT_NAME% - %PROJECT_DESCRIPTION%
echo ===============================================================
echo.

REM Check if Go is installed
echo [*] Checking Go installation...
where go >nul 2>&1
if %errorlevel% neq 0 (
    echo [X] Go is not installed. Please install Go and try again.
    exit /b 1
)

for /f "tokens=*" %%i in ('go version') do set GO_VERSION=%%i
echo [+] Go found: %GO_VERSION%
echo.

REM Detect OS and architecture
echo [*] Detecting platform...
set OS=windows

REM Detect architecture
if "%PROCESSOR_ARCHITECTURE%"=="AMD64" (
    set ARCH=amd64
) else if "%PROCESSOR_ARCHITECTURE%"=="ARM64" (
    set ARCH=arm64
) else if "%PROCESSOR_ARCHITECTURE%"=="x86" (
    set ARCH=386
) else (
    echo [X] Unsupported architecture: %PROCESSOR_ARCHITECTURE%
    exit /b 1
)

echo [+] Platform: %OS%/%ARCH%
echo.

REM Generate build number from current UTC time (MMDDHHMM format)
echo [*] Generating build number...

REM Get UTC time components
for /f "tokens=1-4 delims=/ " %%a in ('wmic path win32_utctime get Month^,Day^,Hour^,Minute /format:table ^| findstr /r "[0-9]"') do (
    set DAY=%%a
    set HOUR=%%b
    set MINUTE=%%c
    set MONTH=%%d
)

REM Pad with zeros if needed
if %MONTH% lss 10 set MONTH=0%MONTH%
if %DAY% lss 10 set DAY=0%DAY%
if %HOUR% lss 10 set HOUR=0%HOUR%
if %MINUTE% lss 10 set MINUTE=0%MINUTE%

set BUILD_NUMBER=%MONTH%%DAY%%HOUR%%MINUTE%

REM Get build time
for /f "tokens=*" %%i in ('powershell -command "Get-Date -Format 'yyyy-MM-dd HH:mm:ss UTC'"') do set BUILD_TIME=%%i

echo [+] Build number: %BUILD_NUMBER%
echo [+] Build time: %BUILD_TIME%
echo.

REM Construct version string
set VERSION_STRING=%VERSION_MAJOR%.%VERSION_MINOR%.%VERSION_REVISION%
set FULL_VERSION=%VERSION_STRING% build %BUILD_NUMBER%

echo [*] Version: %FULL_VERSION%
echo.

REM Construct linker flags for version injection
set LDFLAGS=-s -w
set LDFLAGS=%LDFLAGS% -X "main.AppName=%PROJECT_NAME%"
set LDFLAGS=%LDFLAGS% -X "main.VersionMajor=%VERSION_MAJOR%"
set LDFLAGS=%LDFLAGS% -X "main.VersionMinor=%VERSION_MINOR%"
set LDFLAGS=%LDFLAGS% -X "main.VersionRevision=%VERSION_REVISION%"
set LDFLAGS=%LDFLAGS% -X "main.BuildNumber=%BUILD_NUMBER%"
set LDFLAGS=%LDFLAGS% -X "main.BuildTime=%BUILD_TIME%"

REM Determine output filename
set OUTPUT_NAME=%PROJECT_NAME%-%OS%-%ARCH%.exe

REM Build the application
echo [*] Building %PROJECT_NAME%...
go build -ldflags "%LDFLAGS%" -o "%OUTPUT_NAME%" %MAIN_PACKAGE%

if %errorlevel% equ 0 (
    echo [+] Build successful: %OUTPUT_NAME%

    REM Get file size
    for %%A in ("%OUTPUT_NAME%") do set FILE_SIZE=%%~zA
    echo [+] Binary size: %FILE_SIZE% bytes
) else (
    echo [X] Build failed
    exit /b 1
)
echo.

REM Create copy for easy access
echo [*] Creating convenience copy...
copy /Y "%OUTPUT_NAME%" "%SYMLINK_NAME%" >nul 2>&1
if %errorlevel% equ 0 (
    echo [+] Created copy: %SYMLINK_NAME%
) else (
    echo [!] Failed to create copy: %SYMLINK_NAME%
)
echo.

REM Generate documentation if enabled
if "%GENERATE_DOCS%"=="true" (
    echo [*] Generating documentation...

    REM Create documentation header
    (
        echo # %PROJECT_NAME% - Go Package Documentation
        echo.
        echo **Version:** %FULL_VERSION%
        echo **Generated:** %BUILD_TIME%
        echo.
        echo ---
        echo.
        echo ## Package Overview
        echo.
    ) > "%DOC_OUTPUT%"

    REM Extract package documentation
    go doc -all >> "%DOC_OUTPUT%" 2>&1

    if %errorlevel% equ 0 (
        echo [+] Documentation generated: %DOC_OUTPUT%
    ) else (
        echo [!] Documentation generation had warnings (check %DOC_OUTPUT%^)
    )
    echo.
)

REM Print build summary
echo ===============================================================
echo   Build Complete!
echo ===============================================================
echo.
echo   Project:     %PROJECT_NAME%
echo   Version:     %FULL_VERSION%
echo   Platform:    %OS%/%ARCH%
echo   Binary:      %OUTPUT_NAME%
if "%GENERATE_DOCS%"=="true" (
    echo   Docs:        %DOC_OUTPUT%
)
echo.
echo Run with: %SYMLINK_NAME% --version
echo Usage:    %SYMLINK_NAME% ^<min^> ^<max^>
echo.

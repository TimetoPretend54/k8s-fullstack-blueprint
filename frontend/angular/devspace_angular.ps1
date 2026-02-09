# PowerShell script for Windows - runs Angular watch mode locally
# This builds the app to ./dist/ which DevSpace syncs to the container

# Clear the screen for clean output
Clear-Host

$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location $ScriptDir

# Define colors
$ColorBlue = "DarkBlue"
$ColorGreen = "Green"
$ColorCyan = "Cyan"
$ColorYellow = "Yellow"

# Print styled banner
Write-Host @"
     %########%      
     %###########%       ____                 _____                      
         %#########%    |  _ \   ___ __   __ / ___/  ____    ____   ____ ___ 
         %#########%    | | | | / _ \\ \\ / / \___ \ |  _ \  / _  | / __// _ \\
     %#############%    | |_| |(  __/ \ V /  ____) )| |_) )( (_| |( (__(  __/
     %#############%    |____/  \___|  \_/   \____/ |  __/  \__,_| \___\\\___|
 %###############%                                  |_|
 %###########%$
"@ -ForegroundColor $ColorBlue

Write-Host ""
Write-Host "Welcome to your Angular frontend development container with Devspace!" -ForegroundColor $ColorGreen
Write-Host ""

Write-Host "This is how you can work with it:" -ForegroundColor $ColorBlue
Write-Host "- Files will be synchronized between your local machine and this container via Devspace" -ForegroundColor $ColorYellow
Write-Host "- Some ports will be forwarded, so you can access this container via localhost" -ForegroundColor $ColorYellow
Write-Host "- Run 'npm run watch' to start the application manually" -ForegroundColor $ColorYellow
Write-Host "- Or just let ng serve handle it automatically (it's already running)" -ForegroundColor $ColorYellow
Write-Host ""
Write-Host "Navigate to http://localhost:4200/ to view your 'hello-world' application." -ForegroundColor $ColorCyan
Write-Host "Navigate to http://localhost:4200/appt-booking/home/ to view your 'appt-booking' application." -ForegroundColor $ColorCyan

# Run npm run watch (this typically runs ng build --watch)
npm run watch

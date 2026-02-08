#!/bin/bash
# Bash script for Linux/Mac - runs Angular watch mode locally
# This builds the app to ./dist/ which DevSpace syncs to the container

# Clear the screen for clean output
clear

cd "$(dirname "$0")"

# Colors for output
COLOR_BLUE="\033[0;94m"
COLOR_GREEN="\033[0;92m"
COLOR_CYAN="\033[0;96m"
COLOR_YELLOW="\033[0;93m"
COLOR_RESET="\033[0m"

# Print styled banner
echo -e "${COLOR_BLUE}
     %########%      
     %###########%       ____                 _____                      
         %#########%    |  _ \   ___ __   __ / ___/  ____    ____   ____ ___ 
         %#########%    | | | | / _ \\ \\ / / \___ \ |  _ \  / _  | / __// _ \\
     %#############%    | |_| |(  __/ \ V /  ____) )| |_) )( (_| |( (__(  __/
     %#############%    |____/  \___|  \_/   \____/ |  __/  \__,_| \___\\\___|
 %###############%                                  |_|
 %###########%$
${COLOR_RESET}"

echo -e ""
echo -e "${COLOR_GREEN}Welcome to your Angular frontend development container with Devspace!${COLOR_RESET}"
echo -e ""

echo -e "This is how you can work with it:"
echo -e "- Files will be synchronized between your local machine and this container via Devspace"
echo -e "- Some ports will be forwarded, so you can access this container via localhost"
echo -e "- Run \`${COLOR_GREEN}npm run watch${COLOR_RESET}\` to start the application manually"
echo -e "- Or just let ng serve handle it automatically (it's already running)"
echo -e ""
echo -e "Navigate to ${COLOR_CYAN}http://localhost:4200/${COLOR_RESET} view your application."

# Run npm run watch (this typically runs ng build --watch)
npm run watch

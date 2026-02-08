#!/bin/bash
set +e  # Continue on errors

COLOR_BLUE="\033[0;94m"
COLOR_GREEN="\033[0;92m"
COLOR_RESET="\033[0m"

# Print useful output for user
echo -e "${COLOR_BLUE}
     %########%      
     %###########%       ____                 _____                      
         %#########%    |  _ \   ___ __   __ / ___/  ____    ____   ____ ___ 
         %#########%    | | | | / _ \\ \\ / / \___ \ |  _ \  / _  | / __// _ \\
     %#############%    | |_| |(  __/ \ V /  ____) )| |_) )( (_| |( (__(  __/
     %#############%    |____/  \___|  \_/   \____/ |  __/  \__,_| \___\\\___|
 %###############%                                  |_|
 %###########%${COLOR_RESET}"


echo -e ""
echo -e "Welcome to your Go (Echo) backend development container!"
echo -e ""
echo -e "This is how you can work with it:"
echo -e "- Files will be synchronized between your local machine and this container"
echo -e "- Some ports will be forwarded, so you can access this container via localhost"
echo -e "- Run \`${COLOR_GREEN}go run main.go${COLOR_RESET}\` to start the application manually"
echo -e "- Or just let Air handle it automatically (it's already running)"
echo -e ""

# Set terminal prompt
export PS1="\[${COLOR_BLUE}\]devspace\[${COLOR_RESET}\] ./\\W \[${COLOR_BLUE}\]\\$\[${COLOR_RESET}\] "
if [ -z "$BASH" ]; then export PS1="$ "; fi

# Include project's bin/ folder in PATH
export PATH="./bin:$PATH"

# Install Air if not already installed
if ! command -v air >/dev/null 2>&1; then
    echo -e "Installing Air..."
    go install github.com/cosmtrek/air@v1.51.0
fi

# Download Go dependencies
echo -e ""
echo -e "Downloading Go modules..."
go mod download
echo -e ""

echo -e "Starting Air hot-reload..."
exec air

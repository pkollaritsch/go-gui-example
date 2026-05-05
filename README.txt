Is there an open-source and free of fee GoLang GUI ready for production 
Apps which still maintains the "write once, run anywhere" goal?
A video explaining and demonstrating the example in this repo can be 
found at:  https://youtu.be/1smsqZUAYn4


########## This following has been tested on MacOS: ####################
# download and install go language from: https://go.dev/dl/
echo 'export PATH=$PATH:"/usr/local/go/bin"' >> ~/.bashrc

go get cogentcore.org/core
go get cogentcore.org/core/core@upgrade
go get cogentcore.org/core/events@upgrade
go get cogentcore.org/core/styles@upgrade
go get cogentcore.org/core/styles/units@upgrade
go get cogentcore.org/core/colors@upgrade
go get cogentcore.org/core/paint@upgrade
go get golang.org/x/image/math/fixed@upgrade
go get golang.org/x/image@upgrade
go get cogentcore.org/core/system@upgrade
go get cogentcore.org/core/colors/cam/hct@upgrade
go get cogentcore.org/core/paint/render@upgrade
go get cogentcore.org/core/paint/pimage@upgrade
go get cogentcore.org/core/htmlcore@upgrade
go get github.com/stretchr/testify/assert@upgrade
go get cogentcore.org/core/text/shaped@upgrade
go get cogentcore.org/core/math32@upgrade


############# Init Setup To create executables ################################
go install cogentcore.org/core@latest

#Install homebrew and xCode:
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
echo 'eval "$(/opt/homebrew/bin/brew shellenv)"' >> ~/.bashrc
source ~/.bashrc

#install the Windows cross-compiler:
brew install mingw-w64

#install the docker linux cross-compiler:
#Goto https://www.docker.com/products/docker-desktop/ and download docker desktop for MAC apple silicon:
https://www.docker.com/products/docker-desktop/
Start the docker app and a "Finish setting up docker desktop" window will appear.
Select "advanced settings" 
  * select "User" so it is installed under ~/.docker/bin and so won't need your password.
  * unselect "Allow the default Docker socket to be used (requires password)"
  * unselect "Allow privileged port mapping (requires password)"
  * unselect "Automatically check configuration"
  * click "Finish" button
Add this to ~/.bashrc:
  export PATH="/Users/pkollaritsch/.docker/bin":$PATH   #cannot have spaces around equal sign!!!!
source ~/.bashrc


############# Create executables ####################################################################
###### To cross-compile to windows/amd64:
CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc ~/go/bin/core build windows/amd64
mv cogentcore_example.exe cogentcore_example_win_x86_64.exe
I verified this worked.

###### To compile to MAC OS for apple or intel silicon:
go build .
mv main cogentcore_example_mac.exe
I verified this worked.

###### Build for ARM64 Linux (matches my MACOS architecture)
docker run --rm -v "$PWD":/app -w /app golang:latest \
  sh -c "apt-get update && apt-get install -y libx11-dev libxcursor-dev libxrandr-dev libxinerama-dev libxi-dev && CGO_ENABLED=1 go build -o myapp-linux"
I have not tried this.

######  Build for amd64 Linux (requires cross-compilation toolchain)
docker run --rm -v "$PWD":/app -w /app golang:latest \
  sh -c "dpkg --add-architecture amd64 && apt-get update && apt-get install -y gcc-x86-64-linux-gnu g++-x86-64-linux-gnu libx11-dev:amd64 libxcursor-dev:amd64 libxrandr-dev:amd64 libxinerama-dev:amd64 libxi-dev:amd64 libwayland-dev:amd64 libxkbcommon-dev:amd64 libgl1-mesa-dev:amd64 libglu1-mesa-dev:amd64 libxxf86vm-dev:amd64 && CC=x86_64-linux-gnu-gcc CXX=x86_64-linux-gnu-g++ GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o cogentcore_example_linux_x86_64.exe"

To Test on Ubuntu Linux:
Click the password entry box and Before typing your password in Ubuntu, look for a small gear icon in the bottom right corner of the screen.  Click that gear and select "Ubuntu on Xorg" (this is Ubuntu's name for X11).  Log in, open your terminal, and run ./cogentcore_example_linux_x86_64.exe

I verified this worked.

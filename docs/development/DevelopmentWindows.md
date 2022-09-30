## QuantstopTerminal Development Guide - Windows

### Prerequisites

* [Go](https://go.dev/)
  1. Download Go 1.17.5 from https://go.dev/dl/go1.17.5.windows-amd64.msi
  2. Open the MSI file you downloaded and follow the prompts to install Go. <br>
     By default, the installer will install Go to Program Files or Program Files (x86). <br>
     You can change the location as needed. After installing, you will need to close and reopen any 
     open command prompts so that changes to the environment made by the installer are reflected at the command prompt.
  3. Verify that you've installed Go.
     1. In Windows, click the Start menu.
     2. In the menu's search box, type cmd, then press the Enter key.
     3. In the Command Prompt window that appears, type the following command:
        ```sh 
        go version 
        ```
     4. Confirm that the command prints the installed version of Go like below:
        ```sh 
        go version go1.17.5 windows/amd64 
        ```
* [NodeJS/NPM](https://nodejs.org/en/)
  1. Download NodeJS 16.13.1 LTS from https://nodejs.org/dist/v16.13.1/node-v16.13.1-x64.msi
  2. Open the MSI file you downloaded and follow the prompts to install NodeJS. <br>
  3. Verify that you've installed NodeJS.
     1. In Windows, click the Start menu.
     2. In the menu's search box, type cmd, then press the Enter key.
     3. In the Command Prompt window that appears, type the following command:
        ```sh 
        node -v 
        ```
     4. Confirm that the command prints the installed version of Go like below:
        ```sh 
        v16.13.1
        ```
* [Vue.Js]()
  1. Todo: 
* [Taskfile](https://taskfile.dev/#/installation)
  1. Follow the [instructions](https://taskfile.dev/#/installation)
* [air](https://github.com/cosmtrek/air)
  1. Follow the [instructions](https://github.com/cosmtrek/air#installation)

**Note:** If you want to build the Windows installer for QuantstopTerminal, you must also install [WixToolset.](https://wixtoolset.org/releases/)
This toolset makes a standard .msi installer for the executables within QuantstopTerminal.

### Installation

_Once you have installed the required dependencies above, continue with getting the project 
source files from GitHub. This guide assumes you have a terminal with git installed._

1. Clone the repo
   ```sh
   git clone https://github.com/QuantStop/QuantStopTerminal.git
   ```
2. Run the app directly in development mode:<br>
   ```sh 
     task run
     ```
   Or build production executables: <br>
   _(Note: All executable files will be located in /builds)_ 
   ```sh 
     task build
     ```
3. Profit.

<!--suppress HtmlDeprecatedAttribute -->
<p align="right">(<a href="https://github.com/quantstop/quantstopterminal#top">back to home</a>)</p>

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

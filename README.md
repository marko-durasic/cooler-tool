# Cooler (Go Edition)

A tool written in Go to monitor and cool down your system.

## Installation

1.  **Clone the repository (or download the files).**
    (Assuming the files are in a directory named `cooler`)

2.  **Run the installer.**
    Open a terminal in the `cooler` directory and run:
    ```bash
    ./install.sh
    ```
    The script will guide you through the installation of all necessary dependencies and will compile the `cooler` executable.

## Usage

Once the installation is complete, you can run the tool from the project directory with:
```bash
./cooler
```

## Prerequisites (Managed by Installer)

The installer script will handle the following prerequisites:

-   Go (Golang) compiler
-   Node.js (v18+) & npm
-   Google Gemini CLI
-   `lm-sensors` & `cpufrequtils`
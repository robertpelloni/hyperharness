# Borg Deployment Instructions

_This document contains the latest deployment instructions for the Borg Universal AI Dashboard and Cognitive Control Plane._

## Prerequisites

1.  **Node.js**: >= 22.12.0
2.  **pnpm**: Recommended package manager (`npm install -g pnpm`)
3.  **Git**: For submodule fetching and version control.

## Initial Setup

1.  **Clone the Repository**:
    ```bash
    git clone https://github.com/robertpelloni/borg.git
    cd borg
    ```

2.  **Initialize Submodules**:
    ```bash
    git submodule update --init --recursive
    ```

3.  **Install Dependencies**:
    ```bash
    pnpm install
    ```

4.  **Environment Variables**:
    Copy `.env.example` to `.env` and fill in the required API keys (OpenAI, Anthropic, Gemini, etc.).
    ```bash
    cp .env.example .env
    ```

## Running the Platform

Borg is designed as a long-running service that manages PC memory, CPU, disk, and bandwidth usage.

### Start the Dashboard & Core Server
Use the provided startup scripts:

**Windows**:
```bash
.\start.bat
```

**Linux/macOS**:
```bash
./start.sh
```

Alternatively, run via pnpm:
```bash
pnpm dev
```

This will:
1. Start the Borg Server (core).
2. Start the MCP Router (client/server proxy).
3. Open the Web Dashboard (Next.js) at `http://localhost:3000`.

### Building for Production
```bash
pnpm build
pnpm start
```

## Extension Installation

Once the dashboard is running, navigate to the **Integrations** tab in the WebUI to install:
*   Browser Extensions (Chrome, Firefox).
*   IDE Plugins (VSCode, Cursor, Windsurf).
*   CLI Harnesses.

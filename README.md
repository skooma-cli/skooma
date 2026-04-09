# 🧪 Skooma

![Release Version](https://img.shields.io/github/v/release/mark-rodgers/skooma)
![Go Version](https://img.shields.io/badge/Go-1.26.1%2B-00ADD8?logo=go&logoColor=white)
![License](https://img.shields.io/badge/license-MIT-green)

> _"Khajiit has wares, if you have a project name."_

Skooma is a CLI tool that scaffolds full-stack single-page applications in seconds. One command gives you a Go/Gin backend, a React + TypeScript + Vite + Tailwind CSS frontend, and a `docker-compose.yml` — all wired together and ready to run. Named after the famously suspicious substance from Elder Scrolls lore, because good scaffolding tools should feel a little magical.

---

## What Gets Brewed

Running `skooma brew` conjures the following stack:

| Layer              | Tech                                                                                        |
| ------------------ | ------------------------------------------------------------------------------------------- |
| **Backend**        | Go + [Gin](https://github.com/gin-gonic/gin) + [godotenv](https://github.com/joho/godotenv) |
| **Frontend**       | React 19 + TypeScript + Vite + Tailwind CSS + ESLint                                        |
| **Infrastructure** | Docker Compose (backend, frontend, and database services)                                   |

---

## Installation

**Prerequisites:** Go 1.26.1+, `make`, Git

```bash
git clone https://github.com/mark-rodgers/skooma.git
cd skooma
make build
```

This produces a binary at `bin/skooma` (or `bin/skooma.exe` on Windows). Add it to your `$PATH` to use it anywhere.

---

## Usage

### Interactive mode

Run `brew` with just a project name (or no arguments at all). Skooma will walk you through a TUI form asking for the remaining details.

```bash
skooma brew myapp
```

You'll be prompted for:

- **Project name**
- **Repository URL** — e.g. `github.com/user/myapp`
- **Author** — e.g. `Name <email@example.com>`
- **Database** — Flat File, Microsoft SQL, or PostgreSQL

### Non-interactive mode

Pass all flags upfront to skip the form entirely — useful for scripts and automation.

```bash
skooma brew myapp \
  --repo github.com/you/myapp \
  --author "Jane Doe <jane@example.com>" \
  --database postgres
```

### Flags

| Flag         | Short | Description                                    | Default      |
| ------------ | ----- | ---------------------------------------------- | ------------ |
| `--repo`     | `-r`  | Repository URL (e.g. `github.com/you/myapp`)   | _(prompted)_ |
| `--author`   | `-a`  | Author name and email in `Name <email>` format | _(optional)_ |
| `--database` | `-d`  | Database type: `file`, `mssql`, or `postgres`  | `file`       |

---

## Generated Project Structure

After running `skooma brew myapp`, you'll find a `myapp/` directory with the following layout:

```
myapp/
├── docker-compose.yml
├── backend/
│   ├── go.mod
│   ├── main.go
│   └── Makefile
└── frontend/
    ├── index.html
    ├── package.json
    ├── vite.config.ts
    ├── tsconfig.json
    ├── tsconfig.app.json
    ├── tsconfig.node.json
    ├── eslint.config.js
    ├── .gitignore
    ├── public/
    └── src/
        ├── main.tsx
        ├── App.tsx
        ├── App.css
        └── index.css
```

---

## Database Options

The `--database` flag controls what database service gets configured in `docker-compose.yml` and wired into your backend.

| Value      | Description                                                                                      |
| ---------- | ------------------------------------------------------------------------------------------------ |
| `file`     | Flat file storage — no database container, no extra dependencies. Good for getting started fast. |
| `mssql`    | Microsoft SQL Server container.                                                                  |
| `postgres` | PostgreSQL container. The default in generated `docker-compose.yml` examples.                    |

---

## Building Skooma

If you're hacking on Skooma itself:

```bash
# Build the binary
make build

# Clean build artifacts
make clean
```

The output binary lands in `bin/`.

---

<p align="center">
  <img height="300" src="cmd/templates/frontend/public/khajiit.webp">
</p>

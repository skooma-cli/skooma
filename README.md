# 🧪 Skooma

![Release Version](https://img.shields.io/github/v/release/skooma-cli/skooma?include_prereleases)
![Go Version](https://img.shields.io/badge/Go-1.26.1%2B-00ADD8?logo=go&logoColor=white)
![License](https://img.shields.io/badge/license-MIT-green)

> _"Khajiit has wares, if you have a project name."_

Skooma is a flexible CLI tool that scaffolds full-stack single-page applications in seconds using external template repositories. One command sets up your entire development environment with backend, frontend, and infrastructure — all wired together and ready to run. Named after the famously suspicious substance from Elder Scrolls lore, because good scaffolding tools should feel a little magical.

---

## What Gets Brewed

Running `skooma brew` creates a complete project structure based on your chosen template. Templates define what gets generated - from simple static sites to complex multi-service applications:

| Template Defines   | What You Get                                       |
| ------------------ | -------------------------------------------------- |
| **Project Files**  | Source code, configuration, and build scripts      |
| **Dependencies**   | Package definitions and dependency management      |
| **Tooling**        | Development workflows, linting, and testing setup  |
| **Infrastructure** | Deployment configuration and service orchestration |
| **Variables**      | Custom prompts and inputs collected during a brew  |

---

## Templates

Skooma uses external template repositories to generate projects, allowing for flexible and community-driven project scaffolding. Templates are git repositories containing `.tmpl` files that get processed and customized for each new project.

### Default Template

The [skooma-template-default](https://github.com/skooma-cli/skooma-template-default) provides a full-stack web application setup with:

- **Frontend**: React with TypeScript, Vite build tooling, Tailwind CSS styling, and ESLint
- **Backend**: Go with Gin web framework and environment configuration
- **Infrastructure**: Docker Compose orchestration with database services
- **Database Options**: PostgreSQL, Microsoft SQL Server, or flat file storage

See the template repository for specific versions, detailed setup instructions, and customization options.

### Template Management

Templates can be managed through the Skooma configuration:

```bash
# List available templates
skooma template ls

# Add a custom template repository
skooma template add my-template github.com/user/skooma-template-custom

# Remove a template
skooma template rm my-template

# Use a specific template
skooma brew my-app --template my-template
```

### Creating Custom Templates

Want to create your own template? Templates are git repositories with:

1. **Configuration file** (`skooma.config.json`) defining available variables and prompts
2. **Template files** ending in `.tmpl` (e.g., `index.html.tmpl`)
3. **Template variables** using Go's [text/template](https://pkg.go.dev/text/template) syntax

Template variables are dynamically defined by each template's `skooma.config.json` file. When brewing a project, Skooma will prompt for all variables defined in the template's configuration, then process all `.tmpl` files using Go's templating engine.

See the [default template](https://github.com/skooma-cli/skooma-template-default) repository for a complete example.

---

## Installation

**Prerequisites:** Go 1.26.1+

### Quick Install

```bash
go install github.com/skooma-cli/skooma@latest
```

### Build from Source

If you want to hack on Skooma or build a specific version:

```bash
git clone https://github.com/skooma-cli/skooma.git
cd skooma
make build
```

This produces a binary at `bin/skooma` (or `bin/skooma.exe` on Windows). Add it to your `$PATH` to use it anywhere.

---

## Usage

### Interactive mode

Run `brew` with just a project name (or no arguments at all). Skooma will walk you through a TUI form asking for the remaining details, including template selection.

```bash
skooma brew my-app
```

You'll be prompted for:

- **Project name** — Alphanumeric, no spaces or underscores
- **Template** — Choose from available template repositories
- **Repository URL** — e.g. `github.com/user/repo`
- **Author** — e.g. `Name <email@example.com>`
- **Database** — Flat File, Microsoft SQL Server, or PostgreSQL

### Non-interactive mode

Pass all flags upfront to skip the form entirely — useful for scripts and automation.

```bash
skooma brew my-app \
  --template default \
  --repo github.com/user/repo \
  --author "Jane Doe <jane@example.com>" \
  --database postgres
```

### Flags

| Flag         | Short | Description                                    | Default      |
| ------------ | ----- | ---------------------------------------------- | ------------ |
| `--template` | `-t`  | Template name from configured templates        | _(prompted)_ |
| `--repo`     | `-r`  | Repository URL (e.g. `github.com/user/repo`)   | _(prompted)_ |
| `--author`   | `-a`  | Author name and email in `Name <email>` format | _(optional)_ |
| `--database` | `-d`  | Database type: `file`, `mssql`, or `postgres`  | `file`       |

---

## Related Projects

- **[Skooma Default Template](https://github.com/skooma-cli/skooma-template-default)** - The official default template
- **[Template Gallery](https://github.com/topics/skooma-template)** - Community templates and examples

## License

MIT License - see the [LICENSE](LICENSE) file for details.

---

<p align="center">
  <img height="300" src="cmd/templates/frontend/public/khajiit.webp">
</p>

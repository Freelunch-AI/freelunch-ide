# Freelunch IDE - Contributor Guide

## Codebase Explanation

### Folder & File Structure

### Step-by-step Codebase Tutorial

## Contributing Rules

### Issues

- Follow Issue Template
- Issues must be self-contained and highly descriptive of the problem and solution that should be implemented
- Every Bug Issue closed must create a ./docs/post-mortems/post-mortem_[i].md where i is the number of the issue, which describes the issue, when it happenned, the data evidence for the issue (telemetry, screenshots, user reviews, etc), desription of the solution, data evidence of the solution, who was involved in solving it, time it took to solve it, the process that was used to arrive at the solution, how to avoid and reduce the impact of similar problems in the future, the changes in standard operating protocols imposed.

### PRs

- Follow PR Template
- Code PRs need to resolve a specific issue
- PRs must be self-contained and highly descriptive of the solution and solution implementation
- Prefer multiple straighforward PRs over a single big PR with multiple different stuff



---

# Foundations

The following documentation explains the  code base in plain terms, and walks through getting the repository running on **Linux** or **Windows**.

You do not need to know Go or TypeScript to follow the setup section. Work through it top to bottom and copy the commands exactly.

> **Status of these instructions:** every command in the "Daily commands" section was run and verified on macOS. The Linux and Windows steps are derived from the repository's
> configuration files, not executed on those machines. If something does not work, that is a bug in this document — say so and it will be fixed.

---

## Part 1 — What this pull request does

Five things, in the order they happened.

### 1. Pixi became the single source of truth for tool versions

Everyone building this project must use the **same versions** of Go, Node, and the linters.
If one person has Go 1.24 and another has Go 1.26, code that works for one breaks for the other, and nobody can tell why.

[Pixi](https://pixi.sh) solves this. `pixi.toml` declares the tools; `pixi.lock` records the exact resolved versions. When you run `pixi install`, you get precisely what everyone else has — currently Go 1.26.5, Node 26.6.0, golangci-lint 2.12.2.

**What this means for you:** never install Go or Node yourself. Pixi provides them.

### 2. The Go code was restructured around a ServiceManager

The CLI is built from *services* (logging, command handling) that are registered in one place and wired together at startup. This is so new features can be added without editing unrelated code.

Two rules govern it, and they exist because `freelunch` is a **command-line tool, not a server**. A server starts once and runs for days; a CLI starts fresh on every single command, including `freelunch --help`.

- **`Start()` must not do any I/O.** No network calls, no file reads. Anything expensive is loaded later, on first actual use.
- **`Start()` must never block.** A test (`Test_serviceManagerFinal_StartDoesNotBlock`) enforces this and will fail the build if broken.

If you add a service, follow the existing ones in `src/cli/internal/` closely.

### 3. A broken quality gate was found and fixed

The command `pixi run check` runs every formatter, linter, test, and build. It had been
**failing for four days without anyone noticing**, because nothing runs it automatically.

Two causes:

- **ESLint was scanning `.pixi/`** — the folder where Pixi installs the toolchain. Go ships
  some JavaScript files of its own inside it, and ESLint was reporting **1037 errors** from
  them. None were from our code. ESLint, unlike other tools, does **not** read `.gitignore`,
  so it had to be told explicitly to skip that folder.
- **Prettier was reformatting `pixi.lock`** — a generated file. Rewriting it is pointless
  because Pixi regenerates it, and the check would fail again next time.

**The general lesson**, worth remembering because it will happen again: tools install large
folders full of *other people's code* (`node_modules/` for JavaScript, `.pixi/` for Pixi).
Linters must be told to ignore them, every time a new linter is added.

### 4. The customer monorepo template was added

When a customer runs `freelunch init my-company`, they should get a repository with a
standard folder layout. That layout now exists as a template in `templates/monorepo/`:

```
templates/monorepo/
├── platform/freelunch.yaml     ← platform config, owned by Platform Engineers
├── products/
│   └── example_product/
│       ├── services/           ← always-on workloads
│       └── workflows/          ← placeholder; not supported in the Demo
└── .github/workflows/          ← CI/CD pipeline
```

This is the structure defined in the roadmap, section 1.1.

Two details:

- The folder is called **`example_product`**. It is a placeholder that the `init` command
  will rename to the customer's real product name. We use a normal name rather than
  something like `<product-name>` because angle brackets are **illegal in Windows
  filenames** and would break the template for Windows users.
- **`freelunch init` does not exist yet.** These are the folders it will one day copy. The
  CLI has `install`, `uninstall`, `status` and `version`; scaffolding a monorepo is the
  half of roadmap 1.1 that is still outstanding.

### 5. All source code moved into `src/`

This is the largest change, and the reason most files in this PR look "moved".

**Before** — Go and TypeScript were mixed together at the top level:

```
freelunch-ide/
├── cmd/  internal/  go.mod       ← Go
├── ide/  package.json  tsconfig.json  ← TypeScript
└── node_modules/                  ← 94 MB of dependencies, at the top level
```

**After** — each language is fully self-contained:

```
freelunch-ide/
├── src/
│   ├── cli/          ← ALL Go code (go.mod, cmd/, internal/)
│   └── ide_src/      ← ALL TypeScript code (package.json, node_modules/, core/)
├── templates/        ← the customer monorepo template
├── docs/             ← specifications
├── Taskfile.yml      ← what every command actually runs
├── pixi.toml         ← tool versions
└── .golangci.yml  .goreleaser.yaml  .prettierrc.json
```

The top level now contains **no source code at all** — only configuration.

#### One consequence you should know about

Go identifies a project by a "module path", which must match the folder it lives in. Because
the Go code moved to `src/cli/`, its module path changed:

| | |
|---|---|
| Before | `github.com/Freelunch-AI/freelunch-ide` |
| After | `github.com/Freelunch-AI/freelunch-ide/src/cli` |

Every `import` line in the Go code was updated to match. **If you have older code or notes
referring to the old path, they need updating.**

This also affects release tags. Go requires a submodule's tags to carry the folder prefix —
`src/cli/v1.2.3` rather than plain `v1.2.3`. `branching_strategy.md` still documents the
plain form, so **that document and this layout currently disagree**. It is an open decision,
not something you should work around on your own.

#### Why the configuration files stayed at the top

`.golangci.yml`, `.prettierrc.json`, and `.prettierignore` did **not** move into `src/`.
They configure tools that check the *whole repository* — Prettier formats the YAML files in
`templates/` and `.github/`, which have nothing to do with the IDE frontend. Moving them
would have quietly stopped those files from being checked, which is the same class of bug as
item 3 above.

---

## Part 2 — Setting up your computer

### Which operating system?

| Your OS | What to do |
|---|---|
| **Linux** (Ubuntu, Fedora, Debian…) | Works directly. Go to [Linux setup](#linux-setup). |
| **Windows** | You must use **WSL2**. Go to [Windows setup](#windows-setup). |
| **macOS** | Works directly; follow the Linux steps. |

**Windows cannot be used directly, and this is not optional.** `pixi.toml` declares support
for four platforms only:

```
platforms = ["linux-64", "linux-aarch64", "osx-arm64"]
```

There is no `win-64`. Pixi will refuse to install on native Windows. (`osx-64` was also
dropped — nobody on the team uses an Intel Mac, and it held several packages far below the
versions the other platforms can resolve.) This matches
`tech_stack.md`, which states the development OS is **Linux or WSL2**.

WSL2 is not a workaround or a downgrade — it is real Linux running inside Windows, and it is
the supported path.

---

### Windows setup

#### Step 1 — Install WSL2

Open **PowerShell as Administrator** (right-click the Start button → "Terminal
(Administrator)") and run:

```powershell
wsl --install -d Ubuntu
```

Restart your computer when prompted. On first launch, Ubuntu asks you to create a username
and password. This password is for Linux only — it is not your Windows password, and it will
not show characters as you type. That is normal.

#### Step 2 — Open Ubuntu

Search for "Ubuntu" in the Start menu. Everything from here on happens **inside that Ubuntu
window**, not in PowerShell or Command Prompt.

#### Step 3 — Continue with the Linux steps

Follow [Linux setup](#linux-setup) below. It all applies inside Ubuntu.

> **Important:** put the project inside the Linux home folder (`~/`), **not** in
> `/mnt/c/Users/...`. Working across the Windows/Linux filesystem boundary is very slow and
> causes file-permission problems. Clone into `~/projects/` or similar.

---

### Linux setup

#### Step 1 — Install prerequisites

```bash
sudo apt update && sudo apt install -y curl git
```

(On Fedora use `sudo dnf install curl git`.)

#### Step 2 — Install Pixi

```bash
curl -fsSL https://pixi.sh/install.sh | sh
```

Then **close and reopen your terminal** so the `pixi` command is found. Verify:

```bash
pixi --version
```

You should see a version number. If you get "command not found", the terminal was not
restarted.

#### Step 3 — Get the code

```bash
mkdir -p ~/projects && cd ~/projects
git clone https://github.com/Freelunch-AI/freelunch-ide.git
cd freelunch-ide
```

#### Step 4 — Install the toolchain

```bash
pixi install
```

This downloads Go, Node, and the linters at their pinned versions into a local `.pixi/`
folder. It takes a few minutes the first time and does not touch the rest of your system.

#### Step 5 — Install project dependencies

```bash
pixi run setup
```

This downloads the Go modules and the npm packages, and installs the pinned **k3d** and
**kubectl** into `~/.freelunch/bin`. They are not placed on your `PATH` and not installed
system-wide, so they cannot collide with tools you already use.

It deliberately does **not** download the k3s image bundle, which is about 220MB and only
needed to create a cluster with no network. If you want that, run
`pixi run task setup:airgap-images` once — see
[local-cluster.md](local-cluster.md) for how it works and how to verify it honestly.

#### Step 6 — Verify everything works

```bash
pixi run check
```

This runs every formatter, linter, test, and build. **It should finish with no errors.** If
it does, your machine is correctly configured.

---

## Part 3 — Daily commands

Always run commands through `pixi run`. That is what guarantees you use the pinned tool
versions rather than whatever happens to be installed on your machine.

| Command | What it does |
|---|---|
| `pixi run check` | **Run this before every push.** Formatting, linting, tests, build — everything. |
| `pixi run build` | Build the CLI into `bin/freelunch` and compile the TypeScript. |
| `pixi run test` | Run the Go and TypeScript tests. |
| `pixi run fmt` | Auto-format all code. Run this if `check` complains about formatting. |
| `pixi run lint` | Run the linters only. |
| `pixi run setup` | Re-install dependencies (after someone adds a new one). |
| `pixi run task cluster:up` | Start the local Kubernetes cluster. See [local-cluster.md](local-cluster.md). |
| `pixi run task cluster:down` | Delete the local cluster. |
| `pixi run task setup:airgap-images` | Cache the k3s images so the cluster can be created offline (~220MB, optional). |

Try the built CLI:

```bash
pixi run build
./bin/freelunch version
```

The CLI drives the same cluster as the tasks above, and it is the path a customer actually
gets — they have a `freelunch` binary and no pixi, no checkout and no Taskfile:

```bash
./bin/freelunch install     # create the local Demo cluster
./bin/freelunch status      # is it running, and which nodes
./bin/freelunch uninstall   # tear it down
```

Use whichever is at hand: the tasks print more while you are working on the cluster
itself, the CLI is what we ship. Both read the same pinned k3d and kubectl from
`~/.freelunch/bin`, and both mount the airgap image cache when it exists.

### Finer-grained commands

`Taskfile.yml` defines smaller tasks you can call individually — useful when you only want
to check one language:

```bash
pixi run task test:go      # Go tests only
pixi run task test:ts      # TypeScript tests only
pixi run task lint:go      # Go linter only
pixi run task --list       # show everything available
```

Note the `task` word: `pixi run check` works because it is a shortcut declared in
`pixi.toml`, but the smaller tasks need `pixi run task <name>`.

---

## Part 4 — Rules that are easy to break

These have each caused a real problem already.

1. **Never run `go build ./...` or `go test ./...` with `./...`.** Use the Taskfile
   commands. Some npm packages contain Go files, and the Go tool would try to compile them.
   *(After the `src/` restructure this is less dangerous than it was, but the Taskfile
   commands remain the correct way.)*

2. **Never install Go or Node yourself.** Pixi provides them. If you install your own, you
   will eventually hit a bug nobody else can reproduce.

3. **Always run `pixi run check` before pushing.** Nothing runs it automatically. There is a
   GitHub Action on the repository, but it is an **AI reviewer** — it does not run the tests.
   If you skip the check, broken code reaches the branch and nobody finds out for days. That
   has already happened once.

4. **Do not edit files inside `.pixi/` or `node_modules/`.** They are downloaded
   dependencies, are not committed, and any change is lost on the next install.

5. **The `.md` files in `docs/` are frozen.** They were agreed by the team and are not to be
   edited without asking Marcos or Bruno first.

6. **Use Conventional Commit messages** — `feat:`, `fix:`, `chore:`, `docs:`, `test:`. The
   release tooling reads these prefixes to generate changelogs.

7. **Run the Go linter from `src/cli`, never from the repository root.** From the root it
   prints a confident `0 issues.` while also emitting
   `typechecking error: ... directory prefix internal does not contain main module`. There
   is no Go module at the root, so it checked nothing — the clean result is meaningless.
   `pixi run check` does this correctly; the trap is only when you invoke
   `golangci-lint` by hand.

---

## Part 5 — Troubleshooting

**`pixi: command not found`**
The terminal was not restarted after installing Pixi. Close it and open a new one.

**Pixi says the platform is not supported**
You are on native Windows. Use WSL2 — see [Windows setup](#windows-setup).

**`npm ci` fails with "package.json and package-lock.json are not in sync"**
Run `pixi run task setup:ts` — or, from inside `src/ide_src/`, `npm install` once. This
happens when the lockfile is older than the workspace layout.

**`pixi run check` fails on formatting**
Run `pixi run fmt`, then run the check again. Formatting problems fix themselves.

**Everything is very slow (Windows)**
The project is probably in `/mnt/c/...`. Move it into the Linux home folder (`~/`) and clone
it again there.

**Go complains it cannot find a package after pulling**
The module path changed in this PR. Run `pixi run setup` to refresh dependencies.

---

## Where to look next

| File | What it holds |
|---|---|
| `Taskfile.yml` | The real definition of every build command |
| `pixi.toml` | Which tools and versions are provisioned |
| `docs/specs/demo-global-spec/roadmap.md` | The ordered feature specification — the source of truth for scope |
| `docs/specs/demo-global-spec/tech_stack.md` | The technology choices and the reasoning |
| `docs/contributing/local-cluster.md` | Running the local Kubernetes cluster |
| `src/cli/internal/` | The Go services, and the best examples to copy when adding one |

If any step here fails, ask before working around it. A broken setup usually means the
documentation is wrong, and it is worth fixing for the next person.

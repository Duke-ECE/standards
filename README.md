# Duke-ECE Engineering Standards

This repo is the single source of truth for **how we build services** in the
managed-agents platform. It is written for both humans and AI coding agents
(every service repo's `AGENTS.md` links here).

- `AGENTS.md` — the rules. Read this first.
- `templates/go-service/` — the canonical Go service skeleton. Copy it when
  starting a new service; do not hand-roll a new layout.

When a rule here conflicts with an older habit in a specific repo, this repo
wins — update the code or update the rule via PR, don't fork conventions.

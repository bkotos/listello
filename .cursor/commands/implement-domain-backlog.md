---
name: implement-domain-backlog
description: Implement one command/event pair from docs/domain-model-backlog.md, one PR per layer, waiting for merge between PRs.
---

Follow the project skill **implement-domain-backlog**.

Read `.cursor/skills/implement-domain-backlog/SKILL.md` and `docs/domain-model-backlog.md`.

Take **one** command/event pair (`###` heading). If I named one, use that; otherwise the first heading that still has a `- [ ]`. Walk every remaining layer checkbox under that heading. For each checkbox: TDD (red → stop for my review → green after I approve), mark it `[x]` in `docs/domain-model-backlog.md`, open **one PR**, wait until I merge, then start the next layer on a new branch. Do not put two layers in one PR. Do not start a different command in this invoke.

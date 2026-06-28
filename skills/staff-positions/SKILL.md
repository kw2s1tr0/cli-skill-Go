---
name: staff-positions
description: Use this skill when a user wants to search, list, sort, or find position IDs or position codes in the staff search CLI.
---

# Staff Positions CLI

Run with `staffcli` available on PATH:

```bash
staffcli positions
```

If the command fails with a missing token, unauthorized, or authentication error, do not run login for the user. Ask the user to run:

```bash
staffcli login --email USER_EMAIL --token-name TOKEN_NAME
```

Then retry the position command after the user confirms login is complete.

Optional sort flags:

```bash
staffcli positions --order-by id --order-direction asc
staffcli positions --order-by code --order-direction asc
staffcli positions --order-by name --order-direction desc
```

Allowed values:

- `--order-by`: `id`, `code`, `name`
- `--order-direction`: `asc`, `desc`

Use this command before employee searches when the user gives a position name such as `リーダー`; find the position `ID`, then pass it to `staffcli employees --position-id ID`.

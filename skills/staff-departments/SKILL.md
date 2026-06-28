---
name: staff-departments
description: Use this skill when a user wants to search, list, sort, or find department IDs or department codes in the staff search CLI.
---

# Staff Departments CLI

Run with `staffcli` available on PATH:

```bash
staffcli departments
```

If the command fails with a missing token, unauthorized, or authentication error, do not run login for the user. Ask the user to run:

```bash
staffcli login --email USER_EMAIL --token-name TOKEN_NAME
```

Then retry the department command after the user confirms login is complete.

Optional sort flags:

```bash
staffcli departments --order-by id --order-direction asc
staffcli departments --order-by code --order-direction asc
staffcli departments --order-by name --order-direction desc
```

Allowed values:

- `--order-by`: `id`, `code`, `name`
- `--order-direction`: `asc`, `desc`

Use this command before employee searches when the user gives a department name such as `営業部`; find the department `ID`, then pass it to `staffcli employees --department-id ID`.

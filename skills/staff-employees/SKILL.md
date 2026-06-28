---
name: staff-employees
description: Use this skill when a user wants to search staff employees by keyword, department, position, employment status, employee number, email, name, or related department or position.
---

# Staff Employees CLI

Run with `staffcli` available on PATH:

```bash
staffcli employees
```

If the command fails with a missing token, unauthorized, or authentication error, do not run login for the user. Ask the user to run:

```bash
staffcli login --email USER_EMAIL --token-name TOKEN_NAME
```

Then retry the employee command after the user confirms login is complete.

Search flags:

```bash
staffcli employees --keyword KEYWORD
staffcli employees --department-id DEPARTMENT_ID
staffcli employees --position-id POSITION_ID
staffcli employees --employment-status active
```

Flags can be combined:

```bash
staffcli employees --keyword "山田" --department-id 2 --position-id 1 --employment-status active
```

Allowed employment status values:

- `active`
- `leave`
- `retired`

If the user gives a department or position name instead of an ID, first run `staffcli departments` or `staffcli positions` to find the ID. Example for `営業部の人`: run `staffcli departments --order-by code --order-direction asc`, find `営業部`, then run `staffcli employees --department-id ID`.

The command prints employee number, name, kana, email, employment status, department, and position.

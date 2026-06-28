---
name: staff-me
description: Use this skill when a user wants to see the currently authenticated staff search API user, verify that the saved CLI token works, or inspect the current login identity.
---

# Staff Me CLI

Run with `staffcli` available on PATH:

```bash
staffcli me
```

This command uses the saved access token and prints the authenticated user's `ID`, `Name`, and `Email`.

If it fails with a token or authentication error, do not try to log in for the user. Ask the user to run:

```bash
staffcli login --email USER_EMAIL --token-name TOKEN_NAME
```

Then retry `staffcli me` after the user confirms login is complete.

---
name: staff-login
description: Use this skill when a user wants to log in to the staff search CLI, issue or refresh an API token, authenticate before searching staff data, or check the login command usage.
---

# Staff Login CLI

Do not run the login command on behalf of the user because it prompts for a password. Instead, tell the user to run it themselves with `staffcli` available on PATH:

```bash
staffcli login --email USER_EMAIL --token-name TOKEN_NAME
```

The CLI prompts for the password on stderr and stores the returned access token in the user's home directory. Do not pass passwords as flags and do not print credentials.

If `--email` or `--token-name` is omitted, the API may apply its own validation. Prefer providing both for normal use.

After the user says login is complete, you may verify authentication with:

```bash
staffcli me
```

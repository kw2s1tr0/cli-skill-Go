---
name: staff-logout
description: Use this skill when a user wants to log out of the staff search CLI, revoke the current API token, delete the saved access token, or end an authenticated CLI session.
---

# Staff Logout CLI

Run with `staffcli` available on PATH:

```bash
staffcli logout
```

This command revokes the current token through the API and deletes the locally saved access token.

After logout, search commands should fail until the user logs in again.

If logout fails because there is no saved token or the token is invalid, report that state; do not run login for the user.

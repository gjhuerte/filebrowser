# Troubleshooting

## Session Timeout

By default, user sessions expire after **2 hours**. If you're uploading large files over slower connections, you may need to increase this timeout to prevent sessions from expiring mid-upload.

You can configure the session timeout using the `tokenExpirationTime` setting:

```bash
# CLI flag
filebrowser --tokenExpirationTime 24h

# Environment variable
FB_TOKEN_EXPIRATION_TIME=24h

# Persist to config
filebrowser config set --tokenExpirationTime 24h
```

Valid duration formats: `"2h"`, `"30m"`, `"24h"`, or combinations like `"2h30m"`.

## Permission Denied

If you see permission errors when accessing or uploading files:

1. **Check file ownership** — File Browser runs as a specific system user. Make sure the served directory and its contents are readable (and writable, if uploads are needed) by that user.

2. **Docker PUID/PGID** — When running in Docker, set `PUID` and `PGID` environment variables to match the owner of the files on the host:
   ```bash
   docker run -e PUID=$(id -u) -e PGID=$(id -g) ...
   ```

3. **SELinux/AppArmor** — On systems with mandatory access controls, you may need to set appropriate labels on the served directory.

## Upload Failures

- **Large file uploads** — If uploads fail for large files, check your reverse proxy's upload size limit:
  - nginx: `client_max_body_size 10G;`
  - Caddy: no limit by default
  - Apache: `LimitRequestBody`

- **Timeout during upload** — Increase the session timeout (see above) and check your reverse proxy's timeout settings.

## Can't Log In

- **Default credentials** — The default username and password are both `admin`. Change the password after first login.

- **Forgot password** — Use the CLI to reset:
  ```bash
  filebrowser users update admin --password newpassword
  ```

- **Proxy authentication** — If using proxy-based auth, verify that the correct header (e.g., `X-Forwarded-User`) is being passed through. See [authentication.md](authentication.md).

- **Clear cookies** — Clear your browser's cookies for the File Browser domain and try again.

## Docker Issues

- **Container won't start** — Check logs with `docker logs <container-name>`.

- **Files not visible** — Ensure the volume mount is correct: `-v /path/on/host:/srv`. The path inside the container should match File Browser's `--root` setting (default: `/srv`).

- **Database errors** — Make sure the config/database volume is persistent: `-v /path/to/config:/config`. Without this, the database is recreated on each container restart.

## Reverse Proxy

- **WebSocket errors** — Ensure your reverse proxy passes WebSocket connections. For nginx:
  ```nginx
  location / {
      proxy_pass http://localhost:8080;
      proxy_set_header Upgrade $http_upgrade;
      proxy_set_header Connection "upgrade";
      proxy_set_header Host $host;
  }
  ```

- **Base URL / subpath** — If serving File Browser under a subpath (e.g., `/files`), set the `--baseURL` flag:
  ```bash
  filebrowser --baseURL /files
  ```

See [deployment.md](deployment.md) for full reverse proxy configuration examples.

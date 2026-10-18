## Self-Registration (Signup)

File Browser allows you to enable user self-registration (signup). This can be enabled via **Settings → Global Settings**, or with `filebrowser config set --signup`. Self-registered users inherit the configured **user defaults**, including the scope.

> [!WARNING]
>
> By default, the user scope is the server's root, so a self-registered user could read,
> modify, and delete every file File Browser serves. To prevent this, either:
>
> a. Enable `createUserDir` so each user gets their own directory; or
> b. If users are meant to share files, set the default scope to something other than the root.

## Fail2ban

File Browser has built-in brute-force protection: after **10 failed logins within 15 minutes** from one address, further attempts get `429 Too Many Requests` until the window passes. Password-protected share links have the same limit per address and link. The counters live in memory, so restarting File Browser clears them.

The client address comes from the `X-Real-IP` / `X-Forwarded-For` headers when present, so run File Browser behind a reverse proxy that sets them (Nginx Proxy Manager does). Exposed directly, a client could forge those headers to dodge the limit.

For longer bans, or to block an address at the firewall, add [fail2ban](https://github.com/fail2ban/fail2ban) on top. It works from File Browser's logs; see its [wiki](https://github.com/fail2ban/fail2ban/wiki).

### Filter Configuration

An example filter configuration targeted at matching File Browser's logs.

```ini
[INCLUDES]
before = common.conf

[Definition]
datepattern = `^%%Y\/%%m\/%%d %%H:%%M:%%S`
failregex   = `\/api\/login: (403|429) <HOST> *`
```

### Jail Configuration

An example jail configuration. You should fill it with the path of the logs of File Browser, as well as the port where it is running at.

```ini
[filebrowser]

enabled = true
port = [your_port]
filter = filebrowser
logpath = [your_log_path]
maxretry = 10
bantime = 10m
findtime = 10m
banaction = iptables-allports
banaction_allports = iptables-allports
```

#!/bin/sh
set -e

PUID="${PUID:-1000}"
PGID="${PGID:-1000}"

# Rewrite user/group IDs in passwd/group to match runtime PUID/PGID
sed -i "s/^user:x:[0-9]*:[0-9]*/user:x:${PUID}:${PGID}/" /etc/passwd
sed -i "s/^user:x:[0-9]*/user:x:${PGID}/" /etc/group

# Ensure volumes are owned by the runtime user
chown "$PUID:$PGID" /config /database /srv

# Seed the config on first start
if [ ! -f /config/settings.json ]; then
  cp /defaults/settings.json /config/settings.json
  chown "$PUID:$PGID" /config/settings.json
fi

# Use /config/settings.json unless a -c/--config flag was passed
for arg in "$@"; do
  case "$arg" in
    -c|--config|-c=*|--config=*)
      exec setuidgid user filebrowser "$@"
      ;;
  esac
done

exec setuidgid user filebrowser --config=/config/settings.json "$@"

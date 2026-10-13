#!/bin/sh
set -e

# Seed the config on first start
if [ ! -f /config/settings.json ]; then
  cp /defaults/settings.json /config/settings.json
fi

# Use /config/settings.json unless a -c/--config flag was passed
for arg in "$@"; do
  case "$arg" in
    -c|--config|-c=*|--config=*) exec filebrowser "$@" ;;
  esac
done

exec filebrowser --config=/config/settings.json "$@"

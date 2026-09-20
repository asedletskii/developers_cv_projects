#!/bin/sh
# Пишет рантайм-конфиг фронтенда (window.__CONFIG__) из переменных окружения при каждом старте контейнера,
# поэтому настройку можно поменять в .env и просто перезапустить контейнер — без пересборки образа.
set -eu

CONFIG_PATH="${CONFIG_PATH:-/srv/config.js}"

# В JS попадает только true/false: любое другое значение переменной сюда не доходит
case "${RESPECT_DO_NOT_TRACK:-true}" in
  false | FALSE | False | 0 | no) respect_do_not_track=false ;;
  *) respect_do_not_track=true ;;
esac

cat > "$CONFIG_PATH" <<EOF
window.__CONFIG__ = { respectDoNotTrack: ${respect_do_not_track} }
EOF

exec "$@"

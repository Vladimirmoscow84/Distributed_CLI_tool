
# Порты для узлов
PORTS=(8081 8082 8083)

# Паттерн и файл
PATTERN="Spartak"
FILE="data/example.txt"
QUORUM=2

# Запуск узлов в фоне
PIDS=()
for PORT in "${PORTS[@]}"; do

  PEERS=""
  for P in "${PORTS[@]}"; do
    if [ "$P" != "$PORT" ]; then
      if [ -z "$PEERS" ]; then
        PEERS="localhost:$P"
      else
        PEERS="$PEERS,localhost:$P"
      fi
    fi
  done

  # Запуск узла
  ./mygrep -pattern "$PATTERN" -file "$FILE" -cluster -quorum $QUORUM -peers $PEERS -port $PORT &
  PIDS+=($!)
  echo "Запуск узла на порту $PORT..."
done

echo "Кластер запущен. Логи узлов:"
for PID in "${PIDS[@]}"; do
  wait $PID
done

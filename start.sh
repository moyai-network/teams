until ci & wait; do
    echo "moyai crashed with exit code $?. Respawning.." >&2
    sleep 10
done

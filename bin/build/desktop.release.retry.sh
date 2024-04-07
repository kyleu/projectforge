#!/bin/bash
# Content managed by Project Forge, see [projectforge.md] for details.

function handle_sigint() {
    echo "Received SIGINT (Ctrl-C), exiting."
    exit 1
}

trap 'handle_sigint' SIGINT

MAX_RETRIES=5
attempt=1

while [ $attempt -le $MAX_RETRIES ]; do
    echo "Attempt $attempt of $MAX_RETRIES: $CMD"
    ./desktop.release.sh $1

    if [ $? -eq 0 ]; then
        echo "Desktop build succeeded."
        exit 0
    else
        echo "Desktop build failed with exit code $?."
    fi

    ((attempt++))
done

echo "Desktop build failed after $MAX_RETRIES attempts."
exit 1

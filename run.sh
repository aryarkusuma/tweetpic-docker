#!/bin/bash

# Start Chrome
/usr/bin/google-chrome \
    --no-sandbox \
    --disable-gpu \
    --headless \
    --disable-dev-shm-usage \
    --remote-debugging-address=0.0.0.0 \
    --remote-debugging-port=9222 \
    --user-data-dir=/data &


# Run the compiled go binary application
/project/docker-scraper



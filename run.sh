#!/bin/bash

# Start Chrome in the background
/usr/bin/google-chrome \
    --no-sandbox \
    --disable-gpu \
    --headless \
    --disable-dev-shm-usage \
    --remote-debugging-address=0.0.0.0 \
    --remote-debugging-port=9222 \
    --user-data-dir=/data &


# Run the dockerscraper application
/project/docker-scraper

# Wait for both processes to finish
wait

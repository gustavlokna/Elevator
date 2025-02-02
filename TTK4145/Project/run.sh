#!/bin/bash

ID=$1  # Get ID from argument

# Check if ID is provided, otherwise exit with usage instructions
if [ -z "$ID" ]; then
    echo "Usage: ./run.sh <ID>"
    exit 1
fi

trap 'echo "Ignoring Ctrl+C...";' SIGINT  # Prevents manual termination with ^C

while true; do
    echo "Building the project..."
    go build -o elevator main.go || { echo "Build failed. Retrying..."; sleep 1; continue; }

    echo "Starting elevator program with ID=$ID..."
    ./elevator -id=$ID

    echo "Program crashed or terminal closed. Restarting in a new window..."
    
    sleep 1  # Prevents rapid restart loop

    # Open a new terminal and restart the script with the same ID
    gnome-terminal -- bash -c "cd $(pwd); ./run.sh $ID; exec bash"

    exit  # Exit the current instance so the new terminal takes over
done

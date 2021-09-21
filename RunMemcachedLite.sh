#!/bin/bash
# If this file won't execute, remember to change it's permissions
# Start server
go run server.go &
# Run client test cases
go run client.go
# Remove client.go storage files
rm -f *.ryoost
# Kill the server process
pkill -f go

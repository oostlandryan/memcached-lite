#!/bin/bash
# If this file won't execute, remember to change it's permissions
# Start server
go run server.go -port=9889 &
# Run client test cases
go run client.go -server=localhost:9889
# Remove client.go storage files
rm -f *.ryoost
# Kill the server process
pkill -f go

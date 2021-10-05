#!/bin/bash
# select gcloud project
PROJECTNAME="cloud-computing-327315"
gcloud config set project $PROJECTNAME
# Create server and client instances
gcloud compute instances create ryoost-server --image-family=ubuntu-1804-lts --image-project=ubuntu-os-cloud --zone=us-central1-a
gcloud compute instances create ryoost-client --image-family=ubuntu-1804-lts --image-project=ubuntu-os-cloud --zone=us-central1-a
#gcloud compute config-ssh
# get server internal ip address
SERVERADDRESS=$(gcloud compute instances describe ryoost-server --format='get(networkInterfaces[0].networkIP)')

# ssh into server instance, clone memcached-lite repo, move to repo, install Go, run server
gcloud compute ssh ryoost-server --zone=us-central1-a --command="git clone git://github.com/oostlandryan/memcached-lite.git
cd memcached-lite
yes | sudo apt install golang-go
go run server.go -port=9889
exit" &

# # ssh into client instance, clone memcached-lite repo, move to repo, install Go, run client
gcloud compute ssh ryoost-client --zone=us-central1-a --command="git clone git://github.com/oostlandryan/memcached-lite.git
cd memcached-lite
yes | sudo apt install golang-go
go run client.go -server=$SERVERADDRESS:9889
exit"

# Delete server and client instances
yes | gcloud compute instances delete ryoost-server
yes | gcloud compute instances delete ryoost-client
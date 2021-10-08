#!/bin/bash
# select gcloud project
PROJECTNAME="cloud-computing-327315"
gcloud config set project $PROJECTNAME

# Configure Network
gcloud compute networks create ryoost-network --subnet-mode=auto --bgp-routing-mode=regional --mtu=1460
gcloud compute firewall-rules create ryoost-tcp --network ryoost-network --allow tcp:22,tcp:3389,icmp,tcp:9889
sleep 5

gcloud beta compute addresses create ryoost-reserved-addresses --global --prefix-length=24 --network=ryoost-network --purpose=vpc_peering
gcloud services vpc-peerings connect --service=servicenetworking.googleapis.com --ranges=ryoost-reserved-addresses --network=ryoost-network --project=$PROJECTNAME


# Create server and client instances
gcloud compute instances create ryoost-server --network=ryoost-network --image-family=ubuntu-1804-lts --image-project=ubuntu-os-cloud --zone=us-central1-a
gcloud compute instances create ryoost-client --network=ryoost-network --image-family=ubuntu-1804-lts --image-project=ubuntu-os-cloud --zone=us-central1-a
# Wait a bit to make sure the instances are actually up and running
sleep 10

# Make sure ssh keys are updated (sometimes this fixes issues, sometimes it does not)
gcloud compute config-ssh
# Make sure the ssh keys have propogated
sleep 5

# get server internal ip address
SERVERADDRESS=$(gcloud compute instances describe ryoost-server --format='get(networkInterfaces[0].networkIP)')

# ssh into server instance, clone memcached-lite repo, move to repo, install Go, run server
gcloud compute ssh ryoost-server --zone=us-central1-a --command="git clone git://github.com/oostlandryan/memcached-lite.git
cd memcached-lite
yes | sudo apt install golang-go
go run server.go -port=9889" &
# Since the previous command is run in the background, we don't know if the server is ready when the next command runs, so we'll wait a bit
sleep 10

# # ssh into client instance, clone memcached-lite repo, move to repo, install Go, run client
gcloud compute ssh ryoost-client --zone=us-central1-a --command="git clone git://github.com/oostlandryan/memcached-lite.git
cd memcached-lite
yes | sudo apt install golang-go
go run client.go -server=$SERVERADDRESS:9889"

# Test Google's Memorystore
gcloud memcache instances create ryoost-memcache --node-count=1 --node-cpu=1 --node-memory=1GB --region=us-central1

# Delete server and client instances
yes | gcloud compute instances delete ryoost-server
yes | gcloud compute instances delete ryoost-client
yes | gcloud memcache instances delete ryoost-memcache --region=us-central1
yes | gcloud compute firewall-rules delete ryoost-tcp
yes | gcloud services vpc-peerings delete --service=servicenetworking.googleapis.com --network=ryoost-network
yes | gcloud beta compute addresses delete ryoost-reserved-addresses --global
sleep 5
yes | gcloud compute networks delete ryoost-network
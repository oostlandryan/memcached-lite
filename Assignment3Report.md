# Memcache-Lite
Ryan Oostland  
ryoost
## Running the Server and Client on Google Compute Engine
A bash script, [gcpdeployment.sh](gcpdeployment.sh), starts the server and client instances in the PROJECTNAME project using gcloud. It then uses `gcloud compute ssh` to access the instances, clone the git repo containing my memcached-lite code, install Go, and finally run the programs. Memcached-lite outputs are wrapped in dashes and should be easy to read in the terminal. After the programs have ran, the server and client instances are deleted.
## Assumptions
The [gcpdeployment.sh](gcpdeployment.sh) script assumes that gcloud is installed and that the project is configured with Google's default network settings. It also assumes that the Memorystore API and Service Networking API are enabled for the project.
.
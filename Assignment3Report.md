# Memcache-Lite
Ryan Oostland  
ryoost
## Running the Server and Client on Google Compute Engine
A bash script, [gcpdeployment.sh](gcpdeployment.sh), starts the server and client instances in the PROJECTNAME project using gcloud. It then uses `gcloud compute ssh` to access the instances, clone the git repo containing my memcached-lite code, install Go, and finally run the programs. Memcached-lite outputs are wrapped in dashes and should be easy to read in the terminal. After the programs have run, the server and client instances are deleted.
## Performance
Memcached-lite's performance was slightly better on the gcp VMs than it was locally. I was able to consitently have 1000 concurrent connections, which was did not work as well on my own machine. All other tests performed identically to how they did on my local machine.
## Comparison with Memorystore's Memcached Service
Unfortunately, I was unable to run memorystore's memcached service. Despite following the [online instructions](https://cloud.google.com/memorystore/docs/memcached/establishing-connection) for establishing a private services access connection, memcached failed to start citing that private service access is not enabled. I will be out of town when the assignment is due, so I'm needing to finish it early and have not left myself with enough time to solve this issue. I assumed that using Google's memcached service would be easier than setting up my own, but was clearly mistaken.
## Costs
A compute instance using the standard E2 machine with 2 vCPUs costs $0.067006 per hour to run.  
Memorystore's Memcached service with 1 node located in Iowa, 1 vCPU, and 1GB would cost $0.0544 per hour.  
Network traffic on the internal network, which my instances use to communicte with each other, is free.  


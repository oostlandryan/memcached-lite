# Memcache-Lite
Ryan Oostland  
ryoost
## Running the Server and Client
A bash script, [RunMemcachedLite.sh](RunMemcachedLite.sh), has been included that starts server.go in the background, runs the test cases in client.go, and then deletes the files that were created and kills the server process. It is not currently compatible with off-the-shelf memcached clients, but starting the server and manually sending it commands through something like netcat works fine.
## Server Design
I chose to implement this in Go because it's exactly the type of application Google designed it for. Go handles concurrency through go-routines, which is one of the main features of the language. The entire Go program is run in one process, but go-routines are essentially lightweight processes that Go's scheduler runs concurrently within the process. The server starts listening for TCP connections on port 9887 in an infinite loop, and whenever a new connection comes in it spawns a new go-routine to handle it. The data from the TCP connection is parsed and then calls appropriate functions depending on if it is a get or set command. To store the keys and values on disk, I simply save a new file whose name is the key and contents the value. I chose this approach over storing all of the keys and values in a single file to avoid locking the file with a mutex, slowing down all of the go-routines. Additionally, this method prevents one connection from destroying all of the data if something where to go wrong.

## Performance and Testing
Go-routines are quite efficient and it's possible for one program to spawn hundreds of thousands of them. Since TCP connections are treated as files, the limiting factor in my design is the maximum number of open file descriptors in a single process which is 1024 on most Linux installations. On my local machine I have been able to sustain 1000 concurrent connections, but I've ran into issues trying to repeatedly close and reopen that many. To avoid issues during grading, the concurrency stress test in Client.go opens 500 concurrent connections, sets key value pairs, gets the keys, and then confirms the retrieved value is the same as it was set to.  

Since the key is used as a filename, the maximum key size is 255 bytes, which is just slightly better than memcached's 250 bytes. Unfortunately the maximum value size is not as good. It can consistently store a value as big as 4KB. I'm not entirely sure what causes the error, but I have been able to inconsistently store and retrieve values in the 5KB-10KB range.

## Future Improvements
Adding the standard memcached flags and extensions is something that I would still like to do, even if it won't improve my grade.  
Since concurrency is a large part of the project, I would like to get around the file descriptor limit by running the server concurrently in multiple processes. Go has some support for this, but it doesn't seem to be encouraged.
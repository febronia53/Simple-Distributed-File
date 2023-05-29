# Simple-Distributed-File-in-GO & SQL Database
- This is a simple simulation for google file system.
- There is a server, client and multiple chunks.
- And there are different files containing data stored on small chunks.
- chunks send data data to server, then server make a map on this data to count how many times every char appears in text sent from chunk, and after that stores it in database.
- when server recieves all data from all chunks it make groupby on it to collect them all.
- finally when client connect to server and ask data, server will send it from its database and it will be stored on client device.
- Note: all server, client and chunks should be connected on the same network.
- This code can run on multiple chunks, clients and only one server

# SDCC PROJECT - B3 - Distributed PageRank
This is an implementation of famous algorithm Page-Rank.
It works with Map-Reduce paradigm on different containers managed by Docker-Compose.

To run the programs you must have installed:
1) Go, in particular, the project written with SDK 1.19.3
2) Docker & Docker-Compose

There are three main parts:
1) Mapper, that contains the code to execute Map-job
2) Reducer, that contains the code to execute Reduce-job
3) Master, generate a random graph by using an algorithm similar to Albert-Barabasi and calls the
function offered by mapper and reducer to evaluate page rank.

Configuration json is an extremely important file. It contains all the constants and parameters used.
pagerank.go can launch automatically all the container and it re-generate the docker-compose yml.
Is also possible to launch manually the program by using : "docker-compose up --build" in the main root,
where there is docker-compose.yml.

Every time wants test different configuration set them in configuration.json, then launch pagerank.go.
If configuration doesn't change you can re-launch program from terminal with "docker-compose up".

-POSSIBLE PROBLEM AND HOW TO SOLVE-
--GRPC FILE ARE MISSING--
To fix that must be installed protoc on your device.
If grpc files are missing go in -> ./Mapper/mapper
Generate pb files with
  "protoc --go_out=. .\mapper.proto"
  "protoc --go-grpc_out=. .\mapper.proto"
Another directory called mapper will be created copy the files outside (into ./Mapper/mapper)
and delete its. Now there will be 5 files copy all of them in -> -./Master/mapper

Do the same things for reducer in -> ./Reducer/reducer



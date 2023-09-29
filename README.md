# SDCC PROJECT - B3 - Distributed PageRank
This is an implementation of famous algorithm Page-Rank.<br>
It works with Map-Reduce paradigm on different containers managed by Docker-Compose.<br>

## Install application
To install the application execute:<br>
git clone https://github.com/0334321-LSE/SDCC-Distribuited-Page-Rank.git

## Requirements
To run the programs you must have installed:<br>
- Go, in particular, the project is written with SDK 1.19.3
- Docker & Docker-Compose

## Structure
There are three main parts:<br>
- Mapper, that contains the code to execute Map-job
- Reducer, that contains the code to execute Reduce-job
- Master, generates a random graph by using a simple algorithm and calls the<br>
  function offered by mapper and reducer to evaluate page rank.<br>
  Those three parts communicates by using gRPC.

## Configuration
Configuration json is an extremely important file. <br>
It contains all the constants and parameters used like: <br>
- Number of mapper and reducer
- Seed and other parameters for graph generation
- Constants for pagerank algorithm
- Assigned port for Map and Reduce service

## Run and close the application
To run the application there are different ways. <br>
First thing, start docker daemon with "dockerd" or manually by launch docker compose program<br>
Then, according to needs:
- "pagerank.go" generate docker-compose.yml, so must be used when the number of container changes.<br>
  After the execution, workers will remain up, other parameters can be changed, and new execution can <br>
  be done by using  "docker-compose up app-master" <br>
- Is also possible to launch manually the program by using : "docker-compose up --build" in the same root of .yml file. <br>
  In this case is not possible to modify the number of containers so don't do that. <br>

To close all the container use : "docker-compose stop" <br>
Instead to close only specific one : "docker-compose stop container-name" <br> <br>
<b> IMPORTANT </b> <br>
Every time wants test different configuration in terms of container number and ports number, <br>
change parameters in configuration.json and then launch pagerank.go to create a different docker-compose.yml. <br>

## Application output
The application has a logging system and produces some output, all those things are in output directory. <br>
Version V3 of the program produce output and leave it in output directory. <br />
Version V4 is created to be executed on EC2-Instance, this one save the output into a zip file in a S3 bucket. <br />
In configuration.json there are parameters for bucket name and region, must change if you want use it on your EC2 instance. <br />

## Version
Version 4.1 works on EC2 instances and saves the output into S3 bucket <br />
Version 3.2 works also locally, doesn't save on S3 bucket. It can be used for debugging by running the main inside MasterLocal directory <br />

## Possible problem and how to solve
### Grpc file are missing
To fix that must be installed protoc on your device. <br>
If grpc files are missing go in  ' ./Mapper/mapper ' <br>
Generate pb files with: <br>
- "protoc --go_out=. .\mapper.proto"
    - "protoc --go-grpc_out=. .\mapper.proto"
      Another directory called mapper will be created, copy the files outside (into ./Mapper/mapper).<br>
      Do the same things for reducer in ' ./Reducer/reducer ' <br>



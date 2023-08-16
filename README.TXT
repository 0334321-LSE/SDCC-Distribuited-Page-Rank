This is an implementation of famous algorithm Page-Rank.
It works with Map-Reduce paradigm on different containers managed by Docker-Compose.

To run the programs you must have installed:
1) Go, in particular is written with SDK 1.19.3
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

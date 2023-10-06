# LOCAL MASTER
It can be used outside EC2 instances for testing or debugging. <br>
## How to run
1) Start docker compose
2) As usually set your parameters in /PageRank/configuration.json, in particular, set **locally** at true
3) Start pagerank.go to launch container and services and to copy the configuration file inside the folder.
4) After that, you can launch main.go inside MasterLocale and edit MasterLocale/constants/configuration about graph size <br>
To change the container number, restart pagerank.go, as usually.
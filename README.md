# NS3DockerEmulator
 Network Emulator based on NS3 and Docker
 More details Read the docs: https://chepeftw.github.io/NS3DockerEmulator/


 - Dowload project
 	- git clone -b golang-porting https://github.com/ptrsen/NS3DockerEmulator.git

 - Install needed things (Ns3, docker, etc) 
 	- sudo ./main -op=install
 
 - Create simple scenario -n 5 nodes , -s size 10 mts^2 , -ns node speed 5mts/s , -np node pause 1 sed
 	- sudo ./main -op=create -n=5 -s=10

     		- Its possigle to enter every node env using (do pings , etc) 
			 sudo docker exec -it emu1  bin/sh
			- sudo docker exec -it emu2  bin/sh

 - Start Ns3 simulation ( ns3/tap-wifi-virtual-machine.cc) 
 	- sudo ./main -op=ns3 -n=5 -s=10


 - Destroy everthing 
 	- sudo ./main -op=destroy

 - clean File volumes ( conf, logs ) for each node 
 - local path: container/log, container/cont , Node path: /app/log, /app/conf 
	- sudo ./main -op=clean




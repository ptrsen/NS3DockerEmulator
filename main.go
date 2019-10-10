package main

import (
	"container/list"
	"context"
	"flag"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/homedir"
	"github.com/op/go-logging"
	"github.com/ptrsen/NS3DockerEmulator/tools/cmd"
	"github.com/ptrsen/NS3DockerEmulator/tools/docker"
	"github.com/ptrsen/NS3DockerEmulator/tools/net"
	"github.com/ptrsen/NS3DockerEmulator/tools/ns3"
	"os"
	"strings"
	"time"
)



/************************************************  VARIABLES  ***********************************************************/
/************************************************************************************************************************/

/**********************************************************************************
*	Logger Configuration
***********************************************************************************/

var logFilePath = "main.log"
var log = logging.MustGetLogger("main")
var format = logging.MustStringFormatter(
	`%{time:2006/01/02 15:04:05.999999} %{shortfile} ▶ %{level:.4s} %{id:03x} %{pid}**** < %{message} >`,
)

/***********************************************************************************/


/**********************************************************************************
*	Global Path Variables
***********************************************************************************/

var homePath string
var projectPath string
var ns3Path string

/***********************************************************************************/



/**********************************************************************************
*	Docker Configuration
***********************************************************************************/

var dockerFileName = "Dockerfile"
var dockerFilePath = "docker/minimal"
var baseContainerNameMin = "myminimalbox"


var dockerImagebasename = "emu"
var tapbasename = "tap-"
var bridgebasename = "br-"

var nameList  = list.New()

// Containers Logs
var logsLocalDirectory = "container/log"
var logsContainerDirectory = "/app/log"

// Containers Configuration
var confLocalDirectory = "container/conf"
var confContainerDirectory = "/app/conf"

/***********************************************************************************/


/**********************************************************************************
*	Ns3 Configuration
***********************************************************************************/

var ns3ModuleFileName = "tap-wifi-virtual-machine.cc"

var scenarioSize = 300
var numberOfNodes = 0
var nodeSpeed = 5
var nodePause = 1


/***********************************************************************************/





/************************************************  MAIN  ***********************************************************/
/************************************************************************************************************************/

func main() {

	/**********************************************************************************
	*	Backend Logger Configuration
	***********************************************************************************/

	logFile, err := os.OpenFile(logFilePath, os.O_RDWR | os.O_CREATE|os.O_APPEND, 0666)
	CheckError(err,"Error open log file ")

	defer func() {
		err := logFile.Close()
		CheckError(err,"Error close log file ")
	}()

	backend := logging.NewLogBackend(logFile, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	logging.SetBackend(backendFormatter)

	/**********************************************************************************/


	/**********************************************************************************
	*	Initialization
	***********************************************************************************/

	// Setting Paths
	homePath = homedir.Get()
	projectPath, err = os.Getwd()
	CheckError(err,"Project path not found ")

	ns3Path = homePath + "/Ns3/bake/source"


    // Create context
	ctx := context.Background()


	// Create Docker client
	cli, err := client.NewClientWithOpts(client.FromEnv)
	CheckError(err,"Fail to create Docker client, Check docker installation ")
	cli.NegotiateAPIVersion(ctx)

	/***********************************************************************************/



	/**********************************************************************************
	*	Command-Line Flags Parser
	***********************************************************************************/

	// Operation  -op [install, create, ns3, emulation, destroy, clean, none]
	operationPtr := flag.String("op", "none", " operation to do  [string]  \n" +
		"     install - Install docker and ns3 to run emulator \n" +
		"     create -  Setup everything for network emulation \n" +
		"     ns3 - Run ns3 network emulation \n" +
		"     emulation -Iterates several emulations (not implemented yet) \n" +
		"     destroy - Destroy everything \n" +
		"     clean - Clean nodes conf and logs directories \n" +
		"     ")


	// Ns3 Scenario Size -s 300
	sizePtr := flag.Int("s", 300, "size of the network scenario [int mts^2] - ")
	// Ns3 Number of Nodes -n 0
	numNodesPtr := flag.Int("n", 0, "number of nodes [int] - ")
	// Ns3 Nodes Speed -ns 5
	nodeSpeedPtr := flag.Int("ns", 5, "nodes speed [int mts/s] - ")
	// Ns3 Nodes Pause -np 1
	nodePausePtr := flag.Int("np", 1, "nodes pause [int s] - ")


	flag.Parse()

	// Getting values
	numberOfNodes = *numNodesPtr
	scenarioSize = *sizePtr
	nodeSpeed = *nodeSpeedPtr
	nodePause = *nodePausePtr

	// Generating name list for containers emu1...emuX
	for i := 0; i < numberOfNodes; i++ {
		nameList.PushBack( dockerImagebasename + fmt.Sprintf("%v", i+1) )
	}

	/***********************************************************************************/




	/**********************************************************************************
	*	Program Steps
	***********************************************************************************/

	// commands
	// sudo ./main -op=install
	// sudo ./main -op=create -n=2 -s=10     // nodes 2, Scenario Size 10
	// sudo ./main -op=ns3 -n=2 -s=10
	// sudo ./main -op=destroy
	// sudo ./main -op=clean


	// sudo ./main -op=emulation -n=2 -s=10  // not yet

	switch operation := *operationPtr; operation {
	case "install":
		fmt.Println("-> Install Step ")
		log.Info("Install Step ")
		Install(ctx)

	case "create":
		fmt.Println("-> Create Step ")
		log.Info("Create Step ")
		Create(ctx,cli)
		log.Info("All Setup , ready to run NS3 !!!")

	case "ns3":
		fmt.Println("-> Ns3 running in background ...")
		log.Info("ns3 running background Step...")
		Ns3Run(ctx)
		log.Info("ns3 done !!! ")

	case "emulation":
		fmt.Println("emulation Step ...")  //  next to work ...

	case "destroy":
		fmt.Println("-> Destroy Step")
		log.Info("Destroy Step")
		Destroy(ctx,cli)

	case "clean":
		fmt.Println("-> Clean Step ...")
		fmt.Println("Deleting conf and var directories ...")
		log.Info("Deleting conf and var directories ...")
		Clean(ctx)

	default:
		fmt.Println("run 'sudo ./main -h' for help ...")
	}

	ctx.Done()
	err = cli.Close()
	CheckError(err, "Error closing Docker Client")

}






/************************************************  FUNCTIONS  ***********************************************************/
/************************************************************************************************************************/


/**********************************************************************************
*	CheckError : Function just to check error
***********************************************************************************/

func CheckError(err error, msj string) {
	if err  != nil {
		fmt.Printf(msj + " %v\n", err)
		log.Error(msj, err)
		os.Exit(1)
	}
}

/***********************************************************************************/


/**********************************************************************************
*	Install : Function for Install Step
***********************************************************************************/

func Install(ctx context.Context){

	// Install Docker
	err := cmd.ExecCommand(ctx, "/bin" ,"bash", projectPath + "/install/docker-install.sh")
	CheckError(err, "Unable to install Docker")
	fmt.Println("Docker installed ")
	log.Info("Docker installed")

	// Install drivers
	err = cmd.ExecCommand(ctx, "/bin" ,"bash", projectPath + "/install/driver-install.sh", projectPath)
	CheckError(err, "Unable to install L2Driver")
	fmt.Println("Drivers installed ")
	log.Info("Drivers installed")

	// Install NS3
	err = cmd.ExecCommand(ctx, "/bin" ,"bash","-c","source "+ projectPath + "/install/ns3-install.sh")
	CheckError(err, "Unable to install NS3")
	fmt.Println("NS3 installed ")
	log.Info("NS3 installed")

}


/**********************************************************************************
*	Create : Function for Create Step
***********************************************************************************/

func Create(ctx context.Context, cli client.APIClient){

	var err error
	var msj string

	// Check ns3 installation
	err, ns3Path = ns3.CheckNs3(ctx,ns3Path)
	CheckError(err,ns3Path)

   // Compile Ns3 Module in optimized mode
    err, msj = ns3.BuildModule(ctx, projectPath, ns3Path, ns3ModuleFileName)
	CheckError(err, msj)

    msj = "ns3/tap-vm file up to date! & " + "NS3 optimize build sucess!! | " + time.Now().Format("2006-01-02 15:04:05.0000")
    fmt.Println(msj)
    log.Info(msj)


	// Pull Dockerfile base image
	err, msj = docker.PullImage(ctx, cli, dockerFilePath, dockerFileName, baseContainerNameMin)

	msj = "Dockerfile base image built correctly !! "
	fmt.Println(msj)
	log.Info(msj)


	// Create local directory for containers logs
	if _, err := os.Stat( logsLocalDirectory ); os.IsNotExist(err) {
		err := os.MkdirAll( logsLocalDirectory ,0777)
		CheckError(err,"Unable to create container/log directory")
	}

	// Create local directory for containers configurations
	if _, err := os.Stat( confLocalDirectory ); os.IsNotExist(err) {
		err := os.MkdirAll( confLocalDirectory ,0777)
		CheckError(err,"Unable to create container/conf directory")
	}


	// Creating All
	for nodeName := nameList.Front(); nodeName != nil; nodeName = nodeName.Next(){
		nodeNameStr:= fmt.Sprint(nodeName.Value)

		// Create dir for each container in logs if not exists
		if _, err := os.Stat( logsLocalDirectory + "/" +  nodeNameStr); os.IsNotExist(err) {
			err := os.MkdirAll( logsLocalDirectory + "/" + nodeNameStr,0777)
			CheckError(err,"Unable to create container/log/emuX directory")
		}

		// Create dir for each container in conf if not exists
		if _, err := os.Stat( confLocalDirectory + "/" +  nodeNameStr); os.IsNotExist(err) {
			err := os.MkdirAll( confLocalDirectory + "/" + nodeNameStr,0777)
			CheckError(err,"Unable to create container/conf/emuX directory")
		}

		// Container Volume Path [localDirectory, ContainerDirectory]
		logVolumeDirectory := [2]string{ projectPath + "/" + logsLocalDirectory, logsContainerDirectory}  // logs
		confVolumeDirectory := [2]string{ projectPath + "/" + confLocalDirectory, confContainerDirectory} // Conf

		// Create Docker Network
		err, msj := docker.CreateDockerNetwork(ctx, cli, bridgebasename + nodeNameStr)
		CheckError(err, msj)
		fmt.Println(msj)
		log.Info(msj)

		// Create Docker Container
		err, msj = docker.CreateContainer(ctx, cli, nodeNameStr, baseContainerNameMin, bridgebasename + nodeNameStr, logVolumeDirectory, confVolumeDirectory)
		CheckError(err, msj)
		msj = "Created Container - "+ msj
		fmt.Println(msj)
		log.Info(msj)

		// Create Tap Device and Attach to Docker Network
        err, msj = net.CreateTAP(ctx, tapbasename + nodeNameStr, bridgebasename + nodeNameStr)
		CheckError(err, msj)
		msj = "Created TAP Device - "+ msj
		fmt.Println(msj)
		log.Info(msj)

	}

}
/***********************************************************************************/



/**********************************************************************************
*	Destroy : Function to Destroy Step
***********************************************************************************/

func Destroy (ctx context.Context, cli client.APIClient){

	// Getting a Containers list
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	CheckError(err,"unable to list  all docker containers")

	// Deleting all Containers , Networks and TAP interfaces
	for _, ctainer := range containers {
	  nameContainer := ctainer.Names[0][1:]
	  if  strings.Contains( nameContainer,"emu" ){

		  // Stop Container
		  err := cli.ContainerStop(ctx, nameContainer, nil)
		  CheckError(err,"Error stopping container")

		  // Remove Container
		  removeOptions :=  types.ContainerRemoveOptions{}  // Default
		  err = cli.ContainerRemove(ctx, nameContainer, removeOptions)
		  CheckError(err,"Error removing container")

		  // Deleting Tap device
		  err, msj := net.DeleteTAP(ctx, tapbasename + nameContainer)
		  CheckError(err, msj)

		  // Deleting Networks
		  err = cli.NetworkRemove(ctx, bridgebasename + nameContainer)
		  CheckError(err,"Unable to delete network ..." )

	  }

	}

	// Deleting base Image
	removeImageOptions :=  types.ImageRemoveOptions{}  // Default
    _, err = cli.ImageRemove(ctx, baseContainerNameMin, removeImageOptions )
	CheckError(err,"Error removing base image")

	fmt.Println("All deleted !! ")
	log.Info("All deleted !! ")


	/*
	// Defining Filters for containers and networks
	filter := filters.NewArgs()   // default filter
	//filters.Add("dangling", "true")
	//filters.Add("dangling", "false")

	//filters.Add("dangling", "true")
	//filters.Add("until", "2016-12-15T14:00")

	// Deleting all Containers
	respConPrune, err := cli.ContainersPrune(ctx,filter)
	CheckError(err,"unable to prune Containers ..." )
	fmt.Println(respConPrune.ContainersDeleted)
	fmt.Println("Containers Deleted !!!")
	log.Info("Containers Deleted !!!")


	// Deleting all Networks
	respNetPrune, err := cli.NetworksPrune(ctx,filter)
	CheckError(err,"unable to prune networks ..." )
	fmt.Println(respNetPrune.NetworksDeleted)
	fmt.Println("Networks Deleted !!!")
	log.Info("Networks Deleted !!!")

*/
}

/***********************************************************************************/



/**********************************************************************************
*	Clean : Just to delete Containers conf and var Directories
***********************************************************************************/

func Clean (ctx context.Context)  {

	err:= cmd.ExecCommand(ctx, "/bin" ,"./rm","-rf", projectPath + "/container")
	CheckError(err, "Unable to erase container logs and confs")
	fmt.Println("Container logs and confs files erased")
	log.Info("Container logs and confs files erased")

}

/***********************************************************************************/



/**********************************************************************************
*	NS3Run : Function to Run Ns3 scenario
***********************************************************************************/

func Ns3Run (ctx context.Context) {

	/// Check ns3 installation
	err, ns3Path := ns3.CheckNs3(ctx,ns3Path)
	CheckError(err,ns3Path)

	log.Info("About to start NS3 RUN")

	err, msj := ns3.RunBackground(ctx, ns3Path, scenarioSize, numberOfNodes, nodeSpeed, nodePause)
	CheckError(err,msj)

	log.Info("Finished running NS3 in the background | Date now: " + time.Now().Format("2006-01-02 15:04:05.0000"))

}

/***********************************************************************************/










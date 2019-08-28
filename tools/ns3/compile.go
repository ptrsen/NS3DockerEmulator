package ns3

import (
	"context"
	"fmt"
	"github.com/ptrsen/NS3DockerEmulator/tools/cmd"
	"time"
)



/*********************************************************************************
*	BuildNs3Module :
*			Function to build Ns3 Module .waf optimized
*			return error, output string
**********************************************************************************/

func BuildModule(ctx context.Context, projectDirectoryPath string ,ns3DirectoryPath string, ns3ModuleFilename string) (error, string) {

	// Copy module from  ns3 project folder to $NS3.../scratch
	err:= cmd.ExecCommand(ctx, "/bin", "./cp", projectDirectoryPath + "/ns3/" + ns3ModuleFilename, ns3DirectoryPath + "/scratch/tap-vm.cc")
	if err != nil { return err, "Error copying module " + ns3ModuleFilename + " to Ns3/scratch directory path" }

	// Build module with waf
	// ./waf build -j 1 -d optimized --disable-examples
	err = cmd.ExecCommand(ctx, ns3DirectoryPath , ns3DirectoryPath +"/./waf", "build", "-j", fmt.Sprintf("%v", 1), "-d", "optimized", "--disable-examples")
	if err != nil { return err, "Error Waf optimize build module " + ns3ModuleFilename }

	return err, ""  // Return ""  everything is good
}

/**********************************************************************************/




/*********************************************************************************
*	RunBackground :
*			Function to run Ns3 in background
*			return error, output string
**********************************************************************************/

func RunBackground(ctx context.Context, ns3Path string, scenarioSize int, numberOfNodes int, nodeSpeed int, nodePause int ) (er error, msj string){

	totalEmuTime := (5 * 60) * numberOfNodes
	fmt.Println("About to start NS3 RUN  with total emulation time of ", totalEmuTime)

	err := cmd.ExecCommand(ctx, ns3Path, ns3Path + "/./waf",
		"-j",fmt.Sprintf("%v",1),
		"--run",
		"scratch/tap-vm --NumNodes=" + fmt.Sprintf("%v",numberOfNodes) +
			" --TotalTime=" + fmt.Sprintf("%v",totalEmuTime) +
			" --TapBaseName=emu --SizeX=" + fmt.Sprintf("%v",scenarioSize) +
			" --SizeY=" + fmt.Sprintf("%v",scenarioSize) +
			" --MobilitySpeed=" + fmt.Sprintf("%v",nodeSpeed) +
			" --MobilityPause=" + fmt.Sprintf("%v",nodePause) )


	if err != nil { return err, "Error running NS3 scenario" }

	// missing proc.pid  search later -- add
	time.Sleep(10)

	fmt.Println("Finished running NS3 in the background | Date now: " + time.Now().Format("2006-01-02 15:04:05.0000"))

	return err, ""  // Return ""  everything is good

}

/**********************************************************************************/

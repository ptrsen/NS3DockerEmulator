package net

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/ptrsen/NS3DockerEmulator/tools/cmd"
	"os"
	"strconv"
	"time"
)

var dockerBridgebasename = "br-"
var linuxTapbasename = "tap-"
var sideA="side-int-"
var sideB="side-ext-"



/**********************************************************************************
*	CreateBridgeTAP:
*					Function to create Linux bridge br-emuX with
*					TAP interface tap-emuX for NS3
**********************************************************************************/

func CreateBridgeTAP(ctx context.Context, containerName string) (er error, msj string){

	/*
		- Create Bridge
	  		sudo ip link add br-emu1 type bridge
		- Create TUN/TAP interface
			sudo ip tuntap add mode tap tap-emu1
		- Set TUN/TAP interface to promisc mode
			sudo ip link set tap-emu1 promisc on
		- Active TUN/TAP interface
			sudo ip link set tap-emu1 up
		- Attach TUN/TAP to Bridge
			sudo ip link set tap-emu1 master br-emu1
		- Active Bridge
			sudo ip link set br-emu1 up
	*/

	err := cmd.ExecCommand(ctx,"/sbin","./ip", "link", "add", dockerBridgebasename + containerName,  "type", "bridge")
	if err != nil { return err, "Error creating bridge" }
	err = cmd.ExecCommand(ctx,"/sbin","./ip", "tuntap", "add", "mode", "tap", linuxTapbasename + containerName)
	if err != nil { return err, "Error creating tap" }
	err = cmd.ExecCommand(ctx,"/sbin","./ip", "link", "set", linuxTapbasename + containerName, "promisc", "on")
	if err != nil { return err, "Error setting tap promisc mode on" }
	err = cmd.ExecCommand(ctx,"/sbin","./ip", "link", "set", linuxTapbasename + containerName, "up")
	if err != nil { return err, "Error Active tap" }
	time.Sleep(1)

	err = cmd.ExecCommand(ctx,"/sbin","./ip", "link", "set", linuxTapbasename + containerName, "master", dockerBridgebasename + containerName)
	if err != nil { return err, "Error attach tap to bridge" }
	err = cmd.ExecCommand(ctx,"/sbin","./ip", "link", "set", dockerBridgebasename + containerName, "up")
	if err != nil { return err, "Error Active bridge" }
	time.Sleep(1)

	return err,""
}

/**********************************************************************************/


/**********************************************************************************
*	CreateBridgeTAP:
*					Function to create Linux bridge br-emuX with
*					TAP interface tap-emuX for NS3
**********************************************************************************/

func DeleteBridgeTAP(ctx context.Context, containerName string) (er error, msj string) {

	// De-Active TUN/TAP interface >> sudo ip link set tap-emu1 down
	// Delete TUN/TAP interface >> sudo ip tuntap del mode tap tap-emu1
	err:= cmd.ExecCommand(ctx, "/sbin", "./ip", "link", "set", linuxTapbasename+containerName, "up")
	if err != nil { return err, "Error turning off Tap device" }
	err = cmd.ExecCommand(ctx, "/sbin", "./ip", "tuntap", "del", "mode", "tap", linuxTapbasename+containerName)
	if err != nil { return err, "Error deleting Tap device" }

	// De-Active Bridge >> sudo ip link set br-emu1 down
	// Delete bridge >> sudo ip link del br-emu1 type bridge
	err = cmd.ExecCommand(ctx, "/sbin", "./ip", "link", "set", dockerBridgebasename+containerName, "down")
	if err != nil { return err, "Error turning off bridge" }
	err = cmd.ExecCommand(ctx, "/sbin", "./ip", "link", "del", dockerBridgebasename+containerName, "type", "bridge")
	if err != nil { return err, "Error deleting bridge" }

	return err, ""

}
/**********************************************************************************/



/**********************************************************************************
*	CreateVeth:
*					Function to create Linux bridge br-emuX with
*					TAP interface tap-emuX for NS3
**********************************************************************************/

func CreateVeth(ctx context.Context, containerName string, containerPid string, index int) (er error, msj string){

	Segment3 := strconv.Itoa( index/250 )
	Segment4 :=  strconv.Itoa((index%250)+1 )
	err, MAC := GenerateMAC()
	if err != nil { return err, "Error getting MAC adress" }

	sideA="side-int-"+ containerName
	sideB="side-ext-"+ containerName


	//  Delete container’s network namespace if already exists
	if _, err := os.Stat( "/var/run/netns/"+ containerPid);  os.IsExist(err) {
		err:= cmd.ExecCommand(ctx, "/bin" ,"./rm","-rf", "/var/run/netns/"+ containerPid)
		if err != nil { return err, "Unable to delete /var/run/netns/containerPid directory" }
	}

	//  Create symbolic link to manage container’s network namespace from netns
	// ln -s  /proc/$containerPid/ns/net /var/run/netns/$containerPid
	err = cmd.ExecCommand(ctx, "/bin" ,"./ln","-fs", "/proc/" + containerPid + "/ns/net", "/var/run/netns/"+ containerPid)
	if err != nil { return err, "Unable to create symbolic link for netns container’s network namespace" }


	//  Create veth pair[$VETH (sideA) <--> $PEER (SideB)] to container[$CONTAINER] (containername) in netns[$NETNS] (containerPid)
	// ip link delete sideA || true
	// ip link add sideA type veth peer name sideB
	err = cmd.ExecCommand(ctx, "/sbin", "./ip", "link", "add", sideA, "type", "veth",  "peer", "name", sideB)
	if err != nil { return err, "Error creating veth pair" }


	// Add sideA to bridge
	// ip link set sideA master br-emu1
	// ip link set sideA up
	err = cmd.ExecCommand(ctx, "/sbin", "./ip", "link", "set", sideA, "master", dockerBridgebasename+containerName)
	if err != nil { return err, "Error attaching veth side A to bridge" }
	err = cmd.ExecCommand(ctx, "/sbin", "./ip", "link", "set", sideA , "up")
	if err != nil { return err, "Error activate veth bridge side" }


	// Add sideB to containerPid netns and active

	// ip link set sideB netns containerPid
	err = cmd.ExecCommand(ctx, "/sbin", "./ip", "link", "set", sideB , "netns",  containerPid)
	if err != nil { return err, "Error attaching veth side B to container" }

	// ip netns exec containerPid ip link set dev sideB name eth0
	err = cmd.ExecCommand(ctx, "/sbin", "./ip", "netns","exec", containerPid, "ip", "link", "set", "dev", sideB, "name", "eth0")
	if err != nil { return err, "Error renaming veth side B inside container" }

	// ip netns exec containerPid ip link set eth0 address $MAC_ADDR
	err = cmd.ExecCommand(ctx, "/sbin", "./ip", "netns","exec", containerPid, "ip", "link", "set", "eth0", "address", MAC )
	if err != nil { return err, "Error renaming veth side B inside container" }

	// ip netns exec containerPid ip link set eth0 up
	err = cmd.ExecCommand(ctx, "/sbin", "./ip", "netns","exec", containerPid, "ip", "link", "set", "eth0", "up" )
	if err != nil { return err, "Error renaming veth side B inside container" }


	// ip netns exec containerPid ip addr add 10.12.$SEGMENT3.$SEGMENT4/16 dev eth
	err = cmd.ExecCommand(ctx, "/sbin", "./ip", "netns","exec", containerPid, "ip", "addr", "add", "10.12."+ Segment3 +"."+ Segment4 + "/16", "dev", "eth0" )
	if err != nil { return err, "Error renaming veth side B inside container" }


	return err,""
}
/**********************************************************************************/


/**********************************************************************************
*	DeleteVeth:
*					Function to delete veth pair
**********************************************************************************/
/*
func DeleteVeth(ctx context.Context, containerName string) (er error, msj string){

	sideA="side-int-"+ containerName

	// ip link delete sideA
	// ip link add sideA type veth peer name sideB
	err := cmd.ExecCommand(ctx, "/sbin", "./ip", "link", "delete", sideA)
	if err != nil { return err, "Error deleting veth pair" }


	return err,""
}
 */
/**********************************************************************************/




/**********************************************************************************
*	GenerateMAC:
*					Function to generate random MAC Address
**********************************************************************************/

func GenerateMAC()  (er error, mac string) {

	buf := make([]byte, 6)
	_, err := rand.Read(buf)
	if err != nil { return err , ""}
	// Set the local bit
	buf[0] = (buf[0] | 2) & 0xfe

	return err, fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", buf[0], buf[1], buf[2], buf[3], buf[4], buf[5])
}

/**********************************************************************************/

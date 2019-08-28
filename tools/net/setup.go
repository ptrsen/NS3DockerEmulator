package net

import (
	"context"
	"github.com/ptrsen/NS3DockerEmulator/tools/cmd"
	"time"
)

/**********************************************************************************
*	CreateTAP:
*				Function to create TAP interface tap-emuX used to connect with NS3
*    			and attach it to network bridge br-emuX
*				return error, output string
**********************************************************************************/

func CreateTAP(ctx context.Context, tapName string, bridgeName string) (er error, msj string){

	/*
			- Create TUN/TAP interface
				sudo ip tuntap add mode tap tap-emu1
			- Set TUN/TAP interface to promisc mode
				sudo ip link set tap-emu1 promisc on
			- Active TUN/TAP interface
				sudo ip link set tap-emu1 up

			- Attach TUN/TAP to Bridge
				sudo ip link set tap-emu1 master br-emu1
	*/

	err := cmd.ExecCommand(ctx,"/sbin","./ip", "tuntap", "add", "mode", "tap", tapName)
	if err != nil { return err, "Error creating tap" }
	err = cmd.ExecCommand(ctx,"/sbin","./ip", "link", "set", tapName, "promisc", "on")
	if err != nil { return err, "Error setting tap promisc mode on" }
	err = cmd.ExecCommand(ctx,"/sbin","./ip", "link", "set", tapName, "up")
	if err != nil { return err, "Error Active tap" }
	time.Sleep(1)

	err = cmd.ExecCommand(ctx,"/sbin","./ip", "link", "set", tapName, "master", bridgeName)
	if err != nil { return err, "Error attach tap to bridge" }
	time.Sleep(1)

	return err,""
}

/**********************************************************************************/



/**********************************************************************************
*	DeleteTAP:
*			  	Function to Delete TAP interface tap-emuX used to connect with NS3
*				return error, output string
**********************************************************************************/

func DeleteTAP(ctx context.Context, tapName string) (er error, msj string) {

	// De-Active TUN/TAP interface >> sudo ip link set tap-emu1 down
	// Delete TUN/TAP interface >> sudo ip tuntap del mode tap tap-emu1
	err:= cmd.ExecCommand(ctx, "/sbin", "./ip", "link", "set", tapName , "up")
	if err != nil { return err, "Error turning off Tap device" }
	err = cmd.ExecCommand(ctx, "/sbin", "./ip", "tuntap", "del", "mode", "tap", tapName )
	if err != nil { return err, "Error deleting Tap device" }

	return err, ""
}
/**********************************************************************************/

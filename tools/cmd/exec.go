package cmd

import (
	"context"
	"os"
	"os/exec"
)


/*********************************************************************************
*	ExecCommand :
*			Execute command from Go, Standard Output
*			returns error
**********************************************************************************/

func ExecCommand(ctx context.Context,dir,program string, args ...string) (er error) {

	cmd := exec.CommandContext(ctx, program, args...)

	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	return err
}

/**********************************************************************************/


/*********************************************************************************
*	ExecCommandOutput :
*			Execute command from Go
*			return error, output string
**********************************************************************************/

func ExecCommandOutput(ctx context.Context,dir,program string, args ...string) (er error,msj string) {

	cmd := exec.CommandContext(ctx, program, args...)

	cmd.Dir = dir
	outStr, err := cmd.Output()
	str := string(outStr)

	if err == nil {
		str = str[1 : len(str)-2]
	}

	return err, str
}

/**********************************************************************************/






package docker

import (
	"archive/tar"
	"bytes"
	"context"
	"encoding/binary"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"io"
	"io/ioutil"
	"net"
	"os"
	"strconv"
)


/*********************************************************************************
*	PullBaseImage :
*			Function to pull Docker Image from Dockerfile
*			return error, output string
**********************************************************************************/

func PullImage(ctx context.Context, cli client.APIClient, dockerFilePath string, dockerFileName string, baseContainerNameMin string) (er error,msj string){

	// Dockerfile to tarfile
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)

	defer func() {
		er = tw.Close()
		if er != nil { msj = "Error closing tarfile" }
	}()

	// Open Dockerfile
	dockerFile := dockerFileName
	dockerFileReader, err := os.Open(dockerFilePath + "/" + dockerFileName)
	if err != nil { return err, "Error opening Dockefile" }

	// Read Dockerfile
	readDockerFile, err := ioutil.ReadAll(dockerFileReader)
	if err != nil { return err, "Error reading Dockefile" }

	// Create Tarfile
	tarHeader := &tar.Header{
		Name: dockerFile,
		Size: int64(len(readDockerFile)),
	}
	err = tw.WriteHeader(tarHeader)
	if err != nil { return err, "Error writing tar header" }
	_, err = tw.Write(readDockerFile)
	if err != nil { return err, "Error writing tar body" }
	dockerFileTarReader := bytes.NewReader(buf.Bytes())


	// Docker image options
	imageOptions := types.ImageBuildOptions{
		Tags: []string{baseContainerNameMin},
		Context:    dockerFileTarReader,
		Dockerfile: dockerFile,
		/* Other
		CPUSetCPUs:   "2",
		CPUSetMems:   "12",
		CPUShares:    20,
		CPUQuota:     10,
		CPUPeriod:    30,
		Memory:       256,
		MemorySwap:   512,
		ShmSize:      10,*/
		Remove:     true}

	// Build Docker image
	imageBuildResponse, err := cli.ImageBuild(ctx, dockerFileTarReader,imageOptions)
	if err != nil { return err, "Error building docker image" }


	defer func() {
		er = imageBuildResponse.Body.Close()
		if er != nil { msj =  "Error getting docker image build response" }
	}()

	// Print just to see the response in Stdout
	_, err = io.Copy(os.Stdout, imageBuildResponse.Body)
	if err != nil { return err, "Error reading image build response" }


	return err, ""   // Return ""  everything is good
}

/**********************************************************************************/


/*********************************************************************************
*	CreateDockerNetwork :
*			Function to create Docker Network
*			return error, output string
**********************************************************************************/

func CreateDockerNetwork (ctx context.Context, cli client.APIClient, networkName string) (er error,msj string){

	// check later for "10.12.0.0/16"
	// IPAM Driver Configuration
	ipamConf := network.IPAM{
		Driver: "static",
		Config: []network.IPAMConfig{
			{
				Subnet:  "10.12.0.0/16",    // ipv4     <- check Int2ipv4 funcion
			},
			{
				Subnet: "2001:db8::/64",    // global ipv6 - link ipv6 fe80::/64  <-  check Int2ipv6 function
			},
		},
		Options: make(map[string]string, 0),
	}

	// Network options
	networkCreateOptions := types.NetworkCreate{
		Driver:         	"l2bridge",
		EnableIPv6:     	true,
		IPAM: 				&ipamConf,
		Internal:   		false,
		Attachable:     	true,
		CheckDuplicate :	true,
		Options: map[string]string{
			"com.docker.network.bridge.name":   networkName,
			"l2bridge.gateway": "10.12.255.254",  // ipv4 gateway  check   <- check Int2ipv4 funcion
		},
	}

	respNetwork, err := cli.NetworkCreate(ctx, networkName , networkCreateOptions)

	if err != nil { return err, "Error creating docker network L2bridge" }
	return err, "Docker network Id - " + respNetwork.ID
}


/*********************************************************************************
*	CreateContainer :
*			Function to create Docker Container
*			return error, output string
**********************************************************************************/

func CreateContainer(ctx context.Context, cli client.APIClient, containerName string, imageName string, networkName string, logVolumeDirectory [2]string, confVolumeDirectory [2]string) (er error,msj string){

	// Container configuration
	containerConfig := &container.Config{
		Image:        	imageName,	 // Image Name
		Tty:          	true,     	 // Attach standard streams to a tty.
		AttachStdin:  	true,     	 // Attach the standard input, makes possible user interaction
		AttachStderr: 	true,    	 // Attach the standard error
		AttachStdout: 	true,    	 // Attach the standard output

		//	Hostname:        conf.Hostname,
		//	Domainname:      conf.Domainname,
		//  Cmd:   []string{"ls", "/"},
		//	Env:             conf.Env,
		//	Labels:          conf.Labels,
		//	NetworkDisabled: false,
		//	ExposedPorts:    ports,
	}

	// Host configuration
	hostConfig := &container.HostConfig{
		Privileged: true,
		Sysctls:  map[string]string{
			//"tools.ipv4.ip_forward":   "1",
		},

		//		Binds:      conf.Binds,
		//		CapAdd:        strslice.StrSlice([]string{"NET_ADMIN"}),
		//		RestartPolicy: restartPolicy,
		//		Resources: Container.Resources{
		//			Memory:     conf.Memory,
		//			MemorySwap: conf.MemorySwap,
		//			CPUShares:  conf.CPUShares,
		//		},

		NetworkMode:  container.NetworkMode(networkName), // networkMode -> docker network name to attach container

		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: logVolumeDirectory[0] + "/" + containerName,  //  Container logs -> host local path
				Target: logVolumeDirectory[1],    // Container logs -> path inside Container
			},
			{
				Type:   mount.TypeBind,
				Source: confVolumeDirectory[0] + "/" + containerName,  ///  Container configuration files -> host local path
				Target: confVolumeDirectory[1],    // Container configuration files  -> path inside Container
			},
		},
		//		PortBindings: portBindings,

	}


	n, err := strconv.ParseUint(containerName[3:], 10, 64)  // number of nodes to uint64


	// Network configuration
	netConfig := &network.NetworkingConfig{ //}
			EndpointsConfig: map[string]*network.EndpointSettings{
				networkName : {  // },
					IPAMConfig: &network.EndpointIPAMConfig{
						IPv4Address:  Int2ipv4(uint16(n)).String(),
						IPv6Address:  Int2ipv6(n).String(),
					},
				},
			},
		}

	//  Run configuration
	runOptions :=  types.ContainerStartOptions{}  // Default

	// Create Container
	respContainer, err := cli.ContainerCreate(ctx, containerConfig, hostConfig, netConfig, containerName)
	if err != nil { return err, "Error creating container" }

	// Start Container
	err = cli.ContainerStart(ctx, respContainer.ID, runOptions)
	if err != nil { return err, "Error starting container" }

	return err, "Container ID: " + respContainer.ID   // Return container ID if everything is good
}

/**********************************************************************************/


func Int2ipv4(lo uint16) net.IP {
	ip := make(net.IP, net.IPv4len)  	 //  ipv4 -> 4 bytes (32 bits)  10.12.0.0/16
	ip[0] = 10                           //  10 -> 10.0.0.0
	ip[1] = 12							 //	 12 -> 10.12.0.0
	binary.BigEndian.PutUint16(ip[2:], lo)   //  Max hosts [2^16 -2]  -> broadcast addr 10.12.255.255 , Default Gateway  10.12.255.254
	return ip
}

func Int2ipv6(lo uint64) net.IP {
	ip := make(net.IP, net.IPv6len)  	// ipv4 -> 16 bytes (128 bits)  2001:db8::0/64
	ip[0]= 32							//  0x20 -> 2000::0
	ip[1]= 1                            //  0x01 -> 2001::0
	ip[2]= 13							//  0x0d -> 2001:0d::0
	ip[3]= 184 							//  0xb8 -> 2001:db8::0
	ip[4]= 0
	ip[5]= 0
	ip[6]= 0
	ip[7]= 0
	binary.BigEndian.PutUint64(ip[8:],lo)  //  Max hosts [2^64]
	return ip
}




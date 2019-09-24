#!/bin/bash

# This script install driver plugins for docker
# https://github.com/nategraf/l2bridge-driver
# https://github.com/nategraf/static-ipam-driver
# Execute: bash driver-install.sh projectpath

projectPath=$1

# Update enviroment
cd ~ || exit
echo -e "\n\n Updating enviroment... \n" 
sudo apt update
sudo apt -y dist-upgrade


cd "$projectPath"/tools/driver  || exit


sudo rm -rf /etc/init.d/static-ipam
sudo rm -rf /usr/local/bin/static-ipam

sudo rm -rf /etc/init.d/l2bridge
sudo rm -rf /usr/local/bin/l2bridge


# Static IPAM driver
# Copy service script to init.d
sudo cp static-ipam-sysv.sh /etc/init.d/static-ipam
sudo chmod +x /etc/init.d/static-ipam

# Copy driver to usr/local/bin
sudo cp static-ipam /usr/local/bin/static-ipam
sudo chmod +x /usr/local/bin/static-ipam

# Activate the service
echo -e "\n\n  Active service   ... \n"
sudo update-rc.d static-ipam defaults
sudo service static-ipam start

# Verify that it is running
echo -e "\n\n  Verifying   ... \n"
sudo stat /run/docker/plugins/static.sock


# L2bridge
# Copy service script to init.d
sudo cp l2bridge-sysv.sh /etc/init.d/l2bridge
sudo chmod +x /etc/init.d/l2bridge

# Copy driver to usr/local/bin
sudo cp l2bridge /usr/local/bin/l2bridge
sudo chmod +x /usr/local/bin/l2bridge

# Activate the service
echo -e "\n\n  Active service   ... \n"
sudo update-rc.d l2bridge defaults
sudo service l2bridge start

# Verify that it is running
echo -e "\n\n  Verifying   ... \n"
sudo stat /run/docker/plugins/l2bridge.sock

cd ~  || exit


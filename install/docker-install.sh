#!/bin/bash

# This script install latest docker version , (v18.09.7)
# https://docs.docker.com/install/linux/docker-ce/ubuntu/
# Execute: bash docker-install.sh 

# Update enviroment
cd ~  || exit
echo -e "\n\n Updating enviroment... \n" 
sudo apt update
sudo apt -y dist-upgrade

# Remove old docker version 
echo -e "\n\n Removing old docker version ... \n"
sudo apt -y purge docker docker-engine docker.io containerd runc 

# Remove new versions of Docker
sudo apt -y purge docker-ce docker-ce-cli containerd.io
sudo rm -rf /var/lib/docker



# SETUP REPOSITORY 
# install pre-reqs for Https in apt
echo -e "\n\n  Setting up Docker repository  ... \n"
sudo apt -y install apt-transport-https ca-certificates curl gnupg-agent software-properties-common

# Add GPG key for oficial docker repository
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -

# Check Key
sudo apt-key fingerprint 0EBFCD88

# Add docker latest repository to apt
sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"

# Uptade 
sudo apt update
sudo apt -y dist-upgrade

# Checking recommended docker source and installation version 
echo -e "\n\n  Checking docker repository  ... \n"
apt-cache policy docker-ce


# INSTALL DOCKER 
echo -e "\n\n  Installing latest docker  ... \n"
sudo apt -y install docker-ce docker-ce-cli containerd.io

# Check docker 
echo -e "\n\n  Verifying  Docker  ... \n"
sudo docker run hello-world
sudo docker --version
cd ~  || exit



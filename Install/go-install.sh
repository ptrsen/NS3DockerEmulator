#!/bin/bash

# This script install latest golang version 
# https://golang.org/doc/install
# Execute: bash go-install.sh 

gitUser=ptrsen

# Getting Golang latest version and url 
url="$(wget -qO- https://golang.org/dl/ | grep -oP 'https:\/\/dl\.google\.com\/go\/go([0-9\.]+)\.linux-amd64\.tar\.gz' | head -n 1 )"
latest="$(echo $url | grep -oP 'go[0-9\.]+' | grep -oP '[0-9\.]+' | head -c -2 )"

# Update enviroment
cd ~
echo -e "\n\n Updating enviroment... \n" 
sudo apt update
sudo apt -y dist-upgrade

# Install pre-reqs
echo -e "\n\n Installing Pre-requirements ... \n"
sudo apt -y install python3 git curl wget

# Remove old Golang version 
echo -e "\n\n Removing old golang ... \n"
sudo rm -rf /usr/local/go
sudo rm -rf $HOME/go

# Install latest Golang version 
echo -e "\n\n Installing golang $latest ... \n"
# Download and extract golang , default home directory  path /usr/local/go 
wget $url
sudo tar -xvf go$latest.linux-amd64.tar.gz
rm -rf go$latest.linux-amd64.tar.gz
sudo chown -R root:root ./go
sudo mv go /usr/local

# Setting GO Paths 
echo "export GOROOT=/usr/local/go" >> ~/.bashrc    #use Default, if the installion is done in $HOME/go  GOROOT=$HOME/go" 
echo "export GOPATH=$HOME/go" >> ~/.bashrc
echo "export PATH=$HOME/go/bin:/usr/local/go/bin:$PATH">> ~/.bashrc  
source ~/.bashrc

# Setting Directory 
mkdir -p $HOME/go/src/github.com/$gitUser  # end  
echo -e "\n\n Work directory $HOME/go/src/github.com/$gitUser ... \n"
cd $HOME/go/src/github.com/$gitUser
git clone -b golang-porting https://github.com/ptrsen/NS3DockerEmulator.git
cd ~
source ~/.bashrc
# Checking Version
go version
#cd $HOME/go/src/github.com/$gitUser # project workspace

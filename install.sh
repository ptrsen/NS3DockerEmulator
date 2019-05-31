#!/bin/bash

# This script install all the required packages for The NS3DockerEmulator  https://github.com/chepeftw/NS3DockerEmulator
# Ns3.29 and Docker 18.09
# To run , open Terminal and execute
# source install.sh

echo -e "\n\n Updating enviroment... \n"
sudo apt update
sudo apt -y dist-upgrade

# https://www.nsnam.org/wiki/Installation#Ubuntu.2FDebian.2FMint

echo -e "\n\n Installing tools for Ns3 ... \n"
sudo apt -y install gcc g++ python python3 python-pip python3-pip make cmake tar unzip patch p7zip-full unrar-free
sudo apt -y install python-dev python3-dev python-setuptools python3-setuptools git mercurial qt5-default net-tools
sudo apt -y install gir1.2-goocanvas-2.0 python-gi python-gi-cairo python-pygraphviz python-kiwi python3-gi python3-gi-cairo python3-pygraphviz gir1.2-gtk-3.0 ipython ipython3
sudo apt -y install openmpi-bin openmpi-common openmpi-doc libopenmpi-dev
sudo apt -y install autoconf cvs bzr unrar
sudo apt -y install gdb valgrind
sudo apt -y install uncrustify
sudo apt -y install doxygen graphviz imagemagick
sudo apt -y install texlive texlive-extra-utils texlive-latex-extra texlive-font-utils texlive-lang-portuguese dvipng latexmk
sudo apt -y install python-sphinx dia
sudo apt -y install gsl-bin libgslcblas0:i386 libgslcblas0 libgsl-dev
sudo apt -y install flex bison libfl-dev
sudo apt -y install tcpdump

sudo apt -y install sqlite sqlite3 libsqlite3-dev
sudo apt -y install libxml2 libxml2-dev
sudo apt -y install libc6-dev libc6-dev-i386 libclang-dev llvm-dev automake
sudo apt -y install libgtk2.0-0 libgtk2.0-dev
pip install cxxfilt
pip3 install cxxfilt
sudo apt -y install vtun lxc
sudo apt -y install libboost-signals-dev libboost-filesystem-dev

# https://www.nsnam.org/wiki/Installation#Installation

echo -e "\n\n Setting Ns3 workspace ... \n"
echo -e "\n\n Installing and Setting bake tool ... \n"
mkdir Ns3
cd Ns3
hg clone http://code.nsnam.org/bake
cd bake

echo "export BAKE_HOME=$PWD" >> ~/.bashrc
echo "export PATH=$PATH:$BAKE_HOME:$BAKE_HOME/build/bin" >> ~/.bashrc
echo "export PYTHONPATH=$PYTHONPATH:$BAKE_HOME:$BAKE_HOME/build/lib" >> ~/.bashrc
source ~/.bashrc
echo -e "\n\n Verifying bake and missing packages to start Ns3 installation  ... \n"
python $BAKE_HOME/bake.py check
# rm -rf bakefile.xml
# python3 $BAKE_HOME/bake.py configure -e ns-3.29
python $BAKE_HOME/bake.py configure -e ns-allinone-3.29
python $BAKE_HOME/bake.py show
echo -e "\n\n Downloading and building Ns3   ... \n"
python $BAKE_HOME/bake.py download
python $BAKE_HOME/bake.py build

echo -e "\n\n Verifying Ns3  ... \n"
cd source/ns-3.29
./waf

echo -e "\n\n Recompoling NS3 in optimized mode  ... \n"
./waf distclean
./waf -d optimized configure --disable-examples --disable-tests --disable-python --enable-static --no-task-lines
./waf

echo -e "\n\n Running first Ns3 example  ... \n"
cp examples/tutorial/first.cc scratch/
./waf
./waf --run scratch/first
cd ~


echo -e "\n\n Installing Docker required packages  ... \n"

sudo apt-get install apt-transport-https ca-certificates curl software-properties-common
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
sudo apt-key fingerprint 0EBFCD88
sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"

service lxcfs stop
sudo apt-get remove lxc-common lxcfs lxd lxd-client

sudo apt-get update
sudo apt-get install docker-ce

echo -e "\n\n  Verifying  Docker  ... \n"
sudo docker run hello-world

echo -e "\n\n Installing Network Bridges  ... \n"

sudo apt install bridge-utils 
sudo apt install uml-utilities 

git clone https://github.com/chepeftw/NS3DockerEmulator.git
cd NS3DockerEmulator

# https://docs.docker.com/install/linux/docker-ce/ubuntu/#set-up-the-repository

echo -e "\n\n Installing Docker required packages  ... \n"
sudo apt -y install apt-transport-https ca-certificates curl software-properties-common
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
sudo apt-key fingerprint 0EBFCD88
sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"

service lxcfs stop
sudo apt-get -y remove lxc-common lxcfs lxd lxd-client

sudo apt update
sudo apt -y dist-upgrade
sudo apt -y install docker-ce

echo -e "\n\n  Verifying  Docker  ... \n"
sudo docker run hello-world



echo -e "\n\n Installing Network Bridges  ... \n"
sudo apt -y install bridge-utils
sudo apt -y install uml-utilities
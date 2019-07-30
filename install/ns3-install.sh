#!/bin/bash

# This script install ns3, ns-3.29
# https://www.nsnam.org/wiki/Installation#Ubuntu.2FDebian.2FMint
# Source: source ns3-install.sh 


ns3version=3.29

# Update enviroment
cd ~
echo -e "\n\n Updating enviroment... \n" 
sudo apt update 
sudo apt -y upgrade
sudo apt -y dist-upgrade


# Remove old Ns3 files
echo -e "\n\n Remove old Ns3... \n" 
#sudo rm -rf ~/Ns3

# install pre-reqs for Ns3
echo -e "\n\n Installing pre-reqs for Ns3 ... \n" 
sudo apt -y install gcc g++ 
sudo apt -y install python python3 
sudo apt -y install python-pip python3-pip 
sudo apt -y install python-pip python3-pip 
sudo apt -y install make cmake tar unzip p7zip-full unrar-free patch
sudo apt -y install python-dev python3-dev 
sudo apt -y install python-setuptools python3-setuptools 
sudo apt -y install qt5-default
 
sudo apt -y install git mercurial
sudo apt -y install gir1.2-goocanvas-2.0 python-gi python-gi-cairo 
sudo apt -y install python-pygraphviz python-kiwi 
sudo apt -y install python3-gi python3-gi-cairo python3-pygraphviz 
sudo apt -y install gir1.2-gtk-3.0 ipython ipython3
sudo apt -y install openmpi-bin openmpi-common openmpi-doc libopenmpi-dev
sudo apt -y install autoconf cvs bzr unrar
sudo apt -y install gdb valgrind 
sudo apt -y install uncrustify

sudo apt -y install doxygen graphviz imagemagick
sudo apt -y install texlive texlive-extra-utils 
sudo apt -y install texlive-latex-extra texlive-font-utils 
sudo apt -y install texlive-lang-portuguese dvipng latexmk
sudo apt -y install python-sphinx dia
sudo apt -y install gsl-bin libgslcblas0:i386 libgslcblas0 libgsl-dev
sudo apt -y install flex bison libfl-dev
sudo apt -y install tcpdump net-tools

sudo apt -y install sqlite sqlite3 libsqlite3-dev
sudo apt -y install libxml2 libxml2-dev
sudo apt -y install libc6-dev libc6-dev-i386 automake

sudo apt -y install libtinfo5 libtinfo-dev llvm-6.0-dev
sudo apt -y install llvm-dev
sudo apt -y install libclang-dev

sudo apt -y install libgtk2.0-0 libgtk2.0-dev
pip install cxxfilt
pip3 install cxxfilt

sudo apt -y install qt5-dev-tools 
sudo apt -y install libqt5-dev

sudo apt -y install libboost-signals-dev libboost-filesystem-dev
sudo apt -y install vtun 

sudo apt -y install bridge-utils 
sudo apt -y install uml-utilities

#sudo apt -y install lxc  # Linux Containers




# install Ns3 with Bake
# https://www.nsnam.org/wiki/Installation#Installation

echo -e "\n\n Setting Ns3 workspace ... \n" 
echo -e "\n\n Installing and Setting bake tool ... \n"
cd ~
mkdir Ns3
cd Ns3
git clone https://gitlab.com/nsnam/bake.git
cd bake

echo "export BAKE_HOME=$PWD" >> ~/.bashrc
echo "export PATH=$PATH:$BAKE_HOME:$BAKE_HOME/build/bin" >> ~/.bashrc
echo "export PYTHONPATH=$PYTHONPATH:$BAKE_HOME:$BAKE_HOME/build/lib" >> ~/.bashrc
source ~/.bashrc

echo -e "\n\n Verifying bake and missing packages to start Ns3 installation  ... \n" 

# rm -rf bakefile.xml
# python3 $BAKE_HOME/bake.py configure -e ns-$ns3version # none All in one version
python $BAKE_HOME/bake.py configure -e ns-allinone-$ns3version
python $BAKE_HOME/bake.py check
python $BAKE_HOME/bake.py show

echo -e "\n\n Downloading and building Ns3   ... \n" 
python $BAKE_HOME/bake.py download
python $BAKE_HOME/bake.py build

echo -e "\n\n Verifying and compiling Ns3 in normal mode   ... \n" 
cd source/ns-$ns3version
./waf
cd ~









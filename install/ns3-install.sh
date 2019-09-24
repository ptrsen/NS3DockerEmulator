#!/bin/bash

# This script install ns3 (ns-latest) (latest version)
# https://www.nsnam.org/wiki/Installation#Ubuntu.2FDebian.2FMint
# Source: source ns3-install.sh

url="$(wget -qO- https://www.nsnam.org/release/ | grep -oP 'releases\/ns\-([0-9\-]+)\/.*latest'  | head -n 1 )"
latest="$(echo "$url" | grep -oP 'ns\-[0-9\-]+' | grep -oP '[0-9\-]+' | sed 's/-//'| sed 's/-/./g' | head -c 4  )"

# Update enviroment~/.bashrc
cd ~ || exit
echo -e "\n\n Updating enviroment... \n" 
sudo apt update 
sudo apt -y upgrade
sudo apt -y dist-upgrade

# Remove old Ns3 files
echo -e "\n\n Remove old Ns3... \n" 
sudo rm -rf ~/Ns3

# install pre-reqs for Ns3
echo -e "\n\n Installing pre-reqs for Ns-$latest ... \n"
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

echo -e "\n\n Setting Ns-$latest workspace ... \n"
echo -e "\n\n Installing and Setting bake tool ... \n"
cd ~ || exit
mkdir Ns3
cd Ns3 || exit
git clone https://gitlab.com/nsnam/bake.git
cd bake || exit

bakePath=$PWD

{
echo "export BAKE_HOME=$PWD"
echo "export PATH=$PATH:$BAKE_HOME:$BAKE_HOME/build/bin"
echo "export PYTHONPATH=$PYTHONPATH:$BAKE_HOME:$BAKE_HOME/build/lib"
} >> ~/.bashrc


. ~/.bashrc

echo -e "\n\n Verifying bake and missing packages to start Ns-$latest installation  ... \n"

# rm -rf bakefile.xml
# python3 $BAKE_HOME/bake.py configure -e ns-$ns3version # none All in one version
python "$bakePath"/bake.py configure -e ns-allinone-"$latest"
python "$bakePath"/bake.py check
python "$bakePath"/bake.py show

echo -e "\n\n Downloading and building Ns-$latest   ... \n"
python "$bakePath"/bake.py download
python "$bakePath"/bake.py build

echo -e "\n\n Verifying and compiling Ns-$latest in normal mode   ... \n"
cd source/ns-"$latest" || exit
./waf

# for Optimize mode
echo -e "\n\n Recompoling NS-$latest in optimized mode  ... \n"
./waf distclean
./waf -d optimized configure --disable-examples --disable-tests --disable-python --enable-static --no-task-lines
./waf

echo -e "\n\n Running first Ns-$latest example  ... \n"
cp examples/tutorial/first.cc scratch/
./waf
./waf --run scratch/first
cd ~ || exit










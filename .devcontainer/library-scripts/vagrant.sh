curl -fsSL https://apt.releases.hashicorp.com/gpg | sudo apt-key add -

sudo apt-add-repository "deb [arch=amd64] https://apt.releases.hashicorp.com $(lsb_release -cs) main"


# echo "deb [arch=amd64] https://download.virtualbox.org/virtualbox/debian $(lsb_release -cs) contrib" >> /etc/apt/sources.list 

# wget -q https://www.virtualbox.org/download/oracle_vbox_2016.asc -O- | sudo apt-key add -
# wget -q https://www.virtualbox.org/download/oracle_vbox.asc -O- | sudo apt-key add -

# sudo apt-get update -y
# sudo apt-get install qemu virtualbox-6.1 -y


sudo apt update
apt install fasttrack-archive-keyring
echo "deb https://fasttrack.debian.net/debian-fasttrack/ bullseye-fasttrack main contrib" >> /etc/apt/sources.list
echo "deb https://fasttrack.debian.net/debian-fasttrack/ bullseye-backports-staging main contrib" >> /etc/apt/sources.list
sudo apt update
sudo apt install virtualbox -y
sudo apt-get install vagrant -y
apt install wget curl -y
sudo apt-get install -y  software-properties-common
sudo apt-get install virtualbox-dkms
sudo dpkg-reconfigure virtualbox-dkms 
sudo dpkg-reconfigure virtualbox
sudo apt-get install linux-headers-generic
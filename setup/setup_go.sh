curl -O https://dl.google.com/go/go1.14.3.linux-amd64.tar.gz

tar xvf go1.14.3.linux-amd64.tar.gz

sudo chown -R root:root ./go
sudo mv go /usr/local

export GOPATH=$HOME/rad
export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin

cd /usr/local/go
git clone https://github.com/shudaniel/RuntimeAutoDeploy.git
cd RuntimeAutoDeploy


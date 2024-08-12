REAL_PATH=$( realpath . )
echo $REAL_PATH

cd /tmp
rm -rf teams
git clone --depth=1 https://github.com/moyai-network/teams
cd teams
go build -o $REAL_PATH .

cd $REAL_PATH
./teams
rm -rf teams

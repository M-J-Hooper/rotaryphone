sudo systemctl start ofono
sudo pulseaudio --start

sudo bluetoothctl <<EOF
power on
agent on
default-agent
discoverable on
pairable on
EOF

export PATH=$PATH:/usr/local/go/bin
export GOPATH=$(go env GOPATH)

go run $GOPATH/src/github.com/M-J-Hooper/rotaryphone/run/main.go

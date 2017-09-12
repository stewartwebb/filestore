#!/bin/bash

echo "TCPSocket 3310" > /etc/clamav/clamd.conf
echo "TCPAddr $CLAMAV_SERVER" >> /etc/clamav/clamd.conf

go get ./ && CompileDaemon -command="./src";

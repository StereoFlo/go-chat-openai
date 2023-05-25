#!/bin/bash

if [ "$EUID" -ne 0 ]
  then echo "run sudo as root!"
  exit
fi

cp ./go-chat-tg.service /etc/systemd/system/go-chat-tg.service
systemctl daemon-reload
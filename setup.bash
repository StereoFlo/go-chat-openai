#!/bin/bash

if [ "$EUID" -ne 0 ]
  then echo "run sudo as root!"
  exit
fi

cp ./go-tg-bot.service /etc/systemd/system/go-tg-bot.service
systemctl daemon-reload

#!/bin/bash
echo "running on remote server"
sudo bash
dctl -s 172.16.19.75 login -u admin -p Diamanti@111
dctl cluster status > aa

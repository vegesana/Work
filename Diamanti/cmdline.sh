#!/bin/bash
sudo bash
hn=`hostname`
echo "#########################################"
echo "#########################################"
echo ""
echo "ON SERVER : $hn"
echo ""
echo "#########################################"
echo "#########################################"
dctl -s 172.16.19.75 login -u admin -p Diamanti@111
dctl cluster status > rajuConfig.txt
dctl network list >> rajuConfig.txt

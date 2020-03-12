#!/bin/bash

# PLEASE RUN ON SECURECRT
password=diamanti
if [ "$#" -lt 1 ] ; then
        echo "Must pass at least one server name"
        exit 1
fi
echo "Number of argument passed is $#"
for arg in "$@"; do 
    cat diag.sh | sshpass -p $password ssh diamanti@"$arg"
    cat cmdline.sh | sshpass -p $password ssh diamanti@"$arg"
    mkdir -p "Temp/$arg"
    Dir1="Temp/$arg/${arg}Raju.txt"
    Dir="Temp/$arg/"
    sshpass -p $password scp diamanti@"$arg":/home/diamanti/rajumerged.txt $Dir1
    sshpass -p $password scp diamanti@"$arg":/home/diamanti/ncdutil.log $Dir
    sshpass -p $password scp diamanti@"$arg":/home/diamanti/boardinfo.log $Dir
done


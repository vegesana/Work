#!/bin/bash
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
    file1="Temp/$arg/${arg}Raju.txt"
    file2="Temp/$arg/ncdutil.log"
    file3="Temp/$arg/boardinfo.log"
    sshpass -p $password scp diamanti@"$arg":/home/diamanti/rajumerged.txt $file1
    sshpass -p $password scp diamanti@"$arg":/home/diamanti/ncdutil.log $file2
    sshpass -p $password scp diamanti@"$arg":/home/diamanti/boardinfo.log $file3
done


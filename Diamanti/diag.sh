#!/bin/bash
echo "running on remote server"
sudo bash
SERVICES=("armada" "convoy" "bosun")
for service in ${SERVICES[@]}; do
        status=`systemctl is-active "$service"`
        if [[ ! $status == 'active' ]] ; then
                echo "SERVICE $service Is Not running"
        else
                echo "Happy: SERVICE $service Is running"
        fi
done

echo "Sorting logfiles"
OUTPUTFILE="/home/diamanti/rajumerged.txt"
touch $OUTPUTFILE
chmod 777 $OUTPUTFILE
OUTFILE="/home/diamanti/ncdutil.log"
touch $OUTFILE
chmod 777 $OUTFILE
OUTFILE1="/home/diamanti/boardinfo.log"
touch $OUTFILE1
chmod 777 $OUTFILE1

COREPATH="/var/log/diamanti/core/"
LOGFILES=("armada.log" "convoy.log" "bosun.log")
for logfile in ${LOGFILES[@]}; do
        cat "$COREPATH$logfile"
done > tmp.txt
COREPATH="/var/log/diamanti/embedded/"
LOGFILES=("ncd.log") 
for logfile in ${LOGFILES[@]}; do
        cat "$COREPATH$logfile"
done > tmp1.txt
sudo dstool -c "ncdutil -a" > $OUTFILE
sudo dstool -c "boardinfo" > $OUTFILE1
sort -k1,2 tmp.txt > mytemp.txt   # sorting the file by Timestamp
sort -k1,2 tmp1.txt >> mytemp.txt   # sorting the file by Timestamp
sort -k1,2 mytemp.txt > $OUTPUTFILE   # sorting the file by Timestamp

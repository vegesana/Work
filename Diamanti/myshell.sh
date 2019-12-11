#!/bin/bash
DIR="/tmp/Raju/*"
for file in $DIR; do 
        if [[ $file == *.txt ]]; then 
            n=`awk '/CPSS/{print NR}' $file`  # Getting LIne Number
            l=`awk '/PCL Info/{print NR}' $file` 
            echo "Line n is $n"
            echo "Line l is $l"
            filename=${file}.tmp
            awk 'NR > $n && NR<139 {print}' $file > $filename
        fi
done

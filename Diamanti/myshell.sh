#!/bin/bash
DIR="/tmp/Raju/*"
for file in $DIR; do 
        if [[ $file == *.txt ]]; then 
            n=`awk '/CPSS/{print NR}' $file` 
            l=`awk '/PCL Info/{print NR}' $file` 
            filename=${file}.tmp
            awk 'NR>65 && NR<137 {print}' $file > $filename
        fi
done

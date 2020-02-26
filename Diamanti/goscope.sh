#!/bin/sh

# for Full path
find $PWD -name '*.[ch]' -exec echo \"{}\" \; | sort -u > cscope.files
cscope -bvq

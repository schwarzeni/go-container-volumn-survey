#!/bin/sh

# build project
go build -o build/gv .

# change dir
cd build

# extract rar
[ ! -d "busybox" ] && mkdir busybox
tar -xf ../busybox.rar -C busybox

# create a mount dir
[ ! -d "self_mnt" ] && mkdir self_mnt

# write data to file in self_mnt
echo $(date) > self_mnt/time.txt

# execute program
./gv -r . -mnt ./mnt -m ./self_mnt:/self_mnt

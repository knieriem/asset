#!/bin/sh

exe=$1

a=,,asset.zip

rm -f $a
/usr/bin/zip -r $a assets

cat $a >> $exe
rm -f $a

/usr/bin/zip -A $exe

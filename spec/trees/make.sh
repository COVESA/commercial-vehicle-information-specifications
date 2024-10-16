#! /usr/bin/bash
make -C ../../ $1 VSPECROOT=$2 VALIDATE='-e validate'
#make -C ../../ $1 VSPECROOT=$2


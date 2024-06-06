#!/bin/bash

echo "enter the name of you project: "
read varname

git clone git@github.com:anthdm/gothkit.git _gothkit
cd _gothkit
mv bootstrap ../$varname
cd ..
rm -rf _gothkit
echo ""
echo "project installed successfully" 
echo "your project folder is available => cd $varname"


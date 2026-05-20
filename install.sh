#!/bin/sh

wget https://github.com/z3nnix/chorus/releases/download/1.1.0/chorus.bin
sudo mv chorus.bin /usr/bin/chorus
sudo chmod +x /usr/bin/chorus
rm -rf chorus
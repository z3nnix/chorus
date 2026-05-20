#!/bin/sh

wget https://github.com/z3nnix/chorus/releases/download/1.0.2/chorus.zip
unzip chorus.zip
rm chorus.zip
sudo mv chorus.bin /usr/bin/chorus
sudo chmod +x /usr/bin/chorus
rm -rf chorus
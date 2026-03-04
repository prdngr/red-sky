#!/bin/bash

apt-get update -y
apt-get install -y \
    ca-certificates \
    gnupg \
    make \
    curl \
    wget \
    vim \
    git

mkdir -p /home/kali/.c2
cd /home/kali/.c2

git clone https://github.com/its-a-feature/Mythic --depth 1 mythic
cd /home/kali/.c2/mythic/

./install_docker_kali.sh

make

./mythic-cli install github https://github.com/MythicC2Profiles/http
./mythic-cli start

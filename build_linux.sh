#!/bin/bash

git submodule update --init --recursive

sudo apt install -y libx11-dev xorg-dev libxtst-dev xcb libxcb-xkb-dev x11-xkb-utils libx11-xcb-dev libxkbcommon-x11-dev libxkbcommon-dev xsel xclip

go mod tidy
go build -ldflags "-s -w" ./cmd/sshx

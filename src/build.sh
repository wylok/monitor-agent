#!/usr/bin/env bash
go build -ldflags "-w -s" main.go
mv -f main monitor-agent
upx --brute monitor-agent
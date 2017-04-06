#!/usr/bin/env bash
#ps x|grep 8000|grep server |awk '{print $1}' | xargs kill -9
ps x|grep go-sms-server |grep -v grep|awk '{print $1}' | xargs kill -9
# PRODUCTION
nohup ./go-sms-server TEST >> log.log &
#!/bin/bash
process_name=`pwd | xargs basename`

usage() {
	echo "./start.sh [ start | stop | reload | status | rdoc ]"
}

status() {
 	ps -ef | grep $process_name | grep -v grep
}

start() {
	echo "start $process_name...."
	nohup ./$process_name 2>&1 1>/dev/null &
	status
}

runwithdoc() {
	echo "start $process_name with doc enabled..."
	bee run -gendoc=true -downdoc=true
}

stop() {
	echo "stop $process_name..."
	ps -ef | grep $process_name | grep -v grep | awk '{print $2}' |xargs kill -9
	status
}

reload() {
	echo "reload $process_name..."
	ps -ef | grep $process_name |grep -v grep | awk '{print $2}' |xargs kill -HUP
	status
}

case "$1" in
	start)
		start
		;;
	stop)
		stop
		;;
	reload)
		reload
		;;
	status)
		status
		;;
	rdoc)
		runwithdoc
		;;
	*)
		usage
		;;
esac

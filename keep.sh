#!/bin/sh
cameraIP=camera.u-suke.org

commands () {
	cat <<EOS
start
stop
start_with_keepalive
keepalive
commands
help
EOS
}

getUrlContents () {
	curl "$@"
}

usage () {
	cat <<EOS
$(basename $0) [$(commands | xargs | sed -e 's/ /|/g')]
EOS
}

start () {
	getUrlContents "http://$cameraIP/cam.cgi?mode=accctrl&type=req_acc&value=0&value2=lumix-webcam"
	getUrlContents "http://$cameraIP/cam.cgi?mode=camcmd&value=recmode"
	getUrlContents "http://$cameraIP/cam.cgi?mode=startstream&value=49199&value2=10"
}

stop () {
	getUrlContents "http://$cameraIP/cam.cgi?mode=stopstream"
}

interrupt () {
	stop
	exit
}

keepalive () {
	while true
	do
		getUrlContents "http://$cameraIP/cam.cgi?mode=getstate"
		getUrlContents "http://$cameraIP/cam.cgi?mode=camctrl&type=touch&value=500/500&value2=on"
		sleep 10
	done
}

trap interrupt SIGINT

case $1 in
	start|stop|keepalive)
		$1
		;;
	start_with_keepalive)
		start
		keepalive
		;;
	*)
		usage
		;;
esac

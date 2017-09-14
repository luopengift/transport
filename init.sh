#!/bin/sh
#
# transport - this script start and stop the transport daemon
#
# chkconfig:
# description:
#               
# process name: transport
# config:      config.json
# pidfile:     var/transport.pid

# Source env
#. /etc/profile

# Source function library.
#. /etc/rc.d/init.d/functions

# Source networking configuration.
#. /etc/sysconfig/network

NOW=$(date +%Y-%m-%d.%H:%M:%S)
echo $NOW
DIR=$(pwd) #当前目录
BASEDIR=$(dirname $(dirname $(dirname $DIR))) 
export GOPATH=$GOPATH:$BASEDIR
export GOBIN=$GOPATH/bin/
echo "GOPATH init Finished. GOPATH=$GOPATH"

#################################
APP=$(basename $DIR)
######创建程序运行临时文件#######

mkdir -p $DIR/var
PIDFile=$DIR/var/$APP.pid
LOGFile=$DIR/var/$APP.log
CONFIG=config.json
#编译$2,默认$2未指定时编译main.go
if [ -z $2 ];then #如果$2为空
    main="main.go"
else
    main=$2
fi 
#main="cmd/transport/main"
################################
function check_PID() {
    if [ -f $PIDFile ];then
        PID=$(cat $PIDFile)
        if [ -n $PID ]; then
            running=$(ps -p $PID|grep -v "PID TTY" |wc -l)
            return $running
        fi
    fi
    return 0
}

function build() {
    gofmt -w .
    go build -o $APP $1
    if [ $? -ne 0 ]; then
        exit $?
    fi
    echo "build $APP success."
}

function start() {
    check_PID
    local running=$?
    if [ $running -gt 0 ];then
        echo -n "$APP now is running already, PID="
        cat $PIDFile
        return 1
    fi

    nohup  ./$APP -f $CONFIG >$LOGFile 2>&1 &
    sleep 1
    running=`ps -p $! | grep -v "PID TTY" | wc -l`
    if [ $running -gt 0 ];then
        echo $! > $PIDFile
        echo "$APP started..., PID=$!"
    else
        echo "$APP failed to start!!!"
        return 1
    fi
}

function status() {
    ps -ef |grep $APP|grep -v grep
    check_PID    
    local running=$?
    if [ $running -gt 0 ];then
        echo -n "$APP now is running already, PID="
        cat $PIDFile
        return 1
    else
        echo "$APP is stopped..."
        return 1
    fi
    
}


function debug() {
    go run $1
}

function stop() {
    local PID=$(cat $PIDFile)
    kill $PID
    rm -f $PIDFile
    echo "$APP stoped..."
}

function restart() {
    stop
    sleep 1
    start   
}

function pack() {
    tar -cvf $APP.tar.gz $APP config.json init.sh
}

function tailf() {
   tail -f $LOGFile
}

function help() {
    echo "$0 build|start|stop|kill|restart|reload|run|tail|docs|sslkey"
}


case "$1"  in
    build) build $main ;;
    start) start ;;
    debug) debug $main ;;
    stop) stop ;;
    kill) killall ;;
    restart) restart ;;
    reload) reload ;;
    status) status ;;
    run) run ;;
    pack) pack ;;
    tail) tailf ;;
    docs) docs ;;
    sslkey) sslkey ;;
    *) help ;;
esac




#bee api  cmdb-api -driver=mysql -conn="root:root888@tcp(127.0.0.1:3306)/cmdb"


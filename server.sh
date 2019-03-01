#! /bin/sh
#程序启动脚本
### BEGIN INIT INFO
# Provides:          gonotify
# Required-Start:    $remote_fs $network
# Required-Stop:     $remote_fs $network
# Default-Start:     2 3 4 5
# Default-Stop:      0 1 6
# Short-Description: starts gonotify
# Description:       starts the zcxf api center service
### END INIT INFO
prefix=/www/sites/gonotify
exec_prefix=${prefix}/notify_server
gonotify_PID=${prefix}/notify_server.pid
gonotify_CONF=${prefix}/main.conf

wait_for_pid () {
    try=0
    while test $try -lt 35 ; do
        case "$1" in
            'created')
            if [ -f "$2" ] ; then
                try=''
                break
            fi
            ;;

            'removed')
            if [ ! -f "$2" ] ; then
                try=''
                break
            fi
            ;;
        esac
        echo -n .
        try=`expr $try + 1`
        sleep 1
    done
}

case "$1" in
    start)
        echo -n "Starting gonotify"
        $exec_prefix --pprof --cross --pid $gonotify_PID --config $gonotify_CONF &
        if [ "$?" != 0 ] ; then
            echo " failed"
            exit 1
        fi
        wait_for_pid created $gonotify_PID
        if [ -n "$try" ] ; then
            echo " failed"
            exit 1
        else
            echo " done"
        fi
    ;;

    stop)
        echo -n "Gracefully shutting down gonotify "

        if [ ! -r $gonotify_PID ] ; then
            echo "warning, no pid file found - gonotify is not running ?"
            exit 1
        fi

        kill -QUIT `cat $gonotify_PID`

        wait_for_pid removed $gonotify_PID

        if [ -n "$try" ] ; then
            echo " failed. Use force-quit"
            exit 1
        else
            echo " done"
        fi;
    ;;
esac
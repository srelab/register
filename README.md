# Docker Regsiter

## 安装

```
go get github.com/itchenyi/register
cd $GOPATH/src/github.com/srelab/register

make build
```

## 基于Docker事件关联Consul 和 Gateway

1. `start`, `unpause` 关联注册逻辑

```
                            event
                              |
                              |
                          http check
                              |
                              |
               success <-------------> failed
                  |                      |
                  |                      |
               register              continue
```

2. `pause`, `die` 关联反注册逻辑

```
                            event
                              |
                              |
                           success 
                              |        
                          unRegister
```


## 程序参数

```

#####  ######  ####  #  ####  ##### ###### #####
#    # #      #    # # #        #   #      #    #
#    # #####  #      #  ####    #   #####  #    #
#####  #      #  ### #      #   #   #      #####
#   #  #      #    # # #    #   #   #      #   #
#    # ######  ####  #  ####    #   ###### #    #
NAME:
    start - start a new gateway-register

USAGE:
    start [command options] [arguments...]

OPTIONS:
   --concurrency value      concurrency number (default: 10)
   --docker.endpoint value  Docker Conn EndPoint (default: "unix:///var/run/docker.sock")
   --log.dir value          the log file is written to the path (default: "./")
   --log.level value        valid levels: [debug, info, warn, error, fatal] (default: "info")
   --gateway.host value     gateway server host
   --gateway.port value     gateway server port
   --consul.host value      consul server host
   --consul.port value      consul server port
   --privilege.host value   privilege server host
   --privilege.port value   privilege server port
```




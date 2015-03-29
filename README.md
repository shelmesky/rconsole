### Rconsole

远程管理工具，支持VNC/RDP/SPICE/SSH/TELNET协议。
使用websocket协议传输数据，在HTML5浏览器中使用canvas显示。 

### TODO

+ 连接到Libvirt服务的虚拟机，支持VNC/SPICE协议
+ 连接到XenServer，支持VNC/RDP协议
+ 连接到VMWare ESXI，支持VNC协议


### 协议参数

[协议参数](https://github.com/shelmesky/rconsole/blob/master/PROTOCOLS.md "协议参数")


### 通过URL中指定参数连接

#### VNC:

http://127.0.0.1:9999/connect?type=vnc&hostname=172.31.31.101&port=5900&password=123456&width=1024&height=660&dpi=96

#### RDP:

http://127.0.0.1:9999/connect?type=rdp&hostname=172.31.31.123&port=3389&username=roy&width=1024&height=660&dpi=96


#### SSH:

http://127.0.0.1:9999/connect?type=ssh&hostname=172.31.31.110&port=22&username=roy&width=1024&height=660

#### TELNET:

http://127.0.0.1:9999/connect?type=telnet&hostname=172.31.31.110&port=23&username=roy&width=1024&height=660

#### SPICE:

http://127.0.0.1:9999/connect?type=spice&hostname=127.0.0.1&port=5900&password=123


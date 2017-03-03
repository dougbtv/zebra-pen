# Quagga stuff

Open quagga console...

```
[root@router_a centos]# docker exec -it quagga /bin/bash
root@router_a:/# vtysh
```

[Static routes with quagga](http://www.nongnu.org/quagga/docs/docs-multi/Static-Route-Commands.html)

Non interactive...

```
[root@router_a centos]# docker exec -i quagga vtysh -c 'show ip route'
```

I do not know what the heck I'm doing... Got a little further from the [configuration doc](https://docs.cumulusnetworks.com/display/ROH/Configuring+Cumulus+Quagga)

```
router_a.example.local# configure terminal
router_a.example.local(config)# interface swp1
router_a.example.local(config-if)# route ospf
% Unknown command.
router_a.example.local(config-if)# router ospf
router_a.example.local(config-router)# 
router_a.example.local(config-router)# 
router_a.example.local(config-router)# 
router_a.example.local(config-router)# ip route 192.168.122.31/32 192.168.122.9
router_a.example.local(config)# 
router_a.example.local# show ip route
Codes: K - kernel route, C - connected, S - static, R - RIP,
       O - OSPF, I - IS-IS, B - BGP, T - Table,
       > - selected route, * - FIB route

K>* 0.0.0.0/0 via 192.168.122.1, eth0
C>* 192.168.122.0/24 is directly connected, eth0
S>* 192.168.122.31/32 [1/0] via 192.168.122.9, eth0

router_a.example.local# write 
Building Configuration...
Integrated configuration saved to /etc/quagga/Quagga.conf
[OK]
```

And here's the config after writing... maybe template this?

```
root@router_a:/etc/quagga# cat /etc/quagga/Quagga.conf
hostname zebra
log file /var/log/quagga/zebra.log
hostname bgpd
log file /var/log/quagga/bgpd.log
log timestamp precision 6
username cumulus nopassword
!
service integrated-vtysh-config
!
interface docker0
 link-detect
!
interface eth0
 ipv6 nd suppress-ra
 link-detect
!
interface lo
 link-detect
!
interface swp1
 link-detect
!
ip route 192.168.122.31/32 192.168.122.9
!
ip forwarding
!
line vty
!
```

But, I'm missing something.... is it the [forwarding rule](http://askubuntu.com/questions/227369/how-can-i-set-my-linux-box-as-a-router-to-forward-ip-packets)

YES IT IS!!!!

This worked...

```
[root@router_a centos]# iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE
[root@router_a centos]# iptables -A FORWARD -i eth0 -j ACCEPT
```



# Routing on centos a/b...

```
[root@centos_a centos]# route add -net 192.168.122.31 netmask 255.255.255.255 gw 192.168.122.135
[root@centos_a centos]# route del -net 192.168.122.31 netmask 255.255.255.255 gw 192.168.122.135
```


# Phase 3

Note, I don't know what I'm doing, lolz.

[This config reference looks promising](http://www.brianlinkletter.com/how-to-build-a-network-of-linux-routers-using-quagga/)

That article was good as gold, worked a charm.
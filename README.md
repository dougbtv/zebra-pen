# zebra-pen

A suite of playbooks and roles to create a demo router VNF, in containers.

## Goal

Have box `CentOS A` ping box `CentOS B` via containerized routers running on Router A & Router B

```
CentOS A -> Router A -> Router B -> CentOS B
```

## Single VM, specific network interfaces in containers

![layout diagram](http://i.imgur.com/Tr18ztQ.png)

In this third iteration we're going to use a specialized script, [vethcon](https://github.com/s1061123/vethcon),created by [Tomo](https://github.com/s1061123) which allows us to specifically craft network interfaces for each container. These allow us to create 2 interfaces (one in each container) that are connected. This playbook spins up the diagrammed containers, and creates the interfaces as diagrammed above.

This configuration uses OSPFd in Quagga for routing from Centos A through the two routers (Quagga A & Quagga B) to Centos B, and then back through in the reverse direction.

To kick off this playbook, modify the `./inventory/single_vm.inventory` file and then run:

```
$ ansible-playbook -i inventory/single_vm.inventory vethcon-single-vm.yml
```

Once that has run, you can verify that it is working by checking the running containers on a host, and then entering the `centos_a` container, and pinging the IP address for `centos_b`.

```
[root@centos-host src]# docker ps
CONTAINER ID        IMAGE                                  COMMAND             CREATED             STATUS              PORTS               NAMES
f90fa1659ae1        cumulusnetworks/quagga:xenial-latest   "/bin/bash"         2 minutes ago       Up 2 minutes                            quagga_b
515cf97a9c3d        cumulusnetworks/quagga:xenial-latest   "/bin/bash"         3 minutes ago       Up 3 minutes                            quagga_a
228d07174535        veth_centos                            "/bin/bash"         4 minutes ago       Up 4 minutes                            centos_b
4a28bddbd406        veth_centos                            "/bin/bash"         4 minutes ago       Up 4 minutes                            centos_a
[root@centos-host src]# docker exec -it centos_a /bin/bash
[root@centos_a /]# route -n
Kernel IP routing table
Destination     Gateway         Genmask         Flags Metric Ref    Use Iface
0.0.0.0         172.17.0.1      0.0.0.0         UG    0      0        0 eth0
172.17.0.0      0.0.0.0         255.255.0.0     U     0      0        0 eth0
192.168.2.0     0.0.0.0         255.255.255.0   U     0      0        0 in1
192.168.4.0     192.168.2.101   255.255.255.0   UG    0      0        0 in1
[root@centos_a /]# ping -c 3 192.168.4.101
PING 192.168.4.101 (192.168.4.101) 56(84) bytes of data.
64 bytes from 192.168.4.101: icmp_seq=1 ttl=62 time=0.195 ms
64 bytes from 192.168.4.101: icmp_seq=2 ttl=62 time=0.117 ms
64 bytes from 192.168.4.101: icmp_seq=3 ttl=62 time=0.157 ms

--- 192.168.4.101 ping statistics ---
3 packets transmitted, 3 received, 0% packet loss, time 2000ms
rtt min/avg/max/mdev = 0.117/0.156/0.195/0.033 ms
[root@centos_a /]# 
```

---

Let's specifically look at the interfaces on on each machine. You'll note that they match the above diagram.

### Centos A

```
[root@centos-host src]# docker exec -it centos_a ip a | grep -Pi "(^\d|inet )"
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN qlen 1
    inet 127.0.0.1/8 scope host lo
6: eth0@if7: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc noqueue state UP 
    inet 172.17.0.2/16 scope global eth0
15: in1@if14: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc noqueue state UP 
    inet 192.168.2.100/24 scope global in1
```

### Quagga A

```
[root@centos-host src]# docker exec -it quagga_a ip a | grep -Pi "(^\d|inet )"
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1
    inet 127.0.0.1/8 scope host lo
10: eth0@if11: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc noqueue state UP group default 
    inet 172.17.0.4/16 scope global eth0
14: in2@if15: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc noqueue state UP group default 
    inet 192.168.2.101/24 scope global in2
17: mid1@if16: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc noqueue state UP group default 
    inet 192.168.3.100/24 scope global mid1
```

### Quagga B

```
[root@centos-host src]# docker exec -it quagga_b ip a | grep -Pi "(^\d|inet )"
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1
    inet 127.0.0.1/8 scope host lo
12: eth0@if13: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc noqueue state UP group default 
    inet 172.17.0.5/16 scope global eth0
16: mid2@if17: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc noqueue state UP group default 
    inet 192.168.3.101/24 scope global mid2
19: out1@if18: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc noqueue state UP group default 
    inet 192.168.4.100/24 scope global out1
```

### Centos B

```
[root@centos-host src]# docker exec -it centos_b ip a | grep -Pi "(^\d|inet )"
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN qlen 1
    inet 127.0.0.1/8 scope host lo
8: eth0@if9: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc noqueue state UP 
    inet 172.17.0.3/16 scope global eth0
18: out2@if19: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc noqueue state UP 
    inet 192.168.4.101/24 scope global out2
```


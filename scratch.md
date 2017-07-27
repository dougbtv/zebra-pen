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


## Openshift style.

### Getting a new all in one...

- must use docker from centos (version matching denies use of new naming scheme)
- [this looks eerily similar](https://github.com/opencontainers/runc/issues/1343)
    + tried `grubby --args="user_namespace.enable=1" --update-kernel="$(grubby --default-kernel)"`
    + really failing there, no luck at all.
    + now... 
- trying with [this method](https://github.com/openshift/origin/blob/master/docs/cluster_up_down.md#linux)
    + which has you download 

### REdo! Did it manually.

This worked.

```
[root@ip-172-31-44-254 centos]# history
    1  yum install -y nano
    2  yum update -y
    3  yum install docker
    4  reboot
    5  docke rps
    6  docker ps
    7  systemctl enable docker
    8  systemctl start docker
    9  wget https://github.com/openshift/origin/releases/download/v1.5.0/openshift-origin-client-tools-v1.5.0-031cbe4-linux-64bit.tar.gz
   10  yum install -y wget
   11  nano /etc/sysconfig/docker
   12  systemctl restart docker
   13  wget https://github.com/openshift/origin/releases/download/v1.5.0/openshift-origin-client-tools-v1.5.0-031cbe4-linux-64bit.tar.gz
   14  tar -xzvf openshift-origin-client-tools-v1.5.0-031cbe4-linux-64bit.tar.gz 
   15  cd openshift-origin-client-tools-v1.5.0-031cbe4-linux-64bit
   16  ls
   17  ./oc cluster up
   18  oc login -u system:admin
   19  ./oc login -u system:admin
   20  echo $PATH
   21  cp oc /usr/bin/
   22  exit
   23  oc projects
   24  oc new-project router
   25  oc projects
   26  history
```

Now, spin up a pod...

Wait, need SCC

```
metadata:
  annotations:
    kubernetes.io/description: router privileges - these are fairly loose.
  creationTimestamp: null
  name: router
allowHostDirVolumePlugin: true
allowHostIPC: false
allowHostNetwork: true
allowHostPID: true
allowHostPorts: false
allowPrivilegedContainer: true
allowedCapabilities: null
apiVersion: v1
defaultAddCapabilities: null
fsGroup:
  type: RunAsAny
groups:
- system:cluster-admins
kind: SecurityContextConstraints
priority: 10
readOnlyRootFilesystem: false
requiredDropCapabilities:
- MKNOD
- SYS_CHROOT
runAsUser:
  type: RunAsAny
seLinuxContext:
  type: MustRunAs
supplementalGroups:
  type: RunAsAny
volumes:
- configMap
- downwardAPI
- emptyDir
- persistentVolumeClaim
- secret


```

Create a pod to test...

## OSPF database not up to date over wan?

```
show ip ospf database
```

---

# Using OpenShift playbooks.

---

* Setup non-openshift-side `openshift-vxlan.yml`
* First, create pods `openshift-create-pods.yml`
    - This will fail if you don't delete what's already there.
* Then koko your pods `openshift-koko-pods.yml`

```
quaggab# show running-config 
Building configuration...

Current configuration:
!
!
interface eth0
 link-detect
!
interface lo
 link-detect
!
interface mid2
 ip ospf mtu-ignore
 ip ospf network point-to-point
 ipv6 nd suppress-ra
 link-detect
!
interface out1
 ipv6 nd suppress-ra
 link-detect
!
router ospf
 ospf router-id 3.3.3.3
 redistribute static
 passive-interface out1
 network 2.2.2.0/24 area 0.0.0.0
 network 3.3.3.0/24 area 0.0.0.0
 network 192.168.3.0/24 area 0.0.0.0
 network 192.168.4.0/24 area 0.0.0.0
!
ip forwarding
!
line vty
!
end
quaggab# 

```


```
[root@centos-host centos]# ssh -i ~/.ssh/ajay_aws -R 4789:localhost:4789 ec2-user@54.200.66.51
```


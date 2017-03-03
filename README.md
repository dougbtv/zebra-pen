# zebra-pen

A suite of playbooks and roles to create a demo router VNF, in containers.

## Goal

Have box `CentOS A` ping box `CentOS B` via containerized routers running on Router A & Router B

```
CentOS A -> Router A -> Router B -> CentOS B
```

## Iteration 3: Single VM, specific network interfaces in containers

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



---

# These are previous iterations and can largely be ignored until they're removed...

Thanks! I'm leaving them here in case you want to pick out some pointers / commands from them, for a revision or two.

## Iteration 2: All-in-one-VM

![layout diagram, iteration 2](http://i.imgur.com/99SmRVx.png)

In the second iteration, we run all the containers on a single VM. However, no configuration is made to the routers. Portions of "Iteration 1" are reused.

To develop this, Doug used his own virt-host as specified in "Iteration 1" documentation. To deploy this for help by Ajay to configure the routers, Doug deployed this on Ajay's OpenStack VM instance.

Doug leveraged the `virt-host-setup.yml` playbook (which spins up virtual machines), and then modified the `./vars/all.yml` to change to a single instance.

In order to run the containers against a VM that is already existing (including Ajay's OpenStack VM instance) modify the `./inventory/single_vm.inventory` file to fit your requirements.

And then run the playbook like:

```
ansible-playbook -i inventory/single_vm.inventory containers-on-single-vm.yml 
```

Once the playbook has run, you can see these containers are running:

```
[root@centos-host centos]# docker ps
CONTAINER ID        IMAGE                                  COMMAND             CREATED             STATUS              PORTS               NAMES
5f73798df94d        centos:centos7                         "/bin/bash"         39 minutes ago      Up 39 minutes                           centos_b
c4f7da02d3d8        centos:centos7                         "/bin/bash"         39 minutes ago      Up 39 minutes                           centos_a
f12ea2a547cb        cumulusnetworks/quagga:xenial-latest   "/bin/bash"         41 minutes ago      Up 41 minutes                           quagga_b
56f97629fe5e        cumulusnetworks/quagga:xenial-latest   "/bin/bash"         42 minutes ago      Up 42 minutes                           quagga_a
```

You can enter any container you like with:

```
[root@centos-host centos]# docker exec -it centos_a /bin/bash
[root@c4f7da02d3d8 /]# ping -c 1 4.2.2.2
PING 4.2.2.2 (4.2.2.2) 56(84) bytes of data.
64 bytes from 4.2.2.2: icmp_seq=1 ttl=51 time=53.4 ms

--- 4.2.2.2 ping statistics ---
1 packets transmitted, 1 received, 0% packet loss, time 0ms
rtt min/avg/max/mdev = 53.471/53.471/53.471/0.000 ms
```

And you can go and interactively configure quagga like...

```
[root@centos-host centos]# docker exec -it quagga_a /bin/bash
root@centos-host:/# vtysh

Hello, this is Quagga (version 0.99.23.1+cl3u2).
Copyright 1996-2005 Kunihiro Ishiguro, et al.

centos-host# 
```

---

## How-to -- Iteration 1

There's a few phases here, I would've done this in one shot, but would require some sorcery with dynamic inventories, and to be expedient, I've left it to the user to create a second inventory of hosts. Maybe in the future I'll set it up for a one-shot with a dynamic inventory.

* Phase 1: Configures the virtual machine host and spins up 4 VMs (as above in goal)
* Phase 2: Setup inventory for VM provisioning
* Phase 3: Run playbooks to setup VMs

If you are feeling ambitious, make 4 bare metal machines, and then just use Phase 2 & 3.

### Phase 1: Virthost setup

You'll first need a CentOS 7.3 (with the latest packages) installed on a machine, this machine will need to have virtualization capabilities. This is our virtualization host, which we refer to as virt-host.

Secodarily, you need [ansible installed](http://docs.ansible.com/ansible/intro_installation.html) somewhere, this can also be the virt-host proper if you please. Also, give yourself SSH keys as root to this machine.

Modify the `inventory/virthost.inventory` file to point to the proper IP address for this machine. (On the first line of that file)

There are additional parameters which are flat across VMs, such as amount of RAM and how many CPUs to use, which you can configure if you'd like in the `./vars/all.yml` variables file.

Then you can clone this repo and run the playbook to setup the virt-host with:

```
ansible-playbook -i inventory/virthost.inventory virt-host-setup.yml
```

Take specific note of the section where you get the IP addresses of the virtual machines, we'll use those in phase 2 for provisioning the virtual 
machines.

The section should look approximately like this below, however the IP addresses will likely differ.

```
TASK [vm-spinup : Here are the IPs of the VMs] *********************************
ok: [zebrapen] => {
    "msg": {
        "centos_a": "192.168.122.176", 
        "centos_b": "192.168.122.31", 
        "router_a": "192.168.122.135", 
        "router_b": "192.168.122.9"
    }
}
```

Assuming the playbook has run successfully, you'll see that there's a series of VMs running on the target machine, now.

```
$ virsh list
```

You can of course destroy them if you need to with `virsh destroy $name` and delete with `virsh undefine $name`

To access these, you use the SSH key as generated by this playbook, `/root/.ssh/id_vm_rsa`

### Phase 2: Setup inventories for VMs

Alright, now that you have run phase 1, we've got some info here for us. 

There's 3 keys things you need:

1. The IP address of the virtual machine host as set in `inventory/virthost.inventory` in the first phase.
2. You need the list IP addresses as reported when running the phase 1 playbook.
3. Copy the generated private SSH key from the virtual machine onto the machine you're running ansible from.

Now that you've got those, we're going to modify that `./inventory/vms.inventory` file.

First thing, up at the top there's for boxes, set the IP address for each of those, to the IP address as reported in phase one.

Each line looks about like:

```
centos_a ansible_host=192.168.122.176
```

Just replace the IP address with your own.

Next, we need to configure the jump-host, which is the IP address of the virtual machine host, this goes into the `[all_vms:vars] section`, put it here replacing the IP address that I have committed to the repo.

```
ansible_ssh_common_args='-o ProxyCommand="ssh -W %h:%p root@192.168.1.119"'
```

Lastly, make sure the path to the private key matches the path to where you copied the SSH key down from the virt host, e.g. this line:

```
ansible_ssh_private_key_file=/home/doug/.ssh/id_testvms
```

### Phase 3: Provisioning virtual machines

Now that your inventory is all set, let's go ahead and run the playbook which is going to set up all these hosts...

```
ansible-playbook -i inventory/vms.inventory vm-setup.yml 
```

Note: If you haven't ssh'd to these boxes before you're going to have to say "yes" to the new key fingerprints.

### Let's go verify that it works!

NOTE: There's not a static route yet on CentOS B to go back through the routers, only from CentOS A to statically route requests through CentOS B through the routers. Ping works, but, I'm not sure of the subtleties. Next version for that.

Alright, so first of all, ssh to the virt host. You'll notice there's `/etc/hosts` on the virt-host for each of the machines herein.

Let's ssh to `centos_a` and we'll run a ping there, and we can further this experiment by using `tcpdump` on the routers to see the packets traverse them.

SSH to `centos_a` and run a ping there (remember, we have an ssh key we generated in phase 1 that allows us to ssh there.)

```
[root@virthost ~]# ssh -i .ssh/id_vm_rsa centos@centos_a
[centos@centos_a ~]$ sudo su
[root@centos_a centos]# route -n
Kernel IP routing table
Destination     Gateway         Genmask         Flags Metric Ref    Use Iface
0.0.0.0         192.168.122.1   0.0.0.0         UG    0      0        0 eth0
192.168.122.0   0.0.0.0         255.255.255.0   U     0      0        0 eth0
192.168.122.25  192.168.122.242 255.255.255.255 UGH   0      0        0 eth0

[root@centos_a centos]# ping -c 1 192.168.122.25
PING 192.168.122.25 (192.168.122.25) 56(84) bytes of data.
64 bytes from 192.168.122.25: icmp_seq=1 ttl=62 time=1.38 ms

--- 192.168.122.25 ping statistics ---
1 packets transmitted, 1 received, 0% packet loss, time 0ms
rtt min/avg/max/mdev = 1.385/1.385/1.385/0.000 ms
[root@centos_a centos]# 

```

And if you'd like you can check out the traffic traversing a router, like so:

```
[root@virthost ~]# ssh -i .ssh/id_vm_rsa centos@router_a

[centos@router_a ~]$ sudo su

[root@router_a centos]# tcpdump 'icmp'
tcpdump: verbose output suppressed, use -v or -vv for full protocol decode
listening on eth0, link-type EN10MB (Ethernet), capture size 65535 bytes
21:53:53.761404 IP 192.168.122.196 > 192.168.122.25: ICMP echo request, id 9498, seq 1, length 64
21:53:53.761479 IP router_a.example.local > 192.168.122.25: ICMP echo request, id 9498, seq 1, length 64
21:53:53.762277 IP 192.168.122.25 > router_a.example.local: ICMP echo reply, id 9498, seq 1, length 64
21:53:53.762288 IP 192.168.122.25 > 192.168.122.196: ICMP echo reply, id 9498, seq 1, length 64

4 packets captured
4 packets received by filter
0 packets dropped by kernel

```

And if you'd like to enter quagga and see what's going on there, you can do so by entering an interactive `vtysh` terminal session using:

```
[root@virthost ~]# ssh -i .ssh/id_vm_rsa centos@router_a

[centos@router_a ~]$ sudo su

[root@router_a centos]# docker ps
CONTAINER ID        IMAGE                                  COMMAND             CREATED             STATUS              PORTS               NAMES
e5df9f94de2a        cumulusnetworks/quagga:xenial-latest   "/bin/bash"         20 minutes ago      Up 20 minutes                           quagga

[root@router_a centos]# docker exec -it quagga vtysh

Hello, this is Quagga (version 0.99.23.1+cl3u2).
Copyright 1996-2005 Kunihiro Ishiguro, et al.

router_a.example.local# show ip route
Codes: K - kernel route, C - connected, S - static, R - RIP,
       O - OSPF, I - IS-IS, B - BGP, T - Table,
       > - selected route, * - FIB route

K>* 0.0.0.0/0 via 192.168.122.1, eth0
C>* 192.168.122.0/24 is directly connected, eth0
S>* 192.168.122.25/32 [1/0] via 192.168.122.174, eth0
router_a.example.local# 

```

## Developer notes

Using this article as a [basis for spinning up virtual machines](http://giovannitorres.me/create-a-linux-lab-on-kvm-using-cloud-images.html) from a centos generic cloud image. The meat of the article is [this gist](https://gist.github.com/giovtorres/0049cec554179d96e0a8329930a6d724), which is embedded here as a templated shell script.

# zebra-pen

![nfvpe-bdage](https://img.shields.io/badge/nfvpe-approved-green.svg) ![apache-badge](https://img.shields.io/badge/license-Apache%20v2-blue.svg)

A suite of playbooks and roles to create a demo router VNF, in containers.

## Goal

Have box `CentOS A` ping box `CentOS B` via containerized routers running on Router A & Router B.

```
CentOS A -> Router A -> Router B -> CentOS B
```

## About the architecture

![layout diagram](http://i.imgur.com/KYtqu9R.png)

In the current iteration we're going to use a specialized script, [vethcon](https://github.com/s1061123/vethcon), created by [Tomo](https://github.com/s1061123) which allows us to specifically craft network interfaces for each container. These allow us to create 2 interfaces (one in each container) that are connected. This playbook spins up the diagrammed containers, and creates the interfaces as diagrammed above.

This configuration uses OSPFd in Quagga for routing from Centos A through the two routers (Quagga A & Quagga B) to Centos B, and then back through in the reverse direction.

## Running the playbook

To kick off this playbook, use the inventory file located at `./inventory/single_vm.inventory` as a basis and then run:

```
$ ansible-playbook -i inventory/single_vm.inventory vethcon-single-vm.yml
```

## Verifying the results.

Once that has run, you can verify that it is working by checking the running containers on a host, and then entering the `centos_a` container, and pinging the IP address for `centos_b`.

For example, to ping the IP address for `centos_b` from `centos_a` execute a command like so:

```bash
$ docker exec -it centos_a ping -c2 4.4.4.4
PING 4.4.4.4 (4.4.4.4) 56(84) bytes of data.
64 bytes from 4.4.4.4: icmp_seq=1 ttl=62 time=0.073 ms
64 bytes from 4.4.4.4: icmp_seq=2 ttl=62 time=0.080 ms
```

And vice-versa for a ping to `centos_a` from `centos_b`

```
[root@zebra centos]# docker exec -it centos_b ping -c2 1.1.1.1
PING 1.1.1.1 (1.1.1.1) 56(84) bytes of data.
64 bytes from 1.1.1.1: icmp_seq=1 ttl=62 time=0.098 ms
64 bytes from 1.1.1.1: icmp_seq=2 ttl=62 time=0.065 ms
```

## Further configuration

If you'd like to further configure the routers, without saving the changes to the playbook (and note that these changes may be ephemeral in nature) you can enter the `vtysh` in either the `quagga_a` or `quagga_b` containers, like so:

```
[user@host ~]$ docker exec -it quagga_a vtysh

Hello, this is Quagga (version 0.99.23.1+cl3u2).
Copyright 1996-2005 Kunihiro Ishiguro, et al.

quagga_a# configure terminal 
quagga_a(config)# 
```

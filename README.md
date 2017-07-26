# zebra-pen

![nfvpe-bdage](https://img.shields.io/badge/nfvpe-approved-green.svg) ![apache-badge](https://img.shields.io/badge/license-Apache%20v2-blue.svg)

A suite of playbooks and roles to create a demo router VNF, in containers.

## Goal

Have box `CentOS A` ping box `CentOS B` via containerized routers running on Router A & Router B.

```
CentOS A -> Router A -> Router B -> CentOS B
```

## About the architecture

![layout diagram](http://i.imgur.com/oXQQWjy.png)

In the current iteration we're going to use a specialized script, [koko](https://github.com/redhat-nfvpe/koko), created by [Tomo](https://github.com/s1061123) which allows us to specifically craft network interfaces for each container. These allow us to create 2 interfaces (one in each container) that are connected. This playbook spins up the diagrammed containers, and creates the interfaces as diagrammed above.

This configuration uses OSPFd in Quagga for routing from Centos A through the two routers (Quagga A & Quagga B) to Centos B, and then back through in the reverse direction.

## Running the playbooks

There's basically two styles here -- one of which uses a single host, with all veth connections between docker containers on that single host. The second style uses two VMs and has vxlan between the routers on the two hosts (and veth connections for containers on the same local host).

### Single Host

To kick off this playbook, use the inventory file located at `./inventory/single_vm.inventory` as a basis and then run:

```
$ ansible-playbook -i inventory/single_vm.inventory koko-single-vm.yml
```

### Two-hosts (with VXLAN)

To kick off this playbook, use the inventory file located at `./inventory/vxlan.lab.inventory` as a basis and then run:

```
$ ansible-playbook -i ./inventory/vxlan.lab.inventory vxlan.yml
```

## OpenShift style.

So you wanna use it with openshift?

First, [spin up openshift manually](https://github.com/openshift/origin/blob/master/docs/cluster_up_down.md) -- use `oc cluster up` to make a all-in-one openshift instance for ease-of-use.

Here's how I setup my openshift

```
# Set ip_forward to 1
/sbin/sysctl -w net.ipv4.ip_forward=1

# Install docker (plus a handy wget)
yum install -y docker wget

# Setup docker to allow an "insecure" registry.
sed -i -e "s|\# INSECURE_REGISTRY='--insecure-registry'|INSECURE_REGISTRY='--insecure-registry 172.30.0.0/16'|" /etc/sysconfig/docker

# Start and enable docker.
systemctl daemon-reload
systemctl start docker
systemctl enable docker

# Download the oc command line tool.
wget https://github.com/openshift/origin/releases/download/v3.6.0-rc.0/openshift-origin-client-tools-v3.6.0-rc.0-98b3d56-linux-64bit.tar.gz
tar -xzvf openshift-origin-client-tools-v3.6.0-rc.0-98b3d56-linux-64bit.tar.gz 
cp openshift-origin-client-tools-v3.6.0-rc.0-98b3d56-linux-64bit/oc /usr/bin/
chmod +x /usr/bin/oc

# Check that it's in your path.
oc version

# Bring up the cluster.
oc cluster up

# See that you can get pods (likely nothing there yet)
oc get pods

# Login as admin.
oc login -u system:admin

# Check the cluster (err, AIO) status.
oc status
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

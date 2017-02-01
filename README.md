# zebra-pen

A suite of playbooks and roles to create a demo router VNF.

## Goal

Have box `CentOS A` ping box `CentOS B` via containerized routers running on Router A & Router B

```
CentOS A -> Router A -> Router B -> CentoOS B
```

## How-to

### Phase 1: Virthost setup

You'll first need a CentOS 7.3 (with the latest packages) installed on a machine, this machine will need to have virtualization capabilities. This is our virtualization host, which we refer to as virt-host.

Secodarily, you need [ansible installed](http://docs.ansible.com/ansible/intro_installation.html) somewhere, this can also be the virt-host proper if you please. Also, give yourself SSH keys as root to this machine.

Modify the `inventory/virthost.inventory` file to point to the proper IP address for this machine.

Then you can clone this repo and run the playbook to setup the virt-host with:

```
ansible-playbook -i inventory/virthost.inventory virt-host-setup.yml
```

### Phase 2: Provisioning virtual machines



## Notes from development

Using this article as a [basis for spinning up virtual machines](http://giovannitorres.me/create-a-linux-lab-on-kvm-using-cloud-images.html) from a centos generic cloud image. The meat of the article is [this gist](https://gist.github.com/giovtorres/0049cec554179d96e0a8329930a6d724).
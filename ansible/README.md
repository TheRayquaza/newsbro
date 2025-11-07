# Ansible

## Wireguard Server Setup


Editing the wg variables

```bash
ansible-vault edit group_vars/wireguard_server.yaml
```

Installing wireguard:

```bash
ansible-playbook -i hosts.ini playbook.yaml -v --ask-vault-pass
```
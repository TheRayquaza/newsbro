# Ansible

## Wireguard Server Setup


Editing the wg variables

```bash
ansible-vault edit group_vars/master.yaml
```

Installing wireguard:

```bash
ansible-playbook -i hosts.ini playbook.yaml --tags "wg" -v --ask-vault-pass
```

## DNS

Installing DNS:

```bash
ansible-playbook -i hosts.ini playbook.yaml --tags "dns" -v --ask-vault-pass
```
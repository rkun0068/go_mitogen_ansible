# go_mitogen_ansible
go_mitogen_ansible

## Get Started

### Mitogen

Reference: https://mitogen.networkgenomics.com/ansible_detailed.html

### Ansible

Reference: https://docs.ansible.com/ansible/latest/playbook_guide/playbooks_intro.html#playbook-syntax

### Usage
compile 

```bash
go build -o fast-ansible .
```

optional args(refer to ansible-playbook):
```bash
  -e EXTRA_VARS, --extra-vars EXTRA_VARS
                        set additional variables as key=value or YAML/JSON, if filename prepend with @
  -h, --help            show this help message and exit
  -f HOSTS_FILE, --hosts-file HOSTS_FILE
                        specify hosts file path
  -t TAGS, --tags TAGS  only run plays and tasks tagged with these values               
```
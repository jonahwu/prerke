# Pre Rke

## Why

1. To create new user, sshkey for new user and docker install.
2. disable firewalld and seclinux
3. setup docker permission for new user 

## Todo

1. deploy rr in deploy server(master?)
2. deploy rr for joined k8s cluster

## condition

1. root user
2. root can be no need sshkey, will not affect the root.
3. root need password for remote aceessing
4. pentium user will generate their own sshkey 
5. if deploy server have no pentium user, you need to create it first (put ip of deploy server in precluster.yaml)

## Example

for all new servers, 172.16.155.241 (will be master), and 172.16.155.242 with root user.
where the server's root password is promise shown as followed.


for the precluster.yml

```

#image install:  curl https://releases.rancher.com/install-docker/17.03.sh | sh
password: promise
nodes:
  - address: 172.16.155.241
    info:
      - a
      - b
      - c
    dockerversion: https://releases.rancher.com/install-docker/17.03.sh
  - address: 172.16.155.242
    info:
      - a
      - b
      - c
    dockerversion: https://releases.rancher.com/install-docker/17.03.sh
```

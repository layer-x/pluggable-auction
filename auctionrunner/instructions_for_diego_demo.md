instructions for running diego w/ multi-brain:

**0 - Optional Commands for showing what's running**

- Show running "VMs" (really warden containers running inside of bosh-lite)
- install BOSH-CLI:
 - `gem install bosh_cli`
- run `bosh target 54.67.93.53`
    -  username: `admin`, password: `admin`
-  run `bosh vms`
    * you will see output like the following:

```
Acting as user 'admin' on 'Bosh Lite Director'
Deployment 'cf-warden'


Director task 148

Task 148 done

+------------------------------------+---------+---------------+--------------+
| Job/index                          | State   | Resource Pool | IPs          |
+------------------------------------+---------+---------------+--------------+
...
| uaa_z1/0                           | running | medium_z1     | 10.244.0.130 |
+------------------------------------+---------+---------------+--------------+

VMs total: 12
Deployment `cf-warden-diego'

Director task 149

Task 149 done

+--------------------+---------+------------------+---------------+
| Job/index          | State   | Resource Pool    | IPs           |
+--------------------+---------+------------------+---------------+
| access_z1/0        | running | access_z1        | 10.244.16.6   |
| brain_z1/0         | running | brain_z1         | 10.244.16.134 |
| cc_bridge_z1/0     | running | cc_bridge_z1     | 10.244.16.142 |
| cell_z1/0          | running | cell_z1          | 10.244.16.138 |
| cell_z1/1          | running | cell_z1          | 10.244.16.154 |
| cell_z1/2          | running | cell_z1          | 10.244.16.150 |
| database_z1/0      | running | database_z1      | 10.244.16.130 |
| route_emitter_z1/0 | running | route_emitter_z1 | 10.244.16.146 |
+--------------------+---------+------------------+---------------+

VMs total: 8

```
- see running apps in cf: `cf apps`
- delete running app in cf: `cf delete <app_name>`
- scale running app in cf: `cf scale <app_name> -i <number_of_instances>`
- view a running app through your web browser: go to http://APP_NAME.54.67.93.53.xip.io/ (i think this is a really good idea)

***1 - SSH to Bosh-Lite VM (running on Amazon)***

- create ssh key file (only need to do this once):

* copy and paste the following into a file called `bosh.pem`:
*
```
-----BEGIN RSA PRIVATE KEY-----
MIIEogIBAAKCAQEAhT8+kOqORiTUsGyj+KDZojJ1H+T8MLM+27Hvqo2KGEH+il+ZEypVJ2r/9Vfc
5ooThuFs78uRvvHHwjgTEJBSg+GZVOCF5OGCpV35Wrnpap3DinhF5frdQ45YJLre09PQxldFnltX
UHwSSVOfbbd1YL1ZXhp2Dbgv9MGrABr+G+kCzhG40Fo1s7ye+k/AzsDK2gNjV3hN077r7eUz6FD1
54kcPNfFjQ7+ONrhpb883RnaCsOb4de7TqZWTrvW9fDfyuNeFowakhOwTNkDM7FjPeRgsUvRkX1y
QdEKQQuiOH/HFqjObtKaSr7KLHBzC8JFSEPrplJMh1asrc99kivYIQIDAQABAoIBAHSK6YVEsiXl
xuV8UDqBLXkxGsJHvNA3pq3vRslsvLEU37ZVgQSDTTGJ48/KBprZf9TETEy8R4Cz5l0YQIyHPrS8
2CilrFaRa3yJ4jQZUXAABuyQ38oUDf0tfii6DXVG2V7xLCIikA8ERdY+vr3u7UosswKcsE61n7Q6
w/72nPT5Jw+sHmDyHuf8EpRLx+HCcOEYTp7/J1yIabWB6R+2qvt4bzupNeTNcAcouOljdI7ffbD4
kq57JfpvhMZUNklWJD2MJs9orz69M8fRChwMs3tHV8CkRkpNwPCweoWrDI7zPElnU5nDVZVFNmgz
vGDO3RKjch49QhPVrGI9/8+qkyUCgYEAyrpZnEi411sAyyAx0TDQlB+tz3KhAokylppJh13nAvue
7doZOIxhVtZi6h5MU8agrl3n5DqUd0j12KzhhVtu/phBkE+FwhE9FFLU4t3Ed9AJE3gZcKoaZIvG
nOuiunTHN8uCXYOzXteHV+LKW229BKbpseATkxVzLT9yvjhXVUcCgYEAqELd48dxosHsXefe781W
/2ic+UN5RKfwd1Tnq59SDD5bUA4XozryHXcNy2dc6Sz22jiJweuvDpCtDQnlBEjMbehK3ZW4t8j+
plnqwBmoFy7xnTFs+m1YwhbpnAII1aQljaSJNDJ/x0DhXmO0CGKxK4/CHLgxbEncPm/JBvq1u1cC
gYBmBwlIXVUhlTw9/nLz/CRNF/Bqwh8EXrYmE3pD9V9pIeenfydITWZDxNu9RghV9VYyyzIEq/LC
YebQ6JkLe6vN2CTPEyaXOAPMca+QidnyDrIyqTPsfr+PsMUBfpnESzdj/jkbBUhFyCTmd04uW3lQ
mQxuJ/7R/G6d7Bu8XjCdywKBgDMWKjyQP4ZFDrjsP5nbZICjiJV90QHxY2c31icbdlPVUvAZdz/O
E9iyXvPU7Da3ujNDW0APiNUJRCFjUa9dUwRDtQdWAAF8+yQSxN2SbKCtVhp9+TKHpJ05S7BcRcZn
0icRP78jXfxnTIXWC8FIBbbOLQd/PTI9sqsaUZTW5fp9AoGAH4OnicCxpITR063Ph40Xp7zWmpQT
a0kHkTGyH/yLnpVHrb19ttlWAgPQT3ZKUQgw8vj2WYTAUwWJTBpRtdiLNXJKDpjvosPGspAsYjGj
9ReQXuccqoRit3DHCmVALUSUozeAl/IjzoKeAP2zU3DZrthBgNEo+bGy5on681iK8hw=
-----END RSA PRIVATE KEY-----
```

  - the `cell_z1` vms are the cells; these are not modified
  - the `brain_z1` vm is where the modified auctioneer lives. when we register brains to the auctioneer, we need send the `cURL -X POST` to this ip `10.244.16.134`

* run the command `chmod 400 bosh.pem`

* `ssh -i bosh.pem ubuntu@54.67.93.53`

***2 - Run a brain: (from ssh session)***

- Run the Passive (observe-only) UI Brain: `~/ui_brain`

- Run the Interactive UI Brain Brain: `~/active_ui_brain`

- Run the Diego Brain: `~/original_diego_brain`

- Run the Fenzo Brain: `java -jar ~/fenzo_brain-1.0-SNAPSHOT-jar-with-dependencies.jar`

***2 - Run a brain: (from ssh session)***

- Register the Passive UI Brain with the Auctioneer:
    - run:

    ```
      curl -X curl -X POST 10.244.16.134:3000/Start \
      -d '{"name":"passive_ui_brain","url":"http://172.31.29.198:3333","tags":"ui"}'
    ```

    - this sends a POST request with some json telling the auctioneer the following information about the brain:
        - _name_ = "passive_ui_brain"
        - _url_ = http://172.31.29.198:3333
        - _tags_ = ["ui"]
   - view the passive ui at http://54.67.93.53:3333

- Register the Active UI Brain with the Auctioneer, using the tag "default":
  - run:

    ```
      curl -X curl -X POST 10.244.16.134:3000/Start \
      -d '{"name":"active_ui_brain","url":"http://172.31.29.198:4444","tags":"default"}'
    ```

  - view the active ui at http://54.67.93.53:4444/brain

- Register the Diego Brain with the Auctioneer, using the tags "default", and "diego":
    - run:

    ```
      curl -X curl -X POST 10.244.16.134:3000/Start \
      -d '{"name":"diego_brain","url":"http://172.31.29.198:6666","tags":"default,diego"}'
    ```

- Register the Fenzo Brain with the Auctioneer, using the tag "fenzo":
  - run:

    ```
      curl -X curl -X POST 10.244.16.134:3000/Start \
      -d '{"name":"fenzo_brain","url":"http://172.31.29.198:5555","tags":"fenzo"}'
      ```

***3. Push an App with CloudFoundry (from your laptop)***
- Install the cf cli if you don't already have it:
  - `brew install cf-cli`
- Log in to the cf instance running on our bosh-lite:
  -  `cf login -a api.54.67.93.53.xip.io  -u admin -p admin --skip-ssl-validation`
- Pull the example_cf_apps to your laptop

  `git clone https://github.com/EMC-CMD/cf-example-apps.git`

  `cd example_cf_apps/exampleapp`

- Push an app to the _default brain_:
  - `cf push --no-start`
  - `cf enable_diego exampleapp`
  - `cf start exampleapp`
  - (optional) `cf scale exampleapp -i 3` (replace 3 with # of instances you want)

- Push an app with a tag:
  - edit the `manifest.yml` file in the *exampleapp* directory to look like the following:

```
---
applications:
- name: ANY_NAME_YOU_WANT
  env:
    DIEGO_BRAIN_TAG: ANY_TAG_YOU_WANT
```
  - `cf push --no-start`
  - `cf enable_diego SOME_OTHER_NAME`
  - `cf start SOME_OTHER_NAME`

###Backup Environment:

*1 - In case the environment fails, there is an identical backup environment running on AWS.*

*2 - To use it:*

  a - replace all instances of the ip `54.67.93.53` with `54.183.111.189` (for `ssh ...`, `bosh target ...` and `cf login ...`)

  b - replace the bosh-lite private ip `172.31.29.198` with `172.31.0.196` (for commands that assign Brains to diego)

  c - and replace the ssh key in `bosh.pem` with:

```
-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEAmT+wbEFqGmBglzFJlxid3Vk1zMsznuG3eb32RP6P65JQysh+w65AwYz6BiIx
3hj8Y7P3TsznKCjUY/z8jXve3ydoAEh62rSD/BZNeyYCHmADQEwyNORQU/AplUPKdsqYoVXadwAB
Pmd4HJ84KBxzMZWyFGtzqM0QhbUY0wuxqcy0SwRW8pIo6/BSHWBKKyEMWR/2f9SG6j2g8PR0/iwU
9QXqrGSQVKjR0wDP7zwSGTrdA08I7a/KmbgQyOWLOz8xL5jxwJ7pXPjA4oyfagjAC2SN/9FFm3LX
sy+hS6f8f3zV+fEy29k8o6mH3WGaScrmlXGShnhgcygrt0mwHAl8UwIDAQABAoIBAQCPIN3MbdgE
JIlyDFV36kTezAgkapCezp+G87WDwOF4GiKdEl7ase/HFb0aZ2t9zIZFNHtBPLyUVHXxoQHbvpps
fyhyQz+C7l/q3IWnA9ustO20arXlkmu3ybF8uGDrS9L7s+yjgfynZQnYaZiQVen8oJw+2BCg0k2h
I3+49M4NEC7qxguIorTplcJI7dXpXjOy0MZ4D0QUTOBAKEzx99Gx4dDB0vDI9kznJXx0kXlQw1KU
4gD/KkIIrWDQXOmfNme4s1pn6e0cAaCfEic69MMDZ/NvG0AEvvZZjNPsExoltlqPb62Gh/CNoAlj
XAFQM60Li+uO4rqkSqPPOTSNNjohAoGBAOKPRF5LGm6uYFi/svBhpl8naxy71WDSKGw0bF8gv6ML
xudTuN62R/ggEq8vciNJdhpPlOA4ewFTJXIAJECYe4T5XntvmfSVXd/3I3A223+5DV1QEWKh71wP
aGyyh7wNV7VWReNSvq0zv6giQqALKanPHwhjF8ALcexfixnDqPKjAoGBAK0pqc2oBd9FKLP5FgDR
fxJrpsi5GkeGlxVsOSJMTdGXyrxY/FHFbc/rLGFbyK/eqvQpg2/mWSWpKZtRbaopPZxQafx6ZgeF
cvNrZ9HU7TXqR1XdEj7peWYZKiI/CHHUAbfd+7eiaNIjx26Q+kL89T1dDUDwVIVB052ipZkzNpqR
AoGBAIv7s2WDiABtE7CiOYCXBUHzzBXD5PJex4Ub2v3n8SBfzXTu4OISxGMGBiVh7mbpI+Tb2QO1
QiMuaYuHlN6omGEv5vXjnb9mbstMGwRhkLvY7e4C48sKfSdnicDnikBiChhMBwCPBqtjtv6+tGXI
n+SAyg7XkzwgljJTUlIH96J7AoGBAKoJeWYrEekWXkur0kFndmI+N35u1TFbJkyxAsF9MAUaCsg8
kTgyqAw9IE1R9ZVND43GnfxpsyxaGjMcGJW4/XjbNdfo0PudvSzuUPopHe2NahMUjHAej0kEeO07
/CzaQ/2rCxxdbJS88X7O+hCBmMdy8irMVBKuewAV0IrJUVshAoGAM1Ws17TBurocG8+jNHjuCpJw
pGvwvmJpwojl6D4+2a2yoP9l2zFQRxbZroFM7rFp/7qwyVxIob5oWkJh/NhEN5iW3UqwhZRyFxq1
UD8cA25rhOiV1Nhl+pvr5yRolMvhyaVOqCXNWyfGROOKN2HtHf6+q1yRSIFPbA0S0VvXHkI=
-----END RSA PRIVATE KEY-----
```

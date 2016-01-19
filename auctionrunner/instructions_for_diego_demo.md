instructions for running diego w/ multi-brain:

**0 - Optional Commands for showing what's running**

- Show running "VMs" (really warden containers running inside of bosh-lite)
- install BOSH-CLI:
 - `gem install bosh_cli`
- run `bosh target 54.183.237.92`
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
- view a running app through your web browser: go to http://APP_NAME.54.183.237.92.xip.io/ (i think this is a really good idea)

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

* `ssh -i bosh.pem ubuntu@54.183.237.92`

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
      -d '{"name":"passive_ui_brain","url":"http://172.31.27.58:3333","tags":"ui"}'
    ```

    - this sends a POST request with some json telling the auctioneer the following information about the brain:
        - _name_ = "passive_ui_brain"
        - _url_ = http://172.31.27.58:3333
        - _tags_ = ["ui"]
   - view the passive ui at http://54.183.237.92:3333

- Register the Active UI Brain with the Auctioneer, using the tag "default":
  - run:

    ```
      curl -X curl -X POST 10.244.16.134:3000/Start \
      -d '{"name":"active_ui_brain","url":"http://172.31.27.58:4444","tags":"default"}'
    ```

  - view the active ui at http://54.183.237.92:4444

- Register the Diego Brain with the Auctioneer, using the tags "default", and "diego":
    - run:

    ```
      curl -X curl -X POST 10.244.16.134:3000/Start \
      -d '{"name":"diego_brain","url":"http://172.31.27.58:6666","tags":"default,diego"}'
    ```

- Register the Fenzo Brain with the Auctioneer, using the tag "fenzo":
  - run:

    ```
      curl -X curl -X POST 10.244.16.134:3000/Start \
      -d '{"name":"fenzo_brain","url":"http://172.31.27.58:5555","tags":"fenzo"}'
      ```

***3. Push an App with CloudFoundry (from your laptop)***
- Install the cf cli if you don't already have it:
  - `brew install cf-cli`
- Log in to the cf instance running on our bosh-lite:
  -  `cf login -a api.54.183.237.92.xip.io  -u admin -p admin --skip-ssl-validation`
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

  a - replace all instances of the ip `54.183.237.92` with `54.208.251.28` (for `ssh ...`, `bosh target ...` and `cf login ...`)

  b - replace the bosh-lite private ip `172.31.27.58` with `172.31.5.101` (for commands that assign Brains to diego)

  c - and replace the ssh key in `bosh.pem` with:

```
-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAm4x/UNE5P5LY3ISWUjg6Bm3Vac/OS6QnuVdc12z45o3U28SrSaHYMyyWAy3U
yMHyPbtu4wVHwdRcIWfqxtkVVnDdehFA6srgnAYRW2JSKXDWCqWPrKgHdJ1Xx7jes2tUP54WiVZs
bdkdjjfcHmXls/AOEuVU9G7PM/q2nKMtp6/Wklu/ARWBmycmM7PxOwr70TsDFYg4pNzh8rnpLd1n
gIFeLsw6cGcfGd/GH0DWiTrHIwnaRcg511Sl7bRvzzLfR4HeADTacs7u2EkdZDkUgsGFjXBqFnqK
PKNlxc3WpRPU4V6Yij3GDhhwfDTv7B51W28Ve0vTt+7ayaayj+dG/QIDAQABAoIBAHBIfy8LmNO3
YSvt2cUIKXqyljeHdldL7BDya2Zml1V2VI0/7pV8auCl8rPgxZUVy5OcVXMzQJU+gjLrHKLl2W1I
k9el1MKoKHL4PldFJiIb/aY51PjBYoBfhBn77WZ+t5YkvAfvht3UuG0NDawzyhiV1NL3ENhRlOjk
tiVj9XTxbU4MdKK7lhUVbONq/Z0eC3FNQc8VR4Jlw1Qlqz5Edx4Uihef3R7JpLPF2eVmaMtUlLwR
yrWrYKOM5vCZyXxSOi4BZPRmM06lxOJFFv1cyxfBofu9DjaGx/KJLzGG0E9T40m8+FUZ285E/+os
Dyl+9rcPRgWhwuqROEsVm+0/MNECgYEAx5KmKDC7njAZWjxqap/77NwHpVv5RGB6Qhs3dHNJf8Lk
/8AsvyVlHSY1Hc8owEWRVtVO33ZHzxPnnB/GISYARyOCteArmQcgAtgSvBO3tyC9WUUAQy2Saq8L
Rae/OKfm34+FEXKX6HaJQ3NakTuQGN/xlwgnvCjQhtei0kzBgkMCgYEAx4dUbhpxLKQSDNqxviWt
f+gwH8YeH9dOqHJEyqPD3Sn4MAuQJxP98OjD4xAAHRhG7h/vzROm0zrVz2WjtLqqwYZ8S1qoDDYx
FABU5LBPGgqFj0gTSGA18lm0QU6IKu6g/Cwp6E/viqVt+1HR2yRpkB4JS6+6Q8dfxKEuU/aMnb8C
gYEAiaepXAdhIddjZU5OyITZK6MI0xIBeRxit742Heh3RdyUP6O6OY39lIGKGamOHjDd8trmsFPR
bA/6rUFtU+f2QRtJSVH6QG8dsViAc6HWEkZO1Ig3ih6g410hlUYDK30ETiecTVCRXxKD0zZ5vbsr
xTySUu6ZGbu9OYT7FbtDrikCgYA+bYgsHtfUKM2A+hfsr2s2ftY3ysv4GGyC5aXCZTTOCOifV67V
mzqz2pAXhhUTBVqD/LgRyRlEM79b8agjztfITySqiwXTNE1svaHSH5vQQQSCzQFDft7CIfD1EfYm
wJzb6ZF/HyuKjLH5lSL81sq0jcFIzgWQWVwMcIXHPXfHjwKBgQCw8lQVUeisVSU945Dj2yalnrO8
9JsZmsQ/TWz3wuRNym7o59rZR60Sydc1wJe2mReDfZs/ZnNyUqCWmACfHr6oKgNlIqIFnLQ8afTC
yZJswwIvpubcthCCCTxN/LjsVWdc7WWDabzfAAYKwX3kTGiqMZRSCXamsiO4JSDB2OOvzw==
-----END RSA PRIVATE KEY-----
```

# Klouds.io / Klouds.org

#####We like to leverage automation:

[![wercker status](https://app.wercker.com/status/fcf719ade20c4ab01184d966f4650ee2/s/master "wercker status")](https://app.wercker.com/project/bykey/fcf719ade20c4ab01184d966f4650ee2)  [![Join the chat at https://gitter.im/Superordinate/klouds](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/Superordinate/klouds?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)  [![Code Climate](https://codeclimate.com/github/Superordinate/klouds/badges/gpa.svg)](https://codeclimate.com/github/Superordinate/klouds)  


<img src="http://www.ozzadar.com/klouds.png" align="center"/>



A web framework for dynamically launching and accessing web servers.

Uses rest APIs to manage containers. Uses wercker to automate deployment.  (maybe)


Traefik to handle container routing for front facing web services.

##HOW TO USE:

###Environment Variables

MARATHON_ENDPOINT= 192.168.3.4:8080  	<-- Point this to your own marathon backend
POSTMASTER= postmaster@yourdomain.com  		<-- username for your mail server
POSTMASTER_KEY= password				<-- email for your mail server

All emails sent through port 2525 at the moment. Google containers block sending mail on default ports to prevent spammers. Works with MailGun.

###To build:

Dependencies:
Mercurial

```
	apt-get install mercurial
	go get .
	go build .

```

### To Run

``` 
	POSTMASTER=postmaster@yourdomain.com POSTMASTER_KEY=password MARATHON_ENDPOINT=192.168.3.4:8080 ./klouds

```

OR

```
	export POSTMASTER=postmaster@yourdomain.com
	export POSTMASTER_KEY=password
	export MARATHON_ENDPOINT=192.168.3.4:8080

	./klouds

```


<a href="https://github.com/EmileVauge/traefik"><img src="https://camo.githubusercontent.com/0d83f4ec95b28ecc0353078ca4364bf461b99c2d/687474703a2f2f7472616566696b2e6769746875622e696f2f7472616566696b2e6c6f676f2e737667" align="center" height="96" width="200" ></a><br>

[Traefik on Github](https://github.com/EmileVauge/traefik "Traefik on Github")





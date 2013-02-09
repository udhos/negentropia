
REQUIREMENTS

# Install Go Language
# http://golang.org/doc/install

# Install GIT - required to clone the source tree
# http://git-scm.com/downloads

# Install Mercurial - required for packages such as: go get code.google.com/p/goauth2/oauth
# http://mercurial.selenic.com/downloads/

# Install Redis
# http://redis.io/download
#
# Microsoft Official Redis for Windows at Microsoft Open Tech:
# https://github.com/MSOpenTech/redis
# https://github.com/MSOpenTech/redis/blob/2.4/msvs/bin/release/redisbin64.zip?raw=true
# https://github.com/MSOpenTech/redis/blob/bksavecow/msvs/bin/release/redisbin.zip
#
# Redis 2.4.5 for Windows 7 64-bit:
# https://github.com/dmajkic/redis/downloads


GENERAL BUILDING GUIDELINES:

# 1. Clone the git repository:
git clone https://code.google.com/p/negentropia/

# 2. Set GOPATH to negentropia\webserv

# 3. install goauth2
#
# http://code.google.com/p/goauth2/					OAuth 2.0 for Go clients. Doesn't work with Facebook: http://code.google.com/p/goauth2/issues/detail?id=4 
# https://github.com/robfig/goauth2					A fork of code.google.com/p/goauth2 that supports Facebook
# https://github.com/HairyMezican/goauth2/			This is mostly copied from http://code.google.com/p/goauth2/
#													The original code will fail when contacting facebook; this code fixes that problem
# http://code.google.com/r/jasonmcvetta-goauth2/	This clone contains changes to stock Goauth2 detailed by Ryan.C.K. here: http://code.google.com/p/goauth2/issues/detail?id=4
go get github.com/HairyMezican/goauth2/oauth

# 4. Install redis client library
go get github.com/vmihailenco/redis

# 5. Build and install (to negentropia\webserv\bin)
go install negentropia\webserv

# 6. Start redis
redis-server

# 7. Run
# Under Linux:
negentropia\webserv\bin\webserv
# Under Windows:
negentropia\webserv\bin\webserv.exe


BUILDING UNDER WINDOWS:

## windows git bash:

export DEVEL=/c/tmp/devel
export GOPATH=$DEVEL/negentropia/webserv

mkdir -p $DEVEL
cd $DEVEL

git clone https://code.google.com/p/negentropia/

# fetch from github with git bash
go get github.com/vmihailenco/redis
go get github.com/HairyMezican/goauth2/oauth

# fetch from google code (could be under cmd prompt)
go get code.google.com/p/go.net/websocket

## windows dos prompt:

@rem build
\tmp\devel\negentropia\win-build.cmd


CONFIGURING

The win-run-webserv.cmd script reads config from the following file:
	\tmp\devel\config-webserv.txt

You should start by copying the following example configuration:
	\tmp\devel\negentropia\config-webserv-sample.txt
	
	copy \tmp\devel\negentropia\config-webserv-sample.txt \tmp\devel\config-webserv.txt

Then tweak \tmp\devel\config-webserv.txt

If you want to enable local accounts:
	Provide the following SMTP mail relay server information
	(used to confirm new users' email adresses):
	-smtpAuthUser=user@exampledomain
	-smtpAuthPass=putPasswordHere
	-smtpAuthServer=smtp.exampledomain.com
	-smtpHostPort=smtp.exampledomain.com:587

If you want to enable support for Google login:
	1. Login to https://code.google.com/apis/console/
	2. Create a new project: API Project -> Create
	3. Under "API Access", create a "Client ID for web applications"
	4. Under "Client ID for web applications", notice the fields "Client ID" and "Client secret"
	5. Add the following lines to webserv-config.txt:
	-gId=putGoogleClientIdHere
	-gSecret=putGoogleClientSecretHere
	6. Under "Client ID for web applications", add the following URLs:
	Redirect URIs:		http://localhost:8080/ne/googleCallback
	JavaScript origins:	http://localhost:8080

If you want to enable support for Facebook login:
	1. Login to https://developers.facebook.com/apps
	2. Register as a developer
	3. Create a new app
	4. Select the app on the left menu, then notice the fields "App ID/API Key" and "App Secret"
	5. Add the following lines to webserv-config.txt:	
	-fId=putFacebookAppIdHere
	-fSecret=putFacebookAppSecretHere
	6. Click on "Edit Settings"
	7. Check the box "Website with Facebook Login"
	8. Add http://localhost:8080/ne/facebookCallback to the field "Site URL"


RUNNING / TESTING UNDER WINDOWS

@rem run:
\tmp\devel\negentropia\win-run.cmd

Point your web browser to http://localhost:8080/ne/

--THE END--

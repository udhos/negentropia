
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

# Install Java Runtime - needed for the Dart Editor
# http://java.com/en/download/
#
# Install Dart Editor
# http://www.dartlang.org/tools/editor/


GENERAL BUILDING GUIDELINES:

# 1. Clone the git repository:
export DEVEL=/c/tmp/devel
mkdir -p $DEVEL
cd $DEVEL
git clone https://code.google.com/p/negentropia/

# 2. Set GOPATH to negentropia/webserv
export GOPATH=$DEVEL/negentropia/webserv

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

# 5. Install websocket library
go get code.google.com/p/go.net/websocket

# 6. Build and install server (to negentropia/webserv/bin)
go install negentropia\webserv
go install negentropia\world

# 7. Build client
set DART_SDK=c:\dart\dart-sdk
set DEVEL=c:\tmp\devel
set NEG_DART_SRC=%DEVEL%\negentropia\wwwroot\dart
set NEG_DART_MAIN=%NEG_DART_SRC%\negentropia_home.dart
cd %NEG_DART_SRC%
%DART_SDK%\bin\pub install
%DART_SDK%\bin\pub update
cd \
%DART_SDK%\bin\dart2js -c -o %NEG_DART_MAIN%.js %NEG_DART_MAIN%

# 8. Start redis
redis-server

# 9. Configure server
	copy \tmp\devel\negentropia\config-common-sample.txt \tmp\devel\config-common.txt
	copy \tmp\devel\negentropia\config-webserv-sample.txt \tmp\devel\config-webserv.txt
	copy \tmp\devel\negentropia\config-world-sample.txt \tmp\devel\config-world.txt

# 10. Run server
negentropia\webserv\bin\webserv
negentropia\webserv\bin\world


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

win-run-webserv.cmd script reads config from:
	\tmp\devel\config-common.txt
	\tmp\devel\config-webserv.txt

win-run-world.cmd script reads config from:
	\tmp\devel\config-common.txt
	\tmp\devel\config-world.txt

You should start by copying example configurations:
	\tmp\devel\negentropia\config-common-sample.txt
	\tmp\devel\negentropia\config-webserv-sample.txt
	\tmp\devel\negentropia\config-world-sample.txt
	
	copy \tmp\devel\negentropia\config-common-sample.txt \tmp\devel\config-common.txt
	copy \tmp\devel\negentropia\config-webserv-sample.txt \tmp\devel\config-webserv.txt
	copy \tmp\devel\negentropia\config-world-sample.txt \tmp\devel\config-world.txt

Then tweak:
	\tmp\devel\config-common.txt
	\tmp\devel\config-webserv.txt
	\tmp\devel\config-world.txt

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

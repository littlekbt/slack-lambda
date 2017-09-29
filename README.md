# slack-lambda

slack-lambda build lambda environment for your slack team.

[gif]

If you install slack-lambda for your slack team, you can execute function on slack, register function, call registered function.  
You use reminder and call registered function feature on slack-lambda, You can schedule execute function. (as if cron)

[gif]


[構成図]

## Features
- execute function.(now support golang, ruby)  
- register function.  
- call registered function.  
- show registered function.  
- show list registered function.  
- remove registered function.  

## Build
You have to build three deamon.  
The one is RTM bot for slack.  
The one is proxy server that converts request body to json and posts it to lambda server.  
The one is lambda server that build image and run container.

```
$ git clone https://github.com/littlekbt/slack-lambda
$ cd slack-lambda
$ make build
// create binary proxy server and lambda server
$ make install
// create symbolic link to /etc/systemd/system
$ vim bot/config/vars
// you must write vars used by bot
$ make run
// start! 
```

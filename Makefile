SLACK_LAMBDA_HOME=/usr/local/src/slack-lambda

build:
	sh script/build.sh
install:
	sudo ln -f -s $(SLACK_LAMBDA_HOME)/systemd/slack-lambda-bot.service /etc/systemd/system/
	sudo ln -f -s $(SLACK_LAMBDA_HOME)/systemd/slack-lambda-proxy.service /etc/systemd/system/
	sudo ln -f -s $(SLACK_LAMBDA_HOME)/systemd/slack-lambda-container.service /etc/systemd/system/
run:
	sudo systemctl start slack-lambda-bot
	sudo systemctl start slack-lambda-proxy
	sudo systemctl start slack-lambda-container

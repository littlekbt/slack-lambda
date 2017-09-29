require 'slack-rtm-bot-helper'
require 'faraday'
require 'json'
require 'time'
require 'uri'
require 'pathname'
require "#{ENV['SLACK_LAMBDA_BOT_PATH']}/lib/slack-lambda-bot/user_file"
require "#{ENV['SLACK_LAMBDA_BOT_PATH']}/lib/slack-lambda-bot/message_converter"
require "#{ENV['SLACK_LAMBDA_BOT_PATH']}/lib/slack-lambda-bot/post_proxy"
require "#{ENV['SLACK_LAMBDA_BOT_PATH']}/lib/slack-lambda-bot"

SlackLambdaBot.new(ENV['SLACK_LAMBDA_TOKEN']).run

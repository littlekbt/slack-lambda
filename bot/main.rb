require 'slack-rtm-bot-helper'
require 'faraday'
require 'json'
require 'time'
require 'uri'
require 'pathname'
require './lib/slack-lambda-bot/user_file'
require './lib/slack-lambda-bot/message_converter'
require './lib/slack-lambda-bot/post_proxy'
require './lib/slack-lambda-bot'

Process.daemon
SlackLambdaBot.new(ENV['SLACK_LAMBDA_TOKEN']).run

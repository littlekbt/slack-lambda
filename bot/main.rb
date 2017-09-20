require 'slack-rtm-bot-helper'
require 'faraday'
require 'json'
require 'time'
require 'pry'
require './lib/slack-lambda-bot'
require './lib/slack-lambda-bot/user_file'
require './lib/slack-lambda-bot/post_proxy'

SlackLambdaBot.new(ENV['SLACK_LAMBDA_TOKEN']).run

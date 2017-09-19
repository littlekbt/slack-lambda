require 'slack-rtm-bot-helper'
require 'faraday'
require 'json'
require 'time'

class SlackLambda
  REGEXP = /language:\s?.*\nversion:\s?.*\n```\n?[\s\S]*\n?```/

  def initialize(token)
    @token = token
  end

  def run
    puts "#{[Time.now().to_s]} open web socket connection"
    Slack::Rtm::Bot::Helper.run(token=@token) do |data|
      text = data['text']
      return unless valid(text)
      puts "#{[Time.now().to_s]} #{data}"

      conn = Faraday.new(:url => 'http://localhost:8080') do |faraday|
        faraday.request  :url_encoded             # form-encode POST params
        faraday.adapter  Faraday.default_adapter  # make requests with Net::HTTP
      end

      res = conn.post do |req|
        req.body = text
      end

      puts "#{Time.now().to_s} <- #{res.body}"

      j = JSON.parse(res.body)
      j['error'] != '' ? output_error(j['error']) : j['stdout']
    end
    puts "#{[Time.now().to_s]} close web socket connection"
  end

  private

  def valid(text)
    text.match?(REGEXP)
  end

  def output_error(err_msg)
    <<ERR
Error Found!
#{err_msg}
ERR
  end
end

SlackLambda.new(ENV['SLACK_LAMBDA_TOKEN']).run
# SlackLambda.run

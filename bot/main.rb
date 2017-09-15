require 'slack-rtm-bot-helper'
require 'faraday'
require 'json'

class SlackLambda
  REGEXP = /language:\s?.*\nversion:\s?.*\n```\n?[\s\S]*\n?```/

  def initialize(token)
    @token = token
  end

  def run
    Slack::Rtm::Bot::Helper.run(token=@token) do |data|
      text = data['text']
      return unless valid(text)

      conn = Faraday.new(:url => 'http://localhost:8080') do |faraday|
        faraday.request  :url_encoded             # form-encode POST params
        faraday.adapter  Faraday.default_adapter  # make requests with Net::HTTP
      end

      res = conn.post do |req|
        req.body = text
      end

      j = JSON.parse(res.body)
      j['stdout']
    end
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

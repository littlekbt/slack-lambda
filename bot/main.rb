require 'slack-rtm-bot-helper'
require 'faraday'
require 'json'
require 'pry'

class SlackLambda
  def initialize(token)
    @token = token
  end

#   TEXT =<<EOF
# language: ruby
# version: 2.3.0
# ```
# def hello
#   "hello"
# end
# 
# def world
#   "world"
# end
# 
# puts hello + world
# ```
# EOF

  def run
    Slack::Rtm::Bot::Helper.run(token=@token) do |data|
      text = data['text']
      return unless valid(text)

      json = build_json(text)

      conn = Faraday.new(:url => 'http://localhost:8080') do |faraday|
        faraday.request  :url_encoded             # form-encode POST params
        faraday.adapter  Faraday.default_adapter  # make requests with Net::HTTP
      end

      res = conn.post do |req|
        req.headers['Content-Type'] = 'application/json'
        req.body = json
      end

      res.body
    end
  end

  private

  REGEXP = /language:\s?(.*)\nversion:\s?(.*)\n```\n([\s\S]*)\n```/

  def valid(text)
    text.match?(REGEXP)
  end

=begin
    lang: ruby
    version: 2.3.0
    ```
      def hello
        "hello"
      end

      def world
        "world"
      end

      puts hello + world
    ```
=end
  def build_json(text)
    m = text.match(REGEXP)
    language = m[1]
    version = m[2]
    program = m[3]
    {"language" => language, "version" => version, "program" => program}.to_json
  end
end

SlackLambda.new(ENV['SLACK_LAMBDA_TOKEN']).run

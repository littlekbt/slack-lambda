require 'slack-rtm-bot-helper'
require 'faraday'
require 'json'
require 'time'
require 'pry'

class SlackLambda
  NAME = 'slack-lambda'
  def initialize(token)
    @token = token
    @name  = NAME
  end

  def run
    puts "#{[Time.now().to_s]} open web socket connection"
    Slack::Rtm::Bot::Helper.run(token: @token, name: @name) do |data|
      text = data['text']
      case extract_action(text)
      when :register
        UserFile.register(text)
      when :show
        UserFile.show(text)
      when :list
        UserFile.list(text)
      when :remove
        UserFile.remove(text)
      when :run
        UserFile.run(text)
      when :exec
        PostConverter.post(text)
      when :help
        show_help
      end
    end
    puts "#{[Time.now().to_s]} close web socket connection"
  end

  private
  
  CONVERT_TO_ACTION = {
    "register" => :register,
    "show"     => :show,
    "list"     => :list,
    "remove"   => :remove,
    "run"      => :run,
    "exec"     => :exec,
    "help"     => :help,
  }

  def extract_action(text)
    m = text.match(/>\s*(register|show|list|remove|run|exec|help)/)

    CONVERT_TO_ACTION[m[1]] if m && CONVERT_TO_ACTION[m[1]]
  end

  class UserFile
    REGISTER_REGEXP = /name:\s?(.+)\n(language:\s?.*\nversion:\s?.*\n```\n?[\s\S]*\n?```)/
    SHOW_REGEXP     = />\s*show\s*([a-z|A-Z|0-9|_|-]*)/
    LIST_REGEXP     = />\s*list/
    REMOVE_REGEXP     = />\s*remove\s*([a-z|A-Z|0-9|_|-]*)/
    RUN_REGEXP      = />\s*run\s*([a-z|A-Z|0-9|_|-]*)/
    class << self
      def register(text)
        m = text.match(REGISTER_REGEXP)
        if m
          name = m[1]
          parameter = m[2]

          namem = name.match(/[a-z|A-Z|0-9|_|-]*/)
          return 'you can use only [a-z|A-Z|0-9|_|-] for function name.' unless namem[0].to_s.length == name.to_s.length

          return "#{name} is already exist" if File.exist?("./storage/user_files/#{name}")

          File.open("./storage/user_files/#{name}", "w") do |f|
            f.puts(parameter)
          end

          return "success register #{name}"
        end
      end

      def show(text)
        m = text.match(SHOW_REGEXP)
        if m
          name = m[1]

          return "#{name} not found" unless File.exist?("./storage/user_files/#{name}")

          File.read("./storage/user_files/#{name}") 
        end
      end

      def list(text)
        m = text.match(LIST_REGEXP)
        Dir.glob("./storage/user_files/*").map{ |f| File.basename(f) }.join("\n") if m
      end

      def remove(text)
        m = text.match(REMOVE_REGEXP)
        if m
          name = m[1]

          return "#{name} not found" unless File.exist?("./storage/user_files/#{name}")

          File.delete("./storage/user_files/#{name}")

          "delete #{name} function"
        end
      end

      def run(text)
        m = text.match(RUN_REGEXP)
        if m
          name = m[1]

          return "#{name} not found" unless File.exist?("./storage/user_files/#{name}")

          t = File.read("./storage/user_files/#{name}") 
          PostConverter.post(t)
        end
      end
    end
  end

  class PostConverter
    REGEXP = /language:\s?.*\nversion:\s?.*\n```\n?[\s\S]*\n?```/
    class << self
      def post(text)
        return unless valid?(text)

        puts "#{[Time.now().to_s]} #{text}"

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

      private
      def valid?(text)
        text.match?(REGEXP)
      end

      def output_error(err_msg)
        <<ERR
Error Found!
#{err_msg}
ERR
      end
    end
  end

  def show_help
  <<EOF
slack-lambda is bot that execute function.

# exec
execute function.

@slack-lambda exec
language: [ruby|golang]
version: <ex: 2.3.0>
```
# input your program
```

# register
register function

@slack-lambda register
name: <input function name>
language: [ruby|golang]
version: <ex: 2.3.0>
```
# input your program
```

# show
show registered function

@slack-lambda show <function name>

# list
list functions

@slack-lambda list 

# remove
delete function

@slack-lambda remove <function name>

# run
execute registered function

@slack-lambda run <function name>
EOF
  end
end

SlackLambda.new(ENV['SLACK_LAMBDA_TOKEN']).run
# SlackLambda.run

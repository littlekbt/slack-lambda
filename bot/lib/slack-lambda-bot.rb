class SlackLambdaBot
  NAME = 'slack-lambda'

  CONVERT_TO_ACTION_SYM = {
    "register" => :register,
    "show"     => :show,
    "list"     => :list,
    "remove"   => :remove,
    "run"      => :run,
    "exec"     => :exec,
    "help"     => :help,
  }

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
        PostProxy.post(text)
      when :help
        show_help
      end
    end
    puts "#{[Time.now().to_s]} close web socket connection"
  end

  private

  def extract_action(text)
    m = text.match(/>\s*(register|show|list|remove|run|exec|help)/)

    m && CONVERT_TO_ACTION_SYM[m[1]] ? CONVERT_TO_ACTION_SYM[m[1]] : :help
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

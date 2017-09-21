class SlackLambdaBot
  module MessageConverter
    def convert_program(message)
      message = remove_url_parenthesis(message)
      message
    end

    # urlが<http://hogehoge.com>のように<>で囲われるので、とる。
    def remove_url_parenthesis(message)
      regexp = /```[\s\S]*(<#{URI.regexp}>)[\s\S]*```/
      m = message.match(regexp)
      m ? message.gsub(m[1], m[1][1...-1]) : message
    end
  end
end

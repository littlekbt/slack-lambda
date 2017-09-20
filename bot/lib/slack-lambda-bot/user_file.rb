class SlackLambdaBot
  class UserFile
    REGISTER_REGEXP = /name:\s?(.+)\n(language:\s?.*\nversion:\s?.*\n```\n?[\s\S]*\n?```)/
    SHOW_REGEXP     = />\s*show\s*([a-z|A-Z|0-9|_|-]*)/
    LIST_REGEXP     = />\s*list/
    REMOVE_REGEXP   = />\s*remove\s*([a-z|A-Z|0-9|_|-]*)/
    RUN_REGEXP      = />\s*run\s*([a-z|A-Z|0-9|_|-]*)/

    FILE_PATH = Pathname.new("./storage/user_files/")

    class << self
      def register(text)
        m = text.match(REGISTER_REGEXP)
        if m
          name = m[1]
          parameter = m[2]

          namem = name.match(/[a-z|A-Z|0-9|_|-]*/)
          return 'you can use only [a-z|A-Z|0-9|_|-] for function name.' unless namem[0].to_s.length == name.to_s.length

          return "#{name} is already exist" if File.exist?(to_s_path(name))

          File.open(to_s_path(name), "w") do |f|
            f.puts(parameter)
          end

          return "success register #{name}"
        end
      end

      def show(text)
        m = text.match(SHOW_REGEXP)
        if m
          name = m[1]

          return "#{name} not found" unless File.exist?(to_s_path(name))

          File.read(to_s_path(name)) 
        end
      end

      def list(text)
        m = text.match(LIST_REGEXP)
        Dir.glob("#{FILE_PATH}*").map{ |f| File.basename(f) }.join("\n") if m
      end

      def remove(text)
        m = text.match(REMOVE_REGEXP)
        if m
          name = m[1]

          return "#{name} not found" unless File.exist?(to_s_path(name))

          File.delete(to_s_path(name))

          "delete #{name} function"
        end
      end

      def run(text)
        m = text.match(RUN_REGEXP)
        if m
          name = m[1]

          return "#{name} not found" unless File.exist?(to_s_path(name))

          t = File.read(to_s_path(name)) 
          PostConverter.post(t)
        end
      end

      private
      def to_s_path(name)
        Pathname.new(FILE_PATH).join(name).to_s
      end
    end
  end
end

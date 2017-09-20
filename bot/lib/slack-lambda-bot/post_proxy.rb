class SlackLambdaBot
  class PostProxy
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
end

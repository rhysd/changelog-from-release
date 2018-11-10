def run(cmdline)
  puts "+#{cmdline}"
  system cmdline
end

guard :shell do
  watch /\.go$/ do |m|
    puts "#{Time.now}: #{m[0]}"
    case m[0]
    when /_test\.go$/
      parent = File.dirname m[0]
      sources = Dir["#{parent}/*.go"].reject{|p| p.end_with? '_test.go'}.join(' ')
      run "go test -v #{m[0]} #{sources}"
    else
      run 'go build'
    end
  end
end

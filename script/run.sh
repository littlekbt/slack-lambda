CURRENT=`pwd`
(cd $CURRENT/bot && STORAGE_PATH=`pwd`/storage bundle exec ruby main.rb)
nohup go run proxy-server/main.go &
nohup go run container-server/main.go &

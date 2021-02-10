kill -9 `pgrep BattleSnake2019` # kills the old process

# rebuild the golang server
git pull
go build
go install BattleSnake2019

# start the server on a background process so it will not die
nohup BattleSnake2019 2> logs.txt &
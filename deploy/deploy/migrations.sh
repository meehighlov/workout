go install github.com/pressly/goose/v3/cmd/goose@latest
git clone https://github.com/meehighlov/workout.git workout-tmp
./goose -dir=workout-tmp/migrations sqlite3 workout.db up
./goose -dir=workout-tmp/migrations sqlite3 workout.db status
rm -rf workout-tmp
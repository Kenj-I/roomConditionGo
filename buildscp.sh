GOARCH=arm GOARM=5 GOOS=linux go build -o app
scp ./app raspizerowh:/home/pi/roomCondition
# Go WebOS 📺
This repository was forked from https://github.com/kaperys/go-webos

# Launch YouTube Video with a URL
Authorization
```
go run cmd/auth/main.go 192.168.3.10 > client-id
```
client id for accessing to your TV will be saved to `client-id` file.

Launch a youtube video by using a URL
```
go run cmd/launch_youtube/main.go 192.168.3.10 $(cat client-id) https://www.youtube.com/watch?v=79XaA_4CYj8
```

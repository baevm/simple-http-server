# simple-http-server
Simple HTTP server implemented using TCP listener

1. Start server
```
go run ./app/server.go
```

2. Test server on 0.0.0.0:3000
```
- Example JSON response with gzip encoding: http -v localhost:3000/get-user

GET /get-user HTTP/1.1
Accept: */*
Accept-Encoding: gzip, deflate
Connection: keep-alive
Host: localhost:3000
User-Agent: HTTPie/2.6.0


HTTP/1.1 200 OK
Content-Encoding: gzip
Content-Length: 66
Content-Type: application/json

{
    "Password": "helloworld",
    "User": "test123"
}

- Example getting files from server: http -v localhost:3000/files/static/index.html

GET /files/static/index.html HTTP/1.1
Accept: */*
Accept-Encoding: gzip, deflate
Connection: keep-alive
Host: localhost:3000
User-Agent: HTTPie/2.6.0


HTTP/1.1 200 OK
Content-Length: 218
Content-Type: application/octet-stream

<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Document</title>
</head>
<body>
    <h1>test</h1>
</body>
</html>
```
## Sample HTTP Server

### Start server

    go run main.go
   
    http: 2021/08/11 18:42:05
    Server is starting on port:8080...

#### Send request

    curl --location --request PATCH 'http://127.0.0.1:8080/v1/catalog' \
    --form 'variations=@"<path-to-file"' \
    --form 'file_name="variations"'

#### Response
    PATCH request w/ file variations.csv was sent successfully
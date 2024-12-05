#!/bin/sh

URL=localhost:8080/upload
FILE=mastering-postgresql.pdf
curl -X POST \
     -F upload=@$FILE \
     -H "Content-Type: multipart/form-data" \
     $URL 

# micro-playlist



Curl Requests


Get


curl -v localhost:9090/product | jq
curl -v localhost:9090/product/1 -XPUT -d '{"name":"tea", "description":"nice cup of tea"}' 
curl -F id=3 -F file=@barath.jpg -v localhost:9090/file/multi
curl -v -H "Accept-Encoding:gzip" localhost:9090/file/1/something.txt
curl -v localhost:9090/file/1/something.txt --compressed 
curl -v -H localhost:9090/file/1/barath.png --compressed -o barath.png

# kearch-sp-google

Use Google Custome Search Engine as a specialist search engine

### Running the server

First you need to rewrite ./GoogleCustomSearchAPIKey and ./GoogleCustomSearchEngineID with your own API key and engine ID. You can get APIkey from [https://developers.google.com/custom-search/v1/introduction](https://developers.google.com/custom-search/v1/introduction).

And then

```
go run main.go
```

To run the server in a docker container
```
docker build --network=host -t kearch-sp-google .
```

Once image is built use
```
docker run --rm -p 32500:32500 -it kearch-sp-google 
```



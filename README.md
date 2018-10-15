To build the dev environment image:

```
docker build -t hell .
```

To start the dev environment:

```
docker run -it --mount type=bind,source="$(pwd)",target=/go/src/github.com/gcapizzi/hell --privileged hell
```

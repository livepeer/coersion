# coersion
prototype for video content validation

## To run the app
```shell
git clone https://github.com/mkrufky/coersion
cd coersion
docker-compose up
```

...then navigate to http://localhost:8080

To try out the feature, use the `match2` of `match3` endpoints and provide the arguments:
```
s0 - a link to a video source
s1 - a link to a video source
s2 - a link to a video source
w - a desired width dimension for resizing
h - a desired height dimension for resizing
o - an optional offset from the start of the video in seconds
v - an allowed variance value for use in matching.
```

The smaller the `v` value, the more strictly the pixels must match. For example: 0 for exact, 1 for slight fuzz, 5 for more fuzz, 10 for much fuzz.

I recommend to use samples such as those listed at https://file-examples.com/index.php/sample-video-files/sample-mp4-files/

Choose a resolution such as 360x240.

Use a different rendition for each source, and start with a variance of 1.  You'll notice that they don't match up perfectly.  Increase the variance to 2.  Then 5.  Then 10...

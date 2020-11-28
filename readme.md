# spotify-charts-scraper
A tool to download all time charts from [spotifycharts.com](https://spotifycharts.com).

## Intention

This was just a little friday-night-toy project.
You can read a bit more about it on [my blog](https://ldb.github.io/posts/spotify-scraper).
Maybe someone with a background in data science will find this useful to generate some nice visualizations.

## Building

Just use `go build` to build the project:
```
go build -o spotifyscraper cmd/download/main.go
```

## Running

Running the tool will automatically download the all time charts between the two dates in `main.go` for each region 
available to a directory called `data`: 

```shell script
2020/11/01 17:52:00 1398 days, 66 regions
Peru                 76 / 1398 [>----------------]   9m14s   5 %
Czech Republic      242 / 1398 [==>--------------]   2m32s  17 %
Indonesia           201 / 1398 [=>---------------]    3m3s  14 %
Chile               182 / 1398 [=>---------------]   3m26s  13 %
Slovakia            125 / 1398 [=>---------------]   4m33s   9 %
Thailand            387 / 1398 [====>------------]    1m7s  28 %
Estonia             144 / 1398 [=>---------------]   3m36s  10 %
Turkey              132 / 1398 [=>---------------]   3m58s   9 %
Mexico               67 / 1398 [>----------------]   7m53s   5 %
Hong Kong            38 / 1398 [-----------------]  13m38s   3 %
```

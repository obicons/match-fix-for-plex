# TV Series Match Fixer for Plex
Many times, naming schemes for TV series do not match Plex's naming expectations. For example, you might have:
```
$ ls /plex-media/tv/
Friends Season 1 COMPLETE 1080p/ 
Friends Season 2/
Friends.Season.3/
swcw Season1/
```
With these directory names, Plex will not match your shows correctly.

This tool attempts to correct this problem. After running this program, you'll have this:
```
$ tree /plex-media/tv/
Friends
 | Season 1
 | Season 2
 | Season 3
star wars: the clone wars
 | Season 1
```
And Plex will correctly match your shows. Note that no abbreviations are hard-coded into this app. Just try running, and see if it matches!

## Running
### Usage
```
$ ./match-fix -help
Usage of ./match-fix:
  -matchee string
    	title that needs matched
  -titles string
    	path to file containing TV Show titles
```


### Example
```
./match-fix -matchee /plex-media/tv/Friends\ Season\ 1\ COMPLETE\ 1080p -titles ./show_list
```
*This program never modifies your filesystem without your consent.* You are prompted at every modification to prevent a bad match from ruining your naming scheme. This is just faster than manually repairing yourself.

### Fixing Bad Matches
You can use the `-titles` flag to specify a file containing a bunch of TV show titles. This program matches by:
1. Using heuristics to remove unnecessary data from directory names
2. Applying a fuzzy-search algorithm to find the best match
3. Using popularity as a tie-breaker

If something is constantly mismatched for you, then you can probably fix it by making the show's title appear earlier in the file specified by the `-titles` flag.

## Building
You need to have Go installed. Then, you can just run `make`.

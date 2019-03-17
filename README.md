# gomul
Joins multicast ipv4/ipv6 groups on specified interface (used for debugging multicast issues)

Tested on windows and linux

## Building
Make sure you have [Go](https://golang.org/doc/install) properly installed.
Requires go1.12+

Next, run

 ```
 $ go get github.com/42wim/gomul
 ```

 You'll have the binary 'gomul' in $GOPATH/bin


## Usage
```
Usage of gomul:
  -group="ff02::42:1 239.42.42.1": multicast groups to join (space seperated)
  -interface=0: interface to listen on (number)
  -ip="": use interface where the specified ip is bound on
  -li=false: show available interfaces
```

## Examples
![demo](http://i.snag.gy/fPmeD.jpg)
![result](http://i.snag.gy/IlteH.jpg)

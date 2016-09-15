# show-prs

```
$ ./show-prs --help
Usage: show-prs [-s] [-f value] [-o value] [-p value] [-t value] [parameters ...]
 -f, --from=value  Sender email address
 -o, --org=value   org
 -p, --project=value
                   Project
 -s, --send
 -t, --to=value    Recipient email address
```
EXAMPLE
---------
./show-prs --send --from "user@example.com" --to "group@example.com" -p project1 -p projectN --org githubuser

INSTALL
--------
Use go get to install the program. It will be located in ```$GOPATH/bin```

```
go get github.com/jmrodri/show-prs
```

MANUAL BUILD
---------------
```
git clone https://github.com/jmrodri/show-prs.git
cd show-prs
go build
./show-prs
```

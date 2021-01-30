## ogit

ogit means "open git" and it's suitable to open from current folder the github's project website.

Make sure you have [Golang](https://golang.org/) installed if you want to make any further modifications.

ogit is simple and straightforward. All github project has a internal and hidden folder called *.git*. If you access its content, you can see a file named *config* that contains the url information for the current GitHub website.

### Installation

Git clone the master repository of the project:

```
git clone github.com\inocencio\master\ogit
```

Run the following command to build and execute it:

```
go build #build the ogit into your project's folder.
go install #build it to GOPATH bin folder. <recommended for global use>
go run ogit.go #execute it directly from go
```
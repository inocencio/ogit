## ogit

ogit means _"open git"_, and it's suitable to open from current folder the github's project website.

Make sure you have [Golang](https://golang.org/) installed if you want to make any further modifications.

ogit is simple and straightforward. All github project has a internal and hidden folder called *.git*. If you access its content, you can see a file named *config* that contains the url information for the current GitHub website.

### Installation

Git clone the master repository of the project:

```
git clone https://github.com/inocencio/ogit
```

Run the following command to build and execute it:

```
go build #build the ogit into your project's folder.
go install #build it to GOPATH bin folder. <recommended for global uses>
go run ogit.go #execute it directly from go
```

Create an executable file from `go build` or `go install` and make sure to place it in user's path to access it from anywhere. Just type `ogit` inside a folder with a *.git* directory and the magic goes on and get ready to browse it. ;)

### Run

Type `ogit` in the same folder from a git project.

```bash
ogit
```

You can set the browser you want to be used instead the default browser configured.

```bash
ogit -b <browser's bath>
```

ogit config file is stored in the system config path. Usually in Linux and Mac it's store in `~/.config/.ogit`, in windows it's usually stored in `<user>\AppData\Roaming\.ogit`.
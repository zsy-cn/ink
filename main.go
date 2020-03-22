package main

import (
	"bufio"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"html/template"
	"time"

	"github.com/codegangsta/cli"
)

const (
	VERSION            = "RELEASE 2020-03-21"
	DEFAULT_ROOT       = "blog"
	DATE_FORMAT_STRING = "2006-01-02 15:04:05"
	INDENT             = "  " // 2 spaces
	POST_TEMPLATE      = `title: {{.Title}}
date: {{.DateString}}
author: {{.Author}}
{{- if .Cover}}
cover: {{.Cover}}
{{- end}}
draft: {{.Draft}}
top: {{.Top}}
{{- if .Preview}}
preview: {{.Preview}}
{{- end}}
{{- if .Tags}}
{{.Tags}}
{{- end}}
type: {{.Type}}
hide: {{.Hide}}
---
`
)

var globalConfig *GlobalConfig
var rootPath string

func main() {

	app := cli.NewApp()
	app.Name = "ink"
	app.Usage = "An elegant static blog generator"
	app.Authors = []*cli.Author{
		{Name: "Xavier", Email: "zsy@iotinit.com"},
	}
	app.Version = VERSION
	app.Commands = []*cli.Command{
		{
			Name:  "build",
			Usage: "Generate blog to public folder",
			Action: func(c *cli.Context) error {
				ParseGlobalConfigByCli(c, false)
				Build()
				return nil
			},
		},
		{
			Name:  "preview",
			Usage: "Run in server mode to preview blog",
			Action: func(c *cli.Context) error {
				ParseGlobalConfigByCli(c, true)
				Build()
				Watch()
				Serve()
				return nil
			},
		},
		{
			Name:  "publish",
			Usage: "Generate blog to public folder and publish",
			Action: func(c *cli.Context) error {
				ParseGlobalConfigByCli(c, false)
				Build()
				Publish()
				return nil
			},
		},
		{
			Name:  "serve",
			Usage: "Run in server mode to serve blog",
			Action: func(c *cli.Context) error {
				ParseGlobalConfigByCli(c, true)
				Build()
				Serve()
				return nil
			},
		},
		{
			Name:  "new",
			Usage: "Creates a new article",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:  "hide",
					Usage: "Hides the article",
				},
				&cli.BoolFlag{
					Name:  "top",
					Usage: "Places the article at the top",
				},
				&cli.BoolFlag{
					Name:  "post",
					Usage: "The article is a post",
				},
				&cli.BoolFlag{
					Name:  "page",
					Usage: "The article is a page",
				},
				&cli.BoolFlag{
					Name:  "draft",
					Usage: "The article is a draft",
				},

				&cli.StringFlag{
					Name:  "title",
					Usage: "Article title",
				},
				&cli.StringFlag{
					Name:  "author",
					Usage: "Article author",
				},
				&cli.StringFlag{
					Name:  "cover",
					Usage: "Article cover path",
				},
				&cli.StringFlag{
					Name:  "date",
					Usage: "The date and time on which the article was created (2006-01-02 15:04:05)",
				},
				&cli.StringFlag{
					Name:  "file",
					Usage: "The path of where the article will be stored",
				},

				&cli.StringSliceFlag{
					Name:  "tag",
					Usage: "Adds a tag to the article",
				},
			},
			Action: func(c *cli.Context) error {
				New(c)
				return nil
			},
		},
	}
	app.Run(os.Args)
	os.Exit(exitCode)
}

func ParseGlobalConfigByCli(c *cli.Context, develop bool) {
	if c.Args().Len() > 0 {
		rootPath = c.Args().Slice()[0]
	} else {
		rootPath = "."
	}
	ParseGlobalConfigWrap(rootPath, develop)
	if globalConfig == nil {
		ParseGlobalConfigWrap(DEFAULT_ROOT, develop)
		if globalConfig == nil {
			Fatal("Parse config.yml failed, please specify a valid path")
		}
	}
}

func ParseGlobalConfigWrap(root string, develop bool) {
	rootPath = root
	globalConfig = ParseGlobalConfig(filepath.Join(rootPath, "config.yml"), develop)
	if globalConfig == nil {
		return
	}
}

func New(c *cli.Context) {
	// If source folder does not exist, create
	if _, err := os.Stat("source/"); os.IsNotExist(err) {
		os.Mkdir("source", os.ModePerm)
	}

	var author, blogTitle, fileName string
	var tags []string

	// Default values
	draft := "false"
	top := "false"
	postType := "post"
	hide := "false"
	date := time.Now()

	// Empty string values
	preview := ""
	cover := ""

	// Parse args
	args := c.Args()
	if args.Len() > 0 {
		blogTitle = args.Slice()[0]
	}
	if blogTitle == "" {
		if c.String("title") != "" {
			blogTitle = c.String("title")
		} else {
			Fatal("Please specify the name of the blog post")
		}
	}

	fileName = blogTitle + ".md"
	if c.String("file") != "" {
		fileName = c.String("file")
	}

	if args.Len() > 1 {
		author = args.Slice()[1]
	}
	if author == "" {
		author = c.String("author")
	}

	if c.Bool("post") && c.Bool("page") {
		Fatal("The post and page arguments are mutually exclusive and cannot appear together")
	}
	if c.Bool("post") {
		postType = "post"
	}
	if c.Bool("page") {
		postType = "page"
	}
	if c.Bool("hide") {
		hide = "true"
	}
	if c.Bool("draft") {
		draft = "true"
	}
	if c.Bool("top") {
		top = "true"
	}

	if c.String("preview") != "" {
		preview = c.String("preview")
	}
	if c.String("cover") != "" {
		cover = c.String("cover")
	}

	var filePath = "source/" + fileName
	file, err := os.Create(filePath)
	if err != nil {
		Fatal(err)
	}
	postTemplate, err := template.New("post").Parse(POST_TEMPLATE)
	if err != nil {
		Fatal(err)
	}

	if c.StringSlice("tag") != nil {
		tags = c.StringSlice("tag")
	}

	var tagString string
	if len(tags) > 0 {
		tagString = "tags:"
		for _, tag := range tags {
			tagString += "\n" + INDENT + "- " + tag
		}
	}

	var dateString string
	if c.String("date") != "" {
		dateString = c.String("date")
		_, err = time.Parse(DATE_FORMAT_STRING, dateString)
		if err != nil {
			Fatal("Illegal date string")
		}
	} else {
		dateString = date.Format(DATE_FORMAT_STRING)
	}
	data := map[string]string{
		"Title":      blogTitle,
		"DateString": dateString,
		"Author":     author,
		"Draft":      draft,
		"Top":        top,
		"Type":       postType,
		"Hide":       hide,
		"Preview":    preview,
		"Cover":      cover,
		"Tags":       tagString,
	}
	fileWriter := bufio.NewWriter(file)
	err = postTemplate.Execute(fileWriter, data)
	if err != nil {
		Fatal(err)
	}
	err = fileWriter.Flush()
	if err != nil {
		Fatal(err)
	}
}

func Publish() {
	command := globalConfig.Build.Publish
	// Prepare exec command
	var shell, flag string
	if runtime.GOOS == "windows" {
		shell = "cmd"
		flag = "/C"
	} else {
		shell = "/bin/sh"
		flag = "-c"
	}
	cmd := exec.Command(shell, flag, command)
	cmd.Dir = filepath.Join(rootPath, globalConfig.Build.Output)
	// Start print stdout and stderr of process
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	out := bufio.NewScanner(stdout)
	err := bufio.NewScanner(stderr)
	// Print stdout
	go func() {
		for out.Scan() {
			Log(out.Text())
		}
	}()
	// Print stdin
	go func() {
		for err.Scan() {
			Log(err.Text())
		}
	}()
	// Exec command
	cmd.Run()
}

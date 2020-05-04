package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

const downloadHelp = `Download a paper and upload to remarkable cloud.`

func (cmd *downloadCommand) Name() string      { return "download" }
func (cmd *downloadCommand) Args() string      { return "[OPTIONS] URL TITLE" }
func (cmd *downloadCommand) ShortHelp() string { return downloadHelp }
func (cmd *downloadCommand) LongHelp() string  { return downloadHelp }
func (cmd *downloadCommand) Hidden() bool      { return false }

func (cmd *downloadCommand) Register(fs *flag.FlagSet) {
}

type downloadCommand struct {
}

func (cmd *downloadCommand) Run(ctx context.Context, args []string) error {
	if len(args) < 1 {
		return errors.New("must pass a url for a pdf")
	}
	if len(args) < 2 {
		return errors.New("must pass a title for the pdf")
	}

	if err := createDirectory(args[2]); err != nil {
		return err
	}

	// Download the pdf.
	
	file := args[0]
	if strings.HasPrefix("http", file) {
	  	logrus.WithFields(logrus.Fields{
			"link": args[0],
		}).Debug("downloading paper")

		file := filepath.Join(dataDir, args[1]+".pdf")
		if err := downloadPaper(args[0], file); err != nil {
			return err
		}
	}

	logrus.WithFields(logrus.Fields{
		"link": args[0],
		"file": file,
	}).Info("downloaded paper to file")

	// Sync the file with remarkable cloud.
	if err := rmAPI.SyncFileAndRename(file, args[1], args[2]); err != nil {
		return err
	}

	fmt.Printf("Downloaded %s and renamed to %s in folder %s", args[0], args[1], args[2])
	return nil
}

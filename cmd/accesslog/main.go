package main

import (
	"log"
	"os"

	"github.com/spf13/pflag"
	"github.com/wzshiming/accesslog"
	"github.com/wzshiming/accesslog/nginx"
	"github.com/wzshiming/accesslog/tocsv"
)

var (
	fields    = nginx.AccessLogFormatted{}.Fields()
	condition string
)

func init() {
	pflag.StringSliceVar(&fields, "field", fields, "fields")
	pflag.StringVar(&condition, "condition", condition, "condition")
	pflag.Parse()
}

func main() {
	err := run(condition, fields)
	if err != nil {
		log.Fatal(err)
	}
}

func run(condition string, fields []string) error {
	ch := make(chan nginx.AccessLogFormatted, 128)

	go func() {
		defer close(ch)
		err := accesslog.ProcessEntries[nginx.DefaultAccessLog](os.Stdin, func(entry accesslog.Entry[nginx.DefaultAccessLog], err error) error {
			if err != nil {
				return err
			}

			e := entry.Entry()
			f, err := e.Formatted()
			if err != nil {
				return err
			}
			ch <- f
			return nil
		})
		if err != nil {
			log.Fatal("Error", err)
		}
	}()

	return tocsv.ProcessToCSV[nginx.AccessLogFormatted](os.Stdout, condition, fields, ch)
}

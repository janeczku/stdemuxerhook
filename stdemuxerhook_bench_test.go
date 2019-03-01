package stdemuxerhook

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"testing"
	"time"
)

// smallFields is a small size data set for benchmarking
var smallFields = logrus.Fields{
	"foo":   "bar",
	"baz":   "qux",
	"one":   "two",
	"three": "four",
}

// largeFields is a large size data set for benchmarking
var largeFields = logrus.Fields{
	"foo":       "bar",
	"baz":       "qux",
	"one":       "two",
	"three":     "four",
	"five":      "six",
	"seven":     "eight",
	"nine":      "ten",
	"eleven":    "twelve",
	"thirteen":  "fourteen",
	"fifteen":   "sixteen",
	"seventeen": "eighteen",
	"nineteen":  "twenty",
	"a":         "b",
	"c":         "d",
	"e":         "f",
	"g":         "h",
	"i":         "j",
	"k":         "l",
	"m":         "n",
	"o":         "p",
	"q":         "r",
	"s":         "t",
	"u":         "v",
	"w":         "x",
	"y":         "z",
	"this":      "will",
	"make":      "thirty",
	"entries":   "yeah",
}

var errorFields = logrus.Fields{
	"foo": fmt.Errorf("bar"),
	"baz": fmt.Errorf("qux"),
}

func BenchmarkSmallNopFormatter(b *testing.B) {
	doFormatterBenchmark(b, &NopFormatter{}, smallFields)
}

func BenchmarkLargeNopFormatter(b *testing.B) {
	doFormatterBenchmark(b, &NopFormatter{}, largeFields)
}

func doFormatterBenchmark(b *testing.B, formatter logrus.Formatter, fields logrus.Fields) {
	entry := &logrus.Entry{
		Time:    time.Time{},
		Level:   logrus.InfoLevel,
		Message: "message",
		Data:    fields,
	}
	var d []byte
	var err error
	for i := 0; i < b.N; i++ {
		d, err = formatter.Format(entry)
		if err != nil {
			b.Fatal(err)
		}
		b.SetBytes(int64(len(d)))
	}
}

func BenchmarkLoggerWithoutHook(b *testing.B) {
	logger := logrus.New()
	nullfs, err := os.OpenFile("/dev/null", os.O_WRONLY, 0666)
	if err != nil {
		b.Fatalf("%v", err)
	}

	defer nullfs.Close()
	logger.Out = nullfs
	doLoggerBenchmark(b, logger, smallFields)
}

func BenchmarkLoggerWithHook(b *testing.B) {
	nullfs1, err := os.OpenFile("/dev/null", os.O_WRONLY, 0666)
	if err != nil {
		b.Fatalf("%v", err)
	}
	defer nullfs1.Close()
	nullfs2, err := os.OpenFile("/dev/null", os.O_WRONLY, 0666)
	if err != nil {
		b.Fatalf("%v", err)
	}
	defer nullfs2.Close()

	logger := logrus.New()
	hook := New(logger)
	hook.SetOutput(nullfs1, nullfs2)
	logger.Hooks.Add(hook)

	doLoggerBenchmark(b, logger, smallFields)
}

func doLoggerBenchmark(b *testing.B, logger *logrus.Logger, fields logrus.Fields) {
	entry := logger.WithFields(fields)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			entry.Info("aaa")
		}
	})
}

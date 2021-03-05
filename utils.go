package dotfiles

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var (
	flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "config",
			Aliases: []string{"c"},
			Value:   "config.yml",
			Usage:   "config file path",
			EnvVars: []string{"DOTFILES_CONFIG_FILE"},
		},
	}
	verboseFlag = &cli.BoolFlag{
		Name:    "verbose",
		Aliases: []string{"v"},
		Usage:   "verbose output",
	}
)

func beforeAction() func(c *cli.Context) error {
	return func(c *cli.Context) error {
		if c.Bool("verbose") {
			log.SetLevel(log.DebugLevel)
		}
		if c.String("config") == "" {
			return fmt.Errorf("config path missing")
		}
		return nil
	}
}

const (
	chunkSize     = 64000
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func randString(n int) string {
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func fileExists(filename string) (bool, error) {
	_, err := os.Stat(filename)
	if err != nil {
		if !os.IsNotExist(err) {
			return false, err
		}
		return false, nil
	}
	return true, nil
}

func fileStat(filename string) (os.FileInfo, error) {
	info, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	return info, err
}

func mkdir(path string) error {
	if ok, err := fileExists(path); !ok {
		if err != nil {
			return err
		}
		return os.MkdirAll(path, 0755)
	}
	return nil
}

func fileGetContents(name string) ([]byte, error) {
	if _, err := os.Stat(name); err != nil {
		return nil, err
	}
	return ioutil.ReadFile(name)
}

func fileGets(name string, len int) ([]byte, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf := make([]byte, len)
	_, err = f.Read(buf)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func fileCompare(source string, dest string) (bool, error) {
	src, err := os.Open(source)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	defer src.Close()

	dst, err := os.Open(dest)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	defer dst.Close()

	for {
		b1 := make([]byte, chunkSize)
		_, err1 := src.Read(b1)

		b2 := make([]byte, chunkSize)
		_, err2 := dst.Read(b2)

		if err1 != nil || err2 != nil {
			if err1 == io.EOF && err2 == io.EOF {
				return true, nil
			} else if err1 == io.EOF || err2 == io.EOF {
				return false, nil
			} else if err1 != io.EOF {
				return false, err1
			} else {
				return false, err2
			}
		}

		if !bytes.Equal(b1, b2) {
			return false, nil
		}
	}
}

func copyFile(source string, dest string) error {
	// check is same file
	if same, err := fileCompare(source, dest); err != nil || same {
		if same {
			log.Debug("not modified: ", source)
		}
		return err
	}

	src, err := os.Open(source)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err == nil {
		log.Info("copyed: ", source, " => ", dest)
	}
	return err
}

package app

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
)

// String utils
// const (
// 	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
// 	letterIdxBits = 6                    // 6 bits to represent a letter index
// 	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
// 	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
// )

// func randString(n int) string {
// 	b := make([]byte, n)
// 	// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
// 	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
// 		if remain == 0 {
// 			cache, remain = rand.Int63(), letterIdxMax
// 		}
// 		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
// 			b[i] = letterBytes[idx]
// 			i--
// 		}
// 		cache >>= letterIdxBits
// 		remain--
// 	}

// 	return string(b)
// }

// File utils
// const chunkSize = 64000

func openFile(name string) (*os.File, error) {
	err := os.MkdirAll(filepath.Dir(name), 0755)
	if err != nil {
		return nil, err
	}
	return os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0644)
}

func saveFile(f *os.File, contents []byte) (bool, error) {
	existingContents, err := io.ReadAll(f)
	if err != nil {
		return false, err
	}
	if bytes.Equal(contents, existingContents) {
		return true, nil
	}
	err = f.Truncate(0)
	if err != nil {
		return false, err
	}
	_, err = f.WriteAt(contents, 0)
	return false, err
}

// func fileGets(name string, len int) ([]byte, error) {
// 	f, err := os.Open(name)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer f.Close()

// 	buf := make([]byte, len)
// 	_, err = f.Read(buf)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return buf, nil
// }

// func fileCompare(source string, dest string) (bool, error) {
// 	src, err := os.Open(source)
// 	if err != nil {
// 		if os.IsNotExist(err) {
// 			return false, nil
// 		}
// 		return false, err
// 	}
// 	defer src.Close()

// 	dst, err := os.Open(dest)
// 	if err != nil {
// 		if os.IsNotExist(err) {
// 			return false, nil
// 		}
// 		return false, err
// 	}
// 	defer dst.Close()

// 	for {
// 		b1 := make([]byte, chunkSize)
// 		_, err1 := src.Read(b1)

// 		b2 := make([]byte, chunkSize)
// 		_, err2 := dst.Read(b2)

// 		if err1 != nil || err2 != nil {
// 			if err1 == io.EOF && err2 == io.EOF {
// 				return true, nil
// 			} else if err1 == io.EOF || err2 == io.EOF {
// 				return false, nil
// 			} else if err1 != io.EOF {
// 				return false, err1
// 			} else {
// 				return false, err2
// 			}
// 		}

// 		if !bytes.Equal(b1, b2) {
// 			return false, nil
// 		}
// 	}
// }

// func copyFile(source string, dest string) error {
// 	// check is same file
// 	if same, err := fileCompare(source, dest); err != nil || same {
// 		if same {
// 			log.Debug("not modified: ", source)
// 		}
// 		return err
// 	}

// 	src, err := os.Open(source)
// 	if err != nil {
// 		return err
// 	}
// 	defer src.Close()

// 	dst, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
// 	if err != nil {
// 		return err
// 	}
// 	defer dst.Close()

// 	_, err = io.Copy(dst, src)
// 	return err
// }

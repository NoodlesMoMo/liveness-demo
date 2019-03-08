package healthy

import (
	"fmt"
	"io"
	"os"
	"runtime"
)

func GoRoutingCheck(max int) HealthyCheck {

	return func() error {
		cnt := runtime.NumGoroutine()
		if cnt > max {
			return fmt.Errorf("Toooo many goroutines, count: %d", cnt)
		}
		return nil
	}
}

func FileReadWriteCheck(fileName string) HealthyCheck {
	return func() error {
		f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0664)
		if err != nil {
			return err
		}
		defer f.Close()

		line := `123456789`

		/*
			n, err := f.WriteString(line)
			if n != 9 || err != nil {
				return err
			}

			f.Seek(0, 0)
		*/

		out := make([]byte, 9)
		n, err := f.Read(out)
		if n != 9 || (err != nil && err != io.EOF) {
			return fmt.Errorf("file error: readout failed")
		}

		if string(out) != line {
			return fmt.Errorf("file error: in != out")
		}

		return nil
	}
}

func MysqlPingCheck(dsn string) HealthyCheck {
	return func() error {
		// TODO: implement later ...
		return nil
	}
}

func RedisPingCheck(dst string) HealthyCheck {
	return func() error {
		// TODO: implement later ...
		return nil
	}
}

func DNSResolveCheck(domain string) HealthyCheck {
	return func() error {
		// TODO: implement later ...
		return nil
	}
}

func MemoryMaxCheck(max int64) HealthyCheck {
	return func() error {
		// TODO: implement later ...
		return nil
	}
}


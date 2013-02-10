package configflag

import (
	"io"
	"os"
	"fmt"
	"flag"
	"bufio"
	"strings"
)

func loadFlagsFromFile(config string) ([]string, error) {

	input, err := os.Open(config)
	if err != nil {
		return nil, fmt.Errorf("failure opening flags config file: %s: %s", config, err)
	}

	defer input.Close()

	var flags []string
	var num int = 0

	reader := bufio.NewReader(input)
	for line, pref, fail := reader.ReadLine(); fail != io.EOF; line, pref, fail = reader.ReadLine() {
		if fail != nil {
			err = fmt.Errorf("failure reading line from flags config file: %s: %s", config, err)
			break
		}
		num++
		if pref {
			return nil, fmt.Errorf("very long flags config line at %d", num)
		}
		f := strings.TrimSpace(string(line))
		if f == "" || f[:1] == "#" {
			continue
		}
		flags = append(flags, f)
	}

	return flags, err
}

func Load(fs *flag.FlagSet, config string) error {
	f, err := loadFlagsFromFile(config)
	if err != nil {
		return fmt.Errorf("failure reading config flags: %s", err)
	}

	err = fs.Parse(f)
	if err != nil {
		return fmt.Errorf("failure parsing config flags: %s", err)
	}

	return nil
}

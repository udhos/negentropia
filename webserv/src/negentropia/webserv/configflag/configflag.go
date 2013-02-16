package configflag

import (
	"io"
	"os"
	//"log"
	"fmt"
	"flag"
	"bufio"
	"strings"
)

type FileList []string

// String is the method to get the flag value, part of the flag.Value interface.
func (fl *FileList) String() string {
	return strings.Join(*fl, ",")
}

// Set is the method to set the flag value, part of the flag.Value interface.
// Set's argument is a string to be parsed to set the flag.
// It's a comma-separated list, so we split it.
func (fl *FileList) Set(value string) error {
	*fl = append(*fl, strings.Split(value, ",")...)
	return nil
}

func loadOneFile(config string, flags *[]string) (error) {

	//log.Printf("loadOneFile: %s", config)

	input, err := os.Open(config)
	if err != nil {
		return fmt.Errorf("failure opening flags config file: %s: %s", config, err)
	}

	defer input.Close()

	var num int = 0

	reader := bufio.NewReader(input)
	for line, pref, fail := reader.ReadLine(); fail != io.EOF; line, pref, fail = reader.ReadLine() {
		if fail != nil {
			err = fmt.Errorf("failure reading line from flags config file: %s: %s", config, err)
			break
		}
		num++
		if pref {
			return fmt.Errorf("very long flag from config file %s at line %d", config, num)
		}
		f := strings.TrimSpace(string(line))
		if f == "" || f[:1] == "#" {
			continue
		}
		*flags = append(*flags, f)
	}

	//log.Printf("loadOneFile: %s: loaded %d flags", config, len(*flags))

	return err
}

func Load(fs *flag.FlagSet, configs FileList) error {
	var f []string
	
    for _, cfg := range configs {
        if err := loadOneFile(cfg, &f); err != nil {
			return fmt.Errorf("failure loading config file: %s: %s", cfg, err)
		}
    }

	//log.Printf("Load: %d total flags", len(f))

	if err := fs.Parse(f); err != nil {
		return fmt.Errorf("failure parsing config flags: %s", err)
	}

	return nil
}

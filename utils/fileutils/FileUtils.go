package fileutils

import (
	"bufio"
	"os"
)

func ReadFile(filename string) (map[string]struct{}, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	set := make(map[string]struct{})
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		set[scanner.Text()] = struct{}{}
	}
	return set, scanner.Err()
}

package files

import "os"
import "fmt"

// Exists returns whether the given file or directory exists.
func Exists(path string) (bool, error) {
	if path == nil {
		return false, fmt.Errorf("invalid input path")
	}
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return true, err
		}
	}
	return true, nil
}

// Exists returns whether the given file or directory exists.
func IsFile(path string) (bool, error) {
	if path == nil {
		return false, fmt.Errorf("invalid input path")
	}
	stat, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, err
		}
	}
	return !stat.IsDir(), nil
}

// Exists returns whether the given file or directory exists.
func IsDir(path string) (bool, error) {
	if path == nil {
		return false, fmt.Errorf("invalid input path")
	}
	stat, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, err
		}
	}
	return stat.IsDir(), nil
}

package problem

import (
	"archive/zip"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func DoesFileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

// readDir reads the directory named by dirname and returns
// a list of directory entries sorted by filename.
func readDir(dirname string) ([]string, error) {
	f, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}
	list, err := f.Readdir(-1)
	_ = f.Close()
	if err != nil {
		return nil, err
	}

	sort.Slice(list, func(i, j int) bool { return list[i].Name() < list[j].Name() })

	// Filter all subdirectories that is in the form "Test<xx>"
	var ans []string
	for _, x := range list {
		if !x.IsDir() {
			continue
		}
		tokens := strings.Split(x.Name(), "Test")
		if len(tokens) != 2 || tokens[0] != "" {
			continue
		}
		ans = append(ans, tokens[1])
	}

	if len(ans) == 0 {
		return []string{}, errors.New("no test detected")
	}
	return ans, nil
}

// unzip will decompress a zip archive, moving all files and folders
// within the zip file (parameter 1) to an output directory (parameter 2).
func unzip(src string, dest string) ([]string, error) {
	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {
		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return filenames, err
		}
	}
	return filenames, nil
}

func normaliseProblem(problem *Problem) {
	problem.Code = strings.TrimSpace(problem.Code)
	problem.Name = strings.TrimSpace(problem.Name)
}

func parseProblemFromRows(rows *sql.Rows) (Problem, error) {
	problem := Problem{}
	if err := rows.Scan(
		&problem.Id,
		&problem.Code,
		&problem.Name,
		&problem.TimeLimit,
		&problem.MemoryLimit,
		&problem.MaxScore,
		&problem.TestCount,
	); err != nil {
		return problem, err
	}
	normaliseProblem(&problem)
	return problem, nil
}
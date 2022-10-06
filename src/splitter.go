// Responsible for splitting files into parts.
package main

import (
	"math"
	"os"
	"path/filepath"
	"strconv"
)

// Splits *filename* into parts of size *limit*
// and saves them in *dest*.
func split(filename string, dest string, limit int64) []string {
	file, err := os.Open(filename)
	check(err)

	defer file.Close()

	os.MkdirAll(dest, os.ModePerm)

	var res []string

	file_stats, _ := file.Stat()

	var file_size int64 = file_stats.Size()

	var file_part_size = limit * (1 << 20)

	total_parts := int64(math.Ceil(float64(file_size) / float64(file_part_size)))

	for i := int64(0); i < total_parts; i++ {
		part_size := int(math.Min(float64(file_part_size), float64(file_size-int64(i*file_part_size))))
		part_buffer := make([]byte, part_size)

		file.Read(part_buffer)

		part_file_name := filepath.Join(dest, filepath.Base(filename)+"."+rjust(strconv.FormatInt(i, 10), 3, "0"))
		_, err := os.Create(part_file_name)
		check(err)

		err = os.WriteFile(part_file_name, part_buffer, os.ModeAppend)
		check(err)

		abs_path, err := filepath.Abs(part_file_name)
		check(err)

		res = append(res, abs_path)
	}

	return res
}

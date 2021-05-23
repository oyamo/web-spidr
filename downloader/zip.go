package downloader

import (
	"archive/zip"
	"io"
	"os"
)

func zipFolder(path string, fileName string) {
	zipFile, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}

	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	directory, err := os.ReadDir(path)
	println(path)
	if err != nil {
		panic(err)
	}

	for _, entry := range directory {
		fileInfo, err :=entry.Info()
		if err != nil {
			return
		}

		filePath := "." + string(os.PathSeparator) + folder + string(os.PathSeparator) + fileInfo.Name()
		file, _ := os.Open(filePath)

		info, err := file.Stat()
		if err != nil {
			return
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return
		}
		header.Name = fileInfo.Name()

		header.Method = zip.Deflate

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return
		}
		_, _ = io.Copy(writer, file)

		file.Close()
	}

}

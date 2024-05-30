package data

import (
	"bufio"
	"errors"
	"github.com/spf13/afero"
	"net.dalva.GvawSkinSync/logger"
)

func LoadLines(fname string, fs afero.Fs) (dest []string) {
	in, err := fs.Open(fname)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error opening")
		return nil
	}
	defer in.Close()

	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		dest = append(dest, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		logger.Log.Error().Err(err).Msg("Error reading")
		return nil
	}

	logger.Log.Debug().Str("fname", fname).Msg("Loaded")
	return dest
}

func ListSubDirs(fname string, fs afero.Fs) (subdirs []NamedDir) {
	files, err := afero.ReadDir(fs, fname)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error reading dir")
		return nil
	}
	for _, dir := range files {
		if dir.IsDir() {
			d := NamedDir{
				FullDirName: fname + "/" + dir.Name(),
				Name:        Mixdir(dir.Name()),
			}
			subdirs = append(subdirs, d)
			logger.Log.Trace().Str("dirname", fname+"/"+dir.Name()).Msg("Subdir")
		}
	}
	return subdirs
}

type NamedDir struct {
	FullDirName string
	Name        Mixdir
}

type Mixdir string

func (a *Mixdir) Add(b Mixdir) string {
	return string(*a) + "_" + string(b)
}

func (base *Mixdir) AddPers(personal Mixdir) string {
	return "[" + string(personal) + "] " + string(*base)
}

func ReadOffsetChunk(file afero.File, chunkNumber int64, chunkSize int32) ([]byte, error) {
	buffer := make([]byte, chunkSize)
	read, err := file.ReadAt(buffer, chunkNumber*int64(chunkSize))
	if err != nil && err.Error() != "EOF" {
		return nil, err
	}
	return buffer[:read], nil
}

func WriteOffsetChunk(file afero.File, bytes []byte, chunkNumber int64, chunkSize int32) error {
	n, err := file.WriteAt(bytes, chunkNumber*int64(chunkSize))
	if err != nil {
		return err
	}
	if n != len(bytes) {
		return errors.New("failed to write all bytes")
	}
	return nil
}

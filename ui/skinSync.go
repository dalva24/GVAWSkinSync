package ui

import (
	"context"
	"github.com/mholt/archiver/v4"
	"github.com/spf13/afero"
	alastor "net.dalva.GvawSkinSync/Alastor"
	"net.dalva.GvawSkinSync/logger"
	"os"
)

var chunkSize = 1024 * 1024

func skinSync() {
	logger.Log.Info().Msg("Basedir Prep")
	baseDir := savedGamesDir.Text + "\\Mods\\Tech\\GVAW-SkinSync"
	err := os.MkdirAll(baseDir, 0666)
	if err != nil {
		ShowInfo("Error", "Cannot make "+baseDir)
		return
	}
	var fs = afero.NewBasePathFs(afero.NewOsFs(), baseDir)

	download(fs, "descriptor.7z", ".", ".")
	extractAndCreateFiles("descriptor.7z", fs, true, ".")

	downloadAircrafts(fs)
}

func download(fs afero.Fs, fname string, source string, dest string) {
	err := fs.MkdirAll(dest, 0666)
	if err != nil {
		ShowInfo("Error", "Cannot download - cannot create "+dest)
		logger.Log.Error().Err(err).Msg("Cannot download - cannot create " + dest)
		return
	}
	logger.Log.Info().Str("source", source+"/"+fname).Str("dest", dest+"/"+fname).Msg("Downloading")
	f, err := fs.Create(dest + "/" + fname)
	if err != nil {
		ShowInfo("Error", "Cannot download - cannot create "+dest+"/"+fname)
		logger.Log.Error().Err(err).Msg("Cannot download - cannot create " + dest + "/" + fname)
		return
	}
	fw := &alastor.FlameWeaver{
		AuthCode:    authCode.Text,
		AddressPort: addressPort,
		FileName:    source + "/" + fname,
		ChunkSize:   chunkSize,
		Destination: &f,
	}
	err = fw.Weave()
	if err != nil {
		ShowInfo("Error", "Download failed: "+err.Error())
		logger.Log.Error().Err(err).Msg("Download failed")
		f.Close()
		return
	}
	f.Close()
}

func extractAndCreateFiles(archiveName string, fs afero.Fs, delArchive bool, extractTo string) {
	f, err := fs.Open(archiveName)
	if err != nil {
		ShowInfo("Error", "Cannot open file "+archiveName)
		logger.Log.Error().Err(err).Str("archiveName", archiveName).Msg("Cannot open file")
		return
	}
	defer f.Close()

	err = archiver.SevenZip{}.Extract(context.Background(),
		f, nil,
		func(ctx context.Context, f archiver.File) error {
			logger.Log.Trace().Str("fname", f.Name()).Msg("file")
			createFiles(fs, f, extractTo)
			return nil
		})
	if err != nil {
		logger.Log.Error().Err(err).Msg("Extract")
		ShowInfo("Error", "Cannot extract file "+archiveName)
		return
	}

	if delArchive {
		f.Close()
		err = fs.Remove(archiveName)
		if err != nil {
			logger.Log.Error().Err(err).Msg("Extract")
			ShowInfo("Error", "Cannot delete file "+archiveName)
			return
		}
	}
}

func createFiles(fs afero.Fs, file archiver.File, subdir string) {
	if file.IsDir() {
		err := fs.MkdirAll(file.NameInArchive, 0666)
		if err != nil {
			logger.Log.Error().Err(err).Msg("Mkdir")
			ShowInfo("Error", "Cannot create directory "+err.Error())
			return
		}
	} else {
		open, err := file.Open()
		if err != nil {
			logger.Log.Error().Err(err).Msg("Mkfile")
			ShowInfo("Error", "Cannot create file "+err.Error())
			return
		}
		err = afero.WriteReader(fs, subdir+"/"+file.NameInArchive, open)
		if err != nil {
			logger.Log.Error().Err(err).Msg("Write")
			ShowInfo("Error", "Cannot write file "+err.Error())
			return
		}
	}
}

func downloadAircrafts(fs afero.Fs) {
	acfts, err := alastor.LS(addressPort, authCode.Text, "base_full")
	if err != nil {
		logger.Log.Error().Err(err).Msg("DownloadAircraft")
		ShowInfo("Error", "Cannot Download aircraft list")
		return
	}
	for _, acft := range acfts {
		skins, err := alastor.LS(addressPort, authCode.Text, "base_full/"+acft)
		if err != nil {
			logger.Log.Error().Err(err).Msg("DownloadSkins")
			ShowInfo("Error", "Cannot Download aircraft skin list")
			return
		}
		for _, skin := range skins {
			download(fs, "description.lua",
				serverSkinPath(base, full, acft, skin),
				clientSkinPath(acft, skin),
			)
			download(fs, "textures.7z",
				serverSkinPath(base, full, acft, skin),
				clientSkinPath(acft, skin),
			)
			extractAndCreateFiles(
				clientSkinPath(acft, skin)+"/textures.7z",
				fs, true, clientSkinPath(acft, skin))
		}
	}
}

var base = "base"
var full = "full"
var half = "half"

func serverSkinPath(kind string, res string, aircraft string, skin string) string {
	return kind + "_" + res + "/" + aircraft + "/" + skin
}

func clientSkinPath(aircraft string, skin string) string {
	return "Liveries/" + aircraft + "/" + skin
}

func prepDownloadQuery(fname string, bs int, offset int64) alastor.DataQuery {
	return alastor.DataQuery{
		ApiKey:            authCode.Text,
		RequestedFilename: fname,
		ChunkSize:         int32(bs),
		ChunkOffset:       offset,
	}
}

func extractFile(fname string, destination string) {

}

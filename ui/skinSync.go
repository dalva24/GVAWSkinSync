package ui

import (
	"github.com/spf13/afero"
	"log"
	alastor "net.dalva.GvawSkinSync/Alastor"
	"os"
)

var chunkSize = 1024 * 1024

func skinSync() {
	baseDir := savedGamesDir.Text + "\\Mods\\Tech\\GVAW-SkinSync"
	err := os.MkdirAll(baseDir, 0666)
	if err != nil {
		return
	}
	log.Println("Basedir prep")

	var fs = afero.NewBasePathFs(afero.NewOsFs(), baseDir)
	f, err := fs.Create("descriptor.zip")
	if err != nil {
		return
	}
	log.Println("Creating descriptor.zip")

	fw := &alastor.FlameWeaver{
		AuthCode:    authCode.Text,
		AddressPort: addressPort,
		FileName:    "descriptor.zip",
		ChunkSize:   chunkSize,
		Destination: &f,
	}

	log.Println("weaving...")
	err = fw.Weave()
	if err != nil {
		return
	}

	downloadBaseMod()

	downloadAircrafts()
}

func downloadBaseMod() {

}

func downloadAircrafts() {

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

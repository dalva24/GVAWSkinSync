package ui

import (
	"github.com/spf13/afero"
	alastor "net.dalva.GvawSkinSync/Alastor"
	"net.dalva.GvawSkinSync/logger"
	"os"
	"os/exec"
	"syscall"
)

var chunkSize = 1024 * 1024

var aircrafts []string

func skinSync() {
	logger.Log.Info().Msg("Basedir Prep")
	baseDir := fullPath(".")
	err := os.MkdirAll(baseDir, 0666)
	if err != nil {
		ShowInfo("Error", "Cannot make "+baseDir)
		return
	}
	var fs = afero.NewBasePathFs(afero.NewOsFs(), baseDir)

	UpdateStatusMajor("Connected", 1, 0, 0)

	download(fs, "descriptor.7z", ".", ".", "descriptor")

	extractAndCreateFiles("descriptor.7z", fs, true, ".")

	err = fs.RemoveAll("Liveries")
	if err != nil {
		ShowInfo("Error", "Cannot scrub old skin "+err.Error())
		logger.Log.Error().Err(err).Msg("Cannot scrub old skin")
		return
	}

	downloadAircrafts(fs)

	if persSkin.Checked {
		downloadPersonalSkins(fs)
	}

	UpdateStatusMajor("Done", 4, 0, 0)
}

func download(fs afero.Fs, fname string, source string, dest string, shownName string) {
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
		AuthCode:        authCode.Text,
		AddressPort:     addressPort,
		FileName:        source + "/" + fname,
		ChunkSize:       chunkSize,
		Destination:     &f,
		StatusShownName: shownName,
		StatusMinor:     statusMinor,
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
	cmd := exec.Command("./7za.exe", "x", fullPath(archiveName), "-y", "-bd", "-o"+fullPath(extractTo))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 0x08000000}
	err := cmd.Run()
	if err != nil {
		logger.Log.Error().Err(err).Msg("Extract")
		return
	}

	if delArchive {
		err = fs.Remove(archiveName)
		if err != nil {
			logger.Log.Error().Err(err).Msg("Extract")
			ShowInfo("Error", "Cannot delete file "+archiveName)
			return
		}
	}
}

func downloadAircrafts(fs afero.Fs) {
	logger.Log.Info().Msg("Downloading Aircrafts")
	UpdateStatusMajor("Downloading Aircrafts...", 2, 0, 0)
	var err error
	aircrafts, err = alastor.LS(addressPort, authCode.Text, "base_full")
	if err != nil {
		logger.Log.Error().Err(err).Msg("DownloadAircraft")
		ShowInfo("Error", "Cannot Download aircraft list")
		return
	}
	for i, acft := range aircrafts {
		skins, err := alastor.LS(addressPort, authCode.Text, "base_full/"+acft)
		if err != nil {
			logger.Log.Error().Err(err).Msg("DownloadSkins")
			ShowInfo("Error", "Cannot Download aircraft skin list")
			return
		}
		for j, skin := range skins {
			UpdateStatusMajor("Downloading "+acft, 2, float64(i)/float64(len(aircrafts)), float64(j)/float64(len(skins)))
			download(fs, "description.lua",
				serverSkinPath(base, full, acft, skin),
				clientSkinPath(acft, skin),
				skin,
			)
			download(fs, "textures.7z",
				serverSkinPath(base, full, acft, skin),
				clientSkinPath(acft, skin),
				skin,
			)
			extractAndCreateFiles(
				clientSkinPath(acft, skin)+"/textures.7z",
				fs, true, clientSkinPath(acft, skin))
		}
	}
}

func downloadPersonalSkins(fs afero.Fs) {
	logger.Log.Info().Msg("Downloading User Skin Textures")
	UpdateStatusMajor("Downloading User Skins...", 3, 0, 0)
	for i, acft := range aircrafts {
		skins, err := alastor.LS(addressPort, authCode.Text, "personal_full/"+acft)
		if err != nil {
			logger.Log.Error().Err(err).Msg("DownloadSkins")
			ShowInfo("Error", "Cannot Download aircraft personal skin list")
			return
		}
		for j, skin := range skins {
			UpdateStatusMajor("DL User Skin "+acft, 2, float64(i)/float64(len(aircrafts)), float64(j)/float64(len(skins)))
			download(fs, "textures.7z",
				serverSkinPath(pers, full, acft, skin),
				clientPersonalSkinPath(acft, skin),
				skin,
			)
			extractAndCreateFiles(
				clientPersonalSkinPath(acft, skin)+"/textures.7z",
				fs, true, clientPersonalSkinPath(acft, skin))
		}
	}
	logger.Log.Info().Msg("Downloading User Skin Descriptor")
	UpdateStatusMajor("Downloading User Skin Descriptor...", 3, 0, 0)
	for _, acft := range aircrafts {
		skins, err := alastor.LS(addressPort, authCode.Text, "compiled/"+acft)
		if err != nil {
			logger.Log.Error().Err(err).Msg("DownloadSkins")
			ShowInfo("Error", "Cannot Download aircraft personal skin list")
			return
		}
		for _, skin := range skins {
			download(fs, "description.lua",
				serverCompiledSkinPath(acft, skin),
				clientSkinPath(acft, skin),
				skin,
			)
		}
	}
}

var base = "base"
var pers = "personal"
var full = "full"
var half = "half"

func serverSkinPath(kind string, res string, aircraft string, skin string) string {
	return kind + "_" + res + "/" + aircraft + "/" + skin
}

func serverCompiledSkinPath(aircraft string, skin string) string {
	return "compiled/" + aircraft + "/" + skin
}

func clientSkinPath(aircraft string, skin string) string {
	return "Liveries/" + aircraft + "/" + skin
}

func clientPersonalSkinPath(aircraft string, skin string) string {
	return "Liveries/" + aircraft + "/personal/" + skin
}

func fullPath(path string) string {
	return savedGamesDir.Text + "\\Mods\\Tech\\GVAW-SkinSync\\" + path
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

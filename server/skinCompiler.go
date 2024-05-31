package main

import (
	"bufio"
	"context"
	"github.com/mholt/archiver/v4"
	"net.dalva.GvawSkinSync/data"
	"net.dalva.GvawSkinSync/logger"
	"strings"
)
import "github.com/spf13/afero"

var fs = afero.NewBasePathFs(afero.NewOsFs(), "./server-data/")

func relPathDir(path string) string {
	return ".\\server-data\\" + path + "\\"
}

var (
	dirComp  = data.Mixdir("compiled")
	dirBase  = data.Mixdir("base")
	dirPers  = data.Mixdir("personal")
	resFull  = data.Mixdir("full")
	resHalf  = data.Mixdir("half")
	enabled  = data.Mixdir("on")
	disabled = data.Mixdir("off")
)

func sanityCheck() {
	logger.Log.Debug().
		Bool("descriptor.7z", checkExist("descriptor.7z", true)).
		Bool("aircraft.txt", checkExist("aircraft.txt", true)).
		Bool("base_full", checkExist(dirBase.Add(resFull), true)).
		Bool("personal_full", checkExist(dirPers.Add(resFull), true)).
		Msg("Checking Sanity...")
}

// updateSkins are unused and are
// Deprecated: currently unused since archiving is just too slow
func updateSkins() {
	sanityCheck()
	aircraft := data.LoadLines("aircraft.txt", fs)

	for _, acft := range aircraft {
		baseSkins := data.ListSubDirs("base_full/"+acft, fs)
		logger.Log.Info().Str("acft", acft).Msg("Updating Aircraft")
		for _, bSkin := range baseSkins {
			if checkExist(bSkin.FullDirName, false) {
				logger.Log.Info().Str("bSkin", relPathDir(bSkin.FullDirName)).Msg("Updating Skin")
				files, err := archiver.FilesFromDisk(nil, map[string]string{
					relPathDir(bSkin.FullDirName): "",
				})
				if err != nil {
					logger.Log.Error().Err(err).Str("bSkin", relPathDir(bSkin.FullDirName)).Msg("Updating Skin")
				}

				out, err := fs.Create("base_full/" + acft + ".tar.xz")
				if err != nil {
					logger.Log.Error().Err(err).Str("bSkin", relPathDir(bSkin.FullDirName)).Msg("Updating Skin")
				}

				format := archiver.CompressedArchive{
					Compression: archiver.Xz{},
					Archival:    archiver.Tar{},
				}

				err = format.Archive(context.Background(), out, files)
				if err != nil {
					out.Close()
					logger.Log.Error().Err(err).Str("bSkin", relPathDir(bSkin.FullDirName)).Msg("Updating Skin")
				}
				out.Close()
			} else {
				break
			}
		}
	}
}

func compileSkins(enabled data.Mixdir) {
	sanityCheck()

	err := fs.RemoveAll(dirComp.Add(enabled))
	if err != nil {
		logger.Log.Err(err).Msg("Failed to remove old files")
	}

	aircraft := data.LoadLines("aircraft.txt", fs)

	for _, acft := range aircraft {
		baseSkins := data.ListSubDirs("base_full/"+acft, fs)
		logger.Log.Info().Str("acft", acft).Msg("Processing Personalized Skins")
		personalSkins := data.ListSubDirs("personal_full/"+acft, fs)
		for _, bSkin := range baseSkins {
			//todo check if base skin allows customization
			for _, pSkin := range personalSkins {

				name := string(bSkin.Name)

				if strings.HasPrefix(name, "[GVAW] ") {
					newName := data.Mixdir(name[7:])
					//prepare folder
					compPersDir := dirComp.Add(enabled) + "/" + acft + "/" + newName.AddPers(pSkin.Name)
					logger.Log.Trace().Str("personal", compPersDir).Msg("Creating Skin")
					err := fs.MkdirAll(compPersDir, 0777)
					if err != nil {
						logger.Log.Error().Str("name", string(pSkin.Name)).Err(err).Msg("Error creating directory")
					}

					//prepare desc.lua
					compPersDesc, err := fs.Create(compPersDir + "/description.lua")
					if err != nil {
						logger.Log.Error().Str("name", string(pSkin.Name)).Err(err).Msg("Error creating description.lua")
					}

					//copy contents from base' desc.lua
					base, err := fs.Open(bSkin.FullDirName + "/description.lua")
					if err != nil {
						logger.Log.Error().Str("name", string(pSkin.Name)).Err(err).Msg("Error copying description.lua")
					}
					copyContents(base, compPersDesc, string(bSkin.Name))
					base.Close()

					if enabled == "on" {
						//add additional lines from personal skins
						pers, err := fs.Open(pSkin.FullDirName + "/description.lua")
						if err != nil {
							logger.Log.Error().Str("name", string(pSkin.Name)).Err(err).Msg("Error merging description.lua")
						}
						appendLuaContents(pers, compPersDesc, "personal/"+string(pSkin.Name))
						pers.Close()

						err = compPersDesc.Close()
						if err != nil {
							logger.Log.Error().Str("name", string(pSkin.Name)).Err(err).Msg("Error saving description.lua")
						}
					}
				}
			}
		}
	}
}

func copyContents(from afero.File, to afero.File, sourceFolder string) {
	scanner := bufio.NewScanner(from)
	for scanner.Scan() {
		toWrite := scanner.Text()
		if strings.HasPrefix(strings.TrimSpace(toWrite), "{") {
			split := strings.SplitAfter(toWrite, "\"")
			if len(split) == 5 {
				split[3] = "../" + sourceFolder + "/" + split[3]
			}
			toWrite = ""
			for _, s := range split {
				toWrite = toWrite + s
			}
		}
		_, err := to.WriteString(toWrite + "\n")
		if err != nil {
			logger.Log.Error().Err(err).Msg("Error writing")
		}
	}
}

func appendLuaContents(from afero.File, to afero.File, sourceFolder string) {
	scanner := bufio.NewScanner(from)
	for scanner.Scan() {
		toWrite := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(toWrite, "{") {
			split := strings.SplitAfter(toWrite, "\"")
			if len(split) == 5 {
				split[3] = "../" + sourceFolder + "/" + split[3]
			}
			toWrite = ""
			for _, s := range split {
				toWrite = toWrite + s
			}
		}
		_, err := to.WriteString("livery[#livery+1]=" + toWrite + "\n")
		if err != nil {
			logger.Log.Error().Err(err).Msg("Error writing")
		}
	}
}

func checkExist(path string, fatal bool) bool {
	exist, err := afero.Exists(fs, path)
	if err != nil {
		logger.Log.Err(err).Msg("Error checking existence")
		return false
	}
	if !exist && fatal {
		logger.Log.Fatal().Str("path", path).Msg("File does not exist")
	}
	return exist
}

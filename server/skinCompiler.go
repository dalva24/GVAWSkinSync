package main

import (
	"bufio"
	"net.dalva.GvawSkinSync/data"
	"net.dalva.GvawSkinSync/logger"
	"strings"
)
import "github.com/spf13/afero"

var fs = afero.NewBasePathFs(afero.NewOsFs(), "./server-data/")

var (
	dirComp = data.Mixdir("compiled")
	dirBase = data.Mixdir("base")
	dirPers = data.Mixdir("personal")
	resFull = data.Mixdir("full")
	resHalf = data.Mixdir("half")
)

func compileSkins() {
	logger.Log.Debug().
		Bool("descriptor.zip", checkExist("descriptor.zip", true)).
		Bool("aircraft.txt", checkExist("aircraft.txt", true)).
		Bool("base_full", checkExist(dirBase.Add(resFull), true)).
		Bool("personal_full", checkExist(dirPers.Add(resFull), true)).
		Msg("Checking Sanity...")

	err := fs.RemoveAll(dirComp.Add(resFull))
	if err != nil {
		logger.Log.Err(err).Msg("Failed to remove old files")
	}
	err = fs.RemoveAll(dirComp.Add(resHalf))
	if err != nil {
		logger.Log.Err(err).Msg("Failed to remove old files")
	}

	aircraft := data.LoadLines("aircraft.txt", fs)

	for _, acft := range aircraft {
		baseSkins := data.ListSubDirs("base_full/"+acft, fs)
		logger.Log.Info().Str("acft", acft).Msg("Processing Personalized Skins")
		personalSkins := data.ListSubDirs("personal_full/"+acft, fs)
		for _, bSkin := range baseSkins {
			for _, pSkin := range personalSkins {

				//prepare folder
				compPersDir := dirComp.Add(resFull) + "/" + bSkin.Name.AddPers(pSkin.Name)
				logger.Log.Trace().Str("personal", compPersDir).Msg("Creating Skin")
				err := fs.MkdirAll(dirComp.Add(resFull)+"/"+bSkin.Name.AddPers(pSkin.Name), 0666)
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

package main

import (
	"github.com/jfixby/coinknife"
	"github.com/jfixby/pin"
	"github.com/jfixby/pin/commandline"
	"github.com/jfixby/pin/fileops"
	"path/filepath"
	"strings"
)

func main() {

	set := &coinknife.Settings{
		PathToInputRepo:        `D:\PICFIGHT\src\github.com\jfixby\dcrregtest`,
		PathToOutputRepo:       `D:\PICFIGHT\src\github.com\picfight\pfcregtest`,
		DoNotProcessAnyFiles:   false,
		DoNotProcessSubfolders: false,
		FileNameProcessor:      nameGenerator,
		IsFileProcessable:      processableFiles,
		FileContentProcessor:   fileGenerator,
		IgnoredFiles:           ignoredFiles(),
		InjectorsPath:          filepath.Join("", "code_injections"),
	}

	coinknife.Build(set)
}

func nameGenerator(data string) string {
	//data = coinknife.Replace(data, "decred/dcrd", "picfight/pfcd")
	return data
}

func fileGenerator(data string) string {

	data = coinknife.Replace(data, "Decred regression testing", "PicFight coin regression testing")
	data = coinknife.Replace(data, "jfixby/dcr", "picfight/pfc")
	data = coinknife.Replace(data, "decred/dcrd", "picfight/pfcd")
	data = coinknife.Replace(data, "dcrharness", "pfcharness")
	data = coinknife.Replace(data, "dcrregtest", "pfcregtest")
	data = coinknife.Replace(data, "dcrwallet", "pfcwallet")
	data = coinknife.Replace(data, "DCR", "PFC")
	data = coinknife.Replace(data, "dcrd", "pfcd")

	return data
}

// ignoredFiles
func ignoredFiles() map[string]bool {
	ignore := make(map[string]bool)
	ignore[".git"] = true
	ignore[".github"] = true
	ignore[".idea"] = true
	ignore["pfcregtestbuilder"] = true
	return ignore
}

// processableFiles
func processableFiles(file string) bool {
	if strings.HasSuffix(file, ".png") {
		return false
	}
	if strings.HasSuffix(file, ".jpg") {
		return false
	}
	if strings.HasSuffix(file, ".jpeg") {
		return false
	}
	if strings.HasSuffix(file, ".exe") {
		return false
	}
	if strings.HasSuffix(file, ".svg") {
		return false
	}
	if strings.HasSuffix(file, ".ico") {
		return false
	}
	if strings.HasSuffix(file, ".bin") {
		return false
	}
	if strings.HasSuffix(file, ".bin") {
		return false
	}
	if strings.HasSuffix(file, ".db") {
		return false
	}
	if strings.HasSuffix(file, ".bz2") {
		return false
	}
	if strings.HasSuffix(file, ".gz") {
		return false
	}
	if strings.HasSuffix(file, ".hex") {
		return false
	}
	if strings.HasSuffix(file, ".mp4") {
		return false
	}
	if strings.HasSuffix(file, ".gif") {
		return false
	}
	if strings.HasSuffix(file, ".ttf") {
		return false
	}
	if strings.HasSuffix(file, ".icns") {
		return false
	}
	if strings.HasSuffix(file, ".woff") {
		return false
	}
	if strings.HasSuffix(file, ".woff2") {
		return false
	}
	if strings.HasSuffix(file, ".eot") {
		return false
	}
	if strings.HasSuffix(file, ".sum") {
		return false
	}
	//-
	if strings.HasSuffix(file, "api.proto") {
		return false
	}
	if strings.HasSuffix(file, ".pot") {
		return false
	}
	if strings.HasSuffix(file, ".gyp") {
		return false
	}
	if strings.HasSuffix(file, ".cc") {
		return false
	}
	if strings.HasSuffix(file, ".h") {
		return false
	}
	if strings.HasSuffix(file, "notes.sample") {
		return false
	}
	if strings.HasSuffix(file, ".desktop") {
		return false
	}
	if strings.HasSuffix(file, ".log") {
		return false
	}
	if strings.HasSuffix(file, "pfcd.service") {
		return false
	}
	if strings.HasSuffix(file, ".conf") {
		return false
	}
	if strings.HasSuffix(file, ".json") {
		return false
	}
	if strings.HasSuffix(file, ".py") {
		return false
	}
	if strings.HasSuffix(file, ".tmpl") {
		return false
	}
	if strings.HasSuffix(file, ".js") {
		return false
	}
		if strings.HasSuffix(file, ".css") {
		return false
	}
	if strings.HasSuffix(file, ".lock") {
		return false
	}
	if strings.HasSuffix(file, "LICENSE") {
		return false
	}
	if strings.HasSuffix(file, "CONTRIBUTORS") {
		return false
	}
	if strings.HasSuffix(file, "Dockerfile") {
		return false
	}
	if strings.HasSuffix(file, "Dockerfile.alpine") {
		return false
	}
	if strings.HasSuffix(file, "CHANGES") {
		return false
	}
	if strings.HasSuffix(file, ".iml") {
		return false
	}
	if strings.HasSuffix(file, ".yml") {
		return false
	}
	if strings.HasSuffix(file, ".toml") {
		return false
	}

	if strings.HasSuffix(file, ".xml") {
		return false
	}
	if strings.HasSuffix(file, ".gitignore") {
		return false
	}
	if strings.HasSuffix(file, ".editorconfig") {
		return false
	}
	if strings.HasSuffix(file, ".eslintignore") {
		return false
	}
	if strings.HasSuffix(file, ".stylelintrc") {
		return false
	}
	if strings.HasSuffix(file, "config") {
		return false
	}
	if strings.HasSuffix(file, ".html") {
		return false
	}
	if strings.HasSuffix(file, ".po") {
		return false
	}
	if strings.HasSuffix(file, ".less") {
		return false
	}

	//------------------------------
	if strings.HasSuffix(file, ".mod") {
		return true
	}
	if strings.HasSuffix(file, ".go") {
		return true
	}
	if strings.HasSuffix(file, ".sh") {
		return true
	}
	if strings.HasSuffix(file, ".md") {
		return true
	}

	pin.E("Unknown file type", file)
	return false
}

func fixSecp256k1Checksum(targetProject string) {
	invalidParent := filepath.Join(targetProject, "btcec")
	invalid := filepath.Join(invalidParent, "secp256k1.go")
	fileops.Delete(invalid)

	batName := "checksum_update.bat"
	batTemplate := filepath.Join("assets", batName)
	batData := fileops.ReadFileToString(batTemplate)
	batData = strings.Replace(batData, "#TARGET_FOLDER#", invalidParent, -1)
	batFile := filepath.Join(batName)
	fileops.WriteStringToFile(batFile, batData)

	ext := &commandline.ExternalProcess{
		CommandName: batFile,
	}
	ext.Launch(true)
	ext.Wait()
}

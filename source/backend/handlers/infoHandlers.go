package handlers

import (
	"archive/zip"
	"encoding/json"
	"github.com/zlepper/go-modpack-packer/source/backend/types"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"runtime/debug"
	"strings"
)

func gatherInformation(conn types.WebsocketConnection, data interface{}) {
	modpack := types.CreateSingleModpackData(data)
	gatherInformationAboutMods(path.Join(modpack.InputDirectory, "mods"), conn)
}

func gatherInformationAboutMods(inputDirectory string, conn types.WebsocketConnection) {
	filesAndDirectories, _ := ioutil.ReadDir(inputDirectory)
	files := make([]os.FileInfo, 0)
	for _, f := range filesAndDirectories {
		if f.IsDir() {
			continue
		}
		files = append(files, f)
	}
	conn.Write("total-mod-files", len(files))
	for _, f := range files {
		fullname := path.Join(inputDirectory, f.Name())
		go gatherInformationAboutMod(fullname, conn)
	}
}

func gatherInformationAboutMod(modfile string, conn types.WebsocketConnection) {
	reader, err := zip.OpenReader(modfile)
	if err != nil {
		if err == zip.ErrFormat {
			conn.Log("err: " + modfile + " is not a valid zip file")
			return
		} else {
			log.Panic(err)
		}
	}
	defer reader.Close()

	// Iterate the files in the archive to find the info files
	for _, f := range reader.File {
		// We only need .info and a certain .json file
		if strings.HasSuffix(f.Name, "mod.info") || f.Name == "litemod.json" {
			r, err := f.Open()
			if err != nil {
				log.Fatal(err)
			}
			readInfoFile(r, conn, f.FileInfo().Size(), modfile)
		}
	}
}

func readInfoFile(file io.ReadCloser, conn types.WebsocketConnection, size int64, filename string) {
	content := make([]byte, size)
	_, err := file.Read(content)
	content = []byte(strings.Replace(string(content), "\n", " ", -1))
	if err != nil {
		conn.Log(err.Error() + "\n" + string(debug.Stack()))
	}
	var mod types.ModInfo
	normalMod := make([]types.ModInfo, 0)
	err = json.Unmarshal(content, &normalMod)
	if err != nil {
		// Try with mod version 2, or with litemod
		err = json.Unmarshal(content, &mod)
		if err != nil {
			conn.Log(err.Error() + "\n" + string(content) + "\n" + filename)
		}
		// Handle version 2 mods
		if mod.ModListVersion == 2 {
			createModResponse(conn, mod.ModList[0], filename)
		} else {
			// Handle liteloader mods
			createModResponse(conn, mod, filename)
		}
		return
	}
	if len(normalMod) > 0 {
		createModResponse(conn, normalMod[0], filename)
	} else {
		createModResponse(conn, types.ModInfo{}, filename)
	}

}

func createModResponse(conn types.WebsocketConnection, mod types.ModInfo, filename string) {
	const modDataReadyEvent string = "mod-data-ready"
	modRes := mod.CreateModResponse(filename)
	modRes.NormalizeId()
	conn.Write(modDataReadyEvent, modRes)
}

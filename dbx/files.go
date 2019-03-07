package dbx

import (
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox"
	dbxFiles "github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/files"
)

type Folder struct {
	Name string
	Path string
}

type File struct {
	Name string
	Path string
}

// Returns list of folders and files of user's Dropbox account
func GetList(accessToken, path string) ([]Folder, []File, error) {
	config := dropbox.Config{
		Token: accessToken,
	}

	client := dbxFiles.New(config)

	result, err := client.ListFolder(&dbxFiles.ListFolderArg{
		Path: path,
	})
	if err != nil {
		return nil, nil, err
	}

	var folders []Folder
	var files []File

	for _, entry := range result.Entries {
		switch meta := entry.(type) {
		case *dbxFiles.FolderMetadata:
			folders = append(folders, Folder{
				Name: meta.Name,
				Path: meta.PathLower,
			})
		case *dbxFiles.FileMetadata:
			files = append(files, File{
				Name: meta.Name,
				Path: meta.PathLower,
			})
		}
	}

	return folders, files, nil
}

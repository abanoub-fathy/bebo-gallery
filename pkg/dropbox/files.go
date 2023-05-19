package dropbox

import (
	dbx "github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	dbxFiles "github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
)

type File struct {
	Name string
	Path string
}

type Folder struct {
	Name string
	Path string
}

// List will list files and folders for specific path
func List(accessToken string, path string) ([]Folder, []File, error) {
	dbxConfig := dbx.Config{
		Token:    accessToken,
		LogLevel: dbx.LogInfo, // if needed, set the desired logging level. Default is off
	}

	// create dbx files client
	filesClient := dbxFiles.New(dbxConfig)

	// hit the listFolder endpoint
	result, err := filesClient.ListFolder(&dbxFiles.ListFolderArg{Path: path})
	if err != nil {
		return nil, nil, err
	}

	// define files, folders
	files := []File{}
	folders := []Folder{}

	for _, entry := range result.Entries {
		switch metadata := entry.(type) {
		case *dbxFiles.FileMetadata:
			files = append(files, File{
				Name: metadata.Name,
				Path: metadata.PathLower,
			})
		case *dbxFiles.FolderMetadata:
			folders = append(folders, Folder{
				Name: metadata.Name,
				Path: metadata.PathLower,
			})
		}
	}

	return folders, files, nil
}

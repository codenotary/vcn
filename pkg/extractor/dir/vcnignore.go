/*
 * Copyright (c) 2018-2019 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package dir

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/src-d/go-git.v4/plumbing/format/gitignore"
)

// IgnoreFilename is the name of the ignore file used by this package when processing directories.
//
// The ignore file supports all pattern formats as specified in the gitignore specification:
//  - https://git-scm.com/docs/gitignore
//  - https://github.com/src-d/go-git/blob/master/plumbing/format/gitignore/doc.go
//
// However, this package implementation:
//  - only uses the ignore file within the root directory (nested ignore files will be treated as normal files)
//  - always ignores the manifest file (it cannot be excluded by the ignore file)
//
const IgnoreFilename = ".vcnignore"

// DefaultIgnoreFileContent is the content of ignore file with default patterns.
const DefaultIgnoreFileContent = `# Windows thumbnail cache files
Thumbs.db
Thumbs.db:encryptable
ehthumbs.db
ehthumbs_vista.db

# Windows folder config file
[Dd]esktop.ini

# Windows Recycle Bin used on file shares
$RECYCLE.BIN/

# macOS
.DS_Store
.AppleDouble
.LSOverride

# macOS Thumbnails
._*

# macOS files that might appear in the root of a volume
.DocumentRevisions-V100
.fseventsd
.Spotlight-V100
.TemporaryItems
.Trashes
.VolumeIcon.icns
.com.apple.timemachine.donotpresent

# Directories potentially created on remote AFP share
.AppleDB
.AppleDesktop
Network Trash Folder
Temporary Items
.apdisk

# temporary files which can be created if a process still has a handle open of a deleted file
.fuse_hidden*

# KDE directory preferences
.directory

# Linux trash folder which might appear on any partition or disk
.Trash-*

# .nfs files are created when an open file is removed but is still being accessed
.nfs*
`

const (
	ignorefileCommentPrefix = "#"
	ignorefileEOL           = "\n"
)

func getIgnoreFileData(f *os.File) (data []byte, err error) {
    if f != nil {
	    data, err := ioutil.ReadAll(f)
        return data, err
    }
    data = []byte(DefaultIgnoreFileContent)
    err = nil
    return
}

// newIgnoreFileMatcher reads and parses the ignore file in path and return a gitignore.Matcher.
// If the ignore file does not exists, create a matcher based on our template
func newIgnoreFileMatcher(path string) (m gitignore.Matcher, err error) {
    // Try to open the .vcnignore file
	f, err := os.Open(filepath.Join(path, IgnoreFilename))
	if err != nil {
        // That didn't go well, let's see if it's just because it isn't there
		if os.IsNotExist(err) {
            // File does not exist, let's use the template instead
            f = nil
            err = nil
		} else {
            // Unexpected error.. (possibly permission)
            // Let's not continue. Something may be very wrong
		    return
        }
	}

    if f != nil {
        // Ensure file is closed if we opened it
	    defer f.Close()
    }
    // Create the git ignore pattern
	ps := []gitignore.Pattern{}
    // get the file contents of the .vcnignore file or the hardcoded template
    data, err := getIgnoreFileData( f )
    if err!= nil {
        return
    }
    // Parse the contents of the file or template
    for _, s := range strings.Split(string(data), ignorefileEOL) {
        // Skip over empty lines and lines starting with a comment smbol
        if !strings.HasPrefix(s, ignorefileCommentPrefix) && len(strings.TrimSpace(s)) > 0 {
            // Useful line - add it to the parser
            ps = append(ps, gitignore.ParsePattern(s, nil))
        }
    }
    // Create a matcher from the paterns
	m = gitignore.NewMatcher(ps)
	return
}

// initIgnoreFile writes the default ignore file if it does not exist.
func initIgnoreFile(root string) error {
	filename := filepath.Join(root, IgnoreFilename)

	// create and open the file if not exists
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		if os.IsExist(err) {
			return nil // file exists already
		}
		return err
	}

	// otherwise, write the default content
	_, err = f.WriteString(DefaultIgnoreFileContent)
	return err
}

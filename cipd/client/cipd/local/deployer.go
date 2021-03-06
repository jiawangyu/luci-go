// Copyright 2014 The LUCI Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package local

import (
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"golang.org/x/net/context"

	"go.chromium.org/luci/cipd/common"
	"go.chromium.org/luci/common/data/sortby"
	"go.chromium.org/luci/common/data/stringset"
	"go.chromium.org/luci/common/logging"
)

// TODO(vadimsh): How to handle path conflicts between two packages? Currently
// the last one installed wins.

// TODO(vadimsh): Use some sort of file lock to reduce a chance of corruption.

// File system layout of a site directory <base> for "symlink" install method:
// <base>/.cipd/pkgs/
//   <arbitrary index>/
//     description.json
//     _current -> symlink to fea3ab83440e9dfb813785e16d4101f331ed44f4
//     fea3ab83440e9dfb813785e16d4101f331ed44f4/
//       bin/
//         tool
//         ...
//       ...
// bin/
//    tool -> symlink to ../.cipd/pkgs/<package name digest>/_current/bin/tool
//    ...
//
// Where <arbitrary index> is chosen to be the smallest number available for
// this installation (installing more packages gets higher numbers, removing
// packages and installing new ones will reuse the smallest ones).
//
// Some efforts are made to make sure that during the deployment a window of
// inconsistency in the file system is as small as possible.
//
// For "copy" install method everything is much simpler: files are directly
// copied to the site root directory and .cipd/pkgs/* contains only metadata,
// such as description and manifest files with a list of extracted files (to
// know what to uninstall).

// Deployer knows how to unzip and place packages into site root directory.
type Deployer interface {
	// DeployInstance installs an instance of a package into the given subdir of
	// the root.
	//
	// It unpacks the package into <base>/.cipd/pkgs/*, and rearranges
	// symlinks to point to unpacked files. It tries to make it as "atomic" as
	// possible. Returns information about the deployed instance.
	//
	// Due to a historical bug, if inst contains any files which are intended to
	// be deployed to `.cipd/*`, they will not be extracted and you'll see
	// warnings logged.
	DeployInstance(ctx context.Context, subdir string, inst PackageInstance) (common.Pin, error)

	// CheckDeployed checks whether a given package is deployed at the given
	// subdir.
	//
	// It returns information about installed version (or error if not installed).
	CheckDeployed(ctx context.Context, subdir, packageName string) (common.Pin, error)

	// FindDeployed returns a list of packages deployed to a site root.
	FindDeployed(ctx context.Context) (out common.PinSliceBySubdir, err error)

	// RemoveDeployed deletes a package from a subdir given its name.
	RemoveDeployed(ctx context.Context, subdir, packageName string) error

	// TempFile returns os.File located in <base>/.cipd/tmp/*.
	//
	// The file is open for reading and writing.
	TempFile(ctx context.Context, prefix string) (*os.File, error)

	// CleanupTrash attempts to remove stale files.
	//
	// May return errors if some files are still locked, this is fine.
	CleanupTrash(ctx context.Context) error
}

// NewDeployer return default Deployer implementation.
func NewDeployer(root string) Deployer {
	var err error
	if root == "" {
		err = fmt.Errorf("site root path is not provided")
	} else {
		root, err = filepath.Abs(filepath.Clean(root))
	}
	if err != nil {
		return errDeployer{err}
	}
	trashDir := filepath.Join(root, SiteServiceDir, "trash")
	return &deployerImpl{NewFileSystem(root, trashDir)}
}

////////////////////////////////////////////////////////////////////////////////
// Implementation that returns error on all requests.

type errDeployer struct{ err error }

func (d errDeployer) DeployInstance(context.Context, string, PackageInstance) (common.Pin, error) {
	return common.Pin{}, d.err
}

func (d errDeployer) CheckDeployed(context.Context, string, string) (common.Pin, error) {
	return common.Pin{}, d.err
}

func (d errDeployer) FindDeployed(context.Context) (out common.PinSliceBySubdir, err error) {
	return nil, d.err
}
func (d errDeployer) RemoveDeployed(context.Context, string, string) error { return d.err }
func (d errDeployer) TempFile(context.Context, string) (*os.File, error)   { return nil, d.err }
func (d errDeployer) CleanupTrash(context.Context) error                   { return d.err }

////////////////////////////////////////////////////////////////////////////////
// Real deployer implementation.

// packagesDir is a subdirectory of site root to extract packages to.
const packagesDir = SiteServiceDir + "/pkgs"

// currentSymlink is a name of a symlink that points to latest deployed version.
// Used on Linux and Mac.
const currentSymlink = "_current"

// currentTxt is a name of a text file with instance ID of latest deployed
// version. Used on Windows.
const currentTxt = "_current.txt"

// deployerImpl implements Deployer interface.
type deployerImpl struct {
	fs FileSystem
}

func (d *deployerImpl) DeployInstance(ctx context.Context, subdir string, inst PackageInstance) (common.Pin, error) {
	if err := common.ValidateSubdir(subdir); err != nil {
		return common.Pin{}, err
	}

	pin := inst.Pin()
	logging.Infof(ctx, "Deploying %s into %s(/%s)", pin, d.fs.Root(), subdir)

	// Be paranoid.
	if err := common.ValidatePin(pin); err != nil {
		return common.Pin{}, err
	}
	if _, err := d.fs.EnsureDirectory(ctx, filepath.Join(d.fs.Root(), subdir)); err != nil {
		return common.Pin{}, err
	}

	// Extract new version to the .cipd/pkgs/* guts. For "symlink" install mode it
	// is the final destination. For "copy" install mode it's a temp destination
	// and files will be moved to the site root later (in addToSiteRoot call).
	// ExtractPackageInstance knows how to build full paths and how to atomically
	// extract a package. No need to delete garbage if it fails.
	pkgPath, err := d.packagePath(ctx, subdir, pin.PackageName, true)
	if err != nil {
		return common.Pin{}, err
	}

	svcDir := SiteServiceDir + "/"
	filterCipd := func(f File) bool {
		name := f.Name()
		if strings.HasPrefix(name, svcDir) {
			logging.Warningf(ctx, "[non-fatal] ignoring internal file: %s", name)
			return true
		}
		return false
	}

	// Unzip the package into the final destination inside .cipd/* guts.
	destPath := filepath.Join(pkgPath, pin.InstanceID)
	if err := ExtractInstance(ctx, inst, NewFileSystemDestination(destPath, d.fs), filterCipd); err != nil {
		return common.Pin{}, err
	}

	// We want to cleanup 'destPath' if something is not right with it.
	deleteFailedInstall := true
	defer func() {
		if deleteFailedInstall {
			logging.Warningf(ctx, "Deploy aborted, cleaning up %s", destPath)
			d.fs.EnsureDirectoryGone(ctx, destPath)
		}
	}()

	// Read and sanity check the manifest.
	newManifest, err := d.readManifest(ctx, destPath)
	if err != nil {
		return common.Pin{}, err
	}
	installMode, err := checkInstallMode(newManifest.InstallMode)
	if err != nil {
		return common.Pin{}, err
	}

	// Remember currently deployed version (to remove it later). Do not freak out
	// if it's not there (prevInstanceID == "") or broken (err != nil).
	prevInstanceID, err := d.getCurrentInstanceID(pkgPath)
	prevManifest := Manifest{}
	if err == nil && prevInstanceID != "" {
		prevManifest, err = d.readManifest(ctx, filepath.Join(pkgPath, prevInstanceID))
	}
	if err != nil {
		logging.Warningf(ctx, "Previous version of the package is broken: %s", err)
		prevManifest = Manifest{} // to make sure prevManifest.Files == nil.
	}

	// Install all new files to the site root.
	err = d.addToSiteRoot(ctx, subdir, newManifest.Files, installMode, pkgPath, destPath)
	if err != nil {
		return common.Pin{}, err
	}

	// Mark installed instance as a current one. After this call the package is
	// considered installed and the function must not fail. All cleanup below is
	// best effort.
	if err = d.setCurrentInstanceID(ctx, pkgPath, pin.InstanceID); err != nil {
		return common.Pin{}, err
	}
	deleteFailedInstall = false

	// Wait for async cleanup to finish.
	wg := sync.WaitGroup{}
	defer wg.Wait()

	// When using 'copy' install mode all files (except .cipdpkg/*) are moved away
	// from 'destPath', leaving only an empty husk with directory structure.
	// Remove it to save some inodes.
	if installMode == InstallModeCopy {
		wg.Add(1)
		go func() {
			defer wg.Done()
			removeEmptyTree(destPath, func(string) bool { return true })
		}()
	}

	// Remove old instance directory completely.
	if prevInstanceID != "" && prevInstanceID != pin.InstanceID {
		wg.Add(1)
		go func() {
			defer wg.Done()
			d.fs.EnsureDirectoryGone(ctx, filepath.Join(pkgPath, prevInstanceID))
		}()
	}

	// Remove no longer present files from the site root directory.
	if len(prevManifest.Files) > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			toKeep := map[string]bool{}
			for _, f := range newManifest.Files {
				toKeep[f.Name] = true
			}
			toKill := []FileInfo{}
			for _, f := range prevManifest.Files {
				if !toKeep[f.Name] {
					toKill = append(toKill, f)
				}
			}
			d.removeFromSiteRoot(ctx, subdir, toKill)
		}()
	}

	// Verify it's all right.
	newPin, err := d.CheckDeployed(ctx, subdir, pin.PackageName)
	if err == nil && newPin.InstanceID != pin.InstanceID {
		err = fmt.Errorf("other instance (%s) was deployed concurrently", newPin.InstanceID)
	}
	if err == nil {
		logging.Infof(ctx, "Successfully deployed %s", pin)
	} else {
		logging.Errorf(ctx, "Failed to deploy %s: %s", pin, err)
	}
	return newPin, err
}

func (d *deployerImpl) CheckDeployed(ctx context.Context, subdir, pkg string) (common.Pin, error) {
	if err := common.ValidateSubdir(subdir); err != nil {
		return common.Pin{}, err
	}

	pkgPath, err := d.packagePath(ctx, subdir, pkg, false)
	if err != nil {
		return common.Pin{}, err
	}
	if pkgPath == "" {
		return common.Pin{}, fmt.Errorf("package %s is not installed", pkg)
	}

	current, err := d.getCurrentInstanceID(pkgPath)
	if err != nil {
		return common.Pin{}, err
	}
	if current == "" {
		return common.Pin{}, fmt.Errorf("package %s is not installed", pkg)
	}
	return common.Pin{
		PackageName: pkg,
		InstanceID:  current,
	}, nil
}

func (d *deployerImpl) FindDeployed(ctx context.Context) (common.PinSliceBySubdir, error) {
	// Directories with packages are direct children of .cipd/pkgs/.
	pkgs := filepath.Join(d.fs.Root(), filepath.FromSlash(packagesDir))
	infos, err := ioutil.ReadDir(pkgs)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	found := common.PinMapBySubdir{}
	for _, info := range infos {
		if !info.IsDir() {
			continue
		}
		// Read the description and the 'current' link.
		pkgPath := filepath.Join(pkgs, info.Name())
		desc, err := d.readDescription(ctx, pkgPath)
		if err != nil || desc == nil {
			continue
		}
		currentID, err := d.getCurrentInstanceID(pkgPath)
		if err != nil || currentID == "" {
			continue
		}

		// Ignore duplicate entries, they can appear if someone messes with pkgs/*
		// structure manually.
		if _, ok := found[desc.Subdir][desc.PackageName]; !ok {
			if _, ok := found[desc.Subdir]; !ok {
				found[desc.Subdir] = common.PinMap{}
			}
			found[desc.Subdir][desc.PackageName] = currentID
		}
	}

	return found.ToSlice(), nil
}

func (d *deployerImpl) RemoveDeployed(ctx context.Context, subdir, packageName string) error {
	if err := common.ValidateSubdir(subdir); err != nil {
		return err
	}

	logging.Infof(ctx, "Removing %s from %s(/%s)", packageName, d.fs.Root(), subdir)
	if err := common.ValidatePackageName(packageName); err != nil {
		return err
	}
	pkgPath, err := d.packagePath(ctx, subdir, packageName, false)
	if err != nil {
		return err
	}
	if pkgPath == "" {
		logging.Warningf(ctx, "Package %s not found", packageName)
		return nil
	}

	// Read the manifest of the currently installed version.
	manifest := Manifest{}
	currentID, err := d.getCurrentInstanceID(pkgPath)
	if err == nil && currentID != "" {
		manifest, err = d.readManifest(ctx, filepath.Join(pkgPath, currentID))
	}

	// Warn, but continue with removal anyway. EnsureDirectoryGone call below
	// will nuke everything (even if it's half broken).
	if err != nil {
		logging.Warningf(ctx, "Package %s is in a broken state: %s", packageName, err)
	} else {
		d.removeFromSiteRoot(ctx, subdir, manifest.Files)
	}
	return d.fs.EnsureDirectoryGone(ctx, pkgPath)
}

func (d *deployerImpl) TempFile(ctx context.Context, prefix string) (*os.File, error) {
	dir, err := d.fs.EnsureDirectory(ctx, filepath.Join(d.fs.Root(), SiteServiceDir, "tmp"))
	if err != nil {
		return nil, err
	}
	return ioutil.TempFile(dir, prefix)
}

func (d *deployerImpl) TempDir(ctx context.Context, prefix string, mode os.FileMode) (string, error) {
	dir, err := d.fs.EnsureDirectory(ctx, filepath.Join(d.fs.Root(), SiteServiceDir, "tmp"))
	if err != nil {
		return "", err
	}
	return tempDir(dir, prefix, mode)
}

func (d *deployerImpl) CleanupTrash(ctx context.Context) error {
	return d.fs.CleanupTrash(ctx)
}

////////////////////////////////////////////////////////////////////////////////
// Utility methods.

type numSet sort.IntSlice

func (s *numSet) addNum(n int) {
	idx := sort.IntSlice((*s)).Search(n)
	if idx == len(*s) {
		// it's insertion point is off the end of the slice
		*s = append(*s, n)
	} else if (*s)[idx] != n {
		// it's insertion point is inside the slice, but is not present.
		*s = append(*s, 0)
		copy((*s)[idx+1:], (*s)[idx:])
		(*s)[idx] = n
	}
	// it's already present in the slice
}

func (s *numSet) smallestNewNum() int {
	prev := -1
	for _, n := range *s {
		if n-1 != prev {
			return prev + 1
		}
		prev = n
	}
	return prev + 1
}

// packagePath returns a path to a package directory in .cipd/pkgs/.
//
// This will scan all directories under pkgs, looking for a description.json. If
// an old-style package folder is encountered (e.g. has an instance folder and
// current manifest, but doesn't have a description.json), the description.json
// will be added.
//
// If no suitable path is found and allocate is true, this will create a new
// directory with an accompanying description.json. Otherwise this returns "".
func (d *deployerImpl) packagePath(ctx context.Context, subdir, pkg string, allocate bool) (string, error) {
	if err := common.ValidateSubdir(subdir); err != nil {
		return "", err
	}

	if err := common.ValidatePackageName(pkg); err != nil {
		return "", err
	}

	rel := filepath.FromSlash(packagesDir)
	abs, err := d.fs.RootRelToAbs(rel)
	if err != nil {
		logging.Errorf(ctx, "Can't get absolute path of %q: %s", rel, err)
		return "", err
	}

	seenNumbers, curPkgs := d.resolveValidPackageDirs(ctx, abs)
	if cur, ok := curPkgs[Description{subdir, pkg}]; ok {
		return cur, nil
	}

	if !allocate {
		return "", nil
	}

	// we didn't find one, so we have to make one
	if _, err := d.fs.EnsureDirectory(ctx, abs); err != nil {
		logging.Errorf(ctx, "Cannot ensure packages directory: %s", err)
		return "", err
	}

	// take the last 2 components from the pkg path.
	pkgParts := strings.Split(pkg, "/")
	prefix := ""
	if len(pkgParts) > 2 {
		prefix = strings.Join(pkgParts[len(pkgParts)-2:], "_")
	} else {
		prefix = strings.Join(pkgParts, "_")
	}
	// 0777 allows umask to take effect
	tmpDir, err := d.TempDir(ctx, prefix, 0777)
	if err != nil {
		logging.Errorf(ctx, "Cannot create new pkg tempdir: %s", err)
		return "", err
	}
	defer d.fs.EnsureDirectoryGone(ctx, tmpDir)
	err = d.fs.EnsureFile(ctx, filepath.Join(tmpDir, descriptionName), func(f *os.File) error {
		return writeDescription(&Description{Subdir: subdir, PackageName: pkg}, f)
	})
	if err != nil {
		logging.Errorf(ctx, "Cannot create new pkg description.json: %s", err)
		return "", err
	}

	// now we have to find a suitable index folder for it.
	for attempts := 0; attempts < 3; attempts++ {
		if attempts > 0 {
			// random sleep up to 1s to help avoid collisions between clients.
			time.Sleep(time.Duration(rand.Int31n(1000)) * time.Millisecond)
		}
		n := seenNumbers.smallestNewNum()
		seenNumbers.addNum(n)

		pkgPath := filepath.Join(abs, strconv.Itoa(n))
		// We use os.Rename instead of d.fs.Replace because we want it to fail if
		// the target directory already exists.
		switch err := os.Rename(tmpDir, pkgPath); le := err.(type) {
		case nil:
			return pkgPath, nil

		case *os.LinkError:
			if le.Err != syscall.ENOTEMPTY {
				logging.Errorf(ctx, "Error while creating pkg dir %s: %s", pkgPath, err)
				return "", err
			}

		default:
			logging.Errorf(ctx, "Unknown error while creating pkg dir %s: %s", pkgPath, err)
			return "", err
		}

		// rename failed with ENOTEMPTY, that means that another client wrote this
		// directory.
		description, err := d.readDescription(ctx, pkgPath)
		if err != nil {
			logging.Warningf(ctx, "Skipping %q: %s", pkgPath, err)
			continue
		}
		if description.PackageName == pkg && description.Subdir == subdir {
			return pkgPath, nil
		}
	}

	logging.Errorf(ctx, "Unable to find valid index for package %q in %s!", pkg, abs)
	return "", err
}

type byLenThenAlpha []string

func (b byLenThenAlpha) Len() int      { return len(b) }
func (b byLenThenAlpha) Swap(i, j int) { b[i], b[j] = b[j], b[i] }
func (b byLenThenAlpha) Less(i, j int) bool {
	return sortby.Chain{
		func(i, j int) bool { return len(b[i]) < len(b[j]) },
		func(i, j int) bool { return b[i] < b[j] },
	}.Use(i, j)
}

// resolveValidPackageDirs scans the .cipd/pkgs dir and returns:
//   * a numeric set of all number-style directories seen.
//   * a map of Description (e.g. subdir + pkgname) to the correct pkg folder
//
// This also will delete (EnsureDirectoryGone) any folders or files in the pkgs
// directory which are:
//   * invalid (contain no description.json and no current instance symlink)
//   * duplicate (where multiple directories contain the same description.json)
//
// Duplicate detection always prefers the folder with the shortest path name
// that sorts alphabetically earlier.
func (d *deployerImpl) resolveValidPackageDirs(ctx context.Context, pkgsAbsDir string) (numbered numSet, all map[Description]string) {
	files, err := ioutil.ReadDir(pkgsAbsDir)
	if err != nil && !os.IsNotExist(err) {
		logging.Errorf(ctx, "Can't read packages dir %q: %s", pkgsAbsDir, err)
		return
	}

	allWithDups := map[Description][]string{}

	for _, f := range files {
		fullPkgPath := filepath.Join(pkgsAbsDir, f.Name())
		description, err := d.readDescription(ctx, fullPkgPath)
		if description == nil || err != nil {
			if err == nil {
				err = fmt.Errorf("missing description.json and current instance")
			}
			logging.Warningf(ctx, "removing junk directory: %q (%s)", fullPkgPath, err)
			if err := d.fs.EnsureDirectoryGone(ctx, fullPkgPath); err != nil {
				logging.Warningf(ctx, "while removing junk directory: %q (%s)", fullPkgPath, err)
			}
			continue
		}
		allWithDups[*description] = append(allWithDups[*description], fullPkgPath)
	}

	all = make(map[Description]string, len(allWithDups))
	for desc, possibilities := range allWithDups {
		sort.Sort(byLenThenAlpha(possibilities))

		// keep track of all non-deleted numeric children of .cipd/pkgs
		if n, err := strconv.Atoi(filepath.Base(possibilities[0])); err == nil {
			numbered.addNum(n)
		}

		all[desc] = possibilities[0]

		if len(possibilities) == 1 {
			continue
		}
		for _, extra := range possibilities[1:] {
			logging.Warningf(ctx, "removing duplicate directory: %q", extra)
			if err := d.fs.EnsureDirectoryGone(ctx, extra); err != nil {
				logging.Warningf(ctx, "while removing duplicate directory: %q (%s)", extra, err)
			}
		}
	}

	return
}

// getCurrentInstanceID returns instance ID of currently installed instance
// given a path to a package directory (.cipd/pkgs/<name>).
//
// It returns ("", nil) if no package is installed there.
func (d *deployerImpl) getCurrentInstanceID(packageDir string) (string, error) {
	var current string
	var err error
	if runtime.GOOS == "windows" {
		var bytes []byte
		bytes, err = ioutil.ReadFile(filepath.Join(packageDir, currentTxt))
		if err == nil {
			current = strings.TrimSpace(string(bytes))
		}
	} else {
		current, err = os.Readlink(filepath.Join(packageDir, currentSymlink))
	}
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	if err = common.ValidateInstanceID(current); err != nil {
		return "", fmt.Errorf(
			"pointer to currently installed instance doesn't look like a valid instance id: %s", err)
	}
	return current, nil
}

// setCurrentInstanceID changes a pointer to currently installed instance ID.
//
// It takes a path to a package directory (.cipd/pkgs/<name>) as input.
func (d *deployerImpl) setCurrentInstanceID(ctx context.Context, packageDir, instanceID string) error {
	if err := common.ValidateInstanceID(instanceID); err != nil {
		return err
	}
	if runtime.GOOS == "windows" {
		return EnsureFile(
			ctx, d.fs, filepath.Join(packageDir, currentTxt),
			strings.NewReader(instanceID))
	}
	return d.fs.EnsureSymlink(ctx, filepath.Join(packageDir, currentSymlink), instanceID)
}

// readDescription reads the package description.json given a path to a package
// directory.
//
// As a backwards-compatibility measure, it will also upgrade CIPD < 1.4 folders
// to contain a description.json. Previous to 1.4, package folders only had
// instance subfolders, and the current instances' manifest was used to
// determine the package name. Versions prior to 1.4 also installed all packages
// at the base (subdir ""), hence the implied subdir location here.
//
// Returns (nil, nil) if no description.json exists and there are no instance
// folders present.
func (d *deployerImpl) readDescription(ctx context.Context, pkgDir string) (desc *Description, err error) {
	descriptionPath := filepath.Join(pkgDir, descriptionName)
	r, err := os.Open(descriptionPath)
	switch {
	case os.IsNotExist(err):
		// try fixup
		break
	case err == nil:
		defer r.Close()
		return readDescription(r)
	default:
		return
	}

	// see if this is a pre 1.4 directory
	currentID, err := d.getCurrentInstanceID(pkgDir)
	if err != nil {
		return
	}

	if currentID == "" {
		logging.Warningf(ctx, "No current instance id in %s", pkgDir)
		err = nil
		return
	}

	manifest, err := d.readManifest(ctx, filepath.Join(pkgDir, currentID))
	if err != nil {
		return
	}

	desc = &Description{
		PackageName: manifest.PackageName,
	}
	// To handle the case where some other user owns these directories, all errors
	// from here to the end are treated as warnings.
	err = d.fs.EnsureFile(ctx, descriptionPath, func(f *os.File) error {
		return writeDescription(desc, f)
	})
	if err != nil {
		logging.Warningf(ctx, "Unable to create description.json: %s", err)
		err = nil
	}
	return
}

// readManifest reads package manifest given a path to a package instance
// (.cipd/pkgs/<name>/<instance id>).
func (d *deployerImpl) readManifest(ctx context.Context, instanceDir string) (Manifest, error) {
	manifestPath := filepath.Join(instanceDir, filepath.FromSlash(manifestName))
	r, err := os.Open(manifestPath)
	if err != nil {
		return Manifest{}, err
	}
	defer r.Close()
	manifest, err := readManifest(r)
	if err != nil {
		return Manifest{}, err
	}
	// Older packages do not have Files section in the manifest, so reconstruct it
	// from actual files on disk.
	if len(manifest.Files) == 0 {
		if manifest.Files, err = scanPackageDir(ctx, instanceDir); err != nil {
			return Manifest{}, err
		}
	}
	return manifest, nil
}

// addToSiteRoot moves or symlinks files into the site root directory (depending
// on passed installMode).
func (d *deployerImpl) addToSiteRoot(ctx context.Context, subdir string, files []FileInfo, installMode InstallMode, pkgDir, srcDir string) error {
	for _, f := range files {
		// e.g. bin/tool
		relPath := filepath.FromSlash(f.Name)
		// e.g. <base>/<subdir>/bin/tool
		destAbs, err := d.fs.RootRelToAbs(filepath.Join(subdir, relPath))
		if err != nil {
			logging.Warningf(ctx, "Invalid relative path %q: %s", relPath, err)
			return err
		}
		if installMode == InstallModeSymlink {
			// e.g. <base>/.cipd/pkgs/name/_current/bin/tool
			targetAbs := filepath.Join(pkgDir, currentSymlink, relPath)
			// e.g. ../.cipd/pkgs/name/_current/bin/tool
			// has more `../` depending on subdir
			targetRel, err := filepath.Rel(filepath.Dir(destAbs), targetAbs)
			if err != nil {
				logging.Warningf(
					ctx, "Can't get relative path from %s to %s",
					filepath.Dir(destAbs), targetAbs)
				return err
			}
			if err = d.fs.EnsureSymlink(ctx, destAbs, targetRel); err != nil {
				logging.Warningf(ctx, "Failed to create symlink for %s", relPath)
				return err
			}
		} else if installMode == InstallModeCopy {
			// E.g. <base>/.cipd/pkgs/name/<id>/bin/tool.
			srcAbs := filepath.Join(srcDir, relPath)
			if err := d.fs.Replace(ctx, srcAbs, destAbs); err != nil {
				logging.Warningf(ctx, "Failed to move %s to %s: %s", srcAbs, destAbs, err)
				return err
			}
		} else {
			// Should not happen. ValidateInstallMode checks this.
			return fmt.Errorf("impossible state")
		}
	}
	return nil
}

// removeFromSiteRoot deletes files from the site root directory.
//
// Best effort. Logs errors and carries on.
func (d *deployerImpl) removeFromSiteRoot(ctx context.Context, subdir string, files []FileInfo) {
	dirsToCleanup := stringset.New(0)

	for _, f := range files {
		absPath, err := d.fs.RootRelToAbs(filepath.Join(subdir, filepath.FromSlash(f.Name)))
		if err != nil {
			logging.Warningf(ctx, "Refusing to remove %q: %s", f.Name, err)
			continue
		}
		if err := d.fs.EnsureFileGone(ctx, absPath); err != nil {
			logging.Warningf(ctx, "Failed to remove a file from the site root: %s", err)
		} else {
			dirsToCleanup.Add(filepath.Dir(absPath))
		}
	}

	if dirsToCleanup.Len() != 0 {
		subdirAbs, err := d.fs.RootRelToAbs(subdir)
		if err != nil {
			logging.Warningf(ctx, "Can't resolve relative %q to absolute path: %s", subdir, err)
		} else {
			removeEmptyTrees(ctx, subdirAbs, dirsToCleanup)
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
// Utility functions.

// checkInstallMode validates the install mode and picks the correct default
// if no install mode is given.
func checkInstallMode(im InstallMode) (InstallMode, error) {
	switch {
	case runtime.GOOS == "windows":
		return InstallModeCopy, nil // Windows supports only 'copy' mode
	case im == "":
		return InstallModeSymlink, nil // default on other platforms
	}
	if err := ValidateInstallMode(im); err != nil {
		return "", err
	}
	return im, nil
}

// scanPackageDir finds a set of regular files (and symlinks) in a package
// instance directory and returns them as FileInfo structs (with slash-separated
// paths relative to dir directory). Skips package service directories (.cipdpkg
// and .cipd) since they contain package deployer gut files, not something that
// needs to be deployed.
func scanPackageDir(ctx context.Context, dir string) ([]FileInfo, error) {
	out := []FileInfo{}
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		if rel == packageServiceDir || rel == SiteServiceDir {
			return filepath.SkipDir
		}
		if info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 {
			symlink := ""
			ok := true
			if info.Mode()&os.ModeSymlink != 0 {
				symlink, err = os.Readlink(path)
				if err != nil {
					logging.Warningf(ctx, "Can't readlink %q, skipping: %s", path, err)
					ok = false
				}
			}
			if ok {
				out = append(out, FileInfo{
					Name:       filepath.ToSlash(rel),
					Size:       uint64(info.Size()),
					Executable: (info.Mode().Perm() & 0111) != 0,
					Symlink:    symlink,
				})
			}
		}
		return nil
	})
	return out, err
}

// removeEmptyTrees recursively removes empty directory subtrees after some
// files have been removed.
//
// It tries to avoid enumerating entire directory tree and instead recurses
// only into directories with potentially empty subtrees. They are indicated by
// 'empty' set with absolute paths to directories that had files removed from
// them (so they MAY be empty now, but not necessarily).
//
// All paths are absolute, using native separators.
//
// Best effort, logs errors.
func removeEmptyTrees(ctx context.Context, root string, empty stringset.Set) {
	// If directory 'A/B/C' has potentially empty subtree, then so do 'A/B' and
	// 'A' and '.'. Expand 'empty' set according to these rules. Note that 'root'
	// itself is always is this set.
	verboseEmpty := stringset.New(empty.Len())
	verboseEmpty.Add(root)
	empty.Iter(func(dir string) bool {
		rel, err := filepath.Rel(root, dir)
		if err != nil {
			// Note: this should never really happen, since there are checks outside
			// of this function.
			logging.Warningf(ctx, "Can't compute %q relative to %q - %s", dir, root, err)
			return true
		}

		// Here 'rel' has form 'A/B/C' or is '.' (but this is already handled).
		if rel != "." {
			path := root
			for _, chunk := range strings.Split(rel, string(filepath.Separator)) {
				path = filepath.Join(path, chunk)
				verboseEmpty.Add(path)
			}
		}

		return true
	})

	// Now we recursively walk through the root subtree, skipping trees we know
	// can't be empty.
	_, err := removeEmptyTree(root, func(candidate string) (shouldCheck bool) {
		return verboseEmpty.Has(candidate)
	})
	if err != nil {
		logging.Warningf(ctx, "Failed to cleanup empty directories under %q - %s", root, err)
	}
}

// removeEmptyTree recursively removes an empty directory tree.
//
// 'path' must point to a directory (not a regular file, not a symlink).
//
// Returns true if deleted 'path' along with its (empty) subtree. Stops on first
// encountered error.
func removeEmptyTree(path string, shouldCheck func(string) bool) (deleted bool, err error) {
	if !shouldCheck(path) {
		return false, nil
	}

	// 'Remove' will delete the directory if it is already empty.
	if err := os.Remove(path); err == nil || os.IsNotExist(err) {
		return true, nil
	}

	// Otherwise need to recurse into it.
	fd, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return true, nil // someone deleted it already, this is OK
		}
		return false, err
	}

	closed := false
	defer func() {
		if !closed {
			fd.Close()
		}
	}()

	total := 0
	removed := 0
	for {
		infos, err := fd.Readdir(100)
		if err == io.EOF || len(infos) == 0 {
			break
		}
		if err != nil {
			return false, err
		}
		total += len(infos)
		for _, info := range infos {
			if info.IsDir() {
				abs := filepath.Join(path, info.Name())
				switch rmed, err := removeEmptyTree(abs, shouldCheck); {
				case err != nil:
					return false, err
				case rmed:
					removed++
				}
			}
		}
	}

	// Close directory, because windows won't remove opened directory.
	fd.Close()
	closed = true

	// The directory is definitely not empty, since we skipped some stuff.
	if total != removed {
		return false, nil
	}

	// The directory is most likely empty now, unless someone concurrently put
	// files there. Unfortunately it is not trivial to detect this specific
	// condition in a cross-platform way. So assume Remove() errors (other than
	// IsNotExit) are due to that.
	err = os.Remove(path)
	return err == nil || os.IsNotExist(err), nil
}

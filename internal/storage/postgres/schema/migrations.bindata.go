// Code generated by go-bindata.
// sources:
// migrations/1565894792_users.down.sql
// migrations/1565894792_users.up.sql
// migrations/1569158364_roles.down.sql
// migrations/1569158364_roles.up.sql
// migrations/1570866542_articles.down.sql
// migrations/1570866542_articles.up.sql
// migrations/1571140600_categories.down.sql
// migrations/1571140600_categories.up.sql
// DO NOT EDIT!

package schema

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var __1565894792_usersDownSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x72\x09\xf2\x0f\x50\x08\x71\x74\xf2\x71\x55\xf0\x74\x53\x70\x8d\xf0\x0c\x0e\x09\x56\x28\x2d\x4e\x2d\x2a\xb6\xe6\x02\x04\x00\x00\xff\xff\x2c\x02\x3d\xa7\x1c\x00\x00\x00")

func _1565894792_usersDownSqlBytes() ([]byte, error) {
	return bindataRead(
		__1565894792_usersDownSql,
		"1565894792_users.down.sql",
	)
}

func _1565894792_usersDownSql() (*asset, error) {
	bytes, err := _1565894792_usersDownSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "1565894792_users.down.sql", size: 28, mode: os.FileMode(420), modTime: time.Unix(1566116474, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var __1565894792_usersUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x9c\x90\x41\x4b\x03\x31\x10\x46\xcf\x93\x5f\x31\xc7\x6e\x29\x54\x04\xf1\xd0\x53\xdc\x4e\x31\x98\x4d\x6b\x36\x91\xf6\x14\x82\x09\x34\xd0\xb5\xcb\x66\x45\x7f\xbe\x74\x71\x6b\x11\xa1\xe0\xf5\xe3\xbd\x61\x78\xa5\x26\x6e\x08\x0d\x7f\x90\x84\x62\x85\x6a\x6d\x90\xb6\xa2\x36\x35\xbe\xe7\xd8\x65\x9c\x30\x48\x01\xa0\x26\x2d\xb8\xc4\x8d\x16\x15\xd7\x3b\x7c\xa2\xdd\x8c\x41\x77\x3c\x44\x97\x02\x0a\x65\x06\x51\x59\x29\x67\x0c\x4e\xe2\x9b\x6f\x22\xbc\x70\x5d\x3e\x72\x8d\x93\xbb\x9b\x02\xad\x12\xcf\x96\x2e\xb9\xd8\xf8\x74\x80\x6b\x54\xeb\x73\xfe\x38\x76\xc1\xed\x7d\xde\xff\xc0\xf7\xb7\xc5\x05\xc5\x60\x3e\xc5\x3e\x35\x31\xf7\xbe\x69\x71\x3a\x67\xf0\xda\x45\xdf\xc7\xe0\x7c\x0f\x46\x54\x54\x1b\x5e\x6d\xce\x06\x2e\x69\xc5\xad\x34\x58\x5a\xad\x49\x19\x77\x46\x4e\xff\xb7\xe1\x3f\x26\x2b\x16\x8c\x7d\xf7\x14\x6a\x49\xdb\xbf\x7a\xba\x31\x8e\x4b\xe1\x13\xd7\x6a\xac\x3c\xce\xc5\xe2\xfa\x89\xa1\xdb\x2f\x7f\xd8\x8a\x05\xb2\xaf\x00\x00\x00\xff\xff\x10\xf2\x3e\xc1\xd1\x01\x00\x00")

func _1565894792_usersUpSqlBytes() ([]byte, error) {
	return bindataRead(
		__1565894792_usersUpSql,
		"1565894792_users.up.sql",
	)
}

func _1565894792_usersUpSql() (*asset, error) {
	bytes, err := _1565894792_usersUpSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "1565894792_users.up.sql", size: 465, mode: os.FileMode(420), modTime: time.Unix(1569210888, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var __1569158364_rolesDownSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x72\x09\xf2\x0f\x50\x08\x71\x74\xf2\x71\x55\xf0\x74\x53\x70\x8d\xf0\x0c\x0e\x09\x56\x28\xca\xcf\x49\x2d\xb6\x06\x04\x00\x00\xff\xff\x8c\x8d\x3a\x89\x1b\x00\x00\x00")

func _1569158364_rolesDownSqlBytes() ([]byte, error) {
	return bindataRead(
		__1569158364_rolesDownSql,
		"1569158364_roles.down.sql",
	)
}

func _1569158364_rolesDownSql() (*asset, error) {
	bytes, err := _1569158364_rolesDownSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "1569158364_roles.down.sql", size: 27, mode: os.FileMode(420), modTime: time.Unix(1569163573, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var __1569158364_rolesUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x9c\xcc\xb1\x4e\x86\x30\x10\xc0\xf1\xf9\xfa\x14\x37\x7e\x1f\x21\xc1\xc5\xc9\xe9\xc4\x23\x36\x96\x8a\xd7\xab\x91\x89\x34\xd2\x81\x44\x94\x40\x7d\x7f\x13\x07\x1e\xe0\xdb\xff\xbf\x7f\x2b\x4c\xca\xa8\xf4\xe8\x18\x6d\x87\xfe\x55\x91\x3f\x6c\xd0\x80\xfb\xcf\x57\x3e\xf0\x62\x60\x99\x01\x02\x8b\x25\x87\x83\xd8\x9e\x64\xc4\x17\x1e\x6b\x03\xdf\x69\xcd\xf0\x4e\xd2\x3e\x93\xe0\xe5\xfe\xee\x8a\xd1\xdb\xb7\xc8\xff\x17\x1f\x9d\xab\x8d\x81\xa6\xc2\xb2\xac\xf9\x28\x69\xdd\xb0\x6a\x0c\x7c\xee\x39\x95\x3c\x4f\xa9\x80\xda\x9e\x83\x52\x3f\x9c\x02\x9f\xb8\xa3\xe8\x14\xdb\x28\xc2\x5e\xa7\x33\xa9\x0d\xfc\x6e\xf3\x2d\xd2\x5c\x1f\xfe\x02\x00\x00\xff\xff\x3b\x99\x38\xe3\xe7\x00\x00\x00")

func _1569158364_rolesUpSqlBytes() ([]byte, error) {
	return bindataRead(
		__1569158364_rolesUpSql,
		"1569158364_roles.up.sql",
	)
}

func _1569158364_rolesUpSql() (*asset, error) {
	bytes, err := _1569158364_rolesUpSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "1569158364_roles.up.sql", size: 231, mode: os.FileMode(420), modTime: time.Unix(1569216254, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var __1570866542_articlesDownSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x72\x09\xf2\x0f\x50\x08\x71\x74\xf2\x71\x55\xf0\x74\x53\x70\x8d\xf0\x0c\x0e\x09\x56\x48\x2c\x2a\xc9\x4c\xce\x49\x2d\xb6\xe6\x02\x04\x00\x00\xff\xff\xc8\xa2\xce\x28\x1f\x00\x00\x00")

func _1570866542_articlesDownSqlBytes() ([]byte, error) {
	return bindataRead(
		__1570866542_articlesDownSql,
		"1570866542_articles.down.sql",
	)
}

func _1570866542_articlesDownSql() (*asset, error) {
	bytes, err := _1570866542_articlesDownSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "1570866542_articles.down.sql", size: 31, mode: os.FileMode(420), modTime: time.Unix(1571404816, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var __1570866542_articlesUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x9c\xce\x41\x4b\xc3\x40\x10\x05\xe0\xf3\xec\xaf\x98\x63\x5b\x0a\x05\xa1\x27\x4f\x63\x9c\xe2\xe2\x26\x96\xc9\x44\xda\x53\x58\xb3\x8b\x2c\xb4\xb4\x24\xe3\xa1\xff\x5e\xf4\x50\xcc\xd5\xf3\x7b\xdf\xe3\x55\xc2\xa4\x8c\x4a\x4f\x81\xd1\xef\xb0\x79\x53\xe4\x83\x6f\xb5\xc5\x38\x5a\x19\x4e\x79\xc2\x85\x83\x92\xa0\x65\xf1\x14\x70\x2f\xbe\x26\x39\xe2\x2b\x1f\xd7\x0e\xbe\xa6\x3c\xf6\x25\x81\x6f\xf4\x97\x36\x5d\x08\x6b\x07\x56\xec\x94\xe1\x9d\xa4\x7a\x21\xc1\xc5\xc3\x76\xbb\xfc\x1b\x7f\x5c\xd2\x0d\x94\x0f\x33\x33\x44\xcb\x9f\x97\xf1\xd6\x97\x84\xf3\x39\x07\x9b\x15\x5a\x39\xe7\xc9\xe2\xf9\x8a\xab\x8d\x83\x61\xcc\xd1\x72\xea\xa3\x81\xfa\x9a\x5b\xa5\x7a\x7f\x17\xf8\xcc\x3b\xea\x82\x62\xd5\x89\x70\xa3\xfd\xbd\xf2\xf3\xf8\x9a\xfe\x23\xdd\xf2\xd1\xb9\xef\x00\x00\x00\xff\xff\x9f\x48\x4d\x56\x2d\x01\x00\x00")

func _1570866542_articlesUpSqlBytes() ([]byte, error) {
	return bindataRead(
		__1570866542_articlesUpSql,
		"1570866542_articles.up.sql",
	)
}

func _1570866542_articlesUpSql() (*asset, error) {
	bytes, err := _1570866542_articlesUpSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "1570866542_articles.up.sql", size: 301, mode: os.FileMode(420), modTime: time.Unix(1571404800, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var __1571140600_categoriesDownSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x72\x09\xf2\x0f\x50\x08\x71\x74\xf2\x71\x55\xf0\x74\x53\x70\x8d\xf0\x0c\x0e\x09\x56\x48\x4e\x2c\x49\x4d\xcf\x2f\xca\x4c\x2d\xb6\xe6\x02\x04\x00\x00\xff\xff\xfd\xee\x8d\x6a\x21\x00\x00\x00")

func _1571140600_categoriesDownSqlBytes() ([]byte, error) {
	return bindataRead(
		__1571140600_categoriesDownSql,
		"1571140600_categories.down.sql",
	)
}

func _1571140600_categoriesDownSql() (*asset, error) {
	bytes, err := _1571140600_categoriesDownSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "1571140600_categories.down.sql", size: 33, mode: os.FileMode(420), modTime: time.Unix(1571144203, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var __1571140600_categoriesUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x9c\xcc\xb1\x4e\xc3\x30\x10\x80\xe1\xf9\xee\x29\x6e\x6c\xab\x48\x65\x61\x62\x3a\xc2\x55\x58\x38\xa6\x9c\xcf\x88\x4e\x91\x95\x58\x28\x43\x20\x4a\xcc\xfb\x23\x31\xe4\x01\xd8\xff\xef\x6f\x55\xd8\x84\x8c\x1f\xbd\x90\xbb\x50\x78\x35\x92\x0f\x17\x2d\xd2\x90\x6b\xf9\xfc\x5e\xa7\xb2\xd1\x01\x61\x1a\x21\x8a\x3a\xf6\x74\x55\xd7\xb1\xde\xe8\x45\x6e\x0d\xc2\x57\x9e\x0b\xbc\xb3\xb6\xcf\xac\x74\xb8\xbf\x3b\x52\x0a\xee\x2d\xc9\xdf\x29\x24\xef\x1b\x44\x38\x9f\xa8\x4e\x73\xd9\x6a\x9e\x17\x3a\x9d\x11\x86\xb5\xe4\x5a\xc6\x3e\x57\x30\xd7\x49\x34\xee\xae\xbb\xa0\x27\xb9\x70\xf2\x46\x6d\x52\x95\x60\xfd\x9e\x34\x08\x3f\xcb\xf8\x1f\x89\xc7\x07\xfc\x0d\x00\x00\xff\xff\xba\x6b\x98\x19\xec\x00\x00\x00")

func _1571140600_categoriesUpSqlBytes() ([]byte, error) {
	return bindataRead(
		__1571140600_categoriesUpSql,
		"1571140600_categories.up.sql",
	)
}

func _1571140600_categoriesUpSql() (*asset, error) {
	bytes, err := _1571140600_categoriesUpSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "1571140600_categories.up.sql", size: 236, mode: os.FileMode(420), modTime: time.Unix(1571221426, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"1565894792_users.down.sql": _1565894792_usersDownSql,
	"1565894792_users.up.sql": _1565894792_usersUpSql,
	"1569158364_roles.down.sql": _1569158364_rolesDownSql,
	"1569158364_roles.up.sql": _1569158364_rolesUpSql,
	"1570866542_articles.down.sql": _1570866542_articlesDownSql,
	"1570866542_articles.up.sql": _1570866542_articlesUpSql,
	"1571140600_categories.down.sql": _1571140600_categoriesDownSql,
	"1571140600_categories.up.sql": _1571140600_categoriesUpSql,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}
var _bintree = &bintree{nil, map[string]*bintree{
	"1565894792_users.down.sql": &bintree{_1565894792_usersDownSql, map[string]*bintree{}},
	"1565894792_users.up.sql": &bintree{_1565894792_usersUpSql, map[string]*bintree{}},
	"1569158364_roles.down.sql": &bintree{_1569158364_rolesDownSql, map[string]*bintree{}},
	"1569158364_roles.up.sql": &bintree{_1569158364_rolesUpSql, map[string]*bintree{}},
	"1570866542_articles.down.sql": &bintree{_1570866542_articlesDownSql, map[string]*bintree{}},
	"1570866542_articles.up.sql": &bintree{_1570866542_articlesUpSql, map[string]*bintree{}},
	"1571140600_categories.down.sql": &bintree{_1571140600_categoriesDownSql, map[string]*bintree{}},
	"1571140600_categories.up.sql": &bintree{_1571140600_categoriesUpSql, map[string]*bintree{}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}


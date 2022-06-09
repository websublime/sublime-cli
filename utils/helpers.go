/*
Copyright © 2022 Websublime.dev organization@websublime.dev

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type PackageJson struct {
	Name            string            `json:"name"`
	Version         string            `json:"version"`
	Description     string            `json:"description"`
	Main            string            `json:"main"`
	Typings         string            `json:"typings"`
	Module          string            `json:"module"`
	Scripts         map[string]string `json:"scripts"`
	Keywords        []string          `json:"keywords"`
	Author          string            `json:"author"`
	License         string            `json:"license"`
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

func PathWalk(root string) ([]string, error) {
	var files []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {

		if !info.IsDir() && !strings.HasPrefix(info.Name(), ".") {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking path: %v", err)
	}

	return files, nil
}

func GetMimeType(extension string) string {
	mimes := map[string]string{
		"none":    "application/octet-stream",
		"xpm":     "image/x-xpixmap",
		"7z":      "application/x-7z-compressed",
		"zip":     "application/zip",
		"xlsx":    "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		"docx":    "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		"pptx":    "application/vnd.openxmlformats-officedocument.presentationml.presentation",
		"epub":    "application/epub+zip",
		"jar":     "application/jar",
		"odt":     "application/vnd.oasis.opendocument.text",
		"ott":     "application/vnd.oasis.opendocument.text-template",
		"ods":     "application/vnd.oasis.opendocument.spreadsheet",
		"ots":     "application/vnd.oasis.opendocument.spreadsheet-template",
		"odp":     "application/vnd.oasis.opendocument.presentation",
		"otp":     "application/vnd.oasis.opendocument.presentation-template",
		"odg":     "application/vnd.oasis.opendocument.graphics",
		"otg":     "application/vnd.oasis.opendocument.graphics-template",
		"odf":     "application/vnd.oasis.opendocument.formula",
		"odc":     "application/vnd.oasis.opendocument.chart",
		"sxc":     "application/vnd.sun.xml.calc",
		"pdf":     "application/pdf",
		"fdf":     "application/vnd.fdf",
		"msi":     "application/x-ms-installer",
		"aaf":     "application/octet-stream",
		"msg":     "application/vnd.ms-outlook",
		"xls":     "application/vnd.ms-excel",
		"pub":     "application/vnd.ms-publisher",
		"ppt":     "application/vnd.ms-powerpoint",
		"doc":     "application/msword",
		"ps":      "application/postscript",
		"psd":     "image/vnd.adobe.photoshop",
		"p7s":     "application/pkcs7-signature",
		"ogg":     "application/ogg",
		"oga":     "audio/ogg",
		"ogv":     "video/ogg",
		"png":     "image/png",
		"apng":    "image/vnd.mozilla.apng",
		"jpg":     "image/jpeg",
		"jxl":     "image/jxl",
		"jp2":     "image/jp2",
		"jpf":     "image/jpx",
		"jpm":     "image/jpm",
		"gif":     "image/gif",
		"webp":    "image/webp",
		"exe":     "application/vnd.microsoft.portable-executable",
		"so":      "application/x-sharedlib",
		"a":       "application/x-archive",
		"deb":     "application/vnd.debian.binary-package",
		"tar":     "application/x-tar",
		"xar":     "application/x-xar",
		"bz2":     "application/x-bzip2",
		"fits":    "application/fits",
		"tiff":    "image/tiff",
		"bmp":     "image/bmp",
		"ico":     "image/x-icon",
		"mp3":     "audio/mpeg",
		"flac":    "audio/flac",
		"midi":    "audio/midi",
		"ape":     "audio/ape",
		"mpc":     "audio/musepack",
		"amr":     "audio/amr",
		"wav":     "audio/wav",
		"aiff":    "audio/aiff",
		"au":      "audio/basic",
		"mpeg":    "video/mpeg",
		"mov":     "video/quicktime",
		"mqv":     "video/quicktime",
		"mp4":     "video/mp4",
		"webm":    "video/webm",
		"3gp":     "video/3gpp",
		"3g2":     "video/3gpp2",
		"avi":     "video/x-msvideo",
		"flv":     "video/x-flv",
		"mkv":     "video/x-matroska",
		"asf":     "video/x-ms-asf",
		"aac":     "audio/aac",
		"voc":     "audio/x-unknown",
		"m4a":     "audio/x-m4a",
		"m3u":     "application/vnd.apple.mpegurl",
		"m4v":     "video/x-m4v",
		"rmvb":    "application/vnd.rn-realmedia-vbr",
		"gz":      "application/gzip",
		"class":   "application/x-java-applet",
		"swf":     "application/x-shockwave-flash",
		"crx":     "application/x-chrome-extension",
		"ttf":     "font/ttf",
		"woff":    "font/woff",
		"woff2":   "font/woff2",
		"otf":     "font/otf",
		"ttc":     "font/collection",
		"eot":     "application/vnd.ms-fontobject",
		"wasm":    "application/wasm",
		"shx":     "application/octet-stream",
		"shp":     "application/octet-stream",
		"dbf":     "application/x-dbf",
		"dcm":     "application/dicom",
		"rar":     "application/x-rar-compressed",
		"djvu":    "image/vnd.djvu",
		"mobi":    "application/x-mobipocket-ebook",
		"lit":     "application/x-ms-reader",
		"bpg":     "image/bpg",
		"sqlite":  "application/vnd.sqlite3",
		"dwg":     "image/vnd.dwg",
		"nes":     "application/vnd.nintendo.snes.rom",
		"lnk":     "application/x-ms-shortcut",
		"macho":   "application/x-mach-binary",
		"qcp":     "audio/qcelp",
		"icns":    "image/x-icns",
		"heic":    "image/heic",
		"heif":    "image/heif",
		"hdr":     "image/vnd.radiance",
		"mrc":     "application/marc",
		"mdb":     "application/x-msaccess",
		"accdb":   "application/x-msaccess",
		"zst":     "application/zstd",
		"cab":     "application/vnd.ms-cab-compressed",
		"rpm":     "application/x-rpm",
		"xz":      "application/x-xz",
		"lz":      "application/lzip",
		"torrent": "application/x-bittorrent",
		"cpio":    "application/x-cpio",
		"xcf":     "image/x-xcf",
		"pat":     "image/x-gimp-pat",
		"gbr":     "image/x-gimp-gbr",
		"glb":     "model/gltf-binary",
		"avif":    "image/avif",
		"txt":     "text/plain",
		"html":    "text/html",
		"svg":     "image/svg+xml",
		"xml":     "text/xml",
		"rss":     "application/rss+xml",
		"atom":    "application/atom+xml",
		"x3d":     "model/x3d+xml",
		"kml":     "application/vnd.google-earth.kml+xml",
		"xlf":     "application/x-xliff+xml",
		"dae":     "model/vnd.collada+xml",
		"gml":     "application/gml+xml",
		"gpx":     "application/gpx+xml",
		"tcx":     "application/vnd.garmin.tcx+xml",
		"amf":     "application/x-amf",
		"3mf":     "application/vnd.ms-package.3dmanufacturing-3dmodel+xml",
		"xfdf":    "application/vnd.adobe.xfdf",
		"owl":     "application/owl+xml",
		"php":     "text/x-php",
		"js":      "application/javascript",
		"mjs":     "application/javascript",
		"lua":     "text/x-lua",
		"pl":      "text/x-perl",
		"py":      "application/x-python",
		"json":    "application/json",
		"map":     "application/json",
		"geojson": "application/geo+json",
		"har":     "application/json",
		"ndjson":  "application/x-ndjson",
		"rtf":     "text/rtf",
		"srt":     "application/x-subrip",
		"tcl":     "text/x-tcl",
		"csv":     "text/csv",
		"tsv":     "text/tab-separated-values",
		"vcf":     "text/vcard",
		"ics":     "text/calendar",
		"warc":    "application/warc",
		"vtt":     "text/vtt",
		"md":      "text/markdown",
		"css":     "text/css",
	}

	mimetype := mimes["none"]
	mime := mimes[extension]

	if mime != "" {
		mimetype = mime
	}

	return fmt.Sprintf("%s; charset=utf-8", mimetype)
}

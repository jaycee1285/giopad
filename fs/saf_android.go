// SPDX-License-Identifier: Unlicense OR MIT

//go:build android

package fs

/*
#cgo LDFLAGS: -landroid

#include <jni.h>
#include <stdlib.h>
*/
import "C"
import (
	"errors"
	"strings"

	"gioui.org/app"
	"git.wow.st/gmp/jni"
)

var (
	safClass      jni.Class
	listDirID     jni.MethodID
	listSubDirID  jni.MethodID
	readFileID    jni.MethodID
	writeFileID   jni.MethodID
	getTreeNameID jni.MethodID
	initialized   bool
)

func initSAF(env jni.Env) error {
	if initialized {
		return nil
	}

	cls, err := jni.LoadClass(env, jni.ClassLoaderFor(env, jni.Object(app.AppContext())), "org/gioui/x/explorer/explorer_android")
	if err != nil {
		return err
	}

	safClass = jni.Class(jni.NewGlobalRef(env, jni.Object(cls)))
	listDirID = jni.GetStaticMethodID(env, safClass, "listDir", "(Landroid/content/Context;Ljava/lang/String;)Ljava/lang/String;")
	listSubDirID = jni.GetStaticMethodID(env, safClass, "listSubDir", "(Landroid/content/Context;Ljava/lang/String;Ljava/lang/String;)Ljava/lang/String;")
	readFileID = jni.GetStaticMethodID(env, safClass, "readFile", "(Landroid/content/Context;Ljava/lang/String;)[B")
	writeFileID = jni.GetStaticMethodID(env, safClass, "writeFile", "(Landroid/content/Context;Ljava/lang/String;[B)Z")
	getTreeNameID = jni.GetStaticMethodID(env, safClass, "getTreeName", "(Landroid/content/Context;Ljava/lang/String;)Ljava/lang/String;")

	initialized = true
	return nil
}

// SAFEntry represents a file or directory in SAF
type SAFEntry struct {
	Name  string
	URI   string
	IsDir bool
}

// ListSAFDir lists contents of a SAF tree URI
func ListSAFDir(treeURI string) ([]SAFEntry, error) {
	var entries []SAFEntry

	err := jni.Do(jni.JVMFor(app.JavaVM()), func(env jni.Env) error {
		if err := initSAF(env); err != nil {
			return err
		}

		ctx := jni.Object(app.AppContext())
		uriStr := jni.JavaString(env, treeURI)

		result, err := jni.CallStaticObjectMethod(env, safClass, listDirID, jni.Value(ctx), jni.Value(uriStr))
		if err != nil {
			return err
		}

		resultStr := jni.GoString(env, jni.String(result))
		entries = parseEntries(resultStr)
		return nil
	})

	return entries, err
}

// ListSAFSubDir lists contents of a subdirectory within a SAF tree
func ListSAFSubDir(treeURI, docURI string) ([]SAFEntry, error) {
	var entries []SAFEntry

	err := jni.Do(jni.JVMFor(app.JavaVM()), func(env jni.Env) error {
		if err := initSAF(env); err != nil {
			return err
		}

		ctx := jni.Object(app.AppContext())
		treeStr := jni.JavaString(env, treeURI)
		docStr := jni.JavaString(env, docURI)

		result, err := jni.CallStaticObjectMethod(env, safClass, listSubDirID, jni.Value(ctx), jni.Value(treeStr), jni.Value(docStr))
		if err != nil {
			return err
		}

		resultStr := jni.GoString(env, jni.String(result))
		entries = parseEntries(resultStr)
		return nil
	})

	return entries, err
}

// ReadSAFFile reads file contents from a SAF document URI
func ReadSAFFile(docURI string) ([]byte, error) {
	var content []byte

	err := jni.Do(jni.JVMFor(app.JavaVM()), func(env jni.Env) error {
		if err := initSAF(env); err != nil {
			return err
		}

		ctx := jni.Object(app.AppContext())
		uriStr := jni.JavaString(env, docURI)

		result, err := jni.CallStaticObjectMethod(env, safClass, readFileID, jni.Value(ctx), jni.Value(uriStr))
		if err != nil {
			return err
		}

		if result == 0 {
			return nil
		}

		content = jni.GetByteArrayElements(env, jni.ByteArray(result))
		return nil
	})

	return content, err
}

// GetSAFTreeName gets the display name of a SAF tree root
func GetSAFTreeName(treeURI string) string {
	var name string

	jni.Do(jni.JVMFor(app.JavaVM()), func(env jni.Env) error {
		if err := initSAF(env); err != nil {
			return err
		}

		ctx := jni.Object(app.AppContext())
		uriStr := jni.JavaString(env, treeURI)

		result, err := jni.CallStaticObjectMethod(env, safClass, getTreeNameID, jni.Value(ctx), jni.Value(uriStr))
		if err != nil {
			return err
		}

		name = jni.GoString(env, jni.String(result))
		return nil
	})

	return name
}

func parseEntries(data string) []SAFEntry {
	if data == "" || strings.HasPrefix(data, "ERROR:") {
		return nil
	}

	lines := strings.Split(data, "\n")
	entries := make([]SAFEntry, 0, len(lines))

	for _, line := range lines {
		parts := strings.SplitN(line, "|", 3)
		if len(parts) != 3 {
			continue
		}

		entries = append(entries, SAFEntry{
			IsDir: parts[0] == "d",
			Name:  parts[1],
			URI:   parts[2],
		})
	}

	return entries
}

// WriteSAFFile writes data to a SAF document URI
func WriteSAFFile(docURI string, data []byte) error {
	var success bool

	err := jni.Do(jni.JVMFor(app.JavaVM()), func(env jni.Env) error {
		if err := initSAF(env); err != nil {
			return err
		}

		ctx := jni.Object(app.AppContext())
		uriStr := jni.JavaString(env, docURI)
		dataArr := jni.NewByteArray(env, data)

		result, err := jni.CallStaticBooleanMethod(env, safClass, writeFileID, jni.Value(ctx), jni.Value(uriStr), jni.Value(dataArr))
		if err != nil {
			return err
		}
		success = result
		return nil
	})

	if err != nil {
		return err
	}
	if !success {
		return errors.New("failed to write file")
	}
	return nil
}

// IsSAFURI checks if a path is a SAF content URI
func IsSAFURI(path string) bool {
	return strings.HasPrefix(path, "content://")
}

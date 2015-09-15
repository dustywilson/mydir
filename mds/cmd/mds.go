package main

import (
	"fmt"
	"strings"

	"github.com/dustywilson/mydir/mds"
)

func main() {
	france := mds.NewDirectory("France", nil)

	pants := france.NewFile("Pants")

	pantsV1 := pants.NewVersion()
	fmt.Println(pantsV1.UUID())

	pantsV2 := pants.NewVersion()
	fmt.Println(pantsV2.UUID())

	chance := france.NewDirectory("Plants").NewDirectory("Prance").NewDirectory("Chance")

	trance := chance.NewFile("Trance")

	tranceV1 := trance.NewVersion()
	fmt.Println(tranceV1.UUID())

	tranceV2 := trance.NewVersion()
	fmt.Println(tranceV2.UUID())

	fmt.Println("=====")

	dumpDirectory(france, 0)

	fmt.Println("=====")

	tranceV1.Delete()
	dumpDirectory(france, 0)

	fmt.Println("=====")

	p2d, p2f, err := france.GetByName("Pants")
	if err != nil {
		panic(err)
	}
	if p2d != nil {
		dumpDirectory(p2d, 10)
	}
	if p2f != nil {
		dumpFile(p2f, 10)
	}

	fmt.Println("=====")

	fmt.Printf("Trance: %+v\n", chance.GetFileByUUID(trance.UUID()))
	fmt.Printf("Chance: %+v\n", chance.GetDirectoryByUUID(chance.UUID()))
	fmt.Printf("PantsV2: %+v\n", chance.GetVersionByUUID(pantsV2.UUID()))

	fmt.Println("=====")

	dumpDirectory(chance.GetVersionByUUID(pantsV2.UUID()).File().Directory(), 0)

	fmt.Println("=====")

	dumpDirectory(chance.GetVersionByUUID(tranceV2.UUID()).File().Directory(), 0)
}

func dumpDirectory(d *mds.Directory, level int) {
	fmt.Printf("%s D[%s] %s\n", strings.Repeat(".", level), d.UUID(), d.Name())
	for _, child := range d.Directories() {
		dumpDirectory(child, level+1)
	}
	for _, child := range d.Files() {
		dumpFile(child, level+1)
	}
}

func dumpFile(f *mds.File, level int) {
	fmt.Printf("%s F[%s] %s\n", strings.Repeat(".", level), f.UUID(), f.Name())
	for _, child := range f.Versions() {
		dumpVersion(child, level+1)
	}
}

func dumpVersion(v *mds.Version, level int) {
	fmt.Printf("%s V[%s] version of %s\n", strings.Repeat(".", level), v.UUID(), v.File().Name())
}

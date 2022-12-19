package main

import (
	"fmt"
	"go/types"
	"golang.org/x/tools/go/packages"
	template2 "mgufrone.dev/job-alerts/cmd/generate-domain/template"
	"os"
)

func failErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
func loadPackage(path string) *packages.Package {
	cfg := &packages.Config{Mode: packages.NeedTypes | packages.NeedImports}
	pkgs, err := packages.Load(cfg, path)
	if err != nil {
		failErr(fmt.Errorf("loading packages for inspection: %v", err))
	}
	if packages.PrintErrors(pkgs) > 0 {
		os.Exit(1)
	}

	return pkgs[0]
}
func generate(root string, pkg *packages.Package, str types.Object) {
	//prefix := strings.ReplaceAll(pkg.Types.Path(), "github.com/inbeat/platform-backend/", root)
	//tplPath := root + "/cmd/generate-domain"
	fmt.Println("pkg", pkg.String(), pkg.Name)
	obj := &template2.Entity{
		PkgRoot: root,
		Root:    root,
		Pkg:     pkg,
		Str:     str.Type().Underlying().(*types.Struct),
	}
	obj.GenerateBlockCodes()
}
func main() {
	var (
		domain string
	)
	domain = os.Args[1]
	root := os.Args[2]
	sourceTypeName := "Entity"
	pkg := loadPackage(domain)
	obj := pkg.Types.Scope().Lookup(sourceTypeName)
	if obj == nil {
		failErr(fmt.Errorf("%s not found in declared types of %s",
			sourceTypeName, pkg))
	}

	// 4. We check if it is a declared type
	if _, ok := obj.(*types.TypeName); !ok {
		failErr(fmt.Errorf("%v is not a named type", obj))
	}
	// 5. We expect the underlying type to be a struct
	_, ok := obj.Type().Underlying().(*types.Struct)
	if !ok {
		failErr(fmt.Errorf("type %v is not a struct", obj))
	}

	// 6. Now we can iterate through fields and access tags
	generate(root, pkg, obj)
}

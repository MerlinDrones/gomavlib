package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/alecthomas/kong"
	"github.com/merlindrones/gomavlib/pkg/conversion"
	"github.com/merlindrones/gomavlib/pkg/templates"
)

const (
	_MAVLINK_REPO     = "https://api.github.com/repos/mavlink/mavlink/"
	_MAVLINK_RAW_REPO = "https://raw.githubusercontent.com/mavlink/mavlink/"
)

var cli struct {
	Dialects []string ` help:"names of dialects to generate" default:"common,standard,minimal"`
}

func run(args []string) error {
	parser, err := kong.New(&cli,
		kong.Description("Convert Mavlink dialects from XML format to Go format."),
		kong.UsageOnError())
	if err != nil {
		return err
	}

	_, err = parser.Parse(args)
	if err != nil {
		return err
	}
	cwd, _ := os.Getwd()
	fmt.Println(cwd)
	err = os.RemoveAll(filepath.Join(cwd, "pkg/dialects"))
	//err = shellCommand("rm -rf pkg/dialects/*/")
	if err != nil {
		fmt.Printf("cant remove existing dialects directory: %s\n", err)
		return err
	}

	os.Mkdir(filepath.Join("pkg", "dialects"), 0o755)
	os.Chdir(filepath.Join("pkg", "dialects"))

	var res struct {
		Sha string `json:"sha"`
	}
	err = downloadJSON(fmt.Sprintf("%scommits/master", _MAVLINK_REPO), &res)
	if err != nil {
		return err
	}

	var files []struct {
		Name string `json:"name"`
	}
	err = downloadJSON(fmt.Sprintf("%scontents/message_definitions/v1.0?ref=%s", _MAVLINK_REPO, res.Sha), &files)
	if err != nil {
		return err
	}

	dialectsMap := make(map[string]string)
	for _, d := range cli.Dialects {
		dialectsMap[strings.TrimSuffix(d, filepath.Ext(d))] = ""
	}

	/* @TODO: modify to only grab the dialects we are interested in */
	for _, f := range files {
		//is it in the list of dialects we want to generate
		_, ok := dialectsMap[strings.TrimSuffix(f.Name, filepath.Ext(f.Name))]
		if ok {
			fmt.Printf("Processing: %s\n", f.Name)
			//If not an xml file
			if !strings.HasSuffix(f.Name, ".xml") {
				continue
			}
			name := f.Name[:len(f.Name)-len(".xml")]

			err = processDialect(res.Sha, name)
			if err != nil {
				return err
			}
		}
	}

	err = writeTemplate(
		"package_test.go",
		templates.TplTest,
		map[string]interface{}{})
	if err != nil {
		return err
	}

	return nil
}

func downloadJSON(addr string, data interface{}) error {
	req, err := http.NewRequest(http.MethodGet, addr, nil)
	if err != nil {
		fmt.Printf("cant create http request: %s\n", err)
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("cant do http request: %s\n", err)
		return err
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(data)
	if err != nil {
		b, err := httputil.DumpResponse(res, true)
		fmt.Printf("Addr: %s\n", addr)
		if err != nil {

			log.Fatalln(err)
		}

		fmt.Println(string(b))
		return err
	}
	return nil
}

func processDialect(commit string, name string) error {
	fmt.Printf("[%s]\n", name)

	dialect := fmt.Sprintf("%s/%s/message_definitions/v1.0/%s.xml", _MAVLINK_RAW_REPO, commit, name)
	err := conversion.Convert(dialect, true)
	if err != nil {
		return err
	}

	pkgName := strings.ToLower(strings.ReplaceAll(name, "_", ""))

	err = writeTemplate("./"+pkgName+"/dialect_test.go", templates.TplDialectTest, map[string]interface{}{
		"PkgName": pkgName,
	})
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(pkgName)
	if err != nil {
		return err
	}

	for _, f := range entries {
		if !strings.HasPrefix(f.Name(), "enum_") {
			continue
		}

		buf, err := os.ReadFile(filepath.Join(pkgName, f.Name()))
		if err != nil {
			return err
		}
		str := string(buf)

		if !strings.Contains(str, "MarshalText(") {
			continue
		}

		enumName := f.Name()
		enumName = enumName[len("enum_"):]
		enumName = enumName[:len(enumName)-len(".go")]
		enumName = strings.ToUpper(enumName)

		err = writeTemplate(
			"./"+pkgName+"/"+strings.ReplaceAll(f.Name(), ".go", "_test.go"),
			templates.TplEnumTest,
			map[string]interface{}{
				"PkgName": pkgName,
				"Name":    enumName,
			})
		if err != nil {
			return err
		}
	}

	fmt.Fprintf(os.Stderr, "\n")
	return nil
}

func writeTemplate(fpath string, tpl *template.Template, args map[string]interface{}) error {
	f, err := os.Create(fpath)
	if err != nil {
		return err
	}
	defer f.Close()

	return tpl.Execute(f, args)
}

func shellCommand(cmdstr string) error {
	fmt.Fprintf(os.Stderr, "%s\n", cmdstr)
	cmd := exec.Command("sh", "-c", cmdstr)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func main() {
	err := run(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERR: %s\n", err)
		os.Exit(1)
	}
}

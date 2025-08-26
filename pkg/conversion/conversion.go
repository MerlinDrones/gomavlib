// Package conversion contains functions to convert definitions from XML to Go.
package conversion

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

const DefaultMavlinkRawBase = "https://raw.githubusercontent.com/mavlink/mavlink"
const DefaultMavlinkRef = "9fdcf82c438641984025e5457b4df854faa5a3be"

var (
	reMsgName     = regexp.MustCompile("^[A-Z0-9_]+$")
	reTypeIsArray = regexp.MustCompile(`^(.+?)\[([0-9]+)\]$`)
)

var dialectTypeToGo = map[string]string{
	"double":   "float64",
	"uint64_t": "uint64",
	"int64_t":  "int64",
	"float":    "float32",
	"uint32_t": "uint32",
	"int32_t":  "int32",
	"uint16_t": "uint16",
	"int16_t":  "int16",
	"uint8_t":  "uint8",
	"int8_t":   "int8",
	"char":     "string",
}

// getMavlinkFallbackBase returns the raw URL base for MAVLink includes.
func getMavlinkFallbackBase(ref string) string {
	return fmt.Sprintf("%s/%s/message_definitions/v1.0", DefaultMavlinkRawBase, ref)
}

// resolveInclude turns an <include> path into an absolute local path or absolute URL,
// relative to the current defAddr. Works for both local files and http(s) URLs.
func resolveInclude(defAddr, inc string, isRemote bool) (string, error) {
	if isRemote {
		u, err := url.Parse(defAddr)
		if err != nil {
			return "", fmt.Errorf("invalid base url %q: %w", defAddr, err)
		}
		base := &url.URL{
			Scheme: u.Scheme,
			Host:   u.Host,
			Path:   path.Dir(u.Path) + "/",
		}
		ref, err := url.Parse(inc)
		if err != nil {
			return "", fmt.Errorf("invalid include %q: %w", inc, err)
		}
		return base.ResolveReference(ref).String(), nil
	}

	// local file case
	baseDir := filepath.Dir(defAddr)
	return filepath.Clean(filepath.Join(baseDir, inc)), nil
}

var rePow = regexp.MustCompile(`^\s*(\d+)\s*\*\*\s*(\d+)\s*$`)

func powUint(base, exp uint64) uint64 {
	var res uint64 = 1
	for exp > 0 {
		if exp&1 == 1 {
			res *= base
		}
		base *= base
		exp >>= 1
	}
	return res
}

// parseEnumUint parses MAVLink enum numeric literals, including:
//   - decimal: "123"
//   - hex:     "0x10000"
//   - binary:  "0b0010"
//   - power:   "2**4"
func parseEnumUint(s string) (uint64, error) {
	ss := strings.TrimSpace(s)

	if m := rePow.FindStringSubmatch(ss); m != nil {
		b, err := strconv.ParseUint(m[1], 10, 64)
		if err != nil {
			return 0, err
		}
		e, err := strconv.ParseUint(m[2], 10, 64)
		if err != nil {
			return 0, err
		}
		return powUint(b, e), nil
	}

	switch {
	case strings.HasPrefix(ss, "0b") || strings.HasPrefix(ss, "0B"):
		return strconv.ParseUint(ss[2:], 2, 64)
	case strings.HasPrefix(ss, "0x") || strings.HasPrefix(ss, "0X"):
		return strconv.ParseUint(ss[2:], 16, 64)
	default:
		return strconv.ParseUint(ss, 10, 64)
	}
}

func defAddrToName(pa string) string {
	var b string
	u, err := url.ParseRequestURI(pa)
	if err == nil {
		b = path.Base(u.Path)
	} else {
		b = path.Base(pa)
	}

	b = strings.TrimSuffix(b, path.Ext(b))
	return strings.ToLower(strings.ReplaceAll(b, "_", ""))
}

func dialectNameGoToDef(in string) string {
	re := regexp.MustCompile("([A-Z])")
	in = re.ReplaceAllString(in, "_${1}")
	return strings.ToLower(in[1:])
}

func dialectNameDefToGo(in string) string {
	re := regexp.MustCompile("_[a-z]")
	in = strings.ToLower(in)
	in = re.ReplaceAllStringFunc(in, func(match string) string {
		return strings.ToUpper(match[1:2])
	})
	return strings.ToUpper(in[:1]) + in[1:]
}

func parseDescription(in string) []string {
	var lines []string

	for _, line := range strings.Split(in, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			lines = append(lines, line)
		}
	}

	return lines
}

func uintPow(base, exp uint64) uint64 {
	result := uint64(1)
	for {
		if exp&1 == 1 {
			result *= base
		}
		exp >>= 1
		if exp == 0 {
			break
		}
		base *= base
	}

	return result
}

type outEnumValue struct {
	Value       uint64
	Name        string
	Description []string
}

type outEnum struct {
	DefName     string
	Name        string
	Description []string
	Values      []*outEnumValue
	Bitmask     bool
}

type outField struct {
	Description []string
	Line        string
}

type outMessage struct {
	DefName     string
	OrigName    string
	Name        string
	Description []string
	ID          int
	Fields      []*outField
}

type outDefinition struct {
	Name     string
	Enums    []*outEnum
	Messages []*outMessage
}

// small helper so we don't rely on an undefined isURL()
func isHTTPURL(s string) bool {
	u, err := url.Parse(s)
	return err == nil && (u.Scheme == "http" || u.Scheme == "https")
}

// convertDefinition turns a parsed XML definition into our outDefinition IR.
// defName is derived from the source path/URL (e.g., "common", "standard", "swarmos").
func convertDefinition(defName string, def *definition, version int) (*outDefinition, error) {
	out := &outDefinition{
		Name:     defName,
		Enums:    make([]*outEnum, 0, len(def.Enums)),
		Messages: make([]*outMessage, 0, len(def.Messages)),
	}

	// Enums
	for _, en := range def.Enums {
		oen := &outEnum{
			DefName:     defName,
			Name:        en.Name, // <-- keep snake_case (A_TYPE)
			Description: parseDescription(en.Description),
			Values:      make([]*outEnumValue, 0, len(en.Entries)),
			Bitmask:     en.Bitmask,
		}

		for _, entry := range en.Entries {
			val, err := parseEnumUint(entry.Value)
			if err != nil {
				return nil, fmt.Errorf("enum %q entry %q: %w", en.Name, entry.Name, err)
			}
			oen.Values = append(oen.Values, &outEnumValue{
				Value:       val,
				Name:        entry.Name,
				Description: parseDescription(entry.Description),
			})
		}
		out.Enums = append(out.Enums, oen)
	}

	// Messages
	for _, m := range def.Messages {
		om, err := processMessage(defName, m)
		if err != nil {
			return nil, err
		}
		out.Messages = append(out.Messages, om)
	}

	return out, nil
}

// processDefinition is called recursively to load XML + includes.
func processDefinition(version *string, processedDefs map[string]struct{}, isRemote bool, defAddr string) ([]*outDefinition, error) {
	if _, ok := processedDefs[defAddr]; ok {
		return nil, nil
	}
	processedDefs[defAddr] = struct{}{}

	// open (use existing helpers)
	data, err := getDefinition(isRemote, defAddr)
	if err != nil {
		return nil, err
	}

	// parse XML into the correct type (definition)
	def := &definition{}
	if err := xml.Unmarshal(data, def); err != nil {
		return nil, fmt.Errorf("while parsing %s: %w", defAddr, err)
	}

	// set/keep global version string if available
	if version != nil && *version == "" && def.Version != "" {
		*version = def.Version
	}

	// choose int version for convertDefinition
	vint := 0
	if def.Version != "" {
		if vv, err := strconv.Atoi(def.Version); err == nil {
			vint = vv
		}
	} else if version != nil && *version != "" {
		if vv, err := strconv.Atoi(*version); err == nil {
			vint = vv
		}
	}

	var outDefs []*outDefinition
	defName := defAddrToName(defAddr)
	outDef, err := convertDefinition(defName, def, vint)
	if err != nil {
		return nil, fmt.Errorf("while converting %s: %w", defAddr, err)
	}
	outDefs = append(outDefs, outDef)

	// recurse into includes
	for _, inc := range def.Includes {
		// 1) try relative to current defAddr
		subAddr, err := resolveInclude(defAddr, inc, isRemote)
		if err == nil {
			subDefs, err := processDefinition(version, processedDefs, isHTTPURL(subAddr), subAddr)
			if err == nil {
				outDefs = append(outDefs, subDefs...)
				continue
			}
		}

		// 2) fallback for well-known MAVLink includes
		switch strings.TrimSpace(inc) {
		case "common.xml", "standard.xml", "minimal.xml":
			fb := getMavlinkFallbackBase(DefaultMavlinkRef)
			subAddr = fb + "/" + inc
			subDefs, err := processDefinition(version, processedDefs, true, subAddr)
			if err == nil {
				outDefs = append(outDefs, subDefs...)
				continue
			}
		}

		return nil, fmt.Errorf("unable to resolve include %q in %s", inc, defAddr)
	}

	return outDefs, nil
}

func getDefinition(isRemote bool, defAddr string) ([]byte, error) {
	if isRemote {
		byt, err := download(defAddr)
		if err != nil {
			return nil, fmt.Errorf("unable to download: %w", err)
		}
		return byt, nil
	}

	byt, err := os.ReadFile(defAddr)
	if err != nil {
		return nil, fmt.Errorf("unable to open: %w", err)
	}
	return byt, nil
}

func download(addr string) ([]byte, error) {
	res, err := http.Get(addr)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad return code: %v", res.StatusCode)
	}

	byt, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return byt, nil
}

func processMessage(defName string, msgDef *definitionMessage) (*outMessage, error) {
	if m := reMsgName.FindStringSubmatch(msgDef.Name); m == nil {
		return nil, fmt.Errorf("unsupported message name: %s", msgDef.Name)
	}

	outMsg := &outMessage{
		DefName:     defName,
		OrigName:    msgDef.Name,
		Name:        dialectNameDefToGo(msgDef.Name),
		Description: parseDescription(msgDef.Description),
		ID:          msgDef.ID,
	}

	for _, f := range msgDef.Fields {
		outField, err := processField(f)
		if err != nil {
			return nil, err
		}
		outMsg.Fields = append(outMsg.Fields, outField)
	}

	return outMsg, nil
}

func processField(fieldDef *dialectField) (*outField, error) {
	outF := &outField{
		Description: parseDescription(fieldDef.Description),
	}
	tags := make(map[string]string)

	newname := dialectNameDefToGo(fieldDef.Name)

	// name conversion is not univoque: add tag
	if dialectNameGoToDef(newname) != fieldDef.Name {
		tags["mavname"] = fieldDef.Name
	}

	outF.Line += newname

	typ := fieldDef.Type
	arrayLen := ""

	if typ == "uint8_t_mavlink_version" {
		typ = "uint8_t"
	}

	// string or array
	if matches := reTypeIsArray.FindStringSubmatch(typ); matches != nil {
		// string
		if matches[1] == "char" {
			tags["mavlen"] = matches[2]
			typ = "char"
			// array
		} else {
			arrayLen = matches[2]
			typ = matches[1]
		}
	}

	// extension
	if fieldDef.Extension {
		tags["mavext"] = "true"
	}

	typ = dialectTypeToGo[typ]
	if typ == "" {
		return nil, fmt.Errorf("unknown type: %s", typ)
	}

	outF.Line += " "
	if arrayLen != "" {
		outF.Line += "[" + arrayLen + "]"
	}
	if fieldDef.Enum != "" {
		outF.Line += fieldDef.Enum
		tags["mavenum"] = typ
	} else {
		outF.Line += typ
	}

	if len(tags) > 0 {
		var tmp []string
		for k, v := range tags {
			tmp = append(tmp, fmt.Sprintf("%s:\"%s\"", k, v))
		}
		sort.Strings(tmp)
		outF.Line += " `" + strings.Join(tmp, " ") + "`"
	}
	return outF, nil
}

func writeDialect(
	dir string,
	defName string,
	version string,
	outDefs []*outDefinition,
	enums map[string]*outEnum,
) error {
	var buf bytes.Buffer
	err := tplDialect.Execute(&buf, map[string]interface{}{
		"PkgName": defName,
		"Version": func() int {
			ret, _ := strconv.Atoi(version)
			return ret
		}(),
		"Defs":  outDefs,
		"Enums": enums,
	})
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(dir, "dialect.go"), buf.Bytes(), 0o644)
}

func writeEnum(
	dir string,
	defName string,
	enum *outEnum,
	link bool,
) error {
	var buf bytes.Buffer
	err := tplEnum.Execute(&buf, map[string]interface{}{
		"PkgName": defName,
		"Enum":    enum,
		"Link":    link && defName != enum.DefName,
	})
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(dir, "enum_"+strings.ToLower(enum.Name)+".go"), buf.Bytes(), 0o644)
}

func writeMessage(
	dir string,
	defName string,
	msg *outMessage,
	link bool,
) error {
	var buf bytes.Buffer
	err := tplMessage.Execute(&buf, map[string]interface{}{
		"PkgName": defName,
		"Msg":     msg,
		"Link":    link && defName != msg.DefName,
	})
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(dir, "message_"+strings.ToLower(msg.OrigName)+".go"), buf.Bytes(), 0o644)
}

// Convert converts a XML definition into a Golang definition.
func Convert(path string, link bool) error {
	version := "" // keep as string for writeDialect
	processedDefs := make(map[string]struct{})
	u, err := url.Parse(path)
	isRemote := u != nil && (u.Scheme == "http" || u.Scheme == "https")
	defName := defAddrToName(path)

	_, err = os.Stat(defName)
	if !os.IsNotExist(err) {
		return fmt.Errorf("directory '%s' already exists", defName)
	}

	os.Mkdir(defName, 0o755)

	// parse all definitions recursively
	outDefs, err := processDefinition(&version, processedDefs, isRemote, path)
	if err != nil {
		return err
	}

	// merge enums together
	enums := make(map[string]*outEnum)
	for _, def := range outDefs {
		for _, defEnum := range def.Enums {
			if _, ok := enums[defEnum.Name]; !ok {
				enums[defEnum.Name] = defEnum
			} else {
				enums[defEnum.Name].DefName = defName
				enums[defEnum.Name].Values = append(enums[defEnum.Name].Values, defEnum.Values...)
			}
		}
	}

	err = writeDialect(defName, defName, version, outDefs, enums)
	if err != nil {
		return err
	}

	for _, enum := range enums {
		err = writeEnum(defName, defName, enum, link)
		if err != nil {
			return err
		}
	}

	for _, def := range outDefs {
		for _, msg := range def.Messages {
			err = writeMessage(defName, defName, msg, link)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

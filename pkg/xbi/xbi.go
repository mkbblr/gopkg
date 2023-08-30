// Package xbi (Extended build info) provides functions to access build information embedded in the application. It uses [runtime/debug] package to fetch the same. Apart from the standard build info, it will fetch and display any additional build info injected into the binary. Additional build info can be injected into the binary using either '-ldflags' or [cmd/xbi-gen] command.
package xbi

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"runtime/debug"
	"strings"
)

// keys to be used while constructing input for [WithItem] function.
const (
	//key for git remote i.e. repository url
	X_BI_KEY_GIT_ORIGIN = "X_BI_KEY_GIT_ORIGIN"

	//key for git status
	X_BI_KEY_GIT_STATUS = "X_BI_KEY_GIT_STATUS"

	X_BI_KEY_GIT_LOG           = "X_BI_KEY_GIT_LOG"
	X_BI_KEY_GIT_LOCAL_COMMITS = "X_BI_KEY_GIT_LOCAL_COMMITS"
	X_BI_KEY_BUILD_PATH        = "X_BI_KEY_BUILD_PATH"
	X_BI_KEY_BUILD_TIME        = "X_BI_KEY_BUILD_TIME"
	X_BI_KEY_BUILD_HOST        = "X_BI_KEY_BUILD_HOST"
	X_BI_KEY_BUILD_USER        = "X_BI_KEY_BUILD_USER"
	X_BI_KEY_KV_PAIR           = "X_BI_KEY_KV_PAIR"
)

type XBI struct {
	B debug.BuildInfo
	X map[string]string
}

// construct a new build info object with any additional extended info passed using options.
func NewBuildInfo(opts ...func(*XBI)) XBI {
	var bi XBI
	bi.X = make(map[string]string)

	for _, opt := range opts {
		if opt != nil {
			opt(&bi)
		}
	}

	info, ok := debug.ReadBuildInfo()
	if !ok {
		return bi
	}
	bi.B = *info
	return bi
}

// Return a functional option that sets the extended build info.
// The input should be a string in the format "key=value"
// where key is something like [X_BI_KEY_GIT_ORIGIN].
// This is to be used along with [NewBuildInfo] . e.g.
//
//	xbi.NewXBI(xbi.WithExtendedInfo(main.XBI))
func WithItem(x string) func(*XBI) {
	return func(bi *XBI) {
		kv := strings.SplitN(x, ":", 2)
		if len(kv) < 2 {
			return
		}

		switch kv[0] {
		case X_BI_KEY_KV_PAIR:
			b, err := base64.StdEncoding.DecodeString(kv[1])
			if err != nil {
				return
			}

			pairs := strings.Split(string(b), ",")
			for _, p := range pairs {
				WithItem(p)(bi)
			}
		case X_BI_KEY_GIT_LOG:
			b, err := base64.StdEncoding.DecodeString(kv[1])
			if err != nil {
				return
			}
			bi.X[kv[0]] = string(b)
		case X_BI_KEY_GIT_LOCAL_COMMITS:
			b, err := base64.StdEncoding.DecodeString(kv[1])
			if err != nil {
				return
			}
			bi.X[kv[0]] = string(b)
		case X_BI_KEY_GIT_STATUS:
			b, err := base64.StdEncoding.DecodeString(kv[1])
			if err != nil {
				return
			}
			bi.X[kv[0]] = string(b)
		default:
			bi.X[kv[0]] = kv[1]
		}

	}
}

// Returns a oneliner build information
func (bi *XBI) Oneliner() string {
	ret := fmt.Sprintf("%s| %s| ", bi.B.GoVersion, bi.B.Path)

	for _, s := range bi.B.Settings {
		if s.Key == "vcs.revision" {
			ret += fmt.Sprintf("rev: %s| ", s.Value)
		}
		if s.Key == "vcs.modified" {
			ret += fmt.Sprintf("dirty: %s| ", s.Value)
		}
	}

	for k, v := range bi.X {
		if k == X_BI_KEY_BUILD_HOST {
			ret += fmt.Sprintf("host: %s| ", v)
		}
		// if k == X_BI_KEY_BUILD_USER {
		// 	ret += fmt.Sprintf("usr: %s| ", v)
		// }
		if k == X_BI_KEY_BUILD_TIME {
			ret += fmt.Sprintf("ts: %s| ", v)
		}
	}

	return ret
}

// Returns a text representation of full build information
func (bi *XBI) Text() string {
	ret := bi.B.String() + "\n"

	ret = ret + strings.Repeat("-", 10) + "\n"

	for _, s := range bi.B.Settings {
		ret = ret + s.Key + ":" + s.Value + "\n"
	}

	ret = ret + strings.Repeat("-", 10) + "\n"

	for _, m := range bi.B.Deps {
		ret = ret + "module:\t" + m.Path + "@" + m.Version + "-" + m.Sum + "\n"
	}

	ret = ret + strings.Repeat("-", 10) + "\n"

	for k, v := range bi.X {

		d := strings.Replace(k, "X_BI_KEY_", "", 1)
		d = strings.ToLower(d)
		d = strings.ReplaceAll(d, "_", ".")

		ret = ret + "\n" + d + ":\n" + v + "\n"
	}
	return ret
}

// Returns a json string of full build information
func (bi *XBI) Json() string {

	j := make(map[string]any)

	for _, b := range bi.B.Settings {
		j[b.Key] = b.Value
	}

	for k, v := range bi.X {
		d := strings.Replace(k, "X_BI_KEY_", "xbi.", 1)
		d = strings.ToLower(d)
		d = strings.ReplaceAll(d, "_", ".")

		switch k {
		case X_BI_KEY_GIT_STATUS:
			j[d] = parseGitStatus(v)
			continue
		case X_BI_KEY_GIT_LOG:
			j[d] = parseGitLog(v)
			continue
		case X_BI_KEY_GIT_LOCAL_COMMITS:
			j[d] = parseGitLog(v)
			continue
		default:
			j[d] = v
		}

	}

	ret, err := json.MarshalIndent(j, "", "  ")
	if err != nil {
		return err.Error()
	}

	return string(ret)
}

// parses the output of command `git status --porcelain=v1 -b -uall` and produces a map of key-value pairs
func parseGitStatus(s string) any {

	ret := make(map[string]any)

	lines := strings.Split(s, "\n")
	for _, l := range lines {

		r := []rune(l)
		if len(r) < 4 {
			//ideally should never happen as each line is expected to have 'xx filename'
			continue
		}
		switch {
		// case r[0] == '#':
		// 	ret["branch"] = append(ret["branch"], string(r[3:]))
		case r[0] == 'M':
			if _, ok := ret["modified"]; !ok {
				ret["modified"] = make([]string, 0)
			}
			if v, ok := (ret["modified"]).([]string); ok {
				ret["modified"] = append(v, string(r[3:]))
			}
		case r[0] == '?':
			if _, ok := ret["untracked"]; !ok {
				ret["untracked"] = make([]string, 0)
			}
			if v, ok := (ret["untracked"]).([]string); ok {
				ret["untracked"] = append(v, string(r[3:]))
			}
		}

		switch {
		case r[1] == '#':
			ret["branch"] = string(r[3:])
		case r[1] == '?':
			if _, ok := ret["untracked"]; !ok {
				ret["untracked"] = make([]string, 0)
			}
			if v, ok := (ret["untracked"]).([]string); ok {
				ret["untracked"] = append(v, string(r[3:]))
			}
		case r[1] == 'M':
			if _, ok := ret["modified"]; !ok {
				ret["modified"] = make([]string, 0)
			}
			if v, ok := (ret["modified"]).([]string); ok {
				ret["modified"] = append(v, string(r[3:]))
			}
		case r[1] == 'D':
			if _, ok := ret["deleted"]; !ok {
				ret["deleted"] = make([]string, 0)
			}
			if v, ok := (ret["deleted"]).([]string); ok {
				ret["deleted"] = append(v, string(r[3:]))
			}
		case r[1] == 'C':
			if _, ok := ret["copied"]; !ok {
				ret["copied"] = make([]string, 0)
			}
			if v, ok := (ret["copied"]).([]string); ok {
				ret["copied"] = append(v, string(r[3:]))
			}
		case r[1] == 'R':
			if _, ok := ret["renamed"]; !ok {
				ret["renamed"] = make([]string, 0)
			}
			if v, ok := (ret["renamed"]).([]string); ok {
				ret["renamed"] = append(v, string(r[3:]))
			}
		case r[1] == 'T':
			if _, ok := ret["type-changed"]; !ok {
				ret["type-changed"] = make([]string, 0)
			}
			if v, ok := (ret["type-changed"]).([]string); ok {
				ret["type-changed"] = append(v, string(r[3:]))
			}
		}

	}
	return ret
}

// parse the output of `git log -n 5 --pretty=format:%h: %s`
func parseGitLog(s string) any {

	m := make(map[string]string)

	lines := strings.Split(s, "\n")
	for _, line := range lines {
		l := strings.SplitN(line, ":", 2)
		m[l[0]] = l[1]
	}

	return m
}

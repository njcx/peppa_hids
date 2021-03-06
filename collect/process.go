package collect

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func GetProcessList() (resultData []map[string]string) {
	var dirs []string
	var err error
	dirs, err = dirsUnder("/proc")
	if err != nil || len(dirs) == 0 {
		return
	}
	for _, v := range dirs {
		pid, err := strconv.Atoi(v)
		if err != nil {
			continue
		}
		statusInfo := getStatus(pid)
		ppid,_ := strconv.Atoi(statusInfo["PPid"])
		pstatusInfo := getStatus(ppid)
		command := getcmdline(pid)
		fd := getfd(pid)
		m := make(map[string]string)
		m["pid"] = v
		m["ppid"] = statusInfo["PPid"]
		m["name"] = statusInfo["Name"]

		if len(strings.Fields(statusInfo["Uid"])) == 4 {
			m["uid"] = strings.Fields(statusInfo["Uid"])[0]
			m["euid"] = strings.Fields(statusInfo["Uid"])[1]
			m["suid"] = strings.Fields(statusInfo["Uid"])[2]
			m["fsuid"] =strings.Fields(statusInfo["Uid"])[3]
		}

		if len(strings.Fields(statusInfo["Gid"])) ==4 {
			m["gid"] = strings.Fields(statusInfo["Gid"])[0]
			m["egid"] = strings.Fields(statusInfo["Gid"])[1]
			m["sgid"] = strings.Fields(statusInfo["Gid"])[2]
			m["fsgid"] =strings.Fields(statusInfo["Gid"])[3]
		}

		if len(strings.Fields(pstatusInfo["Uid"])) ==4  {
			m["puid"] = strings.Fields(pstatusInfo["Uid"])[0]
			m["peuid"] = strings.Fields(pstatusInfo["Uid"])[1]
			m["psuid"] = strings.Fields(pstatusInfo["Uid"])[2]
			m["pfsuid"] =strings.Fields(pstatusInfo["Uid"])[3]
		}

		if len(strings.Fields(pstatusInfo["Gid"])) ==4 {
			m["pgid"] = strings.Fields(pstatusInfo["Gid"])[0]
			m["pegid"] = strings.Fields(pstatusInfo["Gid"])[1]
			m["psgid"] = strings.Fields(pstatusInfo["Gid"])[2]
			m["pfsgid"] =strings.Fields(pstatusInfo["Gid"])[3]
		}

		m["fd"] = fd
		m["command"] = command
		resultData = append(resultData, m)
	}
	return
}
func getcmdline(pid int) string {
	cmdlineFile := fmt.Sprintf("/proc/%d/cmdline", pid)
	cmdlineBytes, e := ioutil.ReadFile(cmdlineFile)
	if e != nil {
		return ""
	}
	cmdlineBytesLen := len(cmdlineBytes)
	if cmdlineBytesLen == 0 {
		return ""
	}
	for i, v := range cmdlineBytes {
		if v == 0 {
			cmdlineBytes[i] = 0x20
		}
	}
	return strings.TrimSpace(string(cmdlineBytes))
}



func getStatus(pid int) (status map[string]string) {
	status = make(map[string]string)
	statusFile := fmt.Sprintf("/proc/%d/status", pid)
	var content []byte
	var err error
	content, err = ioutil.ReadFile(statusFile)
	if err != nil {
		return
	}
	for _, line := range strings.Split(string(content), "\n") {
		if strings.Contains(line, ":") {
			kv := strings.SplitN(line, ":", 2)
			status[kv[0]] = strings.TrimSpace(kv[1])
		}
	}
	//fmt.Println(status)
	return
}

func dirsUnder(dirPath string) ([]string, error) {
	fs, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return []string{}, err
	}

	sz := len(fs)
	if sz == 0 {
		return []string{}, nil
	}
	ret := make([]string, 0, sz)
	for i := 0; i < sz; i++ {
		if fs[i].IsDir() {
			name := fs[i].Name()
			if name != "." && name != ".." {
				ret = append(ret, name)
			}
		}
	}
	return ret, nil
}


func getfd(pid int) string {
	fdDir := fmt.Sprintf("/proc/%d/fd", pid)

	dirs, err := dirsFile(fdDir)
	if err != nil || len(dirs) == 0 {
		return ""
	}

	m := []string{}
	for _, v := range dirs {
		fileInfo, err := os.Readlink(v)
		if err != nil {
			continue
		}
		countSplit := strings.Split(v, "/")
		m=append(m,strings.Join(countSplit[3:], "/")+"---"+fileInfo)

	}

	return strings.Join(m, " ")
}

func dirsFile(dirPath string) ([]string, error) {
	fs, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return []string{}, err
	}
	sz := len(fs)
	if sz == 0 {
		return []string{}, nil
	}
	ret := make([]string, 0, sz)
	for i := 0; i < sz; i++ {
		if !fs[i].IsDir() {
			name := dirPath + "/" + fs[i].Name()
			ret = append(ret, name)
		}
	}
	return ret, nil
}
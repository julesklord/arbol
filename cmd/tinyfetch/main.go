package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

type TreeNode struct {
	Text     string
	Children []*TreeNode
}

func parseFlags() (bool, bool, bool, string) {
	noASCII := false
	minimal := false
	noFrame := false
	outputFmt := ""

	for _, arg := range os.Args[1:] {
		if arg == "--no-ascii" {
			noASCII = true
		} else if arg == "--minimal" {
			minimal = true
		} else if arg == "--noframe" {
			noFrame = true
		} else if strings.HasPrefix(arg, "--output=") {
			outputFmt = strings.TrimPrefix(arg, "--output=")
		} else if arg == "--help" || arg == "-h" {
			fmt.Printf("Usage: %s [--no-ascii] [--minimal] [--noframe] [--output=json|xml|txt]\n", os.Args[0])
			os.Exit(0)
		}
	}
	return noASCII, minimal, noFrame, outputFmt
}

func gatherInfo(pluginsDir string) SystemInfo {
	hostname, _ := os.Hostname()
	osName := getOSName()
	kernel := runCommand("uname", "-r")
	uptimeVal := getUptime()
	shellVal := os.Getenv("SHELL")
	if shellVal == "" {
		shellVal = "sh"
	}
	cpuVal := getCPU()

	memRaw := getMemory()
	diskRaw := getDisk()

	var pluginKeys []string
	var pluginVals []string

	// Scan plugins directory
	if entries, err := os.ReadDir(pluginsDir); err == nil {
		type pluginResult struct {
			key string
			val string
			ok  bool
		}
		results := make([]pluginResult, len(entries))
		var wg sync.WaitGroup

		for i, entry := range entries {
			if !entry.IsDir() {
				infoPath := filepath.Join(pluginsDir, entry.Name())
				fileInfo, err := entry.Info()
				if err == nil && (fileInfo.Mode()&0111 != 0) {
					wg.Add(1)
					go func(idx int, path string, name string) {
						defer wg.Done()
						out := runCommandWithTimeout(2*time.Second, path)
						if out != "" {
							lines := strings.Split(out, "\n")
							pluginOut := strings.TrimSpace(lines[0])
							if pluginOut != "" {
								if strings.Contains(pluginOut, ":") {
									parts := strings.SplitN(pluginOut, ":", 2)
									k := parts[0]
									v := strings.TrimSpace(parts[1])
									results[idx] = pluginResult{key: k, val: v, ok: true}
								} else {
									parsedName := name
									if dotIdx := strings.Index(parsedName, "."); dotIdx != -1 {
										parsedName = parsedName[:dotIdx]
									}
									if len(parsedName) > 0 {
										parsedName = strings.ToUpper(parsedName[:1]) + parsedName[1:]
									}
									results[idx] = pluginResult{key: parsedName, val: pluginOut, ok: true}
								}
							}
						}
					}(i, infoPath, entry.Name())
				}
			}
		}
		wg.Wait()
		for _, res := range results {
			if res.ok {
				pluginKeys = append(pluginKeys, res.key)
				pluginVals = append(pluginVals, res.val)
			}
		}
	}

	return SystemInfo{
		Host:   hostname,
		OSName: osName,
		Kernel: kernel,
		Uptime: uptimeVal,
		Shell:  shellVal,
		CPU:    cpuVal,
		Memory: memRaw,
		Disk:   diskRaw,
		Keys:   pluginKeys,
		Vals:   pluginVals,
	}
}

func formatPluginName(filename string) string {
	name := filename
	if idx := strings.Index(name, "."); idx != -1 {
		name = name[:idx]
	}
	parts := strings.Split(name, "_")
	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(part[:1]) + part[1:]
		}
	}
	return strings.Join(parts, " ")
}

func printTree(node *TreeNode, prefixes []string, isLast bool) {
	if len(prefixes) > 0 {
		for _, p := range prefixes[:len(prefixes)-1] {
			fmt.Print("\033[90m" + p + "\033[0m") // Gray branches
		}
		if isLast {
			fmt.Print("\033[90m└── \033[0m")
		} else {
			fmt.Print("\033[90m├── \033[0m")
		}
	}
	fmt.Println(node.Text)

	for i, child := range node.Children {
		var nextPrefixes []string
		if len(prefixes) > 0 {
			nextPrefixes = append(nextPrefixes, prefixes...)
			if isLast {
				nextPrefixes[len(nextPrefixes)-1] = "    "
			} else {
				nextPrefixes[len(nextPrefixes)-1] = "│   "
			}
		}
		nextPrefixes = append(nextPrefixes, "│   ")
		printTree(child, nextPrefixes, i == len(node.Children)-1)
	}
}

func renderOutput(noASCII, minimal, noFrame bool, outputFmt string, infoObj SystemInfo, extPluginsDir string) {
	// Intercept output format flag early
	if outputFmt != "" {
		switch outputFmt {
		case "json":
			printJSON(infoObj)
			os.Exit(0)
		case "xml":
			printXML(infoObj)
			os.Exit(0)
		case "txt":
			printTXT(infoObj)
			os.Exit(0)
		default:
			fmt.Fprintf(os.Stderr, "Unknown output format: %s\n", outputFmt)
			os.Exit(1)
		}
	}

	// Memory & Progress Bar
	memVal := infoObj.Memory
	if strings.Contains(infoObj.Memory, "%") {
		pctPart := strings.Split(infoObj.Memory, "%")[0]
		if pct, err := strconv.Atoi(strings.TrimSpace(pctPart)); err == nil {
			memVal = getBar(pct) + " " + infoObj.Memory
		}
	}

	// Disk & Progress Bar
	diskVal := infoObj.Disk
	if strings.Contains(infoObj.Disk, "%") {
		idx := strings.Index(infoObj.Disk, "%")
		start := idx
		for start > 0 && infoObj.Disk[start-1] >= '0' && infoObj.Disk[start-1] <= '9' {
			start--
		}
		if pctStr := infoObj.Disk[start:idx]; pctStr != "" {
			if pct, err := strconv.Atoi(pctStr); err == nil {
				diskVal = getBar(pct) + " " + infoObj.Disk
			}
		}
	}

	// Styling tokens
	bold := "\033[1m"
	reset := "\033[0m"
	lblue := "\033[94m"
	lcyan := "\033[96m"

	// Build Tree Root
	root := &TreeNode{
		Text: bold + lcyan + "● " + reset + bold + infoObj.Host + reset + " @ " + lblue + infoObj.OSName + reset,
	}

	// Specs category
	specsNode := &TreeNode{Text: lcyan + bold + "specs" + reset}
	specsNode.Children = append(specsNode.Children, &TreeNode{Text: lblue + "kernel: " + reset + infoObj.Kernel})
	specsNode.Children = append(specsNode.Children, &TreeNode{Text: lblue + "uptime: " + reset + infoObj.Uptime})
	specsNode.Children = append(specsNode.Children, &TreeNode{Text: lblue + "shell: " + reset + infoObj.Shell})
	specsNode.Children = append(specsNode.Children, &TreeNode{Text: lblue + "cpu: " + reset + infoObj.CPU})
	root.Children = append(root.Children, specsNode)

	// Resources category
	resourcesNode := &TreeNode{Text: lcyan + bold + "resources" + reset}
	resourcesNode.Children = append(resourcesNode.Children, &TreeNode{Text: lblue + "memory: " + reset + memVal})
	resourcesNode.Children = append(resourcesNode.Children, &TreeNode{Text: lblue + "disk: " + reset + diskVal})
	root.Children = append(root.Children, resourcesNode)

	// Simple Plugins category
	if len(infoObj.Keys) > 0 {
		pluginsNode := &TreeNode{Text: lcyan + bold + "plugins" + reset}
		for i := 0; i < len(infoObj.Keys); i++ {
			key := strings.ToLower(infoObj.Keys[i])
			val := infoObj.Vals[i]
			pluginsNode.Children = append(pluginsNode.Children, &TreeNode{Text: lblue + key + ": " + reset + val})
		}
		root.Children = append(root.Children, pluginsNode)
	}

	// Diagnostics category (extended plugins)
	if !minimal {
		if entries, err := os.ReadDir(extPluginsDir); err == nil {
			type extResult struct {
				name  string
				lines []string
				ok    bool
			}
			results := make([]extResult, len(entries))
			var wg sync.WaitGroup

			for i, entry := range entries {
				if !entry.IsDir() {
					infoPath := filepath.Join(extPluginsDir, entry.Name())
					fileInfo, err := entry.Info()
					if err == nil && (fileInfo.Mode()&0111 != 0) {
						wg.Add(1)
						go func(idx int, path string, filename string) {
							defer wg.Done()
							ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
							out, err := exec.CommandContext(ctx, path).Output()
							cancel()
							if err == nil {
								rawOut := string(out)
								if strings.TrimSpace(rawOut) != "" {
									lines := strings.Split(rawOut, "\n")
									// Remove trailing empty line caused by Split on final newline
									if len(lines) > 0 && lines[len(lines)-1] == "" {
										lines = lines[:len(lines)-1]
									}
									if len(lines) > 0 {
										results[idx] = extResult{
											name:  formatPluginName(filename),
											lines: lines,
											ok:    true,
										}
									}
								}
							}
						}(i, infoPath, entry.Name())
					}
				}
			}
			wg.Wait()

			var diagChildren []*TreeNode
			for _, res := range results {
				if res.ok {
					pluginNode := &TreeNode{Text: lblue + strings.ToLower(res.name) + reset}
					for _, line := range res.lines {
						pluginNode.Children = append(pluginNode.Children, &TreeNode{Text: line})
					}
					diagChildren = append(diagChildren, pluginNode)
				}
			}

			if len(diagChildren) > 0 {
				diagNode := &TreeNode{Text: lcyan + bold + "diagnostics" + reset}
				diagNode.Children = diagChildren
				root.Children = append(root.Children, diagNode)
			}
		}
	}

	// Render the Tree
	printTree(root, []string{}, true)
}

func getPluginsDir() string {
	if env := os.Getenv("TINYFETCH_PLUGINS_DIR"); env != "" {
		return env
	}
	exe, err := os.Executable()
	if err != nil {
		return "./plugins"
	}
	realExe, err := filepath.EvalSymlinks(exe)
	if err != nil {
		realExe = exe
	}
	return filepath.Join(filepath.Dir(realExe), "plugins")
}

func main() {
	noASCII, minimal, noFrame, outputFmt := parseFlags()
	pluginsDir := getPluginsDir()
	extPluginsDir := filepath.Join(pluginsDir, "extended")
	infoObj := gatherInfo(pluginsDir)
	renderOutput(noASCII, minimal, noFrame, outputFmt, infoObj, extPluginsDir)
}

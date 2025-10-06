package main

// splitRmArgs separates rm flags from file operands.
// - everything before a standalone "--" that begins with '-' is treated as a flag.
// - after "--" everything is treated as a file (even if it starts with '-').
// - "-" (stdin) is ignored for backup.
func splitRmArgs(args []string) (flags []string, files []string) {
	sawDashDash := false
	for _, a := range args {
		if a == "--" {
			sawDashDash = true
			continue
		}
		if !sawDashDash {
			if a == "-" {
				// stdin, not a path to backup
				continue
			}
			if len(a) > 0 && a[0] == '-' {
				// option like -f or -rf or --flag (before --)
				flags = append(flags, a)
				continue
			}
		}
		// otherwise treat as filename/operand
		files = append(files, a)
	}
	return flags, files
}

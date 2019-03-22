package main

import (
	"github.com/urfave/cli"
)

func getFlagBuildOutput() cli.Flag {
	return cli.StringFlag{
		Name:  "output, o",
		Usage: "| where <value> is the relative path to the binary",
		Value: DefaultBuildOutput,
	}
}

func getFlagCommandsDelimiter() cli.Flag {
	return cli.StringFlag{
		Name:  "exec-delim",
		Usage: "| where <value> is the delimiter for commands in an execution group",
		Value: DefaultCommandsDelimiter,
	}
}

func getFlagEnvVars() cli.Flag {
	return cli.StringSliceFlag{
		Name:  "env, e",
		Usage: "| where <value> is the relative path to the binary - specify multiple of these to pass in multiple environment variables",
	}
}

func getFlagExecGroups() cli.Flag {
	return cli.StringSliceFlag{
		Name:  "exec",
		Usage: "| where <value> is a comma-delimited set of commands to run in parallel - specify multiple of these to define multiple execution groups",
	}
}

func getFlagFileExtensions() cli.Flag {
	return cli.StringFlag{
		Name:  "exts",
		Usage: "| where <value> is a comma-delimited set of file extensions without the period (.)",
		Value: DefaultFileExtensions,
	}
}

func getFlagIgnoredNames() cli.Flag {
	return cli.StringFlag{
		Name:  "ignore",
		Usage: "| where <value> is a comma-delimited set of file/directory names to not watch",
		Value: DefaultIgnoredNames,
	}
}

func getFlagRate() cli.Flag {
	return cli.DurationFlag{
		Name:  "rate",
		Usage: "| where <value> is a duration",
		Value: DefaultRefreshRate,
	}
}

func getFlagWatchDirectory() cli.Flag {
	return cli.StringFlag{
		EnvVar: "watch",
		Name:   "watch",
		Usage:  "| where <value> is an absolute path to a directory to watch",
		Value:  getCurrentWorkingDirectory(),
	}
}

func getFlagWorkDirectory() cli.Flag {
	return cli.StringFlag{
		EnvVar: "DIR",
		Name:   "dir",
		Usage:  "| where <value> is an absolute path to a directory to use as the current working directory",
		Value:  getCurrentWorkingDirectory(),
	}
}

func getFlagCommit() cli.Flag {
	return cli.BoolFlag{
		EnvVar: "COMMIT",
		Name:   "commit",
		Usage:  "| set the display to only the commit hash",
	}
}

func getFlagSemver() cli.Flag {
	return cli.BoolFlag{
		EnvVar: "SEMVER",
		Name:   "semver",
		Usage:  "| set the display to only the semver version",
	}
}

func getFlagSilent() cli.Flag {
	return cli.BoolFlag{
		EnvVar: "SILENT",
		Name:   "silent, s",
		Usage:  "| silence the logs",
	}
}

func getFlagVerboseLogs() cli.Flag {
	return cli.BoolFlag{
		EnvVar: "VERBOSE",
		Name:   "verbose, vv",
		Usage:  "| print verbose (debug level) logs (use for debugging)",
	}
}

func getFlagSuperVerboseLogs() cli.Flag {
	return cli.BoolFlag{
		EnvVar: "VVERBOSE",
		Name:   "vverbose, vvv",
		Usage:  "| print very verbose (trace level) logs (use for development of godev itself)",
	}
}

package main

import (
	"github.com/urfave/cli"
)

func getGlobalFlags() []cli.Flag {
	return []cli.Flag{
		getFlagBuildOutput(),
		getFlagCommandsDelimiter(),
		getFlagEnvVars(),
		getFlagExecGroups(),
		getFlagFileExtensions(),
		getFlagIgnoredNames(),
		getFlagRate(),
		getFlagSilent(),
		getFlagSuperVerboseLogs(),
		getFlagVerboseLogs(),
		getFlagWatchDirectory(),
		getFlagWorkDirectory(),
	}
}

func getCommandTestFlags() []cli.Flag {
	return []cli.Flag{
		getFlagBuildOutput(),
		getFlagCommandsDelimiter(),
		getFlagEnvVars(),
		getFlagFileExtensions(),
		getFlagIgnoredNames(),
		getFlagRate(),
		getFlagSilent(),
		getFlagSuperVerboseLogs(),
		getFlagVerboseLogs(),
		getFlagWatchDirectory(),
		getFlagWorkDirectory(),
	}
}

func getCommandVersionFlags() []cli.Flag {
	return []cli.Flag{
		getFlagSemver(),
		getFlagCommit(),
	}
}

func getFlagBuildOutput() cli.Flag {
	return cli.StringFlag{
		Usage: "| where <value> is the relative path to the binary",
		Name:  "output, o",
		Value: DefaultBuildOutput,
	}
}

func getFlagCommandsDelimiter() cli.Flag {
	return cli.StringFlag{
		Usage: "| where <value> is the delimiter for commands in an execution group",
		Name:  "exec-delim",
		Value: DefaultCommandsDelimiter,
	}
}

func getFlagEnvVars() cli.Flag {
	return cli.StringSliceFlag{
		Usage: "| where <value> is the relative path to the binary - specify multiple of these to pass in multiple environment variables",
		Name:  "env, e",
	}
}

func getFlagExecGroups() cli.Flag {
	return cli.StringSliceFlag{
		Usage: "| where <value> is a comma-delimited set of commands to run in parallel - specify multiple of these to define multiple execution groups",
		Name:  "exec",
	}
}

func getFlagFileExtensions() cli.Flag {
	return cli.StringFlag{
		Usage: "| where <value> is a comma-delimited set of file extensions without the period (.)",
		Name:  "exts",
		Value: DefaultFileExtensions,
	}
}

func getFlagIgnoredNames() cli.Flag {
	return cli.StringFlag{
		Usage: "| where <value> is a comma-delimited set of file/directory names to not watch",
		Name:  "ignore",
		Value: DefaultIgnoredNames,
	}
}

func getFlagRate() cli.Flag {
	return cli.DurationFlag{
		Usage: "| where <value> is a duration",
		Name:  "rate",
		Value: DefaultRefreshRate,
	}
}

func getFlagWatchDirectory() cli.Flag {
	return cli.StringFlag{
		Usage: "| where <value> is an absolute path to a directory to watch",
		Name:  "watch",
		Value: getCurrentWorkingDirectory(),
	}
}

func getFlagWorkDirectory() cli.Flag {
	return cli.StringFlag{
		Usage: "| where <value> is an absolute path to a directory to use as the current working directory",
		Name:  "dir",
		Value: getCurrentWorkingDirectory(),
	}
}

func getFlagCommit() cli.Flag {
	return cli.BoolFlag{
		Usage: "| set the display to only the commit hash",
		Name:  "commit",
	}
}

func getFlagSemver() cli.Flag {
	return cli.BoolFlag{
		Usage: "| set the display to only the semver version",
		Name:  "semver",
	}
}

func getFlagSilent() cli.Flag {
	return cli.BoolFlag{
		Usage: "| silence the logs",
		Name:  "silent, s",
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

package commands

import (
	"strings"

	"github.com/openinsight-project/grafinsight/pkg/bus"
	"github.com/openinsight-project/grafinsight/pkg/cmd/grafinsight-cli/commands/datamigrations"
	"github.com/openinsight-project/grafinsight/pkg/cmd/grafinsight-cli/logger"
	"github.com/openinsight-project/grafinsight/pkg/cmd/grafinsight-cli/services"
	"github.com/openinsight-project/grafinsight/pkg/cmd/grafinsight-cli/utils"
	"github.com/openinsight-project/grafinsight/pkg/services/sqlstore"
	"github.com/openinsight-project/grafinsight/pkg/setting"
	"github.com/openinsight-project/grafinsight/pkg/util/errutil"
	"github.com/urfave/cli/v2"
)

func runDbCommand(command func(commandLine utils.CommandLine, sqlStore *sqlstore.SQLStore) error) func(context *cli.Context) error {
	return func(context *cli.Context) error {
		cmd := &utils.ContextCommandLine{Context: context}
		debug := cmd.Bool("debug")

		cfg := setting.NewCfg()

		configOptions := strings.Split(cmd.String("configOverrides"), " ")
		if err := cfg.Load(&setting.CommandLineArgs{
			Config:   cmd.ConfigFile(),
			HomePath: cmd.HomePath(),
			Args:     append(configOptions, cmd.Args().Slice()...), // tailing arguments have precedence over the options string
		}); err != nil {
			return errutil.Wrap("failed to load configuration", err)
		}

		if debug {
			cfg.LogConfigSources()
		}

		engine := &sqlstore.SQLStore{}
		engine.Cfg = cfg
		engine.Bus = bus.GetBus()
		if err := engine.Init(); err != nil {
			return errutil.Wrap("failed to initialize SQL engine", err)
		}

		if err := command(cmd, engine); err != nil {
			return err
		}

		logger.Info("\n\n")
		return nil
	}
}

func runPluginCommand(command func(commandLine utils.CommandLine) error) func(context *cli.Context) error {
	return func(context *cli.Context) error {
		cmd := &utils.ContextCommandLine{Context: context}
		if err := command(cmd); err != nil {
			return err
		}

		logger.Info("\nRestart grafinsight after installing plugins . <service grafinsight-server restart>\n\n")
		return nil
	}
}

// Command contains command state.
type Command struct {
	Client utils.ApiClient
}

var cmd Command = Command{
	Client: &services.GrafinsightComClient{},
}

var pluginCommands = []*cli.Command{
	{
		Name:   "install",
		Usage:  "install <plugin id> <plugin version (optional)>",
		Action: runPluginCommand(cmd.installCommand),
	}, {
		Name:   "list-remote",
		Usage:  "list remote available plugins",
		Action: runPluginCommand(cmd.listRemoteCommand),
	}, {
		Name:   "list-versions",
		Usage:  "list-versions <plugin id>",
		Action: runPluginCommand(cmd.listVersionsCommand),
	}, {
		Name:    "update",
		Usage:   "update <plugin id>",
		Aliases: []string{"upgrade"},
		Action:  runPluginCommand(cmd.upgradeCommand),
	}, {
		Name:    "update-all",
		Aliases: []string{"upgrade-all"},
		Usage:   "update all your installed plugins",
		Action:  runPluginCommand(cmd.upgradeAllCommand),
	}, {
		Name:   "ls",
		Usage:  "list all installed plugins",
		Action: runPluginCommand(cmd.lsCommand),
	}, {
		Name:    "uninstall",
		Aliases: []string{"remove"},
		Usage:   "uninstall <plugin id>",
		Action:  runPluginCommand(cmd.removeCommand),
	},
}

var adminCommands = []*cli.Command{
	{
		Name:   "reset-admin-password",
		Usage:  "reset-admin-password <new password>",
		Action: runDbCommand(resetPasswordCommand),
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "password-from-stdin",
				Usage: "Read the password from stdin",
				Value: false,
			},
		},
	},
	{
		Name:  "data-migration",
		Usage: "Runs a script that migrates or cleanups data in your db",
		Subcommands: []*cli.Command{
			{
				Name:   "encrypt-datasource-passwords",
				Usage:  "Migrates passwords from unsecured fields to secure_json_data field. Return ok unless there is an error. Safe to execute multiple times.",
				Action: runDbCommand(datamigrations.EncryptDatasourcePasswords),
			},
		},
	},
}

var Commands = []*cli.Command{
	{
		Name:        "plugins",
		Usage:       "Manage plugins for grafinsight",
		Subcommands: pluginCommands,
	},
	{
		Name:        "admin",
		Usage:       "Grafinsight admin commands",
		Subcommands: adminCommands,
	},
}

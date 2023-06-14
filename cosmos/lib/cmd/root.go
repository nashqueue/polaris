package cmd

import (
	"cosmossdk.io/client/v2/autocli"
	"cosmossdk.io/depinject"
	"cosmossdk.io/log"
	confixcmd "cosmossdk.io/tools/confix/cmd"
	"cosmossdk.io/x/tx/signing"
	cmtcfg "github.com/cometbft/cometbft/config"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/config"
	"github.com/cosmos/cosmos-sdk/client/debug"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/pruning"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/client/snapshot"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	txmodule "github.com/cosmos/cosmos-sdk/x/auth/tx/config"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	"github.com/spf13/cobra"
	ethcryptocodec "pkg.berachain.dev/polaris/cosmos/crypto/codec"
	"pkg.berachain.dev/polaris/cosmos/simapp"
	evmante "pkg.berachain.dev/polaris/cosmos/x/evm/ante"
	evmmepool "pkg.berachain.dev/polaris/cosmos/x/evm/plugins/txpool/mempool"
)

type EncodingConfig struct {
	InterfaceRegistry codectypes.InterfaceRegistry
	AppCodec          codec.Codec
	TxConfig          client.TxConfig
	LegacyAmino       *codec.LegacyAmino
}

type RootBuilder struct {
	root *cobra.Command

	clientCtxFn     func(*EncodingConfig) client.Context
	initCometBftFn  func() *cmtcfg.Config
	initAppConfigFn func() (string, interface{})
	initRootCmdFn   func(*cobra.Command, *EncodingConfig, module.BasicManager)
}

func NewRootBuilder() *RootBuilder {
	return &RootBuilder{
		root:            &cobra.Command{},
		clientCtxFn:     DefaultClientContext,
		initCometBftFn:  DefaultInitCometBFTConfig,
		initAppConfigFn: DefaultInitAppConfig,
	}
}

func (rb *RootBuilder) Build(appCfg depinject.Config) *cobra.Command {
	var (
		encConfig = &EncodingConfig{
			LegacyAmino: &codec.LegacyAmino{},
		}
		autoCliOpts        autocli.AppOptions
		moduleBasicManager module.BasicManager
	)
	if err := depinject.Inject(depinject.Configs(appCfg, depinject.Supply(
		evmmepool.NewPolarisEthereumTxPool(), log.NewNopLogger())),
		&encConfig.InterfaceRegistry,
		&encConfig.AppCodec,
		&encConfig.TxConfig,
		&encConfig.LegacyAmino,
		&autoCliOpts,
		&moduleBasicManager,
	); err != nil {
		panic(err)
	}

	ethcryptocodec.RegisterInterfaces(encConfig.InterfaceRegistry)

	initClientCtx := rb.clientCtxFn(
		encConfig,
	)

	rootCmd := &cobra.Command{
		Use:   "polard",
		Short: "polaris sample app",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			// set the default command outputs
			cmd.SetOut(cmd.OutOrStdout())
			cmd.SetErr(cmd.ErrOrStderr())

			initClientCtx = initClientCtx.WithCmdContext(cmd.Context())
			initClientCtx, err := client.ReadPersistentCommandFlags(initClientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			initClientCtx, err = config.ReadFromClientConfig(initClientCtx)
			if err != nil {
				return err
			}

			// This needs to go after ReadFromClientConfig, as that function
			// sets the RPC client needed for SIGN_MODE_TEXTUAL.
			txConfigOpts := tx.ConfigOptions{
				TextualCoinMetadataQueryFn: txmodule.NewGRPCCoinMetadataQueryFn(initClientCtx),
			}

			// Add a custom sign mode handler for ethereum transactions.
			txConfigOpts.CustomSignModes = []signing.SignModeHandler{evmante.SignModeEthTxHandler{}}
			txConfigWithTextual, err := tx.NewTxConfigWithOptions(
				codec.NewProtoCodec(encConfig.InterfaceRegistry),
				txConfigOpts,
			)
			if err != nil {
				return err
			}

			initClientCtx = initClientCtx.WithTxConfig(txConfigWithTextual)
			if err := client.SetCmdClientContextHandler(initClientCtx, cmd); err != nil {
				return err
			}

			customAppTemplate, customAppConfig := rb.initAppConfigFn()
			customCMTConfig := rb.initCometBftFn()

			return server.InterceptConfigsPreRunHandler(cmd, customAppTemplate, customAppConfig, customCMTConfig)
		},
	}

	initRootCmd(rootCmd, encConfig, moduleBasicManager)

	if err := autoCliOpts.EnhanceRootCommand(rootCmd); err != nil {
		panic(err)
	}

	return rootCmd
}

func initRootCmd(
	rootCmd *cobra.Command,
	newApp appCreator,
	encConfig *EncodingConfig,
	basicManager module.BasicManager,
) {
	cfg := sdk.GetConfig()
	cfg.Seal()

	rootCmd.AddCommand(
		genutilcli.InitCmd(basicManager, simapp.DefaultNodeHome),
		debug.Cmd(),
		confixcmd.ConfigCommand(),
		pruning.Cmd(newApp),
		snapshot.Cmd(newApp),
	)

	server.AddCommands(rootCmd, simapp.DefaultNodeHome, newApp, appExport, addModuleInitFlags)

	// add keybase, auxiliary RPC, query, genesis, and tx child commands
	rootCmd.AddCommand(
		rpc.StatusCommand(),
		genesisCommand(txConfig, basicManager),
		queryCommand(),
		txCommand(),
		keys.Commands(simapp.DefaultNodeHome),
	)
}

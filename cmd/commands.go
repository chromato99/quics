package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/quic-s/quics/pkg/app"
	"github.com/quic-s/quics/pkg/types"
	"github.com/quic-s/quics/pkg/utils"
	"github.com/spf13/cobra"
)

/**
* Commands
*
* `qis`: Root command (meaning quic-s)
*
* `qis start`: Start quic-s server (run with default IP)
* `qis start --ip <server-ip> --port <server-port>`: Start quic-s server (run with custom IP)
* `qis stop`: Stop quic-s server
* `qis listen`: Listen quic-s protocol
* `qis run`: Run quic-s server (combine of start and listen)
*
* `qis password set --pw <password>`: Change password for quic-s server
* `qis password reset`: Reset password for quic-s server
*
* `qis show`: Show quic-s server information (needed options)
* `qis show client --id <client-UUID>`: Show client information
* `qis show client --all`: Show all clients information
* `qis show dir --id <directory-path>`: Show directory information
* `qis show dir --all`: Show all directories information
* `qis show file --id <file-path>`: Show file information
* `qis show file --all`: Show all files information
* `qis show history --id <file-history-key>`: Show history information
* `qis show history --all`: Show all history information
*
* `qis remove`: Initialize quic-s server (needed options)
* `qis remove client --id <client-UUID>`: Initialize client
* `qis remove client --all`: Initialize all clients
* `qis remove dir --id <directory-path>`: Initialize directory
* `qis remove dir --all`: Initialize all directories
* `qis remove file --id <file-path>`: Initialize file
* `qis remove file --all`: Initialize all files
*
* `qis download file --path --version --target`: Download certain file
 */

/**
* Options & Short Options
*
* `--all`: All option
* `-a`: All short option
*
* `--id`: ID option
* `-i`: ID short option
*
* `--path`: Path option
* `-p`: Path short option
*
* `--version`: Version option
* `-v`: Version short option
*
* `--target`: Target(=destination directory) option
* `--t`: Target short option
*
* `--ip`: IP option
* `--port`: Port option
*
* `--password`: Password option
 */

const (
	RootCommand     = "qis"
	StartCommand    = "start"
	StopCommand     = "stop"
	ListenCommand   = "listen"
	RunCommand      = "run"
	PasswordCommand = "password"
	ShowCommand     = "show"
	RemoveCommand   = "remove"
	DownloadCommand = "download"

	SetCommand   = "set"
	ResetCommand = "reset"

	ClientCommand  = "client"
	DirCommand     = "dir"
	FileCommand    = "file"
	HistoryCommand = "history"
)

const (
	// --all, -a
	AllOption      = "all"
	AllShortOption = "a"

	// --id, -i
	IDOption       = "id"
	IDShortCommand = "i"

	// --path, -p
	PathOption       = "path"
	PathShortCommand = "p"

	// --version, -v
	VersionOption       = "version"
	VersionShortCommand = "v"

	// --target, -t
	TargetOption       = "target"
	TargetShortCommand = "t"

	// --addr (not exist short option)
	AddrOption = "addr"

	// --port (not exist short option)
	PortOption = "port"

	// --port3 (not exist short option)
	Port3Option = "port3"

	// --pw (not exist short option)
	PasswordOption = "pw"
)

var (
	all      bool   = false
	id       string = ""
	path     string = ""
	version  uint64 = 0
	target   string = ""
	addr     string = ""
	port     string = ""
	port3    string = ""
	password string = ""
)

var rootCmd = &cobra.Command{
	Use:   RootCommand,
	Short: "qis is a CLI for interacting with the quics server",
}

var (
	startServerCmd   *cobra.Command
	stopServerCmd    *cobra.Command
	listenCmd        *cobra.Command
	runCmd           *cobra.Command
	passwordCmd      *cobra.Command
	passwordSetCmd   *cobra.Command
	passwordResetCmd *cobra.Command
	showCmd          *cobra.Command
	showClientCmd    *cobra.Command
	showDirCmd       *cobra.Command
	showFileCmd      *cobra.Command
	showHistoryCmd   *cobra.Command
	removeCmd        *cobra.Command
	removeClientCmd  *cobra.Command
	removeDirCmd     *cobra.Command
	removeFileCmd    *cobra.Command
	downloadCmd      *cobra.Command
	downloadFileCmd  *cobra.Command
)

// Run initializes and executes commands using cobra library
func Run() int {
	// initialize
	startServerCmd = initStartServerCmd()
	stopServerCmd = initStopServerCmd()
	listenCmd = initListenCmd()
	runCmd = initRunCmd()
	passwordCmd = initPasswordCmd()
	passwordSetCmd = initPasswordSetCmd()
	passwordResetCmd = initPasswordResetCmd()
	showCmd = initShowCmd()
	showClientCmd = initShowClientCmd()
	showDirCmd = initShowDirCmd()
	showFileCmd = initShowFileCmd()
	showHistoryCmd = initShowHistoryCmd()
	removeCmd = initRemoveCmd()
	removeClientCmd = initRemoveClientCmd()
	removeDirCmd = initRemoveDirCmd()
	removeFileCmd = initRemoveFileCmd()
	downloadCmd = initDownloadCmd()
	downloadFileCmd = initDownloadFileCmd()

	// set flags (= options)
	// qis start --addr <server-ip> --port <http-port> --port3 <http3-port>
	startServerCmd.Flags().StringVarP(&addr, AddrOption, "", "", "Start server with custom address")
	startServerCmd.Flags().StringVarP(&port, PortOption, "", "", "Start http rest server with custom port")
	startServerCmd.Flags().StringVarP(&port3, Port3Option, "", "", "Start http3 rest server with custom port")
	// qis run --addr <server-ip> --port <http-port> --port3 <http3-port>
	runCmd.Flags().StringVarP(&addr, AddrOption, "", "", "Start server with custom address")
	runCmd.Flags().StringVarP(&port, PortOption, "", "", "Start http rest server with custom port")
	runCmd.Flags().StringVarP(&port3, Port3Option, "", "", "Start http3 rest server with custom port")
	// qis password set --pw <password>
	passwordSetCmd.Flags().StringVarP(&password, PasswordOption, "", "", "Change password for quic-s server")
	// qis show client --id, qis show client --all
	showClientCmd.Flags().BoolVarP(&all, AllOption, AllShortOption, false, "Show all status")
	showClientCmd.Flags().StringVarP(&id, IDOption, IDShortCommand, "", "Show status by ID")
	// qis show dir --id, qis show dir --all
	showDirCmd.Flags().BoolVarP(&all, AllOption, AllShortOption, false, "Show all status")
	showDirCmd.Flags().StringVarP(&id, IDOption, IDShortCommand, "", "Show status by ID")
	// qis show file --id, qis show file --all
	showFileCmd.Flags().BoolVarP(&all, AllOption, AllShortOption, false, "Show all status")
	showFileCmd.Flags().StringVarP(&id, IDOption, IDShortCommand, "", "Show status by ID")
	// qis show history --id, qis show history --all
	showHistoryCmd.Flags().BoolVarP(&all, AllOption, AllShortOption, false, "Show all status")
	showHistoryCmd.Flags().StringVarP(&id, IDOption, IDShortCommand, "", "Show status by ID")
	// qis remove client --id, qis remove client --all
	removeClientCmd.Flags().BoolVarP(&all, AllOption, AllShortOption, false, "Initialize all data")
	removeClientCmd.Flags().StringVarP(&id, IDOption, IDShortCommand, "", "Initialize by ID")
	// qis remove dir --id, qis remove dir --all
	removeDirCmd.Flags().BoolVarP(&all, AllOption, AllShortOption, false, "Initialize all data")
	removeDirCmd.Flags().StringVarP(&id, IDOption, IDShortCommand, "", "Initialize by ID")
	// qis remove file --id, qis remove file --all
	removeFileCmd.Flags().BoolVarP(&all, AllOption, AllShortOption, false, "Initialize all data")
	removeFileCmd.Flags().StringVarP(&id, IDOption, IDShortCommand, "", "Initialize by ID")
	// qis download file --path --version
	downloadFileCmd.Flags().StringVarP(&path, PathOption, PathShortCommand, "", "Download a file by path")
	downloadFileCmd.Flags().Uint64VarP(&version, VersionOption, VersionShortCommand, 0, "Download a file by version")
	downloadFileCmd.Flags().StringVarP(&target, TargetOption, TargetShortCommand, "", "Download location")

	// add command to root command
	rootCmd.AddCommand(startServerCmd)
	rootCmd.AddCommand(stopServerCmd)
	rootCmd.AddCommand(listenCmd)
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(passwordCmd)
	rootCmd.AddCommand(showCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(downloadCmd)

	// add command to password command
	passwordCmd.AddCommand(passwordSetCmd)
	passwordCmd.AddCommand(passwordResetCmd)

	// add command to show command
	showCmd.AddCommand(showClientCmd)
	showCmd.AddCommand(showDirCmd)
	showCmd.AddCommand(showFileCmd)
	showCmd.AddCommand(showHistoryCmd)

	// add command to remove command
	removeCmd.AddCommand(removeClientCmd)
	removeCmd.AddCommand(removeDirCmd)
	removeCmd.AddCommand(removeFileCmd)

	// add command to download command
	downloadCmd.AddCommand(downloadFileCmd)

	// execute command
	if err := rootCmd.Execute(); err != nil {
		return 1
	}
	return 0
}

// initStartServerCmd start quic-s server (`qis start`)
func initStartServerCmd() *cobra.Command {
	return &cobra.Command{
		Use:   StartCommand,
		Short: "start quic-s server",
		RunE: func(cmd *cobra.Command, args []string) error {
			quicsApp, err := app.New(addr, port, port3)
			if err != nil {
				return err
			}

			err = quicsApp.StartRestServer()
			if err != nil {
				return err
			}
			return nil
		},
	}
}

// initStopServerCmd stop quic-s server (`qis stop`)
func initStopServerCmd() *cobra.Command {
	return &cobra.Command{
		Use:   StopCommand,
		Short: "stop quic-s server",
		RunE: func(cmd *cobra.Command, args []string) error {
			url := "/api/v1/server/stop"

			restClient := NewRestClient()

			_, err := restClient.PostRequest(url, "application/json", nil) // /server/stop
			if err != nil {
				log.Println("quics err: ", err)
				return err
			}

			err = restClient.Close()
			if err != nil {
				log.Println("quics err: ", err)
				return err
			}

			return nil
		},
	}
}

// initListenCmd listen quic-s protocol (`qis listen`)
func initListenCmd() *cobra.Command {
	return &cobra.Command{
		Use:   ListenCommand,
		Short: "listen quic-s protocol",
		RunE: func(cmd *cobra.Command, args []string) error {
			url := "/api/v1/server/listen"

			restClient := NewRestClient()

			_, err := restClient.PostRequest(url, "application/json", nil) // /server/listen
			if err != nil {
				log.Println("quics err: ", err)
				return err
			}

			err = restClient.Close()
			if err != nil {
				log.Println("quics err: ", err)
				return err
			}

			return nil
		},
	}
}

func initRunCmd() *cobra.Command {
	return &cobra.Command{
		Use:   RunCommand,
		Short: "run quic-s server",
		RunE: func(cmd *cobra.Command, args []string) error {
			quicsApp, err := app.New(addr, port, port3)
			if err != nil {
				return err
			}

			err = quicsApp.Run()
			if err != nil {
				return err
			}
			return nil
		},
	}
}

func initPasswordCmd() *cobra.Command {
	return &cobra.Command{
		Use:   PasswordCommand,
		Short: "change password for quic-s server",
	}
}

func initPasswordSetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   SetCommand,
		Short: "change password for quic-s server",
		RunE: func(cmd *cobra.Command, args []string) error {
			if password == "" {
				log.Println("quics: ", "Please enter password")
				cmd.Help()
				return nil
			}

			url := "/api/v1/server/password/set"

			server := &types.Server{
				Password: password,
			}

			body, err := json.Marshal(server)
			if err != nil {
				log.Println("quics err: ", err)
				return err
			}

			restClient := NewRestClient()

			_, err = restClient.PostRequest(url, "application/json", body)
			if err != nil {
				log.Println("quics err: ", err)
				return err
			}

			err = restClient.Close()
			if err != nil {
				log.Println("quics err: ", err)
				return err
			}

			return nil
		},
	}
}

func initPasswordResetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   ResetCommand,
		Short: "reset password for quic-s server",
		RunE: func(cmd *cobra.Command, args []string) error {
			url := "/api/v1/server/password/reset"

			restClient := NewRestClient()

			_, err := restClient.PostRequest(url, "application/json", nil)
			if err != nil {
				log.Println("quics err: ", err)
				return err
			}

			err = restClient.Close()
			if err != nil {
				log.Println("quics err: ", err)
				return err
			}
			return nil
		},
	}
}

func initShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   ShowCommand,
		Short: "show quic-s server data",
	}
}

func initShowClientCmd() *cobra.Command {
	return &cobra.Command{
		Use:   ClientCommand,
		Short: "show client information",
		RunE: func(cmd *cobra.Command, args []string) error {
			validateOptionByCommand(showClientCmd)

			url := "/api/v1/server/logs/clients?uuid=" + id

			restClient := NewRestClient()

			response, err := restClient.GetRequest(url) // /clients
			if err != nil {
				log.Println("quics err: ", err)
				return err
			}

			err = restClient.Close()
			if err != nil {
				log.Println("quics err: ", err)
				return err
			}

			clients := []types.Client{}
			utils.UnmarshalRequestBody(response.Bytes(), clients)

			for _, client := range clients {
				for _, root := range client.Root {
					fmt.Printf("*   UUID: %s   |   ID: %d   |   IP: %s   |   Root Directoreis: %s   *\n", client.UUID, client.Id, client.Ip, root)
				}
			}

			return nil
		},
	}
}

func initShowDirCmd() *cobra.Command {
	return &cobra.Command{
		Use:   DirCommand,
		Short: "show directory information",
		RunE: func(cmd *cobra.Command, args []string) error {
			validateOptionByCommand(showDirCmd)

			url := "/api/v1/server/logs/directories?afterPath=" + path

			restClient := NewRestClient()

			response, err := restClient.GetRequest(url) // /directories
			if err != nil {
				log.Println("quics err: ", err)
				return err
			}

			err = restClient.Close()
			if err != nil {
				log.Println("quics err: ", err)
				return err
			}

			dirs := []types.RootDirectory{}
			utils.UnmarshalRequestBody(response.Bytes(), dirs)
			for _, dir := range dirs {
				for _, UUID := range dir.UUIDs {
					fmt.Printf("*   Root Directory: %s   |   Owner: %s   |   Password: %s   |   UUID: %s   *\n", dir.AfterPath, dir.Owner, dir.Password, UUID)
				}
			}

			return nil
		},
	}
}

func initShowFileCmd() *cobra.Command {
	return &cobra.Command{
		Use:   FileCommand,
		Short: "show file information",
		RunE: func(cmd *cobra.Command, args []string) error {
			validateOptionByCommand(showFileCmd)

			url := "/api/v1/server/logs/files?afterPath=" + path

			restClient := NewRestClient()

			response, err := restClient.GetRequest(url) // /files
			if err != nil {
				log.Println("quics err: ", err)
				return err
			}

			err = restClient.Close()
			if err != nil {
				log.Println("quics err: ", err)
				return err
			}

			files := []types.File{}
			utils.UnmarshalRequestBody(response.Bytes(), files)

			for _, file := range files {
				fmt.Printf("*   File: %s   |   Root Directory: %s   |   LatestHash: %s   |   LatestSyncTimestamp: %d   |   ContentsExisted: %t   |   Metadata: %s   *\n", file.AfterPath, file.RootDirKey, file.LatestHash, file.LatestSyncTimestamp, file.ContentsExisted, file.Metadata.ModTime)
			}

			return nil
		},
	}
}

func initShowHistoryCmd() *cobra.Command {
	return &cobra.Command{
		Use:   HistoryCommand,
		Short: "show history information",
		RunE: func(cmd *cobra.Command, args []string) error {
			validateOptionByCommand(showHistoryCmd)

			url := "/api/v1/server/logs/histories?afterPath=" + path

			restClient := NewRestClient()

			response, err := restClient.GetRequest(url) // /history
			if err != nil {
				log.Println("quics err: ", err)
				return err
			}

			err = restClient.Close()
			if err != nil {
				log.Println("quics err: ", err)
				return err
			}

			histories := []types.FileHistory{}
			utils.UnmarshalRequestBody(response.Bytes(), histories)

			for _, history := range histories {
				fmt.Printf("*   Path: %s   |   Date: %s   |   UUID: %s   |   Timestamp: %d   |   Hash: %s   |*\n", history.BeforePath+history.AfterPath, history.Date, history.UUID, history.Timestamp, history.Hash)
			}

			return nil
		},
	}
}

func initRemoveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   RemoveCommand,
		Short: "initialize quic-s server",
	}
}

func initRemoveClientCmd() *cobra.Command {
	return &cobra.Command{
		Use:   ClientCommand,
		Short: "remove client",
		RunE: func(cmd *cobra.Command, args []string) error {
			validateOptionByCommand(removeClientCmd)

			url := "/api/v1/server/remove/clients?uuid=" + id

			restClient := NewRestClient()

			_, err := restClient.PostRequest(url, "application/json", nil)
			if err != nil {
				log.Println("quics err: ", err)
				return err
			}

			err = restClient.Close()
			if err != nil {
				log.Println("quics err: ", err)
				return err
			}

			return nil
		},
	}
}

func initRemoveDirCmd() *cobra.Command {
	return &cobra.Command{
		Use:   DirCommand,
		Short: "initialize directory",
		RunE: func(cmd *cobra.Command, args []string) error {
			validateOptionByCommand(removeDirCmd)

			url := "/api/v1/server/remove/directories?afterPath=" + path

			restClient := NewRestClient()

			_, err := restClient.PostRequest(url, "application/json", nil)
			if err != nil {
				log.Println("quics err: ", err)
				return err
			}

			err = restClient.Close()
			if err != nil {
				log.Println("quics err: ", err)
				return err
			}

			return nil
		},
	}
}

func initRemoveFileCmd() *cobra.Command {
	return &cobra.Command{
		Use:   FileCommand,
		Short: "initialize file",
		RunE: func(cmd *cobra.Command, args []string) error {
			validateOptionByCommand(removeFileCmd)

			url := "/api/v1/server/remove/files?afterPath=" + path

			restClient := NewRestClient()

			_, err := restClient.PostRequest(url, "application/json", nil)
			if err != nil {
				log.Println("quics err: ", err)
				return err
			}

			err = restClient.Close()
			if err != nil {
				log.Println("quics err: ", err)
				return err
			}

			return nil
		},
	}
}

func initDownloadCmd() *cobra.Command {
	return &cobra.Command{
		Use:   DownloadCommand,
		Short: "download certain file",
	}
}

func initDownloadFileCmd() *cobra.Command {
	return &cobra.Command{
		Use:   FileCommand,
		Short: "download certain file",
		RunE: func(cmd *cobra.Command, args []string) error {
			if path == "" || version == 0 || target == "" {
				log.Println("quics: ", "Please enter both path and version")
				cmd.Help()
				return nil
			}

			url := "/api/v1/server/download/files?afterPath=" + path + "&timestamp=" + fmt.Sprint(version)

			restClient := NewRestClient()

			response, err := restClient.GetRequest(url)
			if err != nil {
				log.Println("quics err: ", err)
				return err
			}

			err = restClient.Close()
			if err != nil {
				log.Println("quics err: ", err)
				return err
			}

			destinationFile, err := os.Create(target)
			if err != nil {
				return err
			}
			defer destinationFile.Close()

			n, err := destinationFile.Write(response.Bytes())
			if err != nil {
				return err
			}
			if n != len(response.Bytes()) {
				return io.ErrShortWrite
			}

			return nil
		},
	}
}

// ********************************************************************************
//                                  Private Logic
// ********************************************************************************

func validateOptionByCommand(command *cobra.Command) {
	if !all && id == "" {
		log.Println("quics: ", "Please enter only one option")
		command.Help()
		return
	}
}

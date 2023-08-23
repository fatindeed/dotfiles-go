package app

import (
	"fmt"
	"io"
	"strings"

	"github.com/c2fo/vfs/v6"
	"github.com/c2fo/vfs/v6/backend"
	"github.com/c2fo/vfs/v6/backend/s3"
	"github.com/c2fo/vfs/v6/vfssimple"
	"github.com/fatindeed/dotfiles-go/services"
	"github.com/spf13/viper"
)

var loadedFs map[string]struct{}

func init() {
	loadedFs = make(map[string]struct{})
}

func getS3Options() vfs.Options {
	opts := s3.Options{
		AccessKeyID:                 viper.GetString("s3.access_key_id"),
		SecretAccessKey:             viper.GetString("s3.secret_access_key"),
		SessionToken:                viper.GetString("s3.session_token"),
		Region:                      viper.GetString("s3.region"),
		Endpoint:                    viper.GetString("s3.endpoint"),
		ACL:                         viper.GetString("s3.acl"),
		ForcePathStyle:              viper.GetBool("s3.force_path_style"),
		DisableServerSideEncryption: viper.GetBool("s3.disable_server_side_encryption"),
	}
	if opts.AccessKeyID == "" && opts.SecretAccessKey == "" && opts.SessionToken == "" {
		return nil
	}
	return opts
}

// func getGsOptions() vfs.Options {
// 	opts := gs.Options{
// 		APIKey:         viper.GetString("gcs.api_key"),
// 		CredentialFile: viper.GetString("gcs.application_credentials"),
// 		Endpoint:       viper.GetString("gcs.endpoint"),
// 		Scopes:         viper.GetStringSlice("gcs.without_authentication"),
// 	}
// 	if opts.APIKey == "" && opts.CredentialFile == "" {
// 		return nil
// 	}
// 	return opts
// }

// func getAzureOptions() vfs.Options {
// 	opts := azure.Options{
// 		AccountName:  viper.GetString("azure.account_name"),
// 		AccountKey:   viper.GetString("azure.account_key"),
// 		TenantID:     viper.GetString("azure.tenant_id"),
// 		ClientID:     viper.GetString("azure.client_id"),
// 		ClientSecret: viper.GetString("azure.client_secret"),
// 		AzureEnvName: viper.GetString("azure.env_name"),
// 	}
// 	if opts.AzureEnvName == "" {
// 		return nil
// 	}
// 	return opts
// }

func autoloadFs(scheme string) error {
	if _, ok := loadedFs[scheme]; ok {
		return nil
	}

	v := backend.Backend(scheme)
	switch fs := v.(type) {
	case *s3.FileSystem:
		if opts := getS3Options(); opts != nil {
			fs.WithOptions(opts)
		}
	// case *gs.FileSystem:
	// 	if opts := getGsOptions(); opts != nil {
	// 		fs.WithOptions(opts)
	// 	}
	// case *azure.FileSystem:
	// 	if opts := getAzureOptions(); opts != nil {
	// 		fs.WithOptions(opts)
	// 	}
	default:
		return fmt.Errorf("%s scheme unsupported", scheme)
	}
	loadedFs[scheme] = struct{}{}
	return nil
}

func getFileContents(uri string) ([]byte, error) {
	scheme := "file"
	if pos := strings.Index(uri, "://"); pos >= 0 {
		scheme = uri[0:pos]
		if scheme == "hcp" {
			parts := strings.Split(uri[pos+3:], "/")
			if len(parts) != 2 {
				return nil, fmt.Errorf("invalid uri: %s", uri)
			}
			return services.GetHcpVaultSecret(parts[0], parts[1])
		}
	}

	err := autoloadFs(scheme)
	if err != nil {
		return nil, err
	}

	f, err := vfssimple.NewFile(uri)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return io.ReadAll(f)
}

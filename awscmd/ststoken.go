package awscmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/spf13/cobra"
	"gopkg.in/ini.v1"
)

var (
	arn string
	duration int64
	awsid string
	awssecret string
	user string
	account string
)

const arnPrefix = "arn:aws:iam::"
const mfaSuffix = ":mfa/"

// StsTokenCommand allow fetch AWS STS token agains supply MFA token and user arn
var StsTokenCommand = &cobra.Command{
	Use: "token",
	Short: "Get aws sts token against send mfa code",
	Long: "Run MFA authentication and get Security token from AWS STS service",

	RunE: func(command *cobra.Command, args []string) error {
		arn := arnPrefix + account + mfaSuffix + user
		if arn == "" {
			return nil 
		}

		profile, err := command.Flags().GetString("profile")

		if err != nil {
			return err
		}

		fmt.Print("MFA Token: ")
		token, err := command.Flags().GetString("token")
		if err != nil {
			return err
		}		
		if token == "" {
			fmt.Scanln(&token)
		}
		

		var ses session.Session
		path, err := makeAwsPath()
		
		if err != nil {
			return err
		}

		_, err = os.Stat(path)
		config := os.IsNotExist(err)
		if config  {
			awsid, err = command.Flags().GetString("awskey")
			
			if err != nil {
				return err
			} 

			awssecret, err = command.Flags().GetString("awssecret")
			
			if err != nil {
				return err
			}

			s, err := session.NewSession(&aws.Config{
				Credentials: credentials.NewStaticCredentials(awsid, awssecret, ""),
			})

			if err != nil {
				return err
			}

			ses = *s

		} else {
			s, err := session.NewSession(&aws.Config{
				Credentials: credentials.NewSharedCredentials(path, user),
			})

			if err != nil {
				return err
			}

			ses = *s
		}

		svc := sts.New(&ses)

		input := sts.GetSessionTokenInput{
			DurationSeconds: aws.Int64(duration),
			SerialNumber: aws.String(arn),
			TokenCode: aws.String(token),
	    }

		result, err := svc.GetSessionToken(&input)
		
		if err != nil {
			return err
		}

		if config {
			return writeNewConfigFile(result, user, profile)
		} 
		
		return updateExistingConfig(result, profile)
	},
}

func init() {
	StsTokenCommand.Flags().StringVarP(&user, "username", "u", "", "MFA user name")
	StsTokenCommand.Flags().StringVarP(&account, "account", "a", "", "Organization account id")
	StsTokenCommand.Flags().Int64VarP(&duration, "duration","d", 129600, "Duration in seconds accepatble 129,600 or 43,200 or 3,600 or 900")
	StsTokenCommand.Flags().StringP("token","t","","MFA token from authenticate device")
}

func makeAwsPath() (string, error) {
	var awsPath  []string

	home, err := os.UserHomeDir()

	if err != nil {
		return "", err
	}
	
	awsPath = append(awsPath, home)
	awsPath = append(awsPath, ".aws")
	path := strings.Join(awsPath, "/")
	
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, 0766)
	}
	
	awsPath = append(awsPath, "credentials")
	path = strings.Join(awsPath, "/")

	return path, nil
}

func writeNewConfigFile(creds *sts.GetSessionTokenOutput, username string, profile string) error {
	
	path, err := makeAwsPath()

	if err != nil {
		return err
	}
	
	file := ini.Empty()

	section := file.Section(username)
	section.NewKey("aws_access_key_id", awsid)
	section.NewKey("aws_secret_access_key", awssecret)

	section = file.Section(profile)
	section.NewKey("aws_access_key_id", *creds.Credentials.AccessKeyId)
	section.NewKey("aws_secret_access_key", *creds.Credentials.SecretAccessKey)
	section.NewKey("aws_session_token", *creds.Credentials.SessionToken)
	section.NewKey("aws_security_token", *creds.Credentials.SessionToken)

	return file.SaveTo(path)
}

func updateExistingConfig(creds *sts.GetSessionTokenOutput, profile string) error {
	path, err := makeAwsPath()

	if err != nil {
		return nil
	}

	file, err := ini.Load(path)

	if err != nil {
		return err
	}

	section, err := file.GetSection(profile)

	if err != nil {
		return err
	}

	section.DeleteKey("aws_access_key_id")
	section.DeleteKey("aws_secret_access_key")
	section.DeleteKey("aws_session_token")
	section.DeleteKey("aws_security_token")

	section.NewKey("aws_access_key_id", *creds.Credentials.AccessKeyId)
	section.NewKey("aws_secret_access_key", *creds.Credentials.SecretAccessKey)
	section.NewKey("aws_session_token", *creds.Credentials.SessionToken)
	section.NewKey("aws_security_token", *creds.Credentials.SessionToken)

	return file.SaveTo(path)
}
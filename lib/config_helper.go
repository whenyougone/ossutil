package lib

import (
	"fmt"
	"os"
	"os/user"
	"strconv"
	"strings"

	configparser "github.com/alyu/configparser"
)

// sections in config file
const (
	CREDSection string = "Credentials"

	BucketEndpointSection string = "Bucket-Endpoint"

	BucketCnameSection string = "Bucket-Cname"

	AkServiceSection string = "AkService"
)

// config items in section AKSerivce
const (
	ItemEcsAk string = "ecsAk"
)

type configOption struct {
	showNames     []string
	cfInteractive bool
	reveal        bool
	helpChinese   string
	helpEnglish   string
}

// CredOptionList is all options in Credentials section
var CredOptionList = []string{
	OptionLanguage,
	OptionEndpoint,
	OptionAccessKeyID,
	OptionAccessKeySecret,
	OptionSTSToken,
	OptionOutputDir,
}

// CredOptionMap allows alias name for options in Credentials section
// name, allow to show in screen
var CredOptionMap = map[string]configOption{
	OptionLanguage:        configOption{[]string{"language", "Language"}, false, true, "", ""},
	OptionEndpoint:        configOption{[]string{"endpoint", "host"}, true, true, "", ""},
	OptionAccessKeyID:     configOption{[]string{"accessKeyID", "accessKeyId", "AccessKeyID", "AccessKeyId", "access_key_id", "access_id", "accessid", "access-key-id", "access-id"}, true, false, "", ""},
	OptionAccessKeySecret: configOption{[]string{"accessKeySecret", "AccessKeySecret", "access_key_secret", "access_key", "accesskey", "access-key-secret", "access-key"}, true, false, "", ""},
	OptionSTSToken:        configOption{[]string{"stsToken", "ststoken", "STSToken", "sts_token", "sts-token"}, true, false, "", ""},
	OptionOutputDir:       configOption{[]string{"outputDir", "output-dir", "output_dir", "output_directory"}, false, true, "ossutil生成的文件的输出目录, ", "the directory to store files generated by ossutil, "},
}

// DecideConfigFile return the config file, if user not specified, return default one
func DecideConfigFile(configFile string) string {
	if configFile == "" {
		configFile = DefaultConfigFile
	}
	usr, _ := user.Current()
	dir := usr.HomeDir
	if len(configFile) >= 2 && strings.HasPrefix(configFile, "~"+string(os.PathSeparator)) {
		configFile = strings.Replace(configFile, "~", dir, 1)
	}
	return configFile
}

// LoadConfig load the specified config file
func LoadConfig(configFile string) (OptionMapType, error) {
	var configMap OptionMapType
	var err error
	configMap, err = readConfigFromFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("Read config file error: %s, please try \"help config\" to set configuration or use \"--config-file\" option", err)
	}
	if err = checkConfig(configMap); err != nil {
		return nil, err
	}
	return configMap, nil
}

func readConfigFromFile(configFile string) (OptionMapType, error) {
	configFile = DecideConfigFile(configFile)

	config, err := configparser.Read(configFile)
	if err != nil {
		return nil, err
	}

	configMap := OptionMapType{}

	// get options in cred section
	credSection, err := config.Section(CREDSection)
	if err != nil {
		return nil, err
	}

	credOptions := credSection.Options()
	for name, option := range credOptions {
		if opName, ok := getOptionNameByStr(strings.TrimSpace(name)); ok {
			configMap[strings.TrimSpace(opName)] = strings.TrimSpace(option)
		}
	}

	// get options in pair sections
	for _, sec := range []string{BucketEndpointSection, BucketCnameSection} {
		if section, err := config.Section(sec); err == nil {
			configMap[sec] = map[string]string{}
			options := section.Options()
			for bucket, host := range options {
				(configMap[sec]).(map[string]string)[strings.TrimSpace(bucket)] = strings.TrimSpace(host)
			}
		}
	}

	// get options in AKService for user-defined GetAk
	sec := AkServiceSection
	if section, err := config.Section(sec); err == nil {
		configMap[sec] = map[string]string{}
		options := section.Options()
		for ecsUrl, strUrl := range options {
			(configMap[sec]).(map[string]string)[strings.TrimSpace(ecsUrl)] = strings.TrimSpace(strUrl)
		}
	} 
	return configMap, nil
}

func getOptionNameByStr(name string) (string, bool) {
	for optionName, option := range CredOptionMap {
		for _, val := range option.showNames {
			if strings.EqualFold(name, val) {
				return optionName, true
			}
		}
	}
	return "", false
}

func checkConfig(configMap OptionMapType) error {
	for name, opval := range configMap {
		if option, ok := OptionMap[name]; ok {
			if option.optionType == OptionTypeInt64 {
				if _, err := strconv.ParseInt(opval.(string), 10, 64); err != nil {
					return fmt.Errorf("error value of option \"%s\", the value is: %s in config file, which needs int64 type", name, opval)
				}
			}
			if option.optionType == OptionTypeAlternative {
				vals := strings.Split(option.minVal, "/")
				if FindPosCaseInsen(opval.(string), vals) == -1 {
					return fmt.Errorf("error value of option \"%s\", the value is: %s in config file, which is not anyone of %s", name, opval, option.minVal)
				}
			}
		}
	}
	return nil
}

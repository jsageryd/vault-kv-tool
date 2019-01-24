package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	vaultapi "github.com/hashicorp/vault/api"
)

func main() {
	log.SetFlags(0)

	root := flag.String("root", "", "Secrets root (required); must be kv version 1")
	write := flag.Bool("write", false, "If set, write to the given root path instead of reading")

	flag.Parse()

	if *root == "" {
		fmt.Fprintln(os.Stderr, "secrets root not specified.")
		fmt.Fprintln(os.Stderr)
		flag.Usage()
		os.Exit(1)
	}

	vault, err := vaultapi.NewClient(nil)
	if err != nil {
		log.Fatalf("error initializing vault client: %v", err)
	}

	logical := vault.Logical()

	var data map[string]interface{}

	if *write {
		if err := json.NewDecoder(os.Stdin).Decode(&data); err != nil {
			log.Fatalf("error decoding JSON input: %v", err)
		}
		if err := writeData(*root, logical, data); err != nil {
			log.Fatal(err)
		}
	} else {
		data, err = readData(*root, logical, nil)
		if err != nil {
			log.Fatal(err)
		}
	}

	json.NewEncoder(os.Stdout).Encode(data)
}

func readData(root string, logical *vaultapi.Logical, data map[string]interface{}) (map[string]interface{}, error) {
	if data == nil {
		data = make(map[string]interface{})
	}

	sec, err := logical.List(root)
	if err != nil {
		return nil, err
	}

	if sec == nil {
		return nil, fmt.Errorf("[%s] no secrets found", root)
	}

	if sec.Data == nil {
		return nil, fmt.Errorf("[%s] data is nil", root)
	}

	keys, ok := sec.Data["keys"]
	if !ok {
		return nil, fmt.Errorf("[%s] unexpected data; not a list?", root)
	}

	is, ok := keys.([]interface{})
	if !ok {
		return nil, fmt.Errorf("[%s] unexpected list type %T, wanted a %T", root, is, []interface{}{})
	}

	for _, k := range is {
		key, ok := k.(string)
		if !ok {
			return nil, fmt.Errorf("[%s] unexpected key %T, wanted a %T", root, key, "")
		}
		absPath := path.Join(root, key)
		if key[len(key)-1] == '/' {
			data[key], err = readData(absPath, logical, nil)
			if err != nil {
				return nil, err
			}
		} else {
			sec, err := logical.Read(absPath)
			if err != nil {
				return nil, err
			}
			data[key] = sec.Data
		}
	}

	return data, nil
}

func writeData(root string, logical *vaultapi.Logical, data map[string]interface{}) error {
	for key, val := range data {
		mVal, ok := val.(map[string]interface{})
		if !ok {
			return fmt.Errorf("[%s] value is a %T, wanted a %T", root, mVal, map[string]interface{}{})
		}
		absPath := path.Join(root, key)
		if key[len(key)-1] == '/' {
			if err := writeData(absPath, logical, mVal); err != nil {
				return err
			}
		} else {
			if _, err := logical.Write(absPath, mVal); err != nil {
				return err
			}
		}
	}
	return nil
}

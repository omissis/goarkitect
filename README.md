# GoArkitect

This project gives developers the ability to describe and check the architecture of a project and check it is respected at any time.

## Example usage

```sh
# validate the default config file (.goarkitect.yaml) and outputs the result in json
goarkitect validate --output=json

# validate the custom .ark.yaml config file
goarkitect validate .ark.yaml

# validate the custom .ark.yaml config file and all the config files found in the .ark/ folder
goarkitect validate .ark.yaml .ark/

# verify that the current folder follows the rules specified in the default config file (.goarkitect.yaml)
goarkitect verify

# verify that the current folder follows the rules specified in the .ark/ folder and outputs the result in json
goarkitect verify .ark/ --output=json
```

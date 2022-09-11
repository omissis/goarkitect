# GoArkitect

<p align="center">
<img src="docs/assets/goarkitect.logo.jpg" alt="crkitect" title="goarkitect" />
</p>

This project gives developers the ability to describe and check  architectural constraints of a project using a composable set of rules described in one or multiple yaml files.

## Example configuration

```yaml
---
rules:
  - name: all inner asdf manifests contain golang, possibly in its 1.19.1 version  # name of the rule, it should tell what the rule is about
    kind: file # name of the matcher to use, which tells what objects it will operate on
    matcher:
      kind: all # the all kind sets the matcher to match all possible files, which will be narrowed down below
    thats: # 'thats' apply filtering to the selected matchers to narrow down the files to operate on
      - kind: end_with # matcher filter that selects only files whose name ends with .tool-versions"
        suffix: .tool-versions
    excepts: # 'excepts' allow to pull some special cases out of the set of file determined by the 'thats' filters
      - kind: this # excepts the one in the root directory
        filePath: ./.tool-versions
    musts: # 'musts' will trigger errors in case the expectation is not respected, which in turn will have goarkitect to exit with status code 1
      - kind: contain_value
        value: golang
    shoulds: # 'shoulds' will trigger warnings, which won't cause error status codes on exit
      - kind: contain_value
        value: golang 1.19.1
    coulds: # 'coulds' will trigger info-level notices, and they can be seen as suggestions
      - kind: contain_value
        value: golangci-lint
    because: "it is needed for the project to compile the source code" # reason for the rule to exists
```

See the [examples folder](./examples/) for more complete information on how to configure goarkitect rules.

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

## Acknowledgements

Goarkitect draws inspiration from [Alessandro Minoccheri](https://alessandrominoccheri.github.io)'s [PHPArkitect](https://github.com/phparkitect/arkitect), a tool that helps you to keep your PHP codebase coherent and solid.

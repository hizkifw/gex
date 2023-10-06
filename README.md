# gex!

Go Hex Editor with vim-like keybindings

![Screenshot](https://github.com/hizkifw/gex/assets/7418049/ebdfa046-b436-4d87-a94d-2e978194b9e3)

## Status

gex! is very early in development and should not be used to load or modify any
important data.

For documentation on how to use gex!, see the [`doc`] folder.

[`doc`]: doc/help.md

## Install

You can download the latest CI build from the [Actions page], or from here.

| [Linux x64] | [Windows x64] | [Darwin x64] |
| ----------- | ------------- | ------------ |

[Actions page]: https://github.com/hizkifw/gex/actions
[Linux x64]:
  https://nightly.link/hizkifw/gex/workflows/ci.yaml/main/gex-linux-amd64.zip
[Windows x64]:
  https://nightly.link/hizkifw/gex/workflows/ci.yaml/main/gex-windows-amd64.zip
[Darwin x64]:
  https://nightly.link/hizkifw/gex/workflows/ci.yaml/main/gex-darwin-amd64.zip

You can also install from source. You'll need **Go 1.20** to build and install
gex!.

```sh
go install github.com/hizkifw/gex/cmd/gex@latest
```

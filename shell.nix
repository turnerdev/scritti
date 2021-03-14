{ pkgs ? import <nixpkgs> {} }:

with pkgs;
let
  drv = callPackage ./default.nix {};
in
  drv.overrideAttrs (attrs: {
    src = null;
    nativeBuildInputs = [ govers delve plantuml jre graphviz gnumake ] ++ attrs.nativeBuildInputs;
    shellHook = ''
      echo 'Entering ${attrs.pname}'
      set -v
      export GRAPHVIZ_DOT=${graphviz}
      export GOPATH="$(pwd)/.go"
      export GOCACHE=""
      export CGO_ENABLED=0
      export GO111MODULE='on'
      go mod init ${attrs.goPackagePath}
      set +v
    '';
  })
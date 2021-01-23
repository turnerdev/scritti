{ pkgs ? import <nixpkgs> {} }:

with pkgs;
let
  drv = callPackage ./default.nix {};
in
  drv.overrideAttrs (attrs: {
    src = null;
    nativeBuildInputs = [ govers delve ] ++ attrs.nativeBuildInputs;
    shellHook = ''
      echo 'Entering ${attrs.pname}'
      set -v
      export GOPATH="$(pwd)/.go"
      export GOCACHE=""
      export CGO_ENABLED=0
      export GO111MODULE='on'
      go mod init ${attrs.goPackagePath}
      set +v
    '';
  })
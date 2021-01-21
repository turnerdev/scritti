{ buildGoModule
, nix-gitignore
}:

buildGoModule {
  pname = "golang-demo";
  version = "0.0.1";
  src = nix-gitignore.gitignoreSource [] ./.;
  goPackagePath = "scritti";
  modSha256 = "1gqn2vm3wrc3gml8mhkf11sap3hkyzhi90qwzw0x93wv6vmm4mcy";
  vendorSha256 = null;
}
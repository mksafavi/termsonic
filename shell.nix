let
  pkgs = import <nixpkgs> {};
in
pkgs.mkShell {
  packages = [
    pkgs.go
    pkgs.gotools
    pkgs.gopls
    pkgs.pkg-config
    pkgs.alsa-lib
  ];
}

{
  description = "A tool for viewing claude code log files";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = {
    self,
    nixpkgs,
    flake-utils,
  }:
    flake-utils.lib.eachDefaultSystem (
      system: let
        pkgs = nixpkgs.legacyPackages.${system};
        cclog = pkgs.callPackage ./default.nix {};
      in {
        packages.cclog = cclog;
        packages.default = cclog;
        devShells.default = pkgs.mkShell {
          inputsFrom = [
            cclog
          ];
        };
      }
    );
}

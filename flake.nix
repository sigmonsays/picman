{
  description = "picture management tool";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  inputs.flake-utils.url = "github:numtide/flake-utils";
  inputs.gomod2nix.url = "github:nix-community/gomod2nix";

  outputs = {
    self,
    nixpkgs,
    flake-utils,
    gomod2nix,
  }: (
    flake-utils.lib.eachDefaultSystem
    (system: let


      pkgs = import nixpkgs {
        inherit system;
        overlays = [gomod2nix.overlays.default];
      };

      bin = pkgs.callPackage ./. {};

      dockerImage = pkgs.dockerTools.buildImage {
        name = "picman";
        tag = "latest";
        copyToRoot = [bin];
        config = {
          Cmd = ["${bin}/bin/picman"];
        };
      };
    in with pkgs; {
      packages = {
        inherit bin dockerImage;
        default = bin;
      };
      devShells.default = import ./shell.nix {inherit pkgs;};
    })
  );
}

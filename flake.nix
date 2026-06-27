{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    go-overlay.url = "github:purpleclay/go-overlay";
  };

  outputs = { nixpkgs, flake-utils, go-overlay, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs {
          inherit system;
          overlays = [ go-overlay.overlays.default ];
        };

        go = pkgs.go-bin.fromGoMod ./go.mod;
      in {
        packages.default = pkgs.buildGoApplication {
          inherit go;
          pname = "ovr";
          version = "0.1.0";
          src = ./.;
          modules = ./govendor.toml;
          subPackages = [ "cmd/ovr" ];
        };

        devShells.default = pkgs.mkShell {
          buildInputs = [ go.withDefaultTools ];
        };
      }
    );
}

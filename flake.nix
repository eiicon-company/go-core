{
  inputs = {
    nixpkgs.url = "https://flakehub.com/f/NixOS/nixpkgs/*";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      nixpkgs,
      flake-utils,
      ...
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        go = pkgs.go_1_24;
        buildGoModule = pkgs.buildGoModule.override {
          inherit go;
        };
        gopkgs =
          with (import nixpkgs {
            inherit system;
            overlays = [ (_: _: { inherit buildGoModule; }) ];
          }); [
            golangci-lint
            gopls
            gotools
          ];

        tools = with pkgs; [
          mozjpeg
          pngquant
          gifsicle

          gnumake

          actionlint

          nil # nix lsp
          nixfmt-rfc-style
        ];

        devShells.default = pkgs.mkShellNoCC {
          packages = [ go ] ++ gopkgs ++ tools;
        };
      in
      {
        legacyPackages = pkgs;
        inherit devShells;
      }
    );
}

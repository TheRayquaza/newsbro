{
  description = "Simple python fhs devshell that works";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = {
    nixpkgs,
    flake-utils,
    ...
  }:
    flake-utils.lib.eachDefaultSystem (
      system: let
        pkgs = import nixpkgs {
          inherit system;
        };
      in {
        devShells.default = let
          python = pkgs.python312;
        in
          (pkgs.buildFHSEnv {
            name = "simple-python-fhs";

            targetPkgs = _: [
              python
              pkgs.uv
              pkgs.zlib
            ];
            profile = ''
              export UV_PYTHON=${python}
              bash
            '';
          }).env;
      }
    );
}

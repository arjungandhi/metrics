{
  description = "Metrics tracking CLI";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        packages.default = pkgs.buildGoModule {
          pname = "metrics";
          version = "0.1.0";
          src = ./.;
          vendorHash = "sha256-2kUo4KtkmntEw39iefQygQLTgyESgaHEfT1TMj4vEt0=";
          subPackages = [ "cmd/metrics" ];

          meta = with pkgs.lib; {
            description = "Metrics tracking CLI";
            homepage = "https://github.com/arjungandhi/metrics";
            license = licenses.mit;
            mainProgram = "metrics";
          };
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            gopls
          ];
        };
      }
    );
}

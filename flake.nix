{
  description = "giopad - Go/Gio markdown editor";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            gopls
            gotools

            # Gio dependencies for Linux
            pkg-config
            libxkbcommon
            wayland
            wayland-protocols
            libGL
            libglvnd
            egl-wayland
            libx11
            libxcursor
            libxfixes
            libxcb
            libxi
            vulkan-headers
            vulkan-loader

            # Android (optional)
            # android-sdk
          ];

          shellHook = ''
            echo "giopad dev shell"
            echo "  go build    - build desktop"
            echo "  gogio -target android . - build android"
          '';
        };

        packages.default = pkgs.buildGoModule {
          pname = "giopad";
          version = "0.1.0";
          src = ./.;
          vendorHash = null; # Update after first build

          nativeBuildInputs = with pkgs; [ pkg-config ];
          buildInputs = with pkgs; [
            libxkbcommon
            wayland
            wayland-protocols
            libGL
            libglvnd
            egl-wayland
            libx11
            libxcursor
            libxfixes
            libxcb
            libxi
            vulkan-headers
            vulkan-loader
          ];
        };
      }
    );
}

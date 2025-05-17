{ pkgs ? import <nixpkgs> {} }:
pkgs.mkShell {
  nativeBuildInputs = with pkgs; [
    xorg.libX11.dev
    xorg.libXft
    xorg.libXcursor
    xorg.libXrandr
    xorg.libXinerama
    xorg.xinput
    xorg.libXi
    xorg.libXxf86vm
    libGL
    libGLU
    mesa
    mesa_glu
    pkg-config
  ];
}

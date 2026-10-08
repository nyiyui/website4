{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachSystem (import ./flake.systems.nix) (system:
      let
        pkgs = import nixpkgs { inherit system; };
        version = self.shortRev or "dirty";

        # The render-once Go binary.
        render = pkgs.buildGoModule {
          pname = "render";
          inherit version;
          src = ./render;
          subPackages = [ "cmd/render-once" "cmd/imgprep" ];
          vendorHash = "sha256-mhgJ182qAKMpmBOfZ0DK/is9a8G9PFQnUjKg1ARwG7w=";

          nativeBuildInputs = [ pkgs.makeWrapper ];
          postInstall = ''
            wrapProgram $out/bin/render-once \
              --prefix PATH : ${pkgs.lib.makeBinPath [ pkgs.typst ]} \
              --set TYPST_FONT_PATHS ${pkgs.atkinson-hyperlegible}/share/fonts
            wrapProgram $out/bin/imgprep \
              --prefix PATH : ${pkgs.lib.makeBinPath [ pkgs.libwebp pkgs.libavif ]}
          '';

          meta.mainProgram = "render-once";
        };

        # The built static site: runs render-once over ./blog, ./pages, and
        # ./render/template, producing ./public.
        site = pkgs.stdenvNoCC.mkDerivation {
          pname = "website";
          inherit version;
          src = ./.;
          nativeBuildInputs = [ render ];

          buildPhase = ''
            runHook preBuild
            export HOME="$TMPDIR"
            render-once
            runHook postBuild
          '';

          installPhase = ''
            runHook preInstall
            cp -r public "$out"
            runHook postInstall
          '';
        };
      in
      {
        packages.render = render;
        packages.site = site;
        packages.default = site;

        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            typst
            tinymist
            atkinson-hyperlegible

            go
            gopls

            libwebp
            libavif
          ];
          TYPST_FEATURES = "html";
          TYPST_FONT_PATHS = "${pkgs.atkinson-hyperlegible}/share/fonts";
        };
      });
}

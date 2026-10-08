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

        # Pre-fetched @preview typst packages, so `typst compile` doesn't
        # need network access (the Nix build sandbox has none). Laid out to
        # match typst's package cache directory structure, so pointing
        # TYPST_PACKAGE_CACHE_PATH at this is enough for typst to find them
        # without attempting a download.
        typstPackageCache = pkgs.linkFarm "typst-package-cache" [
          {
            # cetz-plot internally pins this exact version (src/cetz.typ
            # imports "@preview/cetz:0.4.0"), so the blog post's own cetz
            # import is pinned to match rather than vendoring two copies.
            name = "preview/cetz/0.4.0";
            path = pkgs.fetchzip {
              url = "https://packages.typst.org/preview/cetz-0.4.0.tar.gz";
              sha256 = "1hikxb66s8s0fbc1cdh9939x3bnqfb4n2v5qsypg32mn37y7933z";
            };
          }
          {
            name = "preview/cetz-plot/0.1.2";
            path = pkgs.fetchzip {
              url = "https://packages.typst.org/preview/cetz-plot-0.1.2.tar.gz";
              sha256 = "1rcs373c5jvx2cj5305886azl309xbjwvicp6pxppi1pv92jjnwk";
            };
          }
          {
            # Transitive dependency of cetz (src/deps.typ).
            name = "preview/oxifmt/1.0.0";
            path = pkgs.fetchzip {
              url = "https://packages.typst.org/preview/oxifmt-1.0.0.tar.gz";
              sha256 = "0bqc5ahiavjds966a8v5llw9imqqaa43x89ha9dl5nlp51vqmla6";
            };
          }
        ];

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
              --set TYPST_FONT_PATHS ${pkgs.atkinson-hyperlegible}/share/fonts \
              --set TYPST_PACKAGE_CACHE_PATH ${typstPackageCache}
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
          TYPST_PACKAGE_CACHE_PATH = "${typstPackageCache}";
        };
      });
}

{

  description = "HTTP (webhook) listener to announce messages in IRC channels";
  inputs.flake-utils.url = "github:numtide/flake-utils";

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      with nixpkgs.legacyPackages.${system}; rec {

        packages = flake-utils.lib.flattenTree rec {

          http2irc = buildGoModule rec {

            pname = "http2irc";
            version = "0.1";

            src = ./.;
            vendorSha256 =
              "sha256-k45e6RSIl3AQdOFQysIwJP9nlYsSFeaUznVIXfbYwLA=";
            subPackages = [ "." ];

            meta = with lib; {
              description =
                "HTTP (webhook) listener to announce messages in IRC channels";
              homepage = "https://github.com/pinpox/http2irc";
              license = licenses.gpl3;
              maintainers = with maintainers; [ pinpox ];
              platforms = platforms.linux;
            };
          };

          announce-drone = pkgs.writeScriptBin "announce-drone" ''
            #!${pkgs.stdenv.shell}
            echo "[$DRONE_REPO - $DRONE_COMMIT_REF] $DRONE_BUILD_STATUS: \
            $DRONE_COMMIT_MESSAGE - $DRONE_COMMIT_AUTHOR $DRONE_BUILD_EVENT $DRONE_BUILD_LINK" | \
            curl -H "Content-Type: application/json" -X POST --data-binary @- localhost:8080
          '';

          announce-test = pkgs.writeScriptBin "announce-test" ''
            #!${pkgs.stdenv.shell}
            echo "I'm a bot!"
          '';

        };

        defaultPackage = packages.http2irc;
      });
}
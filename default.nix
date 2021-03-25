with import <nixpkgs>{};

  buildGoModule rec {

    pname = "http2irc";
    version = "0.1";

    src = pkgs.fetchFromGitHub {
      owner = "pinpox";
      repo = "http2irc";
      # rev = "v${version}";
      rev = "main";
      sha256 = "sha256-hq83FdSaHe6xrFPMngILNOUP++lYNrBYiuPRivNLqqc=";
    };

    vendorSha256 = "sha256-k45e6RSIl3AQdOFQysIwJP9nlYsSFeaUznVIXfbYwLA=";
    subPackages = [ "." ];

    meta = with lib; {
      description = "Webhook reciever to annouce in IRC channels";
      homepage = "https://github.com/pinpox/http2irc";
      license = licenses.gpl3;
      maintainers = with maintainers; [ pinpox ];
      platforms = platforms.linux;
    };
  }

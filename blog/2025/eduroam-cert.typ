#set document(
  title: [How to Connect to GT Eduroam with a Client Certificate],
	author: ("Ken Shibata"),
	description: [*description* wow],
	keywords: ("eduroam", "NetworkManager"),
  date: datetime(year: 2025, month: 7, day: 20)
)

= How to Connect to Eduroam using GT's New Certificate-based Authentication

_This is a living document; content will be updated as necessary (after the initial post date)._

Georgia Tech has been rolling out a new TLS client certificate-based (instead of a username and password) authentication method for eduroam.

As a supplement to the #link("https://gatech.service-now.com/technology?id=kb_article_view&sysparm_article=KB0044211")[guide on ServiceNow], this page provides a walkthrough for NetworkManager, using `nm-connection-editor`.

== Generate a Certificate using the SecureW2 JoinNow Advanced Portal

_Note: This section is essentially the same as the ServiceNow page._

Get your Wi-Fi interface's MAC address (or choose a random MAC address to use):
```sh
nmcli
```

You'll get something like this:
```
wlp0s20f3: connected to eduroam (GT cert)
	"Intel Something Something"
	wifi (iwlwifi), AA:AA:AA:AA:AA:AA, hw, mtu 1234
	ip4 default
	inet4 1.2.3.4/23
	route4 1.2.3.4/23 metric 600
	route4 default via 1.2.3.4 metric 600
	inet6 fe80::/64
	route6 fe80::/64 metric 1024

lo: connected (externally) to lo
	"lo"
	loopback (unknown), 00:00:00:00:00:00, sw, mtu 65536
	inet4 127.0.0.1/8
	inet6 ::1/128

DNS configuration:
	servers: 137.82.1.1 142.103.1.1 137.82.1.2 142.103.1.42
	domains: wireless.ubc.ca
	interface: wlp0s20f3

Use "nmcli device show" to get complete information about known devices and
"nmcli connection show" to get an overview on active connection profiles.

Consult nmcli(1) and nmcli-examples(7) manual pages for complete usage details.
```

1. Open the #link("https://wifi.gatech.edu/joinnowadvanced")[SecureW2 JoinNow Advanced Portal]. Note: make sure to enter the correct MAC address.
2. Download the P12 file. (Make sure to keep the passphrase you used here.)
3. Download the `USERTrust RSA Certification Authority`.

== Setup NetworkManager
#table(
columns: 2,
[Name], [Value],
[Security], [WPA/WPA2 Enterprise],
[Authentication], [TLS],
[Identity], [gburdell3\@gatech.edu],
[Domain], [lawn.gatech.edu],
[CA certificate], [usertrustrsaca [jdk].cer #footnote[Downloaded in step 3.]],
[CA certificate password], [None],
[User certificate], [gburdell3\@gatech.edu.p12 #footnote[This is the #link("https://en.wikipedia.org/wiki/PKCS_12")[PKCS #12] file you downloaded in step 2.] <p12_file>],
[User certificate password], [None],
[User private key], [gburdell3\@gatech.edu.p12 @p12_file (auto-filled when entering User certificate)],
[User key password], [passphrase given in SecureW2 portal]
)
Done!

package main

import "testing"

func TestMeshCredentials(t *testing.T) {
	cfg := map[string]string{"wireless.10.usage": "uplink", "wireless.10.ssid": "vport-example", "aaa.10.ssid": "vport-example", "aaa.1.ssid": "users", "aaa.1.wpa.psk": "userpassword"}
	check := func(wantSSID, wantKey string) {
		t.Helper()
		s, k := meshCredentials(cfg)
		if s != wantSSID || k != wantKey {
			t.Fatal("unexpected credential selection")
		}
	}
	check("", "")
	cfg["mesh.essid"] = "vwire-example"
	cfg["mesh.psk"] = "sitepassword"
	check("vwire-example", "sitepassword")
	delete(cfg, "mesh.essid")
	delete(cfg, "mesh.psk")
	cfg["wireless.4.usage"] = "downlink"
	cfg["wireless.4.mode"] = "master"
	cfg["wireless.4.vwire"] = "enabled"
	cfg["wireless.4.ssid"] = "vwire-example"
	cfg["aaa.4.ssid"] = "vwire-example"
	cfg["aaa.4.wpa.psk"] = "sitepassword"
	check("vwire-example", "sitepassword")
}

package yggdrasilapi

import (
	"craapi/packages/encryption"

	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
)

type Ygg_links struct {
	Homepage string `json:"homepage"`
	Register string `json:"register"`
}
type Ygg_Meta_DATA struct {
	ServerName                  string     `json:"serverName"`
	ImplementationName          string     `json:"implementationName"`
	ImplementationVersion       string     `json:"implementationVersion"`
	Links                       *Ygg_links `json:"links"`
	Non_email_login             bool       `json:"feature.non_email_login"`
	Legacy_skin_api             bool       `json:"feature.legacy_skin_api"`
	No_mojang_namespace         bool       `json:"feature.no_mojang_namespace"`
	Enable_mojang_anti_features bool       `json:"feature.enable_mojang_anti_features"`
	Enable_profile_key          bool       `json:"feature.enable_profile_key"`
	Username_check              bool       `json:"feature.username_check"`
}
type Ygg_Meta struct {
	Meta               *Ygg_Meta_DATA `json:"meta"`
	SkinDomains        []string       `json:"skinDomains"`
	SignaturePublickey string         `json:"signaturePublickey"`
}
type Ygg_PUBKEY struct {
	SignaturePublickey string `json:"signaturePublickey"`
}

var metadata []byte
var metakey []byte

func Yggdrasilapiinit(SkinDomain []string, Domain string) {
	METADATA := Ygg_Meta{
		Meta: &Ygg_Meta_DATA{
			ServerName:            "CRA api",
			ImplementationName:    "CRA api",
			ImplementationVersion: "0x01",
			Links: &Ygg_links{
				Homepage: Domain + "/",
				Register: Domain + "/register",
			},
			Non_email_login:             true,
			Legacy_skin_api:             false,
			No_mojang_namespace:         false,
			Enable_mojang_anti_features: false,
			Enable_profile_key:          true,
			Username_check:              false,
		},
		SkinDomains:        SkinDomain,
		SignaturePublickey: string(encryption.Publickey),
	}
	encoded, error := jsoniter.Marshal(METADATA)
	if error != nil {
		panic(error)
	}
	metadata = encoded
	PUBKEY := Ygg_PUBKEY{
		SignaturePublickey: string(encryption.Publickey),
	}
	encoded, error = jsoniter.Marshal(PUBKEY)
	if error != nil {
		panic(error)
	}
	metakey = encoded
}
func Yggdrasil(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json; charset=utf-8")
	ctx.Response.SetStatusCode(200)
	ctx.Write(metadata)
}
func Yggdrasilpubkey(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json; charset=utf-8")
	ctx.Response.SetStatusCode(200)
	ctx.Write(metadata)
}

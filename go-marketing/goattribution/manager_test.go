package goattribution

import (
	"net/url"
	"reflect"
	"testing"
)

type TestAttributeHandler struct {
}

func (h *TestAttributeHandler) SubChannel() string {
	return "CHANNEL_SUB_TEST"
}

func (h *TestAttributeHandler) Channel() string {
	return "CHANNEL_TEST"
}

func (h *TestAttributeHandler) Match(queryParams url.Values) bool {
	val := queryParams.Get("utm_medium")
	return val == h.Channel()
}

func (h *TestAttributeHandler) Handle(queryParams url.Values) (*AttributeInfo, error) {
	return &AttributeInfo{
		Channel: h.Channel(),
	}, nil
}

func TestAttributeManager_DecryptAttribute(t *testing.T) {
	type args struct {
		appCode string
		referer string
	}
	tests := []struct {
		name    string
		args    args
		want    *AttributeInfo
		wantErr bool
	}{
		//{
		//	name: "userdef.CHANNEL_ORGANIC",
		//	args: args{
		//		appCode: "",
		//		referer: "utm_source=google-play\\u0026utm_medium=organic",
		//	},
		//	want: &AttributeInfo{
		//		Channel: "organic",
		//	},
		//	wantErr: false,
		//},
		{
			name: "userdef.CHANNEL_GOOGLE",
			args: args{
				appCode: "",
				referer: "gclid=123456789&utm_medium=referral&utm_source=apps.facebook.com&utm_campaign=fb4a&utm_content=bytedanceglobal_E.C.P.C&facebook_app_id=",
			},
			want: &AttributeInfo{
				Channel: "google",
			},
			wantErr: false,
		},
		//{
		//	name: "userdef.CHANNEL_FACEBOOK",
		//	args: args{
		//		appCode: "walkgain",
		//		referer: "utm_source=apps.facebook.com\u0026utm_campaign=fb4a\u0026utm_content=%7B%22app%22%3A1442207956662894%2C%22t%22%3A1695907010%2C%22source%22%3A%7B%22data%22%3A%2291ffeebf5c95d0b41f2708122d93d1ea74b3dab359836c3ac70d4762404e9e186518220e888713d5df36ee9b2e53efc1625a404717ad4ff6a8619dac18e4fd6e34da71c91cfb4cc27657ff3b16f503043dbf5b884d77273a0803817ba390f7b0f1828e4d72110637f55eb327fe3df418f0d6fc4ba93bab3f850f47a5142ba79076ba27a816a9e15974c1b2d48681ccfb4a7a8af27c976ff486706ac4c33f7425b03d6fba30a5d5888dd64bc04bcdc42d1769876a93da585d9cc1900437f07120617db0768c220c8d2d594b5894923e42a406f4319bff1c217bc26078558ef2b5f472643fb95c3c58856840aee4219081639f2839c2e6668aff41aebb96156259fb97b96c94f2b9e48b09a6e611e66aa8a63dd4d9c19b8806637837c842714f9ae74cc8078b11470fbd6b69645b20bd97bf421361bc485db48602241578a15b1b86dcc90d926837a881ffc6092592ca9f62454fa3cd8ff5a801dc59ac03647b1b3c7133147767a56f1f42431815c5d007d0c3ed79ca77c4a7c0586161a83811463153f6b3c0e9c4e907c32f2f14c0556d3a307d3cd75e55%22%2C%22nonce%22%3A%22708b3b1a2cb59f5aaa6e10d7%22%7D%7D",
		//	},
		//	want: &AttributeInfo{
		//		Channel:         "meta",
		//		CampaignId:      "23858919684670718",
		//		CampaignName:    "walkgain_cpi_fb_ether_zu_002",
		//		AdCostMode:      "",
		//		CampaignChannel: "fb",
		//		CampaignPartner: "ether",
		//	},
		//	wantErr: false,
		//},
		//{
		//	name: "userdef.CHANNEL_MINTEGRAL",
		//	args: args{
		//		appCode: "walkgain",
		//		referer: "af_tranid\\u003dCtGOjOFjSEE-OXAljqyXHQ\\u0026af_c_id\\u003dss_ads4eachs_android_HappyFruit2048_br_1207\\u0026af_adset_id\\u003d1806894662\\u0026pid\\u003dmintegral_int\\u0026af_prt\\u003dads4eachs\\u0026af_adset\\u003dicon_512x512\\u0026af_ad\\u003dicon_512x512\\u0026af_siteid\\u003dmtg1132839114\\u0026af_ad_id\\u003d1806894662\\u0026c\\u003dads4eachs_happyfruit2048_cpi_mintegral_br_1207",
		//	},
		//	want: &AttributeInfo{
		//		Channel:         "mintegral",
		//		CampaignId:      "ss_ads4eachs_android_HappyFruit2048_br_1207",
		//		CampaignName:    "ads4eachs_happyfruit2048_cpi_mintegral_br_1207_icon_512x512",
		//		AdCostMode:      "cpi",
		//		CampaignChannel: "mintegral",
		//		CampaignPartner: "ads4eachs",
		//	},
		//	wantErr: false,
		//}, {
		//	name: "userdef.CHANNEL_organicnot set",
		//	args: args{
		//		appCode: "walkgain",
		//		referer: "utm_source=(not%20set)&utm_medium=(not%20set)",
		//	},
		//	want: &AttributeInfo{
		//		Channel:         "organic",
		//		CampaignId:      "",
		//		CampaignName:    "",
		//		AdCostMode:      "",
		//		CampaignChannel: "",
		//		CampaignPartner: "",
		//	},
		//	wantErr: false,
		//}, {
		//	name: "userdef.CHANNEL_TEST",
		//	args: args{
		//		appCode: "walkgain",
		//		referer: "utm_source=(not%20set)&utm_medium=CHANNEL_TEST",
		//	},
		//	want: &AttributeInfo{
		//		Channel:         "CHANNEL_TEST",
		//		CampaignId:      "",
		//		CampaignName:    "",
		//		AdCostMode:      "",
		//		CampaignChannel: "",
		//		CampaignPartner: "",
		//	},
		//	wantErr: false,
		//}, {
		//	name: "userdef.CHANNEL_JUMP",
		//	args: args{
		//		appCode: "walkgain",
		//		referer: "utm_source=jump&utm_medium=jump",
		//	},
		//	want: &AttributeInfo{
		//		UtmSource:       "jump",
		//		UtmMedium:       "jump",
		//		Channel:         "jump",
		//		CampaignId:      "",
		//		CampaignName:    "",
		//		AdCostMode:      "",
		//		CampaignChannel: "",
		//		CampaignPartner: "",
		//	},
		//	wantErr: false,
		//},
	}
	Init(Config{
		Name: "default",
		DecryptKeys: map[string]string{
			CHANNEL_META: "key",
		},
	})
	m := GetClient()
	m.AddAttributeHandler(&TestAttributeHandler{})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := m.DecryptAttribute(tt.args.referer)
			if (err != nil) != tt.wantErr {
				t.Errorf("AttributeManager.DecryptAttribute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Channel, tt.want.Channel) {
				t.Errorf("AttributeManager.DecryptAttribute() = %v, want %v", got.Channel, tt.want.Channel)
			}
		})
	}
}

func TestDecryptAttribute(t *testing.T) {

	Init(Config{
		Name: "default",
		DecryptKeys: map[string]string{
			CHANNEL_META: "key",
		},
	})
	m := GetClient()
	type args struct {
		appCode string
		referer string
	}
	tests := []struct {
		name    string
		args    args
		want    *AttributeInfo
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := m.DecryptAttribute(tt.args.referer)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecryptAttribute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecryptAttribute() = %v, want %v", got, tt.want)
			}
		})
	}
}

package management

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/authok/authok-go"
)

func TestBrandingManager_Read(t *testing.T) {
	configureHTTPTestRecordings(t)

	branding, err := api.Branding.Read()
	assert.NoError(t, err)
	assert.IsType(t, &Branding{}, branding)
}

func TestBrandingManager_Update(t *testing.T) {
	configureHTTPTestRecordings(t)

	// Save initial branding settings.
	preTestBrandingSettings, err := api.Branding.Read()
	assert.NoError(t, err)

	expected := &Branding{
		Colors: &BrandingColors{
			Primary: authok.String("#ea5323"),
			PageBackgroundGradient: &BrandingPageBackgroundGradient{
				Type:        authok.String("linear-gradient"),
				Start:       authok.String("#000000"),
				End:         authok.String("#ffffff"),
				AngleDegree: authok.Int(35),
			},
		},
		FaviconURL: authok.String("https://mycompany.org/favicon.ico"),
		LogoURL:    authok.String("https://mycompany.org/logo.png"),
		Font: &BrandingFont{
			URL: authok.String("https://mycompany.org/font.otf"),
		},
	}

	err = api.Branding.Update(expected)
	assert.NoError(t, err)

	actual, err := api.Branding.Read()
	assert.NoError(t, err)
	assert.Equal(t, expected.GetColors().GetPrimary(), actual.GetColors().GetPrimary())
	assert.Equal(t, expected.GetFont().GetURL(), actual.GetFont().GetURL())
	assert.Equal(t, expected.GetFaviconURL(), actual.GetFaviconURL())

	// Restore initial branding settings.
	err = api.Branding.Update(preTestBrandingSettings)
	assert.NoError(t, err)
}

func TestBrandingManager_UniversalLogin(t *testing.T) {
	configureHTTPTestRecordings(t)

	givenACustomDomain(t)

	body := `<!DOCTYPE html><html><head>{%- authok:head -%}</head><body>{%- authok:widget -%}</body></html>`
	expectedUL := &BrandingUniversalLogin{
		Body: authok.String(body),
	}

	err := api.Branding.SetUniversalLogin(expectedUL)
	assert.NoError(t, err)

	actualUL, err := api.Branding.UniversalLogin()
	assert.NoError(t, err)
	assert.Equal(t, expectedUL, actualUL)

	t.Cleanup(func() {
		err = api.Branding.DeleteUniversalLogin()
		assert.NoError(t, err)
	})
}

func TestBrandingColors(t *testing.T) {
	var testCases = []struct {
		name   string
		colors *BrandingColors
		expect string
	}{
		{
			name: "PageBackground",
			colors: &BrandingColors{
				Primary:        authok.String("#ea5323"),
				PageBackground: authok.String("#000000"),
			},
			expect: `{"primary":"#ea5323","page_background":"#000000"}`,
		},
		{
			name: "PageBackgroundGradient",
			colors: &BrandingColors{
				Primary: authok.String("#ea5323"),
				PageBackgroundGradient: &BrandingPageBackgroundGradient{
					Type:        authok.String("linear-gradient"),
					Start:       authok.String("#000000"),
					End:         authok.String("#ffffff"),
					AngleDegree: authok.Int(35),
				},
			},
			expect: `{"primary":"#ea5323","page_background":{"type":"linear-gradient","start":"#000000","end":"#ffffff","angle_deg":35}}`,
		},
		{
			name: "PageBackgroundNil",
			colors: &BrandingColors{
				Primary: authok.String("#ea5323"),
			},
			expect: `{"primary":"#ea5323"}`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			b, err := json.Marshal(testCase.colors)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expect, string(b))

			var colors BrandingColors
			err = json.Unmarshal([]byte(testCase.expect), &colors)
			assert.NoError(t, err)
			assert.Equal(t, testCase.colors, &colors)
		})
	}
}

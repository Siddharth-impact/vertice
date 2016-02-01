package alerts

import (
	"os"

	"gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { check.TestingT(t) }

type S struct{}

var _ = check.Suite(&S{})

var ak = os.Getenv("NIL_MAILGUN_APIKEY")

func (s *S) TestMailgunOnboard(c *check.C) {
	if ak == "" {
		c.Skip("-Mailgun (api_key) not provided")
	}
	c.Assert(len(ak) > 0, check.Equals, true)
	m := NewMailgun(map[string]string{"api_key": ak, "domain": "ojamail.megambox.com"})
	c.Assert(m, check.NotNil)
	err := m.Notify(ONBOARD, map[string]string{
		"email":     "nkishore@megam.io",
		"logo":      "vertice.png",
		"nilavu":    "console.megam.io",
		"click_url": "https://console.megam.io/reset?email=test@email.com&resettoken=9090909090",
		"days":      "20",
		"cost":      "$12",
	})
	c.Assert(err, check.IsNil)
}

func (s *S) TestMailgunReset(c *check.C) {
	if ak == "" {
		c.Skip("-Mailgun (api_key) not provided")
	}
	c.Assert(len(ak) > 0, check.Equals, true)
	m := NewMailgun(map[string]string{"api_key": ak, "domain": "ojamail.megambox.com"})
	c.Assert(m, check.NotNil)
	err := m.Notify(RESET, map[string]string{
		"email":     "nkishore@megam.io",
		"logo":      "vertice.png",
		"nilavu":    "console.megam.io",
		"click_url": "https://console.megam.io/reset?email=test@email.com&resettoken=9090909090",
		"days":      "20",
		"cost":      "$12",
	})
	c.Assert(err, check.IsNil)
}

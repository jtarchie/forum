package templates_test

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/jtarchie/forum/services"
	"github.com/jtarchie/forum/templates"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTemplates(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Templates Suite")
}

var _ = Describe("ListForums", func() {
	It("renders a list of forums, with children nested", func() {
		forums := services.Forums{
			services.Forum{
				Description: "A full description",
				ID:          1,
				Name:        "Parent 1",
				ParentID:    0,
			},
			services.Forum{
				Description: "A full description",
				ID:          3,
				Name:        "Parent 2",
				ParentID:    0,
			},
			services.Forum{
				Description: "A full description",
				ID:          2,
				Name:        "Child 1 of 1",
				ParentID:    1,
			},
		}

		output := templates.ListForums(forums)

		doc, err := goquery.NewDocumentFromReader(strings.NewReader(output))
		Expect(err).ToNot(HaveOccurred())

		titles := doc.Find("article header").Map(func(i int, s *goquery.Selection) string {
			return s.Text()
		})

		Expect(titles).To(HaveLen(3))
		Expect(titles).To(Equal([]string{"Parent 1", "Child 1 of 1", "Parent 2"}))

		titles = doc.Find("article > article header").Map(func(i int, s *goquery.Selection) string {
			return s.Text()
		})

		Expect(titles).To(HaveLen(1))
		Expect(titles).To(Equal([]string{"Child 1 of 1"}))
	})
})

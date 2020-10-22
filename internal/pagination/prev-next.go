// ORIGINAL: java/PagingLinksFinder.java

package pagination

import (
	nurl "net/url"
	"path"
	"strconv"
	"strings"

	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/internal/domutil"
	"github.com/markusmobius/go-domdistiller/internal/model"
	"github.com/markusmobius/go-domdistiller/internal/stringutil"
	"golang.org/x/net/html"
)

type pagingLinkScore struct {
	linkText  string
	linkHref  string
	linkIndex int
	score     int
}

// PrevNextFinder finds the next and previous page links for the distilled document. The functionality
// for next page links is migrated from readability.getArticleTitle() in chromium codebase's
// third_party/readability/js/readability.js, and then expanded for previous page links; boilerpipe
// doesn't have such capability.
// First, it determines the prefix URL of the document. Then, for each anchor in the document, its
// href and text are compared to the prefix URL and examined for next- or previous-paging-related
// information. If it passes, its score is then determined by applying various heuristics on its
// href, text, class name and ID, Lastly, the page link with the highest score of at least 50 is
// considered to have enough confidence as the next or previous page link.
type PrevNextFinder struct{}

func NewPrevNextFinder() *PrevNextFinder {
	return &PrevNextFinder{}
}

func (pnf *PrevNextFinder) FindPagination(root *html.Node, pageURL *nurl.URL) model.PaginationInfo {
	return model.PaginationInfo{
		PrevPage: pnf.FindOutlink(root, pageURL, false),
		NextPage: pnf.FindOutlink(root, pageURL, true),
	}
}

func (pnf *PrevNextFinder) FindOutlink(root *html.Node, pageURL *nurl.URL, findNext bool) string {
	// Clean up URL
	tmp, _ := nurl.Parse(pageURL.String())
	tmp.RawQuery = ""
	tmp.Fragment = ""
	tmp.RawFragment = ""

	// Remove trailing '/' from window location href, because it'll be used to compare with
	// other href's whose trailing '/' are also removed.
	tmp.Path = strings.TrimSuffix(tmp.Path, "/")
	tmp.RawPath = tmp.Path
	currentURL := stringutil.UnescapedString(tmp)

	// Create folder URL
	tmp.Path = strings.TrimSuffix(path.Dir(tmp.Path), "/")
	tmp.RawPath = tmp.Path
	folderURL := stringutil.UnescapedString(tmp)

	// Create allowed prefix
	// The trailing "/" is essential to ensure the whole hostname is matched, and not just the
	// prefix of the hostname. It also maintains the requirement of having a "path" in the URL.
	tmp.Path = "/"
	tmp.RawPath = tmp.Path
	allowedPrefix := stringutil.UnescapedString(tmp)
	lenPrefix := len(allowedPrefix)

	// Loop through all links, looking for hints that they may be next- or previous- page links.
	// Things like having "page" in their textContent, className or id, or being a child of a
	// node with a page-y className or id.
	// Also possible: levenshtein distance? longest common subsequence?
	// After we do that, assign each page a score.
	bannedURLs := make(map[string]struct{})
	candidates := make([]pagingLinkScore, 0)
	for i, link := range dom.GetElementsByTagName(root, "a") {
		// Try to convert relative URL in link href to absolute URL
		linkHref := dom.GetAttribute(link, "href")
		linkHref = stringutil.CreateAbsoluteURL(linkHref, pageURL)

		// Make sure the link href is absolute
		_, err := nurl.ParseRequestURI(linkHref)
		if err != nil {
			continue
		}

		// Make sure the href is related with current page
		if !stringutil.HasPrefixIgnoreCase(linkHref, allowedPrefix) {
			continue
		}

		if findNext && !rxNumber.MatchString(linkHref[lenPrefix:]) {
			continue
		}

		// In original dom-distiller they skip invisible links, but we can't do that here
		// since Go can't compute stylesheet. NEED-COMPUTE-CSS.

		// Remove url anchor and then trailing '/' from link's href.
		tmp, _ := nurl.Parse(linkHref)
		tmp.RawQuery = ""
		tmp.Fragment = ""
		tmp.RawFragment = ""
		tmp.Path = strings.TrimSuffix(tmp.Path, "/")
		tmp.RawPath = tmp.Path
		linkHref = stringutil.UnescapedString(tmp)

		// Ignore page link that is the same as current window location.
		// If the page link is same as the folder URL:
		// - next page link: ignore it, since we would already have seen it.
		// - previous page link: don't ignore it, since some sites will simply have the same
		//                       folder URL for the first page.
		if stringutil.EqualsIgnoreCase(linkHref, currentURL) ||
			(findNext && stringutil.EqualsIgnoreCase(linkHref, folderURL)) {
			continue
		}

		// Get link text using inner text
		linkText := domutil.InnerText(link)
		linkText = strings.TrimSpace(linkText)

		// If the linkText looks like it's not the next or previous page, skip it.
		if len(linkText) > 25 {
			continue
		}

		// If the linkText contains banned text, skip it, and also ban other anchors with the
		// same link URL.
		if rxExtraneous.MatchString(linkText) {
			bannedURLs[linkHref] = struct{}{}
			continue
		}

		// For next page link, if the initial part of the URL is identical to the folder URL, but
		// the rest of it doesn't contain any digits, it's certainly not a next page link.
		// However, this doesn't apply to previous page link, because most sites will just have
		// the folder URL for the first page.
		// TODO(kuan): do we need to apply this heuristic to previous page links if current page
		// number is not 2?
		if findNext {
			remainingLinkHref := linkHref
			if strings.HasPrefix(linkHref, folderURL) {
				remainingLinkHref = linkHref[len(folderURL):]
			}

			if !rxNumber.MatchString(remainingLinkHref) {
				continue
			}
		}

		// Prepare link score
		linkObj := pagingLinkScore{
			linkIndex: i,
			linkText:  linkText,
			linkHref:  linkHref,
			score:     0,
		}

		// If the folder URL isn't part of this URL, penalize this link.  It could still be the
		// link, but the odds are lower.
		// Example: http://www.actionscript.org/resources/articles/745/1/JavaScript-and-VBScript-Injection-in-ActionScript-3/Page1.html.
		if !strings.HasPrefix(linkHref, folderURL) {
			linkObj.score -= 25
		}

		// Concatenate the link text with class name and id, and determine the score based on
		// existence of various paging-related words.
		linkData := linkText + " " + dom.GetAttribute(link, "class") + " " + dom.GetAttribute(link, "id")

		if (findNext && rxNextLink.MatchString(linkData)) ||
			(!findNext && rxPrevLink.MatchString(linkData)) {
			linkObj.score += 50
		}

		if rxPagination.MatchString(linkData) {
			linkObj.score += 25
		}

		if rxFirstLast.MatchString(linkData) {
			// -65 is enough to negate any bonuses gotten from a > or » in the text.
			// If we already matched on "next", last is probably fine.
			// If we didn't, then it's bad.  Penalize.
			// Same for "prev".
			if (findNext && !rxNextLink.MatchString(linkText)) ||
				(!findNext && !rxPrevLink.MatchString(linkText)) {
				linkObj.score -= 65
			}
		}

		if rxNegative.MatchString(linkData) || rxExtraneous.MatchString(linkData) {
			linkObj.score -= 50
		}

		if (findNext && rxPrevLink.MatchString(linkData)) ||
			(!findNext && rxNextLink.MatchString(linkData)) {
			linkObj.score -= 200
		}

		// Check if a parent element contains page or paging or paginate.
		positiveMatch := false
		negativeMatch := false
		parent := domutil.GetParentElement(link)
		for parent != nil && (positiveMatch == false || negativeMatch == false) {
			parentData := dom.GetAttribute(parent, "class") + " " + dom.GetAttribute(parent, "id")
			if !positiveMatch && rxPagination.MatchString(parentData) {
				linkObj.score += 25
				positiveMatch = true
			}

			// TODO(kuan): to get 1st page for prev page link, this can't be applied; however,
			// the non-application might be the cause of recursive prev page being returned,
			// i.e. for page 1, it may incorrectly return page 3 for prev page link.
			if !negativeMatch && rxNegative.MatchString(parentData) {
				// If this is just something like "footer", give it a negative.
				// If it's something like "body-and-footer", leave it be.
				if !rxPositive.MatchString(parentData) {
					linkObj.score -= 25
					negativeMatch = true
				}
			}

			parent = domutil.GetParentElement(parent)
		}

		// If the URL looks like it has paging in it, add to the score.
		// Things like /page/2/, /pagenum/2, ?p=3, ?page=11, ?pagination=34.
		if rxLinkPagination.MatchString(linkHref) || rxPagination.MatchString(linkHref) {
			linkObj.score += 25
		}

		// If the URL contains negative values, give a slight decrease.
		if rxExtraneous.MatchString(linkHref) {
			linkObj.score -= 15
		}

		// If the link text is too long, penalize the link.
		if len(linkText) > 10 {
			linkObj.score -= len(linkText)
		}

		// If the link text can be parsed as a number, give it a minor bonus, with a slight bias
		// towards lower numbered pages.  This is so that pages that might not have 'next' in
		// their text can still get scored, and sorted properly by score.
		// TODO(kuan): it might be wrong to assume that it knows about other pages in the
		// document and that it starts on the first page.
		linkTextAsNumber, _ := strconv.Atoi(linkText)
		if linkTextAsNumber > 0 {
			// Punish 1 since we're either already there, or it's probably before what we
			// want anyway.
			if linkTextAsNumber == 1 {
				linkObj.score -= 10
			} else {
				additionalScore := 10 - linkTextAsNumber
				if additionalScore < 0 {
					additionalScore = 0
				}
				linkObj.score += additionalScore
			}
		}

		diff, valid := pnf.getPageDiff(currentURL, linkHref, len(allowedPrefix))
		if valid {
			if (findNext && diff == 1) || (!findNext && diff == -1) {
				linkObj.score += 25
			}
		}

		// Add final score to candidates
		candidates = append(candidates, linkObj)
	} // loop for all links

	// Loop through all of the possible pages from above and find the top candidate for the next
	// page URL. Require at least a score of 50, which is a relatively high confidence that
	// this page is the next link.
	var topPage pagingLinkScore
	for _, pageObj := range candidates {
		if _, exist := bannedURLs[pageObj.linkHref]; exist {
			continue
		}

		if pageObj.score >= 50 && topPage.score < pageObj.score {
			topPage = pageObj
		}
	}

	return topPage.linkHref
}

func (pnf *PrevNextFinder) getPageDiff(pageURL, linkHref string, skip int) (int, bool) {
	maxLimit := len(pageURL)
	if lenHref := len(linkHref); lenHref < maxLimit {
		maxLimit = lenHref
	}

	var commonLen int
	for i := skip; i < maxLimit; i++ {
		if pageURL[i] != linkHref[i] {
			commonLen = i
			break
		}
	}

	var urlAsNumber int
	if str := rxNumberAtStart.FindString(pageURL[commonLen:]); str != "" {
		urlAsNumber, _ = strconv.Atoi(str)
	}

	var linkAsNumber int
	if str := rxNumberAtStart.FindString(linkHref[commonLen:]); str != "" {
		linkAsNumber, _ = strconv.Atoi(str)
	}

	if urlAsNumber > 0 && linkAsNumber > 0 {
		return linkAsNumber - urlAsNumber, true
	}

	return 0, false
}

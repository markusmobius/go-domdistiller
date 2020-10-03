// ORIGINAL: java/OpenGraphProtocolParser.java

package opengraph

import "strings"

type Prefix uint

const (
	OG Prefix = iota
	Profile
	Article
)

type PrefixNameList map[Prefix]string

func (prefixes PrefixNameList) addObjectType(prefix, objType string) {
	if objType == "" {
		prefixes[OG] = prefix
		return
	}

	objType = strings.TrimPrefix(objType, "/")
	if objType == ProfileObjtype {
		prefixes[Profile] = prefix
		return
	}

	if objType == ArticleObjtype {
		prefixes[Article] = prefix
	}
}

func (prefixes PrefixNameList) setDefault() {
	// For any unspecified prefix, use common ones:
	// - "og": http://ogp.me/ns#
	// - "profile": http://ogp.me/ns/profile#
	// - "article": http://ogp.me/ns/article#.
	if _, exist := prefixes[OG]; !exist {
		prefixes[OG] = "og"
	}

	if _, exist := prefixes[Profile]; !exist {
		prefixes[Profile] = ProfileObjtype
	}

	if _, exist := prefixes[Article]; !exist {
		prefixes[Article] = ArticleObjtype
	}
}

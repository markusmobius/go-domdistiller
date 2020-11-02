# Improvements

After using both Readability.js and DOM Distiller, we found that there are several improvements that can be implemented into this port. Besides that, from our experiments we also found some possible bugs that we decided to fix.

These so-called improvements are listed here as historical documentation and to explain the difference between the main branch and stable branch.

## From Readability

- Implement function to check if a HTML element is probably visible or not. This is especially useful since one of the DOM Distiller strategy is to exclude invisible elements by computing the stylesheets (which is impossible to do in Go).
- Exclude form and input element, since in distilled mode we only want to read.
- Skip byline, empty div and unlikely elements by checking its class name, id and role attributes.
- Convert anchors with Javascript URL into an ordinary text node.
- Convert font to span elements. This is done because the font elements is usually only used for styling, so Readability.js decided to convert it.
- Exclude identification and presentational attributes (eg. `id`, `class` and `style`) from each elements.

## From our own experiments

- Make sure figure's caption doesn't contains noscript elements. This is done because noscript in Go is a bit weird, sometimes it detected as HTML element while the other times it detected as plain text, so we need additional schecks to clean it.
- Mark large blocks around main content's tag level as content as well. In original DOM Distiller, they are looking for the most likely main content, then they mark text blocks that exist in the same tag level of the main content as content as well. Unfortunately, we found out that in some sites parts of the article are omitted by DOM Distiller. To fix this, we decided to make the filter more tolerant by checking text blocks in lower and upper tag levels as well.
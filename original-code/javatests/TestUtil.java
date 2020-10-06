// Copyright 2014 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package org.chromium.distiller;

import com.google.gwt.dom.client.AnchorElement;
import com.google.gwt.dom.client.Document;
import com.google.gwt.dom.client.Element;
import com.google.gwt.dom.client.HeadingElement;
import com.google.gwt.dom.client.IFrameElement;
import com.google.gwt.dom.client.ImageElement;
import com.google.gwt.dom.client.MetaElement;
import com.google.gwt.dom.client.Node;
import com.google.gwt.dom.client.Text;
import com.google.gwt.dom.client.TitleElement;
import com.google.gwt.dom.client.NodeList;
import com.google.gwt.user.client.Random;
import com.google.gwt.user.client.Window;

import java.util.ArrayList;
import java.util.Collections;
import java.util.List;

/**
 * A mixed bag of stuff used in tests.
 */
public class TestUtil {
    public static Text createText(String value) {
        return Document.get().createTextNode(value);
    }

    public static ImageElement createImage() {
        return Document.get().createImageElement();
    }

    public static IFrameElement createIframe() {
        return Document.get().createIFrameElement();
    }

    public static Element createSpan(String value) {
        Element s = Document.get().createElement("SPAN");
        s.setInnerHTML(value);
        return s;
    }

    public static Element createParagraph(String value) {
        Element s = Document.get().createElement("P");
        s.setInnerHTML(value);
        return s;
    }

    public static Element createListItem(String value) {
        Element s = Document.get().createElement("LI");
        s.setInnerText(value);
        return s;
    }

    public static String getElementAsString(Element e) {
        Element div = Document.get().createDivElement();
        div.appendChild(e.cloneNode(true));
        return div.getInnerHTML();
    }

    public static String formHrefWithWindowLocationPath(String strToAppend) {
        String noUrlParams = Window.Location.getPath();
        // Append '/' if necessary.
        if (!strToAppend.isEmpty() && !StringUtil.match(noUrlParams, "\\/$")) {
            noUrlParams += "/";
        }
        return noUrlParams + strToAppend;
    }

    public static String removeAllDirAttributes(String originalHtml) {
        return originalHtml.replaceAll(" dir=\\\"(ltr|rtl|inherit|auto)\\\"","");
    }

    public static List<Node> nodeListToList(NodeList nodeList) {
        List<Node> nodes = new ArrayList<>();
        for (int i = 0; i < nodeList.getLength(); i++) {
            nodes.add(nodeList.getItem(i));
        }
        return nodes;
    }

    /**
     * Randomly shuffle the list in-place.
     */
    public static void shuffle(List<?> list) {
        int size = list.size();
        for (int i=size; i>1; i--) {
            Collections.swap(list, i-1, Random.nextInt(i));
        }
    }
}

// Copyright 2014 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package org.chromium.distiller.webdocument;

import org.chromium.distiller.DomDistillerJsTestCase;
import org.chromium.distiller.TestTextDocumentBuilder;
import org.chromium.distiller.document.TextBlock;
import org.chromium.distiller.document.TextDocument;

import com.google.gwt.dom.client.AnchorElement;
import com.google.gwt.dom.client.Document;
import com.google.gwt.dom.client.Element;
import com.google.gwt.dom.client.Node;
import com.google.gwt.dom.client.Text;

import java.util.List;

public class WebDocumentBuilderTest extends DomDistillerJsTestCase {
    public void testRegression0() {
        String html = "<blockquote><p>“There are plenty of instances where provocation comes into" +
            " consideration, instigation comes into consideration, and I will be on the record" +
            " right here on national television and say that I am sick and tired of men" +
            " constantly being vilified and accused of things and we stop there,”" +
            " <a href=\"http://deadspin.com/i-do-not-believe-women-provoke-violence-says-stephen" +
            "-a-1611060016\" target=\"_blank\">Smith said.</a>  “I’m saying, “Can we go a step" +
            " further?” Since we want to dig all deeper into Chad Johnson, can we dig in deep" +
            " to her?”</p></blockquote>";
        Element div = Document.get().createDivElement();
        mBody.appendChild(div);
        div.setInnerHTML(html);
        TextDocument document = TestTextDocumentBuilder.fromPage(div);
        List<TextBlock> textBlocks = document.getTextBlocks();
        assertEquals(1, textBlocks.size());
        TextBlock tb = textBlocks.get(0);
        assertEquals(74, tb.getNumWords());
        assertTrue(0.1 > tb.getLinkDensity());
    }

    public void testRegression1() {
        String html = "<p>\n"
                + "<a href=\"example\" target=\"_top\"><u>More news</u></a> | \n"
                + "<a href=\"example\" target=\"_top\"><u>Search</u></a> | \n"
                + "<a href=\"example\" target=\"_top\"><u>Features</u></a> | \n"
                + "<a href=\"example\" target=\"_top\"><u>Blogs</u></a> | \n"
                + "<a href=\"example\" target=\"_top\"><u>Horse Health</u></a> | \n"
                + "<a href=\"example\" target=\"_top\"><u>Ask the Experts</u></a> | \n"
                + "<a href=\"example\" target=\"_top\"><u>Horse Breeding</u></a> | \n"
                + "<a href=\"example\" target=\"_top\"><u>Forms</u></a> | \n"
                + "<a href=\"example\" target=\"_top\"><u>Home</u></a> </p>\n";
        Element div = Document.get().createDivElement();
        mBody.appendChild(div);
        div.setInnerHTML(html);
        TextDocument document = TestTextDocumentBuilder.fromPage(div);
        List<TextBlock> textBlocks = document.getTextBlocks();
        assertEquals(1, textBlocks.size());
        TextBlock tb = textBlocks.get(0);
        assertEquals(14, tb.getNumWords());
        assertEquals(1.0, tb.getLinkDensity(), 0.01);
    }
}

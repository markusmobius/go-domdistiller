// Copyright 2014 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package org.chromium.distiller;

import org.chromium.distiller.webdocument.WebTable;

import com.google.gwt.core.client.JsArray;
import com.google.gwt.dom.client.Document;
import com.google.gwt.dom.client.Element;
import com.google.gwt.dom.client.ImageElement;
import com.google.gwt.dom.client.Node;
import com.google.gwt.dom.client.NodeList;

import java.util.Map;
import java.util.List;

public class DomUtilTest extends DomDistillerJsTestCase {
    public void testGetAttributes() {
        Element e = Document.get().createDivElement();
        e.setInnerHTML("<div style=\"width:50px; height:100px\" id=\"f\" class=\"sdf\"></div>");
        e = Element.as(e.getChildNodes().getItem(0));
        JsArray<Node> jsAttrs = DomUtil.getAttributes(e);
        assertEquals(3, jsAttrs.length());
        assertEquals("style", jsAttrs.get(0).getNodeName());
        assertEquals("width:50px; height:100px", jsAttrs.get(0).getNodeValue());
        assertEquals("id", jsAttrs.get(1).getNodeName());
        assertEquals("f", jsAttrs.get(1).getNodeValue());
        assertEquals("class", jsAttrs.get(2).getNodeName());
        assertEquals("sdf", jsAttrs.get(2).getNodeValue());
    }

    public void testGetFirstElementWithClassName() {
        Element rootDiv = TestUtil.createDiv(0);

        Element div1 = TestUtil.createDiv(1);
        div1.addClassName("abcd");
        rootDiv.appendChild(div1);

        Element div2 = TestUtil.createDiv(2);
        div2.addClassName("test");
        div2.addClassName("xyz");
        rootDiv.appendChild(div2);

        Element div3 = TestUtil.createDiv(2);
        div3.addClassName("foobar foo");
        rootDiv.appendChild(div3);

        assertEquals(div1, DomUtil.getFirstElementWithClassName(rootDiv, "abcd"));
        assertEquals(div2, DomUtil.getFirstElementWithClassName(rootDiv, "test"));
        assertEquals(div2, DomUtil.getFirstElementWithClassName(rootDiv, "xyz"));
        assertEquals(null, DomUtil.getFirstElementWithClassName(rootDiv, "bc"));
        assertEquals(null, DomUtil.getFirstElementWithClassName(rootDiv, "t xy"));
        assertEquals(null, DomUtil.getFirstElementWithClassName(rootDiv, "tes"));
        assertEquals(div3, DomUtil.getFirstElementWithClassName(rootDiv, "foo"));
    }

    public void testHasRootDomain() {
        // Positive tests.
        assertTrue(DomUtil.hasRootDomain("http://www.foo.bar/foo/bar.html", "foo.bar"));
        assertTrue(DomUtil.hasRootDomain("https://www.m.foo.bar/foo/bar.html", "foo.bar"));
        assertTrue(DomUtil.hasRootDomain("https://www.m.foo.bar/foo/bar.html", "www.m.foo.bar"));
        assertTrue(DomUtil.hasRootDomain("http://localhost/foo/bar.html", "localhost"));
        assertTrue(DomUtil.hasRootDomain("https://www.m.foo.bar.baz", "foo.bar.baz"));
        // Negative tests.
        assertFalse(DomUtil.hasRootDomain("https://www.m.foo.bar.baz", "x.foo.bar.baz"));
        assertFalse(DomUtil.hasRootDomain("https://www.foo.bar.baz", "foo.bar"));
        assertFalse(DomUtil.hasRootDomain("http://foo", "m.foo"));
        assertFalse(DomUtil.hasRootDomain("https://www.badfoobar.baz", "foobar.baz"));
        assertFalse(DomUtil.hasRootDomain("", "foo"));
        assertFalse(DomUtil.hasRootDomain("http://foo.bar", ""));
        assertFalse(DomUtil.hasRootDomain(null, "foo"));
        assertFalse(DomUtil.hasRootDomain("http://foo.bar", null));
    }

    public void testSplitUrlParams() {
        Map<String, String> result = DomUtil.splitUrlParams("param1=apple&param2=banana");
        assertEquals(2, result.size());
        assertEquals("apple", result.get("param1"));
        assertEquals("banana", result.get("param2"));

        result = DomUtil.splitUrlParams("123=abc");
        assertEquals(1, result.size());
        assertEquals("abc", result.get("123"));

        result = DomUtil.splitUrlParams("");
        assertEquals(0, result.size());

        result = DomUtil.splitUrlParams(null);
        assertEquals(0, result.size());
    }

    public void testNodeDepth() {
        Element div = TestUtil.createDiv(1);

        Element div2 = TestUtil.createDiv(2);
        div.appendChild(div2);

        Element div3 = TestUtil.createDiv(3);
        div2.appendChild(div3);

        assertEquals(2, DomUtil.getNodeDepth(div3));
    }

    public void testZeroOrNoNodeDepth() {
        Element div = TestUtil.createDiv(0);
        assertEquals(0, DomUtil.getNodeDepth(div));
        assertEquals(-1, DomUtil.getNodeDepth(null));
    }

    public void testIsVisibleByOffsetParentDisplayNone() {
        String html =
            "<div style=\"display: none;\">" +
                "<div></div>" +
            "</div>";
        mBody.setInnerHTML(html);
        Element child = mBody.getFirstChildElement().getFirstChildElement();
        assertFalse(DomUtil.isVisibleByOffset(child));
    }

    public void testIsVisibleByOffsetChildDisplayNone() {
        String html =
            "<div>" +
                "<div style=\"display: none;\"></div>" +
            "</div>";
        mBody.setInnerHTML(html);
        Element child = mBody.getFirstChildElement().getFirstChildElement();
        assertFalse(DomUtil.isVisibleByOffset(child));
    }

    public void testIsVisibleByOffsetDisplayBlock() {
        String html =
            "<div>" +
                "<div></div>" +
            "</div>";
        mBody.setInnerHTML(html);
        Element child = mBody.getFirstChildElement().getFirstChildElement();
        assertTrue(DomUtil.isVisibleByOffset(child));
    }

    public void testOnlyProcessArticleElement() {
        final String htmlArticle =
            "<h1></h1>" +
            "<article>a</article>";

        String expected = "<article>a</article>";

        Element result = getArticleElement(htmlArticle);
        assertEquals(expected, result.getString());
    }

    public void testOnlyProcessArticleElementWithHiddenArticleElement() {
        final String htmlArticle =
            "<h1></h1>" +
            "<article>a</article>" +
            "<article style=\"display:none\">b</article>";

        String expected = "<article>a</article>";

        Element result = getArticleElement(htmlArticle);
        assertEquals(expected, result.getString());
    }

    public void testOnlyProcessArticleElementWithZeroAreaElement() {
        final String htmlArticle =
                "<h1></h1>" +
                        "<article>a</article>" +
                        "<article style=\"width: 0px\">b</article>";

        String expected = "<article>a</article>";

        Element result = getArticleElement(htmlArticle);
        assertEquals(expected, result.getString());
    }

    public void testOnlyProcessArticleElementMultiple() {
        final String htmlArticle =
            "<h1></h1>" +
            "<article>a</article>" +
            "<article>b</article>";

        // The existence of multiple articles disables the fast path.
        assertNull(getArticleElement(htmlArticle));
    }

    public void testOnlyProcessSchemaOrgArticle() {
        final String htmlArticle =
            "<h1></h1>" +
            "<div itemscope itemtype=\"http://schema.org/Article\">a" +
            "</div>";

        final String expected =
            "<div itemscope=\"\" " +
                "itemtype=\"http://schema.org/Article\">a" +
            "</div>";

        Element result = getArticleElement(htmlArticle);
        assertEquals(expected, result.getString());
    }

    public void testOnlyProcessSchemaOrgArticleWithHiddenArticleElement() {
        final String htmlArticle =
            "<h1></h1>" +
            "<div itemscope itemtype=\"http://schema.org/Article\">a" +
            "</div>" +
            "<div itemscope itemtype=\"http://schema.org/Article\" " +
                "style=\"display:none\">b" +
            "</div>";

        String expected =
            "<div itemscope=\"\" itemtype=\"http://schema.org/Article\">a" +
            "</div>";

        Element result = getArticleElement(htmlArticle);
        assertEquals(expected, result.getString());
    }

    public void testOnlyProcessSchemaOrgArticleNews() {
        final String htmlArticle =
            "<h1></h1>" +
            "<div itemscope itemtype=\"http://schema.org/NewsArticle\">a" +
            "</div>";

        final String expected =
            "<div itemscope=\"\" " +
                "itemtype=\"http://schema.org/NewsArticle\">a" +
            "</div>";

        Element result = getArticleElement(htmlArticle);
        assertEquals(expected, result.getString());
    }

    public void testOnlyProcessSchemaOrgArticleBlog() {
        final String htmlArticle =
            "<h1></h1>" +
            "<div itemscope itemtype=\"http://schema.org/BlogPosting\">a" +
            "</div>";

        final String expected =
            "<div itemscope=\"\" " +
                "itemtype=\"http://schema.org/BlogPosting\">a" +
            "</div>";

        Element result = getArticleElement(htmlArticle);
        assertEquals(expected, result.getString());
    }

    public void testOnlyProcessSchemaOrgPostal() {
        final String htmlArticle =
            "<h1></h1>" +
            "<div itemscope itemtype=\"http://schema.org/PostalAddress\">a" +
            "</div>";

        Element result = getArticleElement(htmlArticle);
        assertNull(result);
    }

    public void testOnlyProcessSchemaOrgArticleNested() {
        final String htmlArticle =
            "<h1></h1>" +
            "<div itemscope itemtype=\"http://schema.org/Article\">a" +
                "<div itemscope itemtype=\"http://schema.org/Article\">b" +
                "</div>" +
            "</div>";

        final String expected =
            "<div itemscope=\"\" itemtype=\"http://schema.org/Article\">a" +
                "<div itemscope=\"\" itemtype=\"http://schema.org/Article\">b" +
               "</div>" +
            "</div>";

        Element result = getArticleElement(htmlArticle);
        assertEquals(expected, result.getString());
    }

    public void testOnlyProcessSchemaOrgArticleNestedWithNestedHiddenArticleElement() {
        final String htmlArticle =
            "<h1></h1>" +
            "<div itemscope itemtype=\"http://schema.org/Article\">a" +
                "<div itemscope itemtype=\"http://schema.org/Article\">b" +
                "</div>" +
                "<div itemscope itemtype=\"http://schema.org/Article\" " +
                    "style=\"display:none\">c" +
                "</div>" +
            "</div>";

        final String expected =
            "<div itemscope=\"\" itemtype=\"http://schema.org/Article\">a" +
                "<div itemscope=\"\" itemtype=\"http://schema.org/Article\">b" +
                "</div>" +
                "<div itemscope=\"\" itemtype=\"http://schema.org/Article\" " +
                    "style=\"display:none\">c" +
                "</div>" +
            "</div>";

        Element result = getArticleElement(htmlArticle);
        assertEquals(expected, result.getString());
    }

    public void testOnlyProcessSchemaOrgArticleNestedWithHiddenArticleElement() {
        final String paragraph = "<p></p>";

        final String htmlArticle =
            "<h1></h1>" +
            "<div itemscope itemtype=\"http://schema.org/Article\">a" +
                "<div itemscope itemtype=\"http://schema.org/Article\">b" +
                "</div>" +
            "</div>" +
            "<div itemscope itemtype=\"http://schema.org/Article\" " +
                "style=\"display:none\">c" +
            "</div>";

        final String expected =
            "<div itemscope=\"\" itemtype=\"http://schema.org/Article\">a" +
                "<div itemscope=\"\" itemtype=\"http://schema.org/Article\">b" +
                "</div>" +
            "</div>";

        Element result = getArticleElement(htmlArticle);
        assertEquals(expected, result.getString());
    }

    public void testOnlyProcessSchemaOrgNonArticleMovie() {
        final String htmlArticle =
            "<h1></h1>" +
            "<div itemscope itemtype=\"http://schema.org/Movie\">a" +
            "</div>";

        // Non-article schema.org types should not use the fast path.
        Element result = getArticleElement(htmlArticle);
        assertNull(result);
    }

    private Element getArticleElement(String html) {
        mBody.setInnerHTML(html);
        return DomUtil.getArticleElement(mRoot);
    }

    public void testGetArea() {
        String elements =
            "<div style=\"width: 200px; height: 100px\">w</div>" +
            "<div style=\"width: 300px;\">" +
                "<div style=\"width: 300px; height: 200px\"></div>" +
            "</div>" +
            "<div style=\"width: 400px; height: 100px\">" +
                "<div style=\"height: 100%\"></div>" +
            "</div>";
        mBody.setInnerHTML(elements);

        Element element = mBody.getFirstChildElement();
        assertEquals(200*100, DomUtil.getArea(element));

        element = element.getNextSiblingElement();
        assertEquals(300*200, DomUtil.getArea(element));

        element = element.getNextSiblingElement();
        assertEquals(400*100, DomUtil.getArea(element));

        element = element.getFirstChildElement();
        assertEquals(400*100, DomUtil.getArea(element));
    }
}

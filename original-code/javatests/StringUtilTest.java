// Copyright 2014 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package org.chromium.distiller;

import com.google.gwt.regexp.shared.RegExp;

import org.chromium.distiller.StringUtil.WordCounter;

import java.util.ArrayList;
import java.util.List;

public class StringUtilTest extends JsTestCase {
    public void testCountWords() {
        StringUtil.setWordCounter("");
        assertEquals(2, StringUtil.countWords("two words"));
        assertEquals(0, StringUtil.countWords("어"));
        StringUtil.setWordCounter("어");
        assertEquals(1, StringUtil.countWords("어"));
        assertEquals(0, StringUtil.countWords("字"));
        StringUtil.setWordCounter("字");
        assertEquals(1, StringUtil.countWords("字"));
        // Make sure the internal WordCounter is restored to FullWordCounter in the end.
    }

    public void testIsWhitespace() {
        assertTrue(StringUtil.isWhitespace(' '));
        assertTrue(StringUtil.isWhitespace('\t'));
        assertTrue(StringUtil.isWhitespace('\n'));
        assertTrue(StringUtil.isWhitespace('\u00a0'));
        assertFalse(StringUtil.isWhitespace('a'));
        assertFalse(StringUtil.isWhitespace('$'));
        assertFalse(StringUtil.isWhitespace('_'));
        assertFalse(StringUtil.isWhitespace('\u0460'));
    }

    public void testIsStringAllWhitespace() {
        assertTrue(StringUtil.isStringAllWhitespace(""));
        assertTrue(StringUtil.isStringAllWhitespace(" \t\r\n"));
        assertTrue(StringUtil.isStringAllWhitespace(" \u00a0     \t\t\t"));
        assertFalse(StringUtil.isStringAllWhitespace("a"));
        assertFalse(StringUtil.isStringAllWhitespace("     a  "));
        assertFalse(StringUtil.isStringAllWhitespace("\u00a0\u0460"));
        assertFalse(StringUtil.isStringAllWhitespace("\n\t_ "));
    }

    public void testFindAndReplace() {
        assertEquals("", StringUtil.findAndReplace("sdf", ".", ""));
        assertEquals("abc", StringUtil.findAndReplace(" a\tb  c ", "\\s", ""));
    }

    private RegExp toRegex(String s) {
        return RegExp.compile(StringUtil.regexEscape(s));
    }

    public void testRegexEscape() {
        assertTrue(toRegex(".*").test(".*"));
        assertFalse(toRegex(".*").test("test"));
        assertFalse(toRegex("[a-z]+").test("az"));
        assertFalse(toRegex("[a-z]+").test("[a-z]"));
        assertTrue(toRegex("[a-z]+").test("[a-z]+"));
        assertTrue(toRegex("\t\n\\\\d[").test("\t\n\\\\d["));
    }

    public void testIsDigit() {
        assertTrue(StringUtil.isDigit('1'));
        assertTrue(StringUtil.isDigit('0'));
        assertFalse(StringUtil.isDigit(' '));
        assertFalse(StringUtil.isDigit('a'));
        assertFalse(StringUtil.isDigit('$'));
        assertFalse(StringUtil.isDigit('_'));
        assertFalse(StringUtil.isDigit('\u0460'));
    }

    public void testIsStringAllDigits() {
        assertTrue(StringUtil.isStringAllDigits("0"));
        assertTrue(StringUtil.isStringAllDigits("018"));
        assertFalse(StringUtil.isStringAllDigits(""));
        assertFalse(StringUtil.isStringAllDigits("a0"));
        assertFalse(StringUtil.isStringAllDigits("0a"));
        assertFalse(StringUtil.isStringAllDigits(" "));
        assertFalse(StringUtil.isStringAllDigits(" 8"));
        assertFalse(StringUtil.isStringAllDigits("8 "));
        assertFalse(StringUtil.isStringAllDigits("'8_"));
        assertFalse(StringUtil.isStringAllDigits("\u00a0\u0460"));
    }

    public void testContainsDigit() {
        assertTrue(StringUtil.containsDigit("0"));
        assertTrue(StringUtil.containsDigit("018"));
        assertTrue(StringUtil.containsDigit("a0"));
        assertTrue(StringUtil.containsDigit("0a"));
        assertTrue(StringUtil.containsDigit(" 8"));
        assertTrue(StringUtil.containsDigit("8 "));
        assertTrue(StringUtil.containsDigit("'8_"));
        assertFalse(StringUtil.containsDigit(""));
        assertFalse(StringUtil.containsDigit(" "));
        assertFalse(StringUtil.containsDigit("\u00a0\u0460"));
        assertFalse(StringUtil.containsDigit("abc"));
        assertFalse(StringUtil.containsDigit("$"));
        assertFalse(StringUtil.containsDigit("_"));
    }

    public void testToNumber() {
        assertEquals(0, StringUtil.toNumber("0"));
        assertEquals(18, StringUtil.toNumber("018"));
        assertEquals(-1, StringUtil.toNumber("a0"));
        assertEquals(-1, StringUtil.toNumber("0a"));
        assertEquals(-1, StringUtil.toNumber(" 8"));
        assertEquals(-1, StringUtil.toNumber("8 "));
        assertEquals(-1, StringUtil.toNumber("'8_"));
        assertEquals(-1, StringUtil.toNumber(""));
        assertEquals(-1, StringUtil.toNumber(" "));
        assertEquals(-1, StringUtil.toNumber("\u00a0\u0460"));
        assertEquals(-1, StringUtil.toNumber("abc"));
        assertEquals(-1, StringUtil.toNumber("$"));
        assertEquals(-1, StringUtil.toNumber("_"));
    }

    public void testJsSplit() {
        assertArrayEquals(new String[]{""}, StringUtil.jsSplit("", ","));
        assertArrayEquals(new String[]{"1", " 2", " 3"}, StringUtil.jsSplit("1, 2, 3", ","));
        assertArrayEquals(new String[]{"1", "2", "3,4"}, StringUtil.jsSplit("1, 2, 3,4", ", "));

        // Separator is not regex in jsSplit:
        assertArrayEquals(new String[]{"1", "2"}, StringUtil.jsSplit("1.*2", ".*"));
        assertArrayEquals(new String[]{}, "1.*2".split(".*"));

        assertArrayEquals(new String[]{"1", "2", ""}, StringUtil.jsSplit("1,2,", ","));
        // Note the different behavior of Java String.split().
        assertArrayEquals(new String[]{"1", "2"}, "1,2,".split(","));

        assertArrayEquals(new String[]{"1", " 2", " "}, StringUtil.jsSplit("1, 2, ", ","));
        // Note the same behavior of Java String.split().
        assertArrayEquals(new String[]{"1", " 2", " "}, "1, 2, ".split(","));
    }

    public void testJoin() {
        assertEquals("", StringUtil.join(new String[]{}, ""));
        assertEquals("abc", StringUtil.join(new String[]{"abc"}, "def"));
        assertEquals("", StringUtil.join(new String[]{}, "def"));
        assertEquals("1, 2, 3", StringUtil.join(new String[]{"1", "2", "3"}, ", "));
        assertEquals("1, 2, ", StringUtil.join(new String[]{"1", "2", ""}, ", "));
        assertEquals(" , ,  ", StringUtil.join(new String[]{" ", "", " "}, ", "));
        assertEquals("123", StringUtil.join(new String[]{"1", "2", "3"}, ""));
        assertEquals("abc123def", StringUtil.join(new String[]{"abc", "def"}, "123"));
    }
}

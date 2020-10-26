// ORIGINAL: javatest/TextDocumentConstructionTest.java

package testutil_test

import (
	"strings"
	"testing"

	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/internal/stringutil"
	"github.com/markusmobius/go-domdistiller/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func Test_TestUtil_TDB_TextDocumentConstruction(t *testing.T) {
	root := testutil.CreateHTML()
	body := dom.QuerySelector(root, "body")
	dom.SetInnerHTML(body, tdbSourceHTML)

	doc := testutil.NewTextDocumentFromPage(root, stringutil.FastWordCounter{}, nil)
	debugString := doc.DebugString()

	left := strings.Split(tdbExpectedDebug, "\n")
	right := strings.Split(debugString, "\n")
	assert.Equal(t, len(left), len(right))

	for i := 0; i < len(left); i++ {
		assert.Equal(t, left[i], right[i], i)
	}

	assert.Equal(t, tdbExpectedDebug, doc.DebugString())
}

const tdbSourceHTML = `` +
	`<!-- ========= START OF TOP NAVBAR ======= -->` +
	`<div class="topNav"><a name="navbar_top">` +
	`<!--   -->` +
	`</a><a href="#skip-navbar_top" title="Skip navigation links"></a><a name="navbar_top_firstrow">` +
	`<!--   -->` +
	`</a>` +
	`<ul class="navList" title="Navigation">` +
	`<li><a href="overview-summary.html">Overview</a></li>` +
	`<li>Package</li>` +
	`<li>Class</li>` +
	`<li>Use</li>` +
	`<li><a href="overview-tree.html">Tree</a></li>` +
	`<li><a href="deprecated-list.html">Deprecated</a></li>` +
	`<li><a href="index-all.html">Index</a></li>` +
	`<li class="navBarCell1Rev">Help</li>` +
	`</ul>` +
	`<div class="aboutLanguage"><em>GWT 2.5.1</em></div>` +
	`</div>` +
	`<div class="subNav">` +
	`<ul class="navList">` +
	`<li>Prev</li>` +
	`<li>Next</li>` +
	`</ul>` +
	`<ul class="navList">` +
	`<li><a href="index.html?help-doc.html" target="_top">Frames</a></li>` +
	`<li><a href="help-doc.html" target="_top">No Frames</a></li>` +
	`</ul>` +
	`<ul class="navList" id="allclasses_navbar_top">` +
	`<li><a href="allclasses-noframe.html">All Classes</a></li>` +
	`</ul>` +
	`<div>` +
	`</div>` +
	`<a name="skip-navbar_top">` +
	`<!--   -->` +
	`</a></div>` +
	`<!-- ========= END OF TOP NAVBAR ========= -->` +
	`<div class="header">` +
	`<h1 class="title">How This API Document Is Organized</h1>` +
	`<div class="subTitle">This API (Application Programming Interface) document has pages corresponding to the items in the navigation bar, described as follows.</div>` +
	`</div>` +
	`<div class="contentContainer">` +
	`<ul class="blockList">` +
	`<li class="blockList">` +
	`<h2>Overview</h2>` +
	`<p>The <a href="overview-summary.html">Overview</a> page is the front page of this API document and provides a list of all packages with a summary for each.  This page can also contain an overall description of the set of packages.</p>` +
	`</li>` +
	`<li class="blockList">` +
	`<h2>Package</h2>` +
	`<p>Each package has a page that contains a list of its classes and interfaces, with a summary for each. This page can contain six categories:</p>` +
	`<ul>` +
	`<li>Interfaces (italic)</li>` +
	`<li>Classes</li>` +
	`<li>Enums</li>` +
	`<li>Exceptions</li>` +
	`<li>Errors</li>` +
	`<li>Annotation Types</li>` +
	`</ul>` +
	`</li>` +
	`<li class="blockList">` +
	`<h2>Class/Interface</h2>` +
	`<p>Each class, interface, nested class and nested interface has its own separate page. Each of these pages has three sections consisting of a class/interface description, summary tables, and detailed member descriptions:</p>` +
	`<ul>` +
	`<li>Class inheritance diagram</li>` +
	`<li>Direct Subclasses</li>` +
	`<li>All Known Subinterfaces</li>` +
	`<li>All Known Implementing Classes</li>` +
	`<li>Class/interface declaration</li>` +
	`<li>Class/interface description</li>` +
	`</ul>` +
	`<ul>` +
	`<li>Nested Class Summary</li>` +
	`<li>Field Summary</li>` +
	`<li>Constructor Summary</li>` +
	`<li>Method Summary</li>` +
	`</ul>` +
	`<ul>` +
	`<li>Field Detail</li>` +
	`<li>Constructor Detail</li>` +
	`<li>Method Detail</li>` +
	`</ul>` +
	`<p>Each summary entry contains the first sentence from the detailed description for that item. The summary entries are alphabetical, while the detailed descriptions are in the order they appear in the source code. This preserves the logical groupings established by the programmer.</p>` +
	`</li>` +
	`<li class="blockList">` +
	`<h2>Annotation Type</h2>` +
	`<p>Each annotation type has its own separate page with the following sections:</p>` +
	`<ul>` +
	`<li>Annotation Type declaration</li>` +
	`<li>Annotation Type description</li>` +
	`<li>Required Element Summary</li>` +
	`<li>Optional Element Summary</li>` +
	`<li>Element Detail</li>` +
	`</ul>` +
	`</li>` +
	`<li class="blockList">` +
	`<h2>Enum</h2>` +
	`<p>Each enum has its own separate page with the following sections:</p>` +
	`<ul>` +
	`<li>Enum declaration</li>` +
	`<li>Enum description</li>` +
	`<li>Enum Constant Summary</li>` +
	`<li>Enum Constant Detail</li>` +
	`</ul>` +
	`</li>` +
	`<li class="blockList">` +
	`<h2>Use</h2>` +
	`<p>Each documented package, class and interface has its own Use page.  This page describes what packages, classes, methods, constructors and fields use any part of the given class or package. Given a class or interface A, its Use page includes subclasses of A, fields declared as A, methods that return A, and methods and constructors with parameters of type A.  You can access this page by first going to the package, class or interface, then clicking on the "Use" link in the navigation bar.</p>` +
	`</li>` +
	`<li class="blockList">` +
	`<h2>Tree (Class Hierarchy)</h2>` +
	`<p>There is a <a href="overview-tree.html">Class Hierarchy</a> page for all packages, plus a hierarchy for each package. Each hierarchy page contains a list of classes and a list of interfaces. The classes are organized by inheritance structure starting with <code>java.lang.Object</code>. The interfaces do not inherit from <code>java.lang.Object</code>.</p>` +
	`<ul>` +
	`<li>When viewing the Overview page, clicking on "Tree" displays the hierarchy for all packages.</li>` +
	`<li>When viewing a particular package, class or interface page, clicking "Tree" displays the hierarchy for only that package.</li>` +
	`</ul>` +
	`</li>` +
	`<li class="blockList">` +
	`<h2>Deprecated API</h2>` +
	`<p>The <a href="deprecated-list.html">Deprecated API</a> page lists all of the API that have been deprecated. A deprecated API is not recommended for use, generally due to improvements, and a replacement API is usually given. Deprecated APIs may be removed in future implementations.</p>` +
	`</li>` +
	`<li class="blockList">` +
	`<h2>Index</h2>` +
	`<p>The <a href="index-all.html">Index</a> contains an alphabetic list of all classes, interfaces, constructors, methods, and fields.</p>` +
	`</li>` +
	`<li class="blockList">` +
	`<h2>Prev/Next</h2>` +
	`<p>These links take you to the next or previous class, interface, package, or related page.</p>` +
	`</li>` +
	`<li class="blockList">` +
	`<h2>Frames/No Frames</h2>` +
	`<p>These links show and hide the HTML frames.  All pages are available with or without frames.</p>` +
	`</li>` +
	`<li class="blockList">` +
	`<h2>All Classes</h2>` +
	`<p>The <a href="allclasses-noframe.html">All Classes</a> link shows all classes and interfaces except non-static nested types.</p>` +
	`</li>` +
	`<li class="blockList">` +
	`<h2>Serialized Form</h2>` +
	`<p>Each serializable or externalizable class has a description of its serialization fields and methods. This information is of interest to re-implementors, not to developers using the API. While there is no link in the navigation bar, you can get to this information by going to any serialized class and clicking "Serialized Form" in the "See also" section of the class description.</p>` +
	`</li>` +
	`<li class="blockList">` +
	`<h2>Constant Field Values</h2>` +
	`<p>The <a href="constant-values.html">Constant Field Values</a> page lists the static final fields and their values.</p>` +
	`</li>` +
	`</ul>` +
	`<em>This help file applies to API documentation generated using the standard doclet.</em></div>` +
	`<!-- ======= START OF BOTTOM NAVBAR ====== -->` +
	`<div class="bottomNav"><a name="navbar_bottom">` +
	`<!--   -->` +
	`</a><a href="#skip-navbar_bottom" title="Skip navigation links"></a><a name="navbar_bottom_firstrow">` +
	`<!--   -->` +
	`</a>` +
	`<ul class="navList" title="Navigation">` +
	`<li><a href="overview-summary.html">Overview</a></li>` +
	`<li>Package</li>` +
	`<li>Class</li>` +
	`<li>Use</li>` +
	`<li><a href="overview-tree.html">Tree</a></li>` +
	`<li><a href="deprecated-list.html">Deprecated</a></li>` +
	`<li><a href="index-all.html">Index</a></li>` +
	`<li class="navBarCell1Rev">Help</li>` +
	`</ul>` +
	`<div class="aboutLanguage"><em>GWT 2.5.1</em></div>` +
	`</div>` +
	`<div class="subNav">` +
	`<ul class="navList">` +
	`<li>Prev</li>` +
	`<li>Next</li>` +
	`</ul>` +
	`<ul class="navList">` +
	`<li><a href="index.html?help-doc.html" target="_top">Frames</a></li>` +
	`<li><a href="help-doc.html" target="_top">No Frames</a></li>` +
	`</ul>` +
	`<ul class="navList" id="allclasses_navbar_bottom">` +
	`<li><a href="allclasses-noframe.html">All Classes</a></li>` +
	`</ul>` +
	`<div>` +
	`</div>` +
	`<a name="skip-navbar_bottom">` +
	`<!--   -->` +
	`</a></div>`

const tdbExpectedDebug = "" +
	"[0/0;tl=6;nw=1;ld=1.000;]	boilerplate,de.l3s.boilerpipe/LI\n" +
	"   Overview \n" +
	"[1/1;tl=5;nw=1;ld=0.000;]	boilerplate,de.l3s.boilerpipe/LI\n" +
	"Package\n" +
	"[2/2;tl=5;nw=1;ld=0.000;]	boilerplate,de.l3s.boilerpipe/LI\n" +
	"Class\n" +
	"[3/3;tl=5;nw=1;ld=0.000;]	boilerplate,de.l3s.boilerpipe/LI\n" +
	"Use\n" +
	"[4/4;tl=6;nw=1;ld=1.000;]	boilerplate,de.l3s.boilerpipe/LI\n" +
	" Tree \n" +
	"[5/5;tl=6;nw=1;ld=1.000;]	boilerplate,de.l3s.boilerpipe/LI\n" +
	" Deprecated \n" +
	"[6/6;tl=6;nw=1;ld=1.000;]	boilerplate,de.l3s.boilerpipe/LI\n" +
	" Index \n" +
	"[7/7;tl=5;nw=1;ld=0.000;]	boilerplate,de.l3s.boilerpipe/LI\n" +
	"Help\n" +
	"[8/8;tl=4;nw=2;ld=0.000;]	boilerplate,\n" +
	"GWT 2.5.1\n" +
	"[9/9;tl=5;nw=1;ld=0.000;]	boilerplate,de.l3s.boilerpipe/LI\n" +
	"Prev\n" +
	"[10/10;tl=5;nw=1;ld=0.000;]	boilerplate,de.l3s.boilerpipe/LI\n" +
	"Next\n" +
	"[11/11;tl=6;nw=1;ld=1.000;]	boilerplate,de.l3s.boilerpipe/LI\n" +
	" Frames \n" +
	"[12/12;tl=6;nw=2;ld=1.000;]	boilerplate,de.l3s.boilerpipe/LI\n" +
	" No Frames \n" +
	"[13/13;tl=6;nw=2;ld=1.000;]	boilerplate,de.l3s.boilerpipe/LI\n" +
	" All Classes \n" +
	"[14/14;tl=4;nw=6;ld=0.000;]	boilerplate,de.l3s.boilerpipe/H1,de.l3s.boilerpipe/HEADING\n" +
	"How This API Document Is Organized\n" +
	"[15/15;tl=4;nw=19;ld=0.000;]	boilerplate,\n" +
	"This API (Application Programming Interface) document has pages corresponding to the items in the navigation bar, described as follows.\n" +
	"[16/16;tl=5;nw=37;ld=0.027;]	boilerplate,de.l3s.boilerpipe/LI\n" +
	"OverviewThe  Overview  page is the front page of this API document and provides a list of all packages with a summary for each.  This page can also contain an overall description of the set of packages.\n" +
	"[17/23;tl=5;nw=34;ld=0.000;]	boilerplate,de.l3s.boilerpipe/LI\n" +
	"PackageEach package has a page that contains a list of its classes and interfaces, with a summary for each. This page can contain six categories:Interfaces (italic)ClassesEnumsExceptionsErrorsAnnotation Types\n" +
	"[24/38;tl=5;nw=105;ld=0.000;]	boilerplate,de.l3s.boilerpipe/LI\n" +
	"Class/InterfaceEach class, interface, nested class and nested interface has its own separate page. Each of these pages has three sections consisting of a class/interface description, summary tables, and detailed member descriptions:Class inheritance diagramDirect SubclassesAll Known SubinterfacesAll Known Implementing ClassesClass/interface declarationClass/interface descriptionNested Class SummaryField SummaryConstructor SummaryMethod SummaryField DetailConstructor DetailMethod DetailEach summary entry contains the first sentence from the detailed description for that item. The summary entries are alphabetical, while the detailed descriptions are in the order they appear in the source code. This preserves the logical groupings established by the programmer.\n" +
	"[39/44;tl=5;nw=28;ld=0.000;]	boilerplate,de.l3s.boilerpipe/LI\n" +
	"Annotation TypeEach annotation type has its own separate page with the following sections:Annotation Type declarationAnnotation Type descriptionRequired Element SummaryOptional Element SummaryElement Detail\n" +
	"[45/49;tl=5;nw=22;ld=0.000;]	boilerplate,de.l3s.boilerpipe/LI\n" +
	"EnumEach enum has its own separate page with the following sections:Enum declarationEnum descriptionEnum Constant SummaryEnum Constant Detail\n" +
	"[50/50;tl=5;nw=85;ld=0.000;]	boilerplate,de.l3s.boilerpipe/LI\n" +
	"UseEach documented package, class and interface has its own Use page.  This page describes what packages, classes, methods, constructors and fields use any part of the given class or package. Given a class or interface A, its Use page includes subclasses of A, fields declared as A, methods that return A, and methods and constructors with parameters of type A.  You can access this page by first going to the package, class or interface, then clicking on the \"Use\" link in the navigation bar.\n" +
	"[51/53;tl=5;nw=80;ld=0.025;]	boilerplate,de.l3s.boilerpipe/LI\n" +
	"Tree (Class Hierarchy)There is a  Class Hierarchy  page for all packages, plus a hierarchy for each package. Each hierarchy page contains a list of classes and a list of interfaces. The classes are organized by inheritance structure starting with java.lang.Object. The interfaces do not inherit from java.lang.Object.When viewing the Overview page, clicking on \"Tree\" displays the hierarchy for all packages.When viewing a particular package, class or interface page, clicking \"Tree\" displays the hierarchy for only that package.\n" +
	"[54/54;tl=5;nw=42;ld=0.048;]	boilerplate,de.l3s.boilerpipe/LI\n" +
	"Deprecated APIThe  Deprecated API  page lists all of the API that have been deprecated. A deprecated API is not recommended for use, generally due to improvements, and a replacement API is usually given. Deprecated APIs may be removed in future implementations.\n" +
	"[55/55;tl=5;nw=15;ld=0.067;]	boilerplate,de.l3s.boilerpipe/LI\n" +
	"IndexThe  Index  contains an alphabetic list of all classes, interfaces, constructors, methods, and fields.\n" +
	"[56/56;tl=5;nw=16;ld=0.000;]	boilerplate,de.l3s.boilerpipe/LI\n" +
	"Prev/NextThese links take you to the next or previous class, interface, package, or related page.\n" +
	"[57/57;tl=5;nw=18;ld=0.000;]	boilerplate,de.l3s.boilerpipe/LI\n" +
	"Frames/No FramesThese links show and hide the HTML frames.  All pages are available with or without frames.\n" +
	"[58/58;tl=5;nw=15;ld=0.133;]	boilerplate,de.l3s.boilerpipe/LI\n" +
	"All ClassesThe  All Classes  link shows all classes and interfaces except non-static nested types.\n" +
	"[59/59;tl=5;nw=63;ld=0.000;]	boilerplate,de.l3s.boilerpipe/LI\n" +
	"Serialized FormEach serializable or externalizable class has a description of its serialization fields and methods. This information is of interest to re-implementors, not to developers using the API. While there is no link in the navigation bar, you can get to this information by going to any serialized class and clicking \"Serialized Form\" in the \"See also\" section of the class description.\n" +
	"[60/60;tl=5;nw=16;ld=0.188;]	boilerplate,de.l3s.boilerpipe/LI\n" +
	"Constant Field ValuesThe  Constant Field Values  page lists the static final fields and their values.\n" +
	"[61/61;tl=3;nw=12;ld=0.000;]	boilerplate,\n" +
	"This help file applies to API documentation generated using the standard doclet.\n" +
	"[62/62;tl=6;nw=1;ld=1.000;]	boilerplate,de.l3s.boilerpipe/LI\n" +
	"   Overview \n" +
	"[63/63;tl=5;nw=1;ld=0.000;]	boilerplate,de.l3s.boilerpipe/LI\n" +
	"Package\n" +
	"[64/64;tl=5;nw=1;ld=0.000;]	boilerplate,de.l3s.boilerpipe/LI\n" +
	"Class\n" +
	"[65/65;tl=5;nw=1;ld=0.000;]	boilerplate,de.l3s.boilerpipe/LI\n" +
	"Use\n" +
	"[66/66;tl=6;nw=1;ld=1.000;]	boilerplate,de.l3s.boilerpipe/LI\n" +
	" Tree \n" +
	"[67/67;tl=6;nw=1;ld=1.000;]	boilerplate,de.l3s.boilerpipe/LI\n" +
	" Deprecated \n" +
	"[68/68;tl=6;nw=1;ld=1.000;]	boilerplate,de.l3s.boilerpipe/LI\n" +
	" Index \n" +
	"[69/69;tl=5;nw=1;ld=0.000;]	boilerplate,de.l3s.boilerpipe/LI\n" +
	"Help\n" +
	"[70/70;tl=4;nw=2;ld=0.000;]	boilerplate,\n" +
	"GWT 2.5.1\n" +
	"[71/71;tl=5;nw=1;ld=0.000;]	boilerplate,de.l3s.boilerpipe/LI\n" +
	"Prev\n" +
	"[72/72;tl=5;nw=1;ld=0.000;]	boilerplate,de.l3s.boilerpipe/LI\n" +
	"Next\n" +
	"[73/73;tl=6;nw=1;ld=1.000;]	boilerplate,de.l3s.boilerpipe/LI\n" +
	" Frames \n" +
	"[74/74;tl=6;nw=2;ld=1.000;]	boilerplate,de.l3s.boilerpipe/LI\n" +
	" No Frames \n" +
	"[75/75;tl=6;nw=2;ld=1.000;]	boilerplate,de.l3s.boilerpipe/LI\n" +
	" All Classes \n"

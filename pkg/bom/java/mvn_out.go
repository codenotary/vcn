/*
 * Copyright (c) 2018-2020 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package java

import "encoding/xml"

type GraphML struct {
	Graph Graph `xml:"graph"`
}

type Graph struct {
	XMLName xml.Name `xml:"graph"`
	Nodes   []Node   `xml:"node"`
}

type Node struct {
	XMLName xml.Name `xml:"node"`
	Data    Data     `xml:"data"`
}

type Data struct {
	XMLName   xml.Name  `xml:"data"`
	ShapeNode ShapeNode `xml:"ShapeNode"`
}

type ShapeNode struct {
	XMLName   xml.Name `xml:"ShapeNode"`
	NodeLabel string   `xml:"NodeLabel"`
}
